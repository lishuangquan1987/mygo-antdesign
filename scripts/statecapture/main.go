// Command statecapture 把控件在各交互状态下渲染成 PNG，输出到 images/，
// 供逐状态与 antd 规范比对（配合 .workbuddy/skills/visual-audit 使用）。
//
// 每个变体输出：默认态（<file>.png），并可模拟交互输出
// hover / pressed / focus 三张（<file>-hover.png 等）。
// 禁用、加载中等静态态不模拟交互，只有默认态一张。
//
// 用法：
//
//	go run ./scripts/statecapture                 # 全部控件，亮色主题，输出 images/
//	go run ./scripts/statecapture -dark           # 暗色主题
//	go run ./scripts/statecapture -out .audit     # 自定义输出目录
//
// 新增控件：在 specs() 中追加一组 spec（仿照 buttonSpecs 的写法），
// 并在本文件头部注释登记，保持与 widgets/ 的控件清单同步。
package main

import (
	"bytes"
	"flag"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

// spec 描述一个「控件变体」的截图任务。
type spec struct {
	file  string            // 文件名主干，如 button-primary
	label string            // 交互定位文本（hover/pressed/focus 用）；空串表示只截默认态
	view  func(*ui.Context) // 构造该变体（创建即挂载，与 samples 写法一致）
}

// specs 当前覆盖的控件清单。新增控件时在此追加。
func specs() []spec {
	list := []spec{}
	list = append(list, buttonSpecs()...)
	list = append(list, inputSpecs()...)
	return list
}

// buttonSpecs Button：五种 type × danger，以及加载中/禁用静态态。
func buttonSpecs() []spec {
	label := "保存"
	mk := func(name string, opts ...widgets.ButtonOption) spec {
		return spec{
			file:  "button-" + name,
			label: label,
			view: func(c *ui.Context) {
				widgets.Button(c, label, opts...)
			},
		}
	}
	// static 用于禁用/加载中：不模拟 hover/pressed/focus。
	static := func(name string, opts ...widgets.ButtonOption) spec {
		s := mk(name, opts...)
		s.label = ""
		return s
	}
	return []spec{
		mk("default"),
		mk("primary", widgets.Primary()),
		mk("dashed", widgets.Dashed()),
		mk("text", widgets.Text()),
		mk("link", widgets.Link()),
		mk("danger", widgets.Danger()),
		mk("primary-danger", widgets.Primary(), widgets.Danger()),
		mk("text-danger", widgets.Text(), widgets.Danger()),
		static("primary-loading", widgets.Primary(), widgets.Loading(true)),
		static("primary-disabled", widgets.Primary(), widgets.Disabled(true)),
		static("default-disabled", widgets.Disabled(true)),
	}
}

// inputSpecs Input：默认、禁用（label 作为布局锚点，placeholder 作为交互探针）。
func inputSpecs() []spec {
	view := func(c *ui.Context) {
		v := ""
		widgets.Input(c, &v, widgets.InputLabel("邮箱"), widgets.PlaceholderText("请输入邮箱"))
	}
	disabled := func(c *ui.Context) {
		v := ""
		widgets.Input(c, &v, widgets.InputLabel("邮箱"), widgets.PlaceholderText("禁用状态"), widgets.InputDisabled(true))
	}
	return []spec{
		{file: "input-default", label: "请输入邮箱", view: view},
		{file: "input-disabled", label: "", view: disabled},
	}
}

func main() {
	out := flag.String("out", "images", "输出目录")
	dark := flag.Bool("dark", false, "使用暗色主题")
	flag.Parse()

	tk := theme.Light()
	if *dark {
		tk = theme.Dark()
	}
	if err := os.MkdirAll(*out, 0o755); err != nil {
		log.Fatal(err)
	}

	n := 0
	for _, s := range specs() {
		n += capture(s, tk, *out)
	}
	log.Printf("已输出 %d 张状态截图到 %s（theme=%s）", n, *out, map[bool]string{false: "light", true: "dark"}[*dark])
}

// capture 渲染一个变体：默认态必出；label 非空时模拟 hover/pressed/focus。
func capture(s spec, tk *theme.AntToken, out string) int {
	view := func(c *ui.Context) {
		theme.Apply(c, tk)
		ui.Row(c).Padding(48).AlignItems(ui.Center).Children(func() {
			s.view(c)
		})
	}
	tt := ui.NewTester(view, 480, 240)
	tt.Frame()
	if tt.Image() == nil {
		log.Fatalf("%s: 无渲染帧", s.file)
	}
	save(tt, filepath.Join(out, s.file+".png"))
	n := 1

	if s.label == "" {
		return n
	}
	r, ok := tt.Find(s.label)
	if !ok {
		log.Printf("[跳过交互态] %s: 未找到文本 %q", s.file, s.label)
		return n
	}
	x, y := r.X+r.W/2, r.Y+r.H/2

	tt.Move(x, y)
	tt.Frame()
	save(tt, filepath.Join(out, s.file+"-hover.png"))
	n++

	tt.Press(x, y)
	tt.Frame()
	save(tt, filepath.Join(out, s.file+"-pressed.png"))
	n++
	tt.CancelPress()

	tt.ClickAt(x, y)
	tt.Frame()
	save(tt, filepath.Join(out, s.file+"-focus.png"))
	n++

	return n
}

func save(tt *ui.Tester, path string) {
	img := tt.Image()
	if img == nil {
		log.Fatalf("%s: 无渲染帧", path)
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		log.Fatalf("%s: 编码 PNG 失败: %v", path, err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		log.Fatalf("%s: 写入失败: %v", path, err)
	}
	log.Printf("已输出 %s", path)
}
