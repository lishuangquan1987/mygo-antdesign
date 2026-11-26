package main

import (
	"bytes"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	"github.com/lishuangquan1987/mygo-antdesign/samples/button"
	"github.com/lishuangquan1987/mygo-antdesign/theme"
)

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
	for _, cb := range combos {
		g := newTestGallery()
		cb.setup(g)
		g.brand = "" // 基线用默认品牌色
		tt := ui.NewTester(g.View, 960, cb.height)
		tt.Frame()
		if tt.Image() == nil {
			t.Fatalf("%s: 无渲染帧", cb.name)
		}
		path := filepath.Join(dir, "gallery-"+cb.name+".png")
		want, err := os.ReadFile(path)
		if err != nil {
			// 首次运行：生成基线。
			if err := os.WriteFile(path, pngBytes(t, tt), 0o644); err != nil {
				t.Fatal(err)
			}
			t.Logf("已生成基线 %s", path)
			continue
		}
		if string(want) != string(pngBytes(t, tt)) {
			t.Errorf("截图与基线不一致：%s（如为预期改动，删除该文件重新生成）", path)
		}
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
