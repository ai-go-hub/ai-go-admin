package crud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
)

// 上游流式请求配置
const (
	// streamMaxAttempts 总尝试次数（含首次请求）
	streamMaxAttempts = 3
	// streamRetryBase 每次重试等待的基础时长（按尝试次数递增）
	streamRetryBase = 500 * time.Millisecond
)

// streamHTTPClient 上游流式请求客户端
// 不设置总超时，避免长连接流式输出中途被掐断；仅控制建连、响应头等阶段超时
var streamHTTPClient = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ResponseHeaderTimeout: 30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		IdleConnTimeout:       90 * time.Second,
	},
}

// Stream AI 对话流式输出
// 将对话上下文转发到上游 OpenAI Responses 兼容接口，返回响应体供逐段读取（SSE）
// 仅在尚未开始流式输出前对网络错误、5xx、429 重试；一旦开始输出则不再重试
func (s *Service) Stream(ctx context.Context, req *dto.ChatRequest) (io.ReadCloser, error) {
	cfg, err := s.configRepo.GetConfigsByGroup(ctx, "ai")
	if err != nil {
		return nil, err
	}

	url := strings.TrimSpace(cfg["ai_api_url"])
	if url == "" {
		return nil, errors.New("系统配置 > AI 配置的 API URL 未配置")
	}
	key := cfg["ai_api_key"]
	if key == "" {
		return nil, errors.New("系统配置 > AI 配置的 API Key 未配置")
	}
	model := req.Model
	if model == "" {
		return nil, errors.New("缺少模型")
	}
	if len(req.Messages) == 0 {
		return nil, errors.New("缺少对话内容")
	}

	body := map[string]any{
		"model":  model,
		"input":  req.Messages,
		"stream": true,
	}
	if req.Temperature != nil {
		body["temperature"] = *req.Temperature
	}
	if req.TopP != nil {
		body["top_p"] = *req.TopP
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := streamWithRetry(ctx, url, key, payload)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("上游接口错误 %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	return resp.Body, nil
}

// streamWithRetry 带重试的上游流式请求
// 重试条件: 网络错误（未取得响应头）、HTTP 5xx、429; 4xx 为请求错误不重试
func streamWithRetry(ctx context.Context, url, key string, payload []byte) (*http.Response, error) {
	var lastErr error
	for attempt := range streamMaxAttempts {
		if attempt > 0 {
			select {
			case <-time.After(streamRetryBase * time.Duration(attempt)):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+key)
		req.Header.Set("Content-Type", "application/json")

		resp, err := streamHTTPClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests {
			_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
			_ = resp.Body.Close()
			lastErr = fmt.Errorf("上游响应 %d", resp.StatusCode)
			continue
		}
		return resp, nil
	}
	return nil, lastErr
}
