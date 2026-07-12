package captcha

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/config"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/database"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/pkg/filesystem"
	"github.com/ai-go-hub/ai-go-admin/pkg/random"

	"github.com/google/uuid"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"gorm.io/gorm"
)

// ==================== 初始化 ====================

var (
	bootstrapOnce sync.Once
	captchaCfg    config.CaptchaConfig
	chineseRunes  []rune
	uppercasePool []rune
	icons         []iconMeta
	fontOnce      sync.Once
	fontFace      font.Face
	fontErr       error
)

// bootstrap 引导初始化，Create 每次调用均执行过期清理，首次额外加载图标与中文字符池
func bootstrap() {
	cleanExpired()

	bootstrapOnce.Do(func() {
		pool := "们以我到他会作时要动国产的是工就年阶义发成部民可出能方进在和有大这主中为来分生对于学级地用同行面说种过度革而多子后自社加小机也经力线本电高量长党得实家定深法表着水理化争现所起好十战无农使前等反合斗路图把结第里正新开论之物从当两些还天资事队点重其思与间内去因件利相由压员气业代全组数果期导平各基或月然如应形想制心样都向变关问比展那它最及外没看治提五解系林者米群头意只明四道马认次文通但条较克又公孔领军流接席位情运器并飞原油放立题质指建区验活众很教决特此常石强极已根共直团统式转别造九你取西持总料连任志观调么山程百报更见必真保热委手改管处己将修支识象先老光专什六型具示复安带每东增则完风回南劳轮科北打积车计给节做务被整联步类集号列温装即毫知轴研单坚据速防史拉世设达尔场织历花求传断况采精金界品判参层止边清至万确究书术状须离再目海权且青才证低越际八试规斯近注办布门铁需走议县兵固除般引齿胜细影济白格效置推空配叶率述今选养德话查差半敌始片施响收华觉备名红续均药标记难存测身派准斤角降维板许破述技消底床田势端感往便贺村构照容非亚磨族段算适讲按值美态易彪服早班麦削信排台声该击素张密侯何树肥继右属市严径螺检左页抗苏显苦英快称移巴材省黑武培著河帝仅针怎植京助升王眼她苗副杂普谈围食源例致酸旧却充足短划剂宣环落首尺波承粉践府鱼随考刻靠够满夫失包住促枝局菌杆周护岩师举曲春元超负砂封换太模贫减阳扬江析亩木言球朝医校古呢稻宋听唯输滑站另卫字鼓刚写刘微略范供阿块某功友限项余倒卷创律雨让骨远帮初皮播优占圈伟季训激找叫云互跟粮粒练塞钢顶策双留误础阻故寸盾晚丝女散焊功株亲院冷彻弹错散商视艺版烈零室轻倍缺厘泵察绝富城冲壤简否柱李望盘磁雄似困巩益洲脱投送侧润盖挥距触星松送获兴独官混纪依未突架宽冬章偏纹吃执阀矿寨责熟稳夺价努翻奇甲预职评读背协损棉侵灰虽矛厚罗泥辟告箱掌氧恩爱停曾溶营终纲孟钱待尽俄缩沙退陈讨奋械载胞旋征槽倒握担仍呀鲜吧卡介钻逐弱脚怕盐末丰雾冠丙街莱贝辐付吉渗瑞惊顿挤秒悬姆森糖圣凹陶词迟蚕亿矩康遵牧遭幅园腔订香肉屋敏恢忘编印蜂急拿扩飞露核缘游振操央伍域甚迅辉异序免纸夜乡久隶念兰映沟乙吗儒汽磷艰晶埃燃欢铁补咱芽永瓦倾阵碳演威附牙芽永瓦斜灌欧献顺猪洋腐请透司括脉宜笑若尾束壮暴企菜穗楚汉愈绿拖牛份染既秋遍锻玉夏疗尖井费州访吹荣铜沿替滚客召旱悟刺脑措贯藏敢令隙炉壳硫煤迎铸粘探临薄旬善福纵择礼愿伏残雷延烟句纯渐耕跑泽慢栽鲁赤繁境潮横掉锥希池败船假亮谓托伙哲怀摆贡呈劲财仪沉炼麻祖息车穿货销齐鼠抽画饲龙库守筑房歌寒喜哥洗蚀废纳腹乎录镜脂庄擦险赞钟摇典柄辩竹谷乱虚桥奥伯赶垂途额壁网截野遗静谋弄挂课镇妄盛耐扎虑键归符庆聚绕摩忙舞遇索顾胶羊湖钉仁音迹碎伸灯避泛答勇频皇柳哈揭甘诺概宪浓岛袭谁洪谢炮浇斑讯懂灵蛋闭孩释巨徒私银伊景坦累匀霉杜乐勒隔弯绩招绍胡呼峰零柴簧午跳居尚秦稍追梁折耗碱殊岗挖氏刃剧堆赫荷胸衡勤膜篇登驻案刊秧缓凸役剪川雪链渔啦脸户洛孢勃盟买杨宗焦赛旗滤硅炭股坐蒸凝竟枪黎救冒暗洞犯筒您宋弧爆谬涂味津臂障褐陆啊健尊豆拔莫抵桑坡缝警挑冰柬嘴啥饭塑寄赵喊垫丹渡耳虎笔稀昆浪萨茶滴浅拥覆吨浸袖珠雌妈紫戏塔锤震岁貌洁锋疑霸闪埔猛诉刷忽闹乔唐漏闻沈熔氯荒凡抢像浆旁玻亦忠唱蒙予纷捕锁尤乘乌智淡允叛畜俘摸锈扫毕璃宝芯爷鉴秘净蒋钙肩腾枯抛轨堂拌爸循诱祝励肯酒绳塘燥袋朗喂铝软渠颗惯贸综墙趋彼届墨碍启逆卸航孙龄岭休借"
		chineseRunes = []rune(pool)
		uppercasePool = []rune("ABCDEFGHJKMNPQRSTUVWXYZ")

		captchaCfg = config.Get().Captcha
		for _, ic := range captchaCfg.Icons {
			icons = append(icons, iconMeta{
				EnName: ic.EnName,
				CnName: ic.CnName,
				Path:   captchaCfg.IconDir + "/" + ic.EnName + ".png",
			})
		}
	})
}

