package ant_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func TestWatermarkAndMasonryRender(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Watermark(c, "内部资料", func() {
			widgets.Masonry(c, 2, []func(*ui.Context){
				func(c *ui.Context) { widgets.Card(c, "A", func() { widgets.Paragraph(c, "卡片 A") }) },
				func(c *ui.Context) { widgets.Card(c, "B", func() { widgets.Paragraph(c, "卡片 B") }) },
				func(c *ui.Context) { widgets.Card(c, "C", func() { widgets.Paragraph(c, "卡片 C") }) },
			})
		})
	}
	tt := ui.NewTester(view, 640, 420)
	tt.Frame()
	for _, want := range []string{"卡片 A", "卡片 B", "卡片 C"} {
		if !tt.HasText(want) {
			t.Errorf("瀑布流未渲染 %q", want)
		}
	}
	if tt.FramePending() {
		t.Error("水印/瀑布流持续请求帧")
	}
}

func TestMentionsInserts(t *testing.T) {
	reply := ""
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Mentions(c, &reply, []string{"张三", "李四"})
	}
	tt := ui.NewTester(view, 480, 260)
	tt.Frame()
	tt.ClickAt(120, 16)
	tt.Type("@张")
	tt.Frame()
	if !tt.HasText("张三") {
		t.Fatal("提及候选未显示")
	}
	if err := tt.Click("张三"); err != nil {
		t.Fatalf("点击候选失败: %v", err)
	}
	tt.Frame()
	if reply != "@张三 " {
		t.Errorf("插入结果 = %q, want \"@张三 \"", reply)
	}
}

func TestTourSteps(t *testing.T) {
	open := true
	current := 0
	steps := []widgets.TourStep{
		{Title: "第一步", Description: "欢迎使用"},
		{Title: "第二步", Description: "继续"},
	}
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Tour(c, &open, steps, &current)
	}
	tt := ui.NewTester(view, 640, 480)
	tt.Frame()
	if !tt.HasText("第一步") {
		t.Fatal("引导卡片未显示")
	}
	if err := tt.Click("下一步"); err != nil {
		t.Fatalf("点击下一步失败: %v", err)
	}
	tt.Frame()
	if current != 1 || !tt.HasText("第二步") {
		t.Error("未进入第二步")
	}
	if err := tt.Click("完成"); err != nil {
		t.Fatalf("点击完成失败: %v", err)
	}
	tt.Frame()
	if open {
		t.Error("完成后应关闭引导")
	}
}

func TestNotificationAndFloatButton(t *testing.T) {
	open := true
	clicked := 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Notification(c, &open, "部署完成", "v0.3.0 已上线")
		widgets.FloatButton(c, "↑", func() { clicked++ })
	}
	tt := ui.NewTester(view, 640, 400)
	tt.Frame()
	if !tt.HasText("部署完成") || !tt.HasText("v0.3.0 已上线") {
		t.Fatal("通知卡片未显示")
	}
	if err := tt.Click("↑"); err != nil {
		t.Fatalf("点击悬浮按钮失败: %v", err)
	}
	if clicked != 1 {
		t.Errorf("悬浮按钮回调 %d 次, want 1", clicked)
	}
}

func TestAnchorSelects(t *testing.T) {
	active := 0
	selected := ""
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Anchor(c, []widgets.AnchorItem{
			{ID: "intro", Label: "介绍"},
			{ID: "usage", Label: "用法"},
		}, &active, func(id string) { selected = id })
	}
	tt := ui.NewTester(view, 320, 240)
	tt.Frame()
	if err := tt.Click("用法"); err != nil {
		t.Fatalf("点击锚点失败: %v", err)
	}
	tt.Frame()
	if active != 1 || selected != "usage" {
		t.Errorf("active=%d selected=%q, want 1/usage", active, selected)
	}
}

func TestBorderBeamRenders(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Dark())
		widgets.BorderBeam(c, func() {
			widgets.Paragraph(c, "流光边框内容")
		})
	}
	tt := ui.NewTester(view, 420, 240)
	tt.Frame()
	if !tt.HasText("流光边框内容") {
		t.Error("BorderBeam 内容未渲染")
	}
	if tt.FramePending() {
		t.Error("BorderBeam 持续请求帧")
	}
}
