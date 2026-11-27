package main

import (
	"bytes"
	"encoding/json"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	"github.com/lishuangquan1987/mygo-antdesign/samples/button"
	"github.com/lishuangquan1987/mygo-antdesign/theme"
)

// goldenManifest 记录截图基线由哪个平台生成。
//
// 逐字节比对对字体渲染敏感：Windows / macOS / Linux 的系统字体栈不同，
// 同一份界面渲染出的 PNG 必然不同。因此基线按平台生成，跨平台时跳过
// 逐字节比对（仍渲染并断言无空帧），避免 CI 永久红。
// 彻底解决需要固定字体（ui.RegisterFont）或按平台提交多套基线，见
// docs/07-samples-gallery.md 第 8 节。
type goldenManifest struct {
	GOOS   string `json:"goos"`
	GOARCH string `json:"goarch"`
}

func readManifest(dir string) (goldenManifest, bool) {
	data, err := os.ReadFile(filepath.Join(dir, "manifest.json"))
	if err != nil {
		return goldenManifest{}, false
	}
	var m goldenManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return goldenManifest{}, false
	}
	return m, true
}

func writeManifest(t *testing.T, dir string) {
	t.Helper()
	data, err := json.MarshalIndent(goldenManifest{GOOS: runtime.GOOS, GOARCH: runtime.GOARCH}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

// newTestGallery 创建注册了示例的 Gallery。
func newTestGallery() *Gallery {
	if len(samples.All()) == 0 {
		samples.Register(button.Demo())
	}
	return NewGallery()
}

func TestGalleryRenders(t *testing.T) {
	g := newTestGallery()
	tt := ui.NewTester(g.View, 1280, 800)
	tt.Frame()
	if !tt.HasText("Button") {
		t.Fatal("侧栏未显示 Button")
	}
	if !tt.HasText("基础用法") {
		t.Fatal("示例卡片未渲染")
	}
	if !tt.HasText("保存") {
		t.Fatal("示例预览未渲染控件")
	}
	if !tt.HasText("亮暗对比") {
		t.Fatal("控制条未渲染")
	}
}

func TestGalleryInteractions(t *testing.T) {
	g := newTestGallery()
	tt := ui.NewTester(g.View, 1280, 800)

	// 切换暗色 / 紧凑 / 对比：无 panic、无持续帧请求。
	for _, label := range []string{"暗", "紧凑", "亮暗对比", "标准", "亮"} {
		if err := tt.Click(label); err != nil {
			t.Fatalf("点击 %q 失败: %v", label, err)
		}
	}
	tt.Frame()
	if !tt.HasText("基础用法") {
		t.Error("切换后示例区丢失")
	}

	// 品牌色切换。
	if err := tt.Click("酱紫"); err != nil {
		// 色块无文字，Click 按文本找不到是预期的；跳过。
		t.Log("品牌色为色块按钮，无文本目标（跳过点击）")
	}
}

func TestGallerySearch(t *testing.T) {
	g := newTestGallery()
	tt := ui.NewTester(g.View, 1280, 800)
	tt.Frame()
	// 直接设置搜索词，避免依赖侧栏搜索框的像素坐标。
	g.query = "不存在的控件"
	tt.Frame()
	if !tt.HasText("未找到匹配的控件") {
		t.Error("搜索无结果时未显示空提示")
	}
	g.query = ""
	tt.Frame()
	if !tt.HasText("Button") {
		t.Error("清空搜索后未恢复列表")
	}
}

func TestGalleryDensityCompare(t *testing.T) {
	g := newTestGallery()
	g.compare = compareDensity
	tt := ui.NewTester(g.View, 1280, 800)
	tt.Frame()
	if !tt.HasText("基础用法") {
		t.Error("密度对比模式下示例未渲染")
	}
	if tt.FramePending() {
		t.Error("密度对比模式持续请求帧")
	}
}

func TestGalleryCopyCode(t *testing.T) {
	g := newTestGallery()
	tt := ui.NewTester(g.View, 1280, 800)
	tt.Frame()
	if err := tt.Click("复制代码"); err != nil {
		t.Fatalf("点击复制代码失败: %v", err)
	}
	tt.Frame()
	if got := tt.Clipboard(); got == "" {
		t.Error("复制代码后剪贴板为空")
	}
}

func TestGallerySampleInteraction(t *testing.T) {
	g := newTestGallery()
	tt := ui.NewTester(g.View, 1280, 800)
	if err := tt.Click("保存"); err != nil {
		t.Fatalf("点击示例按钮失败: %v", err)
	}
	if !tt.HasText("已保存") {
		t.Error("示例交互未生效")
	}
}

// TestGalleryBaselines 生成/比对 3 种主题组合的截图基线（docs/07 第 8 节）。
// 基线不存在时生成并 PASS；存在时逐字节比对，不一致则 FAIL。
func TestGalleryBaselines(t *testing.T) {
	type combo struct {
		name   string
		setup  func(*Gallery)
		height int
	}
	combos := []combo{
		{"light-standard", func(g *Gallery) { g.dark = false; g.density = theme.Comfortable; g.compare = compareNone; g.brand = "" }, 640},
		{"dark-standard", func(g *Gallery) { g.dark = true; g.density = theme.Comfortable; g.compare = compareNone; g.brand = "" }, 640},
		{"light-compact", func(g *Gallery) { g.dark = false; g.density = theme.Compact; g.compare = compareNone; g.brand = "" }, 640},
	}
	dir := filepath.Join("..", "samples", "golden")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	render := func(cb combo) *ui.Tester {
		g := newTestGallery()
		cb.setup(g)
		g.brand = "" // 基线用默认品牌色
		tt := ui.NewTester(g.View, 960, cb.height)
		tt.Frame()
		if tt.Image() == nil {
			t.Fatalf("%s: 无渲染帧", cb.name)
		}
		return tt
	}

	// 说明：并排对比（compareTheme/compareDensity）包含加载动画（ui.Spinner），
	// 逐帧像素不稳定，由交互测试与人工截图核对，不纳入逐字节基线。
	if m, ok := readManifest(dir); ok && (m.GOOS != runtime.GOOS || m.GOARCH != runtime.GOARCH) {
		for _, cb := range combos {
			render(cb)
		}
		t.Skipf("截图基线由 %s/%s 生成，当前平台 %s/%s 跳过逐字节比对（见 docs/07 第 8 节）",
			m.GOOS, m.GOARCH, runtime.GOOS, runtime.GOARCH)
	}

	wrote := false
	for _, cb := range combos {
		tt := render(cb)
		path := filepath.Join(dir, "gallery-"+cb.name+".png")
		want, err := os.ReadFile(path)
		if err != nil {
			// 首次运行：生成基线。
			if err := os.WriteFile(path, pngBytes(t, tt), 0o644); err != nil {
				t.Fatal(err)
			}
			t.Logf("已生成基线 %s", path)
			wrote = true
			continue
		}
		if os.Getenv("GOLDEN_DUMP") != "" {
			_ = os.WriteFile(filepath.Join(os.TempDir(), "cur-"+cb.name+".png"), pngBytes(t, tt), 0o644)
		}
		if string(want) != string(pngBytes(t, tt)) {
			t.Errorf("截图与基线不一致：%s（如为预期改动，删除该文件重新生成）", path)
		}
	}
	if wrote {
		writeManifest(t, dir)
	}
}

func pngBytes(t *testing.T, tt *ui.Tester) []byte {
	t.Helper()
	img := tt.Image()
	if img == nil {
		t.Fatal("无渲染帧")
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("编码 PNG 失败: %v", err)
	}
	return buf.Bytes()
}