// ==================== 数据结构 ====================

// Result 创建验证码的返回结果
type Result struct {
	Key         string   `json:"key"`
	Elements    []string `json:"elements"`
	ImageWidth  int      `json:"image_width"`
	ImageHeight int      `json:"image_height"`
	ImageBase64 string   `json:"image_base64"`
}

// Request 用户点击验证请求
type Request struct {
	Key            string       `json:"key"`
	Clicks         []ClickPoint `json:"clicks"`
	RenderedWidth  int          `json:"rendered_width"`
	RenderedHeight int          `json:"rendered_height"`
}

// ClickPoint 用户点击坐标
type ClickPoint struct {
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Element string `json:"element"`
}

// storedInfo 存入数据库 Info 字段
type storedInfo struct {
	ImageWidth  int           `json:"image_width"`
	ImageHeight int           `json:"image_height"`
	Elements    []elementInfo `json:"elements"`
}

type elementInfo struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
	W     int    `json:"w"`
	H     int    `json:"h"`
}

type elementType string

const (
	elemChinese   elementType = "chinese"
	elemUppercase elementType = "uppercase"
	elemIcon      elementType = "icon"
)

type captchaElement struct {
	Type  elementType
	Value string // 原始值，icon 为英文文件名，chinese/uppercase 为文字本身
}

type iconMeta struct {
	EnName string // 英文名
	CnName string // 中文名
	Path   string // 完整文件路径
}

// ==================== 公开方法 ====================

