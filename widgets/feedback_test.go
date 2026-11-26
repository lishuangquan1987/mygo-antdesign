package ant_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func TestAlertVariantsAndClose(t *testing.T) {
	closed := 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		col := ui.Column(c).Gap(8).FillWidth()
		col.Children(func() {
			widgets.Alert(c, "信息提示", widgets.AlertTypeOpt(widgets.AlertInfo),
				widgets.AlertDescription("这是一条信息"))
			widgets.Alert(c, "保存成功", widgets.AlertTypeOpt(widgets.AlertSuccess))
			widgets.Alert(c, "磁盘空间不足", widgets.AlertTypeOpt(widgets.AlertWarning))
			widgets.Alert(c, "连接失败", widgets.AlertTypeOpt(widgets.AlertError))
			widgets.Alert(c, "可关闭的提示", widgets.AlertClosable(func() { closed++ }))
		})
	}
	tt := ui.NewTester(view, 480, 400)
	tt.Frame()
	for _, want := range []string{"信息提示", "保存成功", "磁盘空间不足", "连接失败", "可关闭的提示"} {
		if !tt.HasText(want) {
			t.Errorf("提示条 %q 未渲染", want)
		}
	}
	if err := tt.Click("×"); err != nil {
		t.Fatalf("点击关闭失败: %v", err)
	}
	tt.Frame()
	if closed != 1 {
		t.Errorf("AlertClosable 触发 %d 次, want 1", closed)
	}
	// 关闭后提示条消失。
	if tt.HasText("可关闭的提示") {
		t.Error("关闭后提示条仍显示")
	}
}

func TestModalOpenAndActions(t *testing.T) {
	open, okCount, cancelCount := true, 0, 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Modal(c, &open, "确认操作",
			func() {
				widgets.Paragraph(c, "此操作不可撤销，是否继续？")
			},
			widgets.OnOk(func() { okCount++ }),
			widgets.OnCancel(func() { cancelCount++ }))
	}
	tt := ui.NewTester(view, 480, 320)
	tt.Frame()
	if !tt.HasText("此操作不可撤销，是否继续？") {
		t.Fatal("弹窗内容未渲染")
	}
	if err := tt.Click("确定"); err != nil {
		t.Fatalf("点击确定失败: %v", err)
	}
	tt.Frame()
	if okCount != 1 || open {
		t.Errorf("OnOk=%d open=%v, want OnOk=1 open=false", okCount, open)
	}
}

func TestModalCancel(t *testing.T) {
	open, cancelCount := true, 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Modal(c, &open, "取消测试", func() {}, widgets.OnCancel(func() { cancelCount++ }))
	}
	tt := ui.NewTester(view, 480, 320)
	tt.Frame()
	if err := tt.Click("取消"); err != nil {
		t.Fatalf("点击取消失败: %v", err)
	}
	tt.Frame()
	if cancelCount != 1 || open {
		t.Errorf("OnCancel=%d open=%v, want 1/false", cancelCount, open)
	}
}

func TestProgressLine(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		col := ui.Column(c).Gap(8).FillWidth()
		col.Children(func() {
			widgets.ProgressLine(c, 0)
			widgets.ProgressLine(c, 68)
			widgets.ProgressLine(c, 100)
		})
	}
	tt := ui.NewTester(view, 400, 200)
	tt.Frame()
	if !tt.HasText("0%") || !tt.HasText("68%") || !tt.HasText("100%") {
		t.Error("进度百分比文本未渲染完整")
	}
}

func TestFeedbackFollowsThemeAndDensity(t *testing.T) {
	for _, tk := range []*theme.AntToken{
		theme.Light(),
		theme.Dark(),
		theme.Light(theme.WithDensity(theme.Compact)),
		theme.Dark(theme.WithDensity(theme.Compact), theme.WithPrimary(ui.Hex("#7c3aed"))),
	} {
		token := tk
		open := true
		view := func(c *ui.Context) {
			theme.Apply(c, token)
			col := ui.Column(c).Gap(8).FillWidth()
			col.Children(func() {
				widgets.Alert(c, "提示", widgets.AlertTypeOpt(widgets.AlertWarning))
				widgets.ProgressLine(c, 42)
				widgets.Modal(c, &open, "弹窗", func() { widgets.Paragraph(c, "内容") })
			})
		}
		tt := ui.NewTester(view, 400, 320)
		tt.Frame()
		if tt.FramePending() {
			t.Errorf("密度 %v 下反馈控件持续请求帧", token.Density)
		}
	}
}
