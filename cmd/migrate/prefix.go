package migrate

import (
	"io"
	"strings"

	"github.com/golang-migrate/migrate/v4/source"
)

// prefixDriver 包装 source.Driver，在 ReadUp / ReadDown 返回的 SQL 中将
// __PREFIX__ 占位符替换为 `config/config.yaml` 配置的表前缀
type prefixDriver struct {
	source.Driver
	prefix string
}

func (d *prefixDriver) ReadUp(version uint) (io.ReadCloser, string, error) {
	r, id, err := d.Driver.ReadUp(version)
	if err != nil {
		return nil, "", err
	}
	return d.replace(r), id, nil
}

func (d *prefixDriver) ReadDown(version uint) (io.ReadCloser, string, error) {
	r, id, err := d.Driver.ReadDown(version)
	if err != nil {
		return nil, "", err
	}
	return d.replace(r), id, nil
}

func (d *prefixDriver) replace(r io.ReadCloser) io.ReadCloser {
	data, _ := io.ReadAll(r)
	r.Close()
	replaced := strings.ReplaceAll(string(data), "__PREFIX__", d.prefix)
	return io.NopCloser(strings.NewReader(replaced))
}