// Create 创建点选验证码
func Create() (*Result, error) {
	bootstrap()

	bgPaths, err := filesystem.ReadDir(captchaCfg.BackgroundDir)
	if err != nil {
		return nil, fmt.Errorf("read background dir: %w", err)
	}
	if len(bgPaths) == 0 {
		return nil, fmt.Errorf("no background images in %s", captchaCfg.BackgroundDir)
	}

	bgPath := bgPaths[random.Int(0, len(bgPaths))]
	bg, err := loadImage(bgPath)
	if err != nil {
		return nil, fmt.Errorf("load background: %w", err)
	}

	bounds := bg.Bounds()
	bgW, bgH := bounds.Dx(), bounds.Dy()

	allElements := generateElements(captchaCfg.Length+captchaCfg.ConfusionLength, captchaCfg.Elements)
	correctElements := allElements[:captchaCfg.Length]
	shuffleElements(allElements)

	var placedRects []image.Rectangle
	elementBounds := make([]elementInfo, 0, len(allElements))
	for _, elem := range allElements {
		info, rect, drawErr := drawElement(bg, elem, bgW, bgH, placedRects)
		if drawErr != nil {
			return nil, fmt.Errorf("draw %s(%s): %w", elem.Type, elem.Value, drawErr)
		}
		elementBounds = append(elementBounds, info)
		placedRects = append(placedRects, rect)
	}

	correctValues := make([]string, len(correctElements))
	for i, e := range correctElements {
		correctValues[i] = formatElement(e)
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, bg); err != nil {
		return nil, fmt.Errorf("encode png: %w", err)
	}

	infoJSON, _ := json.Marshal(storedInfo{ImageWidth: bgW, ImageHeight: bgH, Elements: elementBounds})
	codeJSON, _ := json.Marshal(correctValues)
	codeHash := md5Hex(codeJSON)

	record := &model.Captcha{
		Key:       uuid.New().String(),
		Code:      codeHash,
		Info:      string(infoJSON),
		ExpiredAt: time.Now().Add(time.Duration(captchaCfg.Expire) * time.Second),
		CreatedAt: time.Now(),
	}

	if err := gorm.G[model.Captcha](database.DB()).Create(context.Background(), record); err != nil {
		return nil, fmt.Errorf("store captcha: %w", err)
	}

	return &Result{
		Key:         record.Key,
		Elements:    correctValues,
		ImageWidth:  bgW,
		ImageHeight: bgH,
		ImageBase64: "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()),
	}, nil
}

// Check 检查点选验证码，deleteOnSuccess 控制验证成功后是否删除记录
func Check(req Request, deleteOnSuccess bool) (bool, error) {
	tx := gorm.G[model.Captcha](database.DB())

	record, err := tx.Where("key = ?", req.Key).First(context.Background())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, fmt.Errorf("验证码不存在或已过期")
		}
		return false, fmt.Errorf("查询验证码: %w", err)
	}

	if time.Now().After(record.ExpiredAt) {
		tx.Where("key = ?", record.Key).Delete(context.Background())
		return false, fmt.Errorf("验证码已过期")
	}

	userValues := make([]string, len(req.Clicks))
	for i, c := range req.Clicks {
		userValues[i] = c.Element
	}
	userJSON, _ := json.Marshal(userValues)
	if md5Hex(userJSON) != record.Code {
		return false, fmt.Errorf("点击元素不正确")
	}

	var info storedInfo
	if err := json.Unmarshal([]byte(record.Info), &info); err != nil {
		return false, fmt.Errorf("解析验证码坐标: %w", err)
	}

	scaleX := float64(req.RenderedWidth) / float64(info.ImageWidth)
	scaleY := float64(req.RenderedHeight) / float64(info.ImageHeight)

	const tolerance = 0
	for i, click := range req.Clicks {
		ei := findElement(info.Elements, click.Element)
		if ei == nil {
			return false, fmt.Errorf("未找到元素 %s 的坐标信息", click.Element)
		}

		ex := int(float64(ei.X) * scaleX)
		ey := int(float64(ei.Y) * scaleY)
		r := tolerance + max(int(float64(ei.W)*scaleX), int(float64(ei.H)*scaleY))/2

		dx := click.X - ex
		dy := click.Y - ey
		if dx*dx+dy*dy > r*r {
			return false, fmt.Errorf("第 %d 个元素点击位置不正确", i+1)
		}
	}

	if deleteOnSuccess {
		tx.Where("key = ?", record.Key).Delete(context.Background())
	}

	return true, nil
}

