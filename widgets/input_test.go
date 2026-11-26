package ant_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func TestInputRendersAndTypes(t *testing.T) {
	// 渲染不 panic（占位文本是否可见由截图基线覆盖，HasText 不统计
	// 原语内部占位文本）。
	v := ""
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Input(c, &v, widgets.PlaceholderText("请输入姓名"))
	}
	tt := ui.NewTester(view, 320, 120)
	tt.Frame()
	if tt.FramePending() {
		t.Error("输入框持续请求帧")
	}
}

func TestInputOnChangedAndOnSubmit(t *testing.T) {
	changed, submitted := 0, 0
	v := ""
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Input(c, &v,
			widgets.PlaceholderText("姓名"),
			widgets.OnChanged(func() { changed++ }),
			widgets.OnSubmit(func() { submitted++ }))
	}
	tt := ui.NewTester(view, 320, 120)
	tt.Frame()
	// 点击输入框区域聚焦，再键入。
	tt.ClickAt(160, 16)
	tt.Type("go")
	tt.Frame()
	if changed == 0 {
		t.Error("OnChanged 未触发")
	}
	if submitted != 0 {
		t.Errorf("未按 Enter 时 OnSubmit 不应触发")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if submitted != 1 {
		t.Errorf("OnSubmit 触发 %d 次, want 1", submitted)
	}
}

func TestInputDisabled(t *testing.T) {
	changed := 0
	v := ""
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Input(c, &v, widgets.InputDisabled(true),
			widgets.OnChanged(func() { changed++ }))
	}
	tt := ui.NewTester(view, 320, 120)
	tt.Frame()
	tt.ClickAt(160, 16)
	tt.Type("x")
	tt.Frame()
	if changed != 0 {
		t.Error("禁用输入框不应触发 OnChanged")
	}
}

func TestTextAreaAndSearchAndNumber(t *testing.T) {
	bio, query := "", ""
	num := 8080.0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.TextArea(c, &bio, widgets.Rows(3), widgets.PlaceholderText("介绍一下自己"))
		widgets.SearchInput(c, &query, widgets.PlaceholderText("搜索…"))
		widgets.NumberInputV(c, &num, 1, 65535, 1)
	}
	tt := ui.NewTester(view, 360, 280)
	tt.Frame()
	if tt.FramePending() {
		t.Error("录入控件持续请求帧")
	}
	// 数字输入修改值应回调 OnChanged。
	changed := 0
	view2 := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.NumberInputV(c, &num, 0, 100, 1,
			widgets.InputLabel("端口"), widgets.OnChanged(func() { changed++ }))
	}
	tt2 := ui.NewTester(view2, 240, 100)
	tt2.Frame()
	r, ok := tt2.Find("端口")
	if !ok && r.W == 0 {
		t.Fatal("未找到数字输入")
	}
	tt2.ClickAt(r.X+20, r.Y+r.H/2)
	tt2.Key(ui.Cmd, ui.KeyA)
	tt2.Type("42")
	tt2.Frame()
	if changed == 0 {
		t.Error("NumberInput 修改值未触发 OnChanged")
	}
}

func TestInputFollowsThemeAndDensity(t *testing.T) {
	// 4 种主题组合渲染不 panic、不持续请求帧；视觉由截图基线覆盖。
	for _, tk := range []*theme.AntToken{
		theme.Light(),
		theme.Dark(),
		theme.Light(theme.WithDensity(theme.Compact)),
		theme.Dark(theme.WithDensity(theme.Compact), theme.WithPrimary(ui.Hex("#7c3aed"))),
	} {
		token := tk
		v := ""
		n := 1.0
		view := func(c *ui.Context) {
			theme.Apply(c, token)
			widgets.Input(c, &v, widgets.PlaceholderText("标准"))
			widgets.Input(c, &v, widgets.InputSize(widgets.Small), widgets.PasswordInput())
			widgets.TextArea(c, &v, widgets.Rows(2))
			widgets.SearchInput(c, &v)
			widgets.NumberInputV(c, &n, 0, 100, 1)
		}
		tt := ui.NewTester(view, 360, 300)
		tt.Frame()
		if tt.FramePending() {
			t.Errorf("密度 %v 下录入控件持续请求帧", token.Density)
		}
	}
}
