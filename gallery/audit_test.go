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
	"github.com/lishuangquan1987/mygo-antdesign/samples/display"
	"github.com/lishuangquan1987/mygo-antdesign/samples/divider"
	extrademo "github.com/lishuangquan1987/mygo-antdesign/samples/extrademo"
	"github.com/lishuangquan1987/mygo-antdesign/samples/feedback"
	"github.com/lishuangquan1987/mygo-antdesign/samples/ide"
	"github.com/lishuangquan1987/mygo-antdesign/samples/input"
	"github.com/lishuangquan1987/mygo-antdesign/samples/layout"
	"github.com/lishuangquan1987/mygo-antdesign/samples/native"
	"github.com/lishuangquan1987/mygo-antdesign/samples/nav"
	"github.com/lishuangquan1987/mygo-antdesign/samples/selectdemo"
	"github.com/lishuangquan1987/mygo-antdesign/samples/table"
	"github.com/lishuangquan1987/mygo-antdesign/samples/typography"
)

// registerAllDemos 与 gallery/main.go 的注册清单保持一致（新增控件时同步追加）。
func registerAllDemos() {
	if len(samples.All()) > 0 {
		return
	}
	samples.Register(button.Demo())
	samples.Register(typography.Demo())
	samples.Register(divider.Demo())
	samples.Register(layout.Demo())
	samples.Register(input.Demo())
	samples.Register(selectdemo.Demo())
	samples.Register(display.Demo())
	samples.Register(table.Demo())
	samples.Register(feedback.Demo())
	samples.Register(nav.Demo())
	samples.Register(extrademo.Demo())
	samples.Register(ide.Demo())
	samples.Register(native.Demo())
}

// TestGalleryAuditScreenshots 在 GALLERY_CAPTURE=1 时把 Gallery 首页与
// 每个控件页（samples.Pages() 全量）渲染成 PNG，输出到 images/gallery/，
// 供视觉审查（.agents/skills/visual-audit 的工作流 C）。
// 注意必须按 Pages() 的控件页名选页：demo 分组名（如「输入框」）不是页名，
// 用它设置 g.selected 会匹配失败、回退到首页，导致多页截图内容重复。
//
// 用法：
//
//	GALLERY_CAPTURE=1 go test ./gallery -run TestGalleryAuditScreenshots
//
// 可用 GALLERY_CAPTURE_DIR 覆盖输出目录（默认 images/gallery）。
func TestGalleryAuditScreenshots(t *testing.T) {
	if os.Getenv("GALLERY_CAPTURE") == "" {
		t.Skip("设置 GALLERY_CAPTURE=1 以生成审查截图")
	}
	registerAllDemos()

	out := os.Getenv("GALLERY_CAPTURE_DIR")
	if out == "" {
		out = filepath.Join("..", "images", "gallery")
	}
	if err := os.MkdirAll(out, 0o755); err != nil {
		t.Fatal(err)
	}

	type combo struct {
		name string
		dark bool
	}
	combos := []combo{{"light", false}, {"dark", true}}

	width, height := 1280, 800
	pages := 0
	for _, cb := range combos {
		// 首页（未选中任何控件）。
		g := NewGallery()
		g.dark = cb.dark
		g.brand = ""
		pages += capturePage(t, g, "gallery-home-"+cb.name, out, width, height)

		// 每个控件页（samples.Pages 展开后的全部页面）。
		for _, d := range samples.Pages() {
			g := NewGallery()
			g.dark = cb.dark
			g.brand = ""
			g.selected = d.Name
			pages += capturePage(t, g, "gallery-"+d.Name+"-"+cb.name, out, width, height)
		}
	}
	t.Logf("已输出 %d 张 Gallery 审查截图到 %s", pages, out)
}

func capturePage(t *testing.T, g *Gallery, name, out string, w, h int) int {
	t.Helper()
	tt := ui.NewTester(g.View, w, h)
	tt.Frame()
	img := tt.Image()
	if img == nil {
		t.Fatalf("%s: 无渲染帧", name)
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("%s: 编码 PNG 失败: %v", name, err)
	}
	path := filepath.Join(out, name+".png")
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	return 1
}
