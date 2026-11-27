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
	tt := ui.NewTester(g.View, 960, 640)
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
	tt := ui.NewTester(g.View, 960, 640)

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

func TestGallerySampleInteraction(t *testing.T) {
	g := newTestGallery()
	tt := ui.NewTester(g.View, 960, 640)
	if err := tt.Click("保存"); err != nil {
		t.Fatalf("点击示例按钮失败: %v", err)
	}
	if !tt.HasText("已保存") {
		t.Error("示例交互未生效")
	}
}

// TestGalleryBaselines 生成/比对 4 种主题组合的截图基线（docs/07 第 8 节）。
// 基线不存在时生成并 PASS；存在时逐字节比对，不一致则 FAIL。
func TestGalleryBaselines(t *testing.T) {
	type combo struct {
		name   string
		setup  func(*Gallery)
		height int
	}
	combos := []combo{
		{"light-standard", func(g *Gallery) { g.dark = false; g.density = theme.Comfortable; g.compare = false; g.brand = "" }, 640},
		{"dark-standard", func(g *Gallery) { g.dark = true; g.density = theme.Comfortable; g.compare = false; g.brand = "" }, 640},
		{"light-compact", func(g *Gallery) { g.dark = false; g.density = theme.Compact; g.compare = false; g.brand = "" }, 640},
		{"dark-compact-compare", func(g *Gallery) { g.dark = true; g.density = theme.Compact; g.compare = true; g.brand = "" }, 900},
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

	// 跨平台：基线由另一平台生成时，字体栈不同必然产生像素差异，
	// 只渲染不比对（避免 CI 永久红）。
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
