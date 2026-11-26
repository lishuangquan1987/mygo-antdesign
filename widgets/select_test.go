package ant_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func TestCheckboxToggleAndChange(t *testing.T) {
	agree, fired := false, 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Checkbox(c, &agree, "同意协议", widgets.OnChange(func() { fired++ }))
	}
	tt := ui.NewTester(view, 320, 120)
	tt.Frame()
	if err := tt.Click("同意协议"); err != nil {
		t.Fatalf("点击失败: %v", err)
	}
	tt.Frame()
	if !agree {
		t.Error("点击后未选中")
	}
	if fired != 1 {
		t.Errorf("OnChange 触发 %d 次, want 1", fired)
	}
	if err := tt.Click("同意协议"); err != nil {
		t.Fatalf("再次点击失败: %v", err)
	}
	tt.Frame()
	if agree {
		t.Error("再次点击后未取消选中")
	}
}

func TestCheckboxDisabled(t *testing.T) {
	agree, fired := false, 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Checkbox(c, &agree, "锁定", widgets.ControlDisabled(true), widgets.OnChange(func() { fired++ }))
	}
	tt := ui.NewTester(view, 320, 120)
	tt.Frame()
	_ = tt.Click("锁定")
	tt.Frame()
	if agree || fired != 0 {
		t.Error("禁用复选框不应改变状态或触发回调")
	}
}

func TestRadioGroupSelection(t *testing.T) {
	size := "middle"
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.RadioGroup(c, func() {
			widgets.Radio(c, &size, "small", "小")
			widgets.Radio(c, &size, "middle", "中")
			widgets.Radio(c, &size, "large", "大")
		})
	}
	tt := ui.NewTester(view, 320, 160)
	tt.Frame()
	if err := tt.Click("大"); err != nil {
		t.Fatalf("点击失败: %v", err)
	}
	tt.Frame()
	if size != "large" {
		t.Errorf("选择 = %q, want large", size)
	}
}

func TestSwitchToggle(t *testing.T) {
	on, fired := false, 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		// 固定内边距，让开关位置可预测（中心约 62, 51）。
		ui.Column(c).Padding(40).Children(func() {
			widgets.Switch(c, &on, widgets.OnChange(func() { fired++ }))
		})
	}
	tt := ui.NewTester(view, 200, 140)
	tt.Frame()
	tt.ClickAt(62, 51)
	tt.Frame()
	if !on {
		t.Error("点击后开关未打开")
	}
	if fired != 1 {
		t.Errorf("OnChange 触发 %d 次, want 1", fired)
	}
}

func TestRateAndSegmented(t *testing.T) {
	score, tab, fired := 0, 1, 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		col := ui.Column(c).Gap(8)
		col.Children(func() {
			widgets.Rate(c, &score, 5, widgets.OnChange(func() { fired++ }))
			widgets.Segmented(c, &tab, []string{"全部", "进行中", "已完成"})
		})
	}
	tt := ui.NewTester(view, 360, 160)
	tt.Frame()
	// 点击分段控制的「全部」。
	if err := tt.Click("全部"); err != nil {
		t.Fatalf("点击失败: %v", err)
	}
	tt.Frame()
	if tab != 0 {
		t.Errorf("Segmented 选择 = %d, want 0", tab)
	}
}

func TestSliderSetValue(t *testing.T) {
	v, fired := 20.0, 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		// 固定内边距，滑轨从 x=40 延伸到 x=280，中心 y≈50。
		ui.Column(c).Padding(40).Children(func() {
			widgets.Slider(c, &v, 0, 100, widgets.OnChange(func() { fired++ })).FillWidth()
		})
	}
	tt := ui.NewTester(view, 320, 100)
	tt.Frame()
	tt.ClickAt(160, 50)
	tt.Frame()
	if fired == 0 {
		t.Error("点击滑轨未触发 OnChange")
	}
}

func TestSelectFollowsThemeAndDensity(t *testing.T) {
	for _, tk := range []*theme.AntToken{
		theme.Light(),
		theme.Dark(),
		theme.Light(theme.WithDensity(theme.Compact)),
		theme.Dark(theme.WithDensity(theme.Compact), theme.WithPrimary(ui.Hex("#7c3aed"))),
	} {
		token := tk
		checked, on, score, tab := false, false, 0, 0
		low, high := 10.0, 90.0
		view := func(c *ui.Context) {
			theme.Apply(c, token)
			col := ui.Column(c).Gap(8)
			col.Children(func() {
				widgets.Checkbox(c, &checked, "复选")
				widgets.CheckboxAll(c, new(ui.CheckState), "全选")
				widgets.RadioGroup(c, func() {
					widgets.Radio(c, &tab, 0, "甲")
					widgets.Radio(c, &tab, 1, "乙")
				})
				widgets.Switch(c, &on)
				widgets.SwitchWithLabel(c, &on, "启用")
				widgets.Slider(c, &low, 0, 100)
				widgets.RangeSliderDouble(c, &low, &high, 0, 100, 1)
				widgets.Rate(c, &score, 5)
				widgets.Segmented(c, &tab, []string{"一", "二"})
			})
		}
		tt := ui.NewTester(view, 400, 400)
		tt.Frame()
		if tt.FramePending() {
			t.Errorf("密度 %v 下选择控件持续请求帧", token.Density)
		}
	}
}