// ==================== 内部函数 ====================

func cleanExpired() {
	gorm.G[model.Captcha](database.DB()).Where("expired_at < ?", time.Now()).Delete(context.Background())
}

func generateElements(count int, elementTypes []string) []captchaElement {
	available := make([]elementType, 0, len(elementTypes))
	for _, e := range elementTypes {
		available = append(available, elementType(e))
	}

	seen := make(map[string]bool, count)
	elements := make([]captchaElement, 0, count)
	for range count {
		for range 50 {
			t := available[random.Int(0, len(available))]
			elem := captchaElement{Type: t}
			switch t {
			case elemChinese:
				elem.Value = randomChinese()
			case elemUppercase:
				elem.Value = randomUppercase()
			case elemIcon:
				idx := random.Int(0, len(icons))
				elem.Value = icons[idx].EnName
			}
			if !seen[elem.Value] {
				seen[elem.Value] = true
				elements = append(elements, elem)
				break
			}
		}
	}
	return elements
}

func randomChinese() string {
	return string(chineseRunes[random.Int(0, len(chineseRunes))])
}

func randomUppercase() string {
	return string(uppercasePool[random.Int(0, len(uppercasePool))])
}

func formatElement(e captchaElement) string {
	switch e.Type {
	case elemIcon:
		for _, ic := range icons {
			if ic.EnName == e.Value {
				return "<" + ic.CnName + ">"
			}
		}
		return "<" + e.Value + ">"
	default:
		return "<" + e.Value + ">"
	}
}

func loadImage(path string) (*image.RGBA, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, rgba.Bounds(), img, bounds.Min, draw.Src)
	return rgba, nil
}

// ==================== 绘制元素（含碰撞检测） ====================

const (
	collisionPadding = 15
	maxPlaceAttempts = 30
)

func drawElement(bg *image.RGBA, elem captchaElement, bgW, bgH int, placed []image.Rectangle) (elementInfo, image.Rectangle, error) {
	switch elem.Type {
	case elemChinese, elemUppercase:
		return drawText(bg, elem, bgW, bgH, placed)
	case elemIcon:
		return drawIcon(bg, elem, bgW, bgH, placed)
	}
	return elementInfo{}, image.Rectangle{}, fmt.Errorf("unknown type: %s", elem.Type)
}

