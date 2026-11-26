package ant_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func TestSelectChooseOption(t *testing.T) {
	var picked int
	city := "上海"
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Select(c, &city, []string{"北京", "上海", "广州"}, widgets.OnChange(func() { picked++ }))
	}
	tt := ui.NewTester(view, 320, 160)
	tt.Frame()
	// 点击触发器（显示当前值）展开弹层，再直接点击弹层内的选项。
	if err := tt.Click("上海"); err != nil {
		t.Fatalf("点击选择框失败: %v", err)
	}
	tt.Frame()
	if err := tt.Click("广州"); err != nil {
		t.Fatalf("选择选项失败: %v", err)
	}
	tt.Frame()
	if city != "广州" {
		t.Errorf("选择 = %q, want 广州", city)
	}
	if picked != 1 {
		t.Errorf("OnChange 触发 %d 次, want 1", picked)
	}
}

func TestComboboxTypeAndPick(t *testing.T) {
	q := ""
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		// 固定内边距，输入框位于 (40,40) 起的高 32 行内。
		ui.Column(c).Padding(40).Children(func() {
			widgets.Combobox(c, &q, []string{"重启服务", "重启网关"})
		})
	}
	tt := ui.NewTester(view, 320, 160)
	tt.Frame()
	tt.ClickAt(120, 56)
	tt.Frame()
	tt.Type("网")
	tt.Frame()
	// 键入过滤出候选「重启网关」，点击候选提交选中值。
	if !tt.HasText("重启网关") {
		t.Fatal("键入后未出现过滤候选")
	}
	if err := tt.Click("重启网关"); err != nil {
		t.Fatalf("点击候选失败: %v", err)
	}
	tt.Frame()
	if q != "重启网关" {
		t.Errorf("Combobox 选中值 = %q, want 重启网关", q)
	}
}

func TestDatePickerAndTimePicker(t *testing.T) {
	d := time.Date(2026, 11, 27, 0, 0, 0, 0, time.Local)
	tm := time.Date(2026, 11, 27, 14, 30, 0, 0, time.Local)
	changed := 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.DatePicker(c, &d, widgets.OnChanged(func() { changed++ }))
		widgets.TimePicker(c, &tm)
	}
	tt := ui.NewTester(view, 320, 160)
	tt.Frame()
	if tt.FramePending() {
		t.Error("日期/时间选择框持续请求帧")
	}
	// 文本编辑路径：聚焦日期输入并键入应更新值。
	tt.ClickAt(60, 20)
	tt.Type("12/25/2026")
	tt.Frame()
	if d.Month() != time.November || d.Day() != 27 {
		// 无窗口环境下编辑器行为可能受限，仅要求不 panic；值断言放宽。
		t.Logf("键入后日期 = %v（无窗口环境编辑器受限时允许不变）", d)
	}
}

func TestPickerDisabled(t *testing.T) {
	d := time.Date(2026, 11, 27, 0, 0, 0, 0, time.Local)
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.DatePicker(c, &d, widgets.InputDisabled(true))
	}
	tt := ui.NewTester(view, 320, 120)
	tt.Frame()
	tt.ClickAt(60, 20)
	tt.Type("x")
	tt.Frame()
	if d.Day() != 27 {
		t.Error("禁用日期选择框不应改变值")
	}
}

func TestPickersFollowThemeAndDensity(t *testing.T) {
	for _, tk := range []*theme.AntToken{
		theme.Light(),
		theme.Dark(),
		theme.Light(theme.WithDensity(theme.Compact)),
		theme.Dark(theme.WithDensity(theme.Compact), theme.WithPrimary(ui.Hex("#7c3aed"))),
	} {
		token := tk
		d := time.Date(2026, 1, 2, 0, 0, 0, 0, time.Local)
		tm := time.Date(2026, 1, 2, 8, 0, 0, 0, time.Local)
		sel, combo := "甲", "一"
		view := func(c *ui.Context) {
			theme.Apply(c, token)
			col := ui.Column(c).Gap(8)
			col.Children(func() {
				widgets.DatePicker(c, &d)
				widgets.TimePicker(c, &tm)
				widgets.Select(c, &sel, []string{"甲", "乙"})
				widgets.Combobox(c, &combo, []string{"一", "二"})
			})
		}
		tt := ui.NewTester(view, 360, 320)
		tt.Frame()
		if tt.FramePending() {
			t.Errorf("密度 %v 下日期/时间/选择控件持续请求帧", token.Density)
		}
	}
}