func drawText(bg *image.RGBA, elem captchaElement, bgW, bgH int, placed []image.Rectangle) (elementInfo, image.Rectangle, error) {
	face, err := getFontFace()
	if err != nil {
		return elementInfo{}, image.Rectangle{}, err
	}

	bounds, _ := font.BoundString(face, elem.Value)
	tw := (bounds.Max.X - bounds.Min.X).Ceil()
	th := (bounds.Max.Y - bounds.Min.Y).Ceil()

	if tw >= bgW || th >= bgH {
		return elementInfo{}, image.Rectangle{}, fmt.Errorf("text too large for background")
	}

	var x, y int
	var rect image.Rectangle
	ok := false
	for range maxPlaceAttempts {
		x = 5 + random.Int(0, bgW-tw-10)
		y = th + random.Int(0, bgH-th-10)

		rect = image.Rect(
			x+bounds.Min.X.Ceil(),
			y+bounds.Min.Y.Ceil(),
			x+bounds.Max.X.Ceil(),
			y+bounds.Max.Y.Ceil(),
		)

		if !hasCollision(placed, rect, collisionPadding) {
			ok = true
			break
		}
	}
	if !ok {
		return elementInfo{}, image.Rectangle{}, fmt.Errorf("cannot place text without collision after %d attempts", maxPlaceAttempts)
	}

	drawer := &font.Drawer{
		Dst:  bg,
		Src:  image.NewUniform(random.RGB("light")),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	drawer.DrawString(elem.Value)

	info := elementInfo{
		Type:  string(elem.Type),
		Value: formatElement(elem),
		X:     rect.Min.X + rect.Dx()/2,
		Y:     rect.Min.Y + rect.Dy()/2,
		W:     rect.Dx(),
		H:     rect.Dy(),
	}
	return info, rect, nil
}

func drawIcon(bg *image.RGBA, elem captchaElement, bgW, bgH int, placed []image.Rectangle) (elementInfo, image.Rectangle, error) {
	ic := findIcon(elem.Value)
	if ic == nil {
		return elementInfo{}, image.Rectangle{}, fmt.Errorf("icon not found: %s", elem.Value)
	}

	f, err := os.Open(ic.Path)
	if err != nil {
		return elementInfo{}, image.Rectangle{}, err
	}
	defer f.Close()

	iconImg, err := png.Decode(f)
	if err != nil {
		return elementInfo{}, image.Rectangle{}, err
	}

	iw := iconImg.Bounds().Dx()
	ih := iconImg.Bounds().Dy()

	if iw >= bgW || ih >= bgH {
		return elementInfo{}, image.Rectangle{}, fmt.Errorf("icon too large for background")
	}

	var x, y int
	var rect image.Rectangle
	ok := false
	for range maxPlaceAttempts {
		x = 5 + random.Int(0, bgW-iw-10)
		y = 5 + random.Int(0, bgH-ih-10)
		rect = image.Rect(x, y, x+iw, y+ih)

		if !hasCollision(placed, rect, collisionPadding) {
			ok = true
			break
		}
	}
	if !ok {
		return elementInfo{}, image.Rectangle{}, fmt.Errorf("cannot place icon without collision after %d attempts", maxPlaceAttempts)
	}

	draw.Draw(bg, rect, iconImg, iconImg.Bounds().Min, draw.Over)

	info := elementInfo{
		Type:  string(elem.Type),
		Value: formatElement(elem),
		X:     rect.Min.X + rect.Dx()/2,
		Y:     rect.Min.Y + rect.Dy()/2,
		W:     rect.Dx(),
		H:     rect.Dy(),
	}
	return info, rect, nil
}

func hasCollision(placed []image.Rectangle, r image.Rectangle, padding int) bool {
	padded := r.Inset(-padding)
	return slices.ContainsFunc(placed, padded.Overlaps)
}

func getFontFace() (font.Face, error) {
	fontOnce.Do(func() {
		fontData, err := os.ReadFile(captchaCfg.FontPath)
		if err != nil {
			fontErr = fmt.Errorf("read font: %w", err)
			return
		}
		f, err := opentype.Parse(fontData)
		if err != nil {
			fontErr = fmt.Errorf("parse font: %w", err)
			return
		}
		fontFace, fontErr = opentype.NewFace(f, &opentype.FaceOptions{
			Size:    28,
			DPI:     72,
			Hinting: font.HintingFull,
		})
	})
	return fontFace, fontErr
}

func findElement(elements []elementInfo, value string) *elementInfo {
	for i := range elements {
		if elements[i].Value == value {
			return &elements[i]
		}
	}
	return nil
}

func findIcon(name string) *iconMeta {
	for i := range icons {
		if icons[i].EnName == name {
			return &icons[i]
		}
	}
	return nil
}

func shuffleElements(elements []captchaElement) {
	for i := len(elements) - 1; i > 0; i-- {
		j := random.Int(0, i+1)
		elements[i], elements[j] = elements[j], elements[i]
	}
}

func md5Hex(data []byte) string {
	h := md5.Sum(data)
	h = md5.Sum([]byte(fmt.Sprintf("%x", h) + "ai-go-click-captcha"))
	return fmt.Sprintf("%x", h)
}
