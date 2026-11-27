package ant_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func TestNativeWrappersRender(t *testing.T) {
	date := time.Date(2026, 11, 27, 0, 0, 0, 0, time.Local)
	tokens := []string{"Go"}
	amount := "12.30"
	query := ""
	count := float64(3)
	name := "旧名字"
	open := true
	general := false
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		col := ui.Column(c).Gap(8).FillWidth()
		col.Children(func() {
			widgets.Hyperlink(c, "查看文档", "https://example.com")
			widgets.DecimalInput(c, &amount, "0", "1000000", "0.01")
			widgets.SearchField(c, &query)
			widgets.Stepper(c, &count, 0, 100, 1)
			widgets.Calendar(c, &date)
			widgets.TokenField(c, &tokens, []string{"Go", "GUI"})
			widgets.EditableText(c, &name)
			widgets.Meter(c, 72, 0, 100, &ui.MeterLevels{Warning: 60, Critical: 85})
			widgets.Collapsible(c, "折叠项", &open, func() {
				widgets.Paragraph(c, "折叠内容")
			})
			widgets.Accordion(c, func() {
				widgets.AccordionItem(c, "常规", &general, func() {
					widgets.Paragraph(c, "常规内容")
				})
			})
		})
	}
	tt := ui.NewTester(view, 640, 900)
	tt.Frame()
	for _, want := range []string{"查看文档", "折叠项", "折叠内容", "常规"} {
		if !tt.HasText(want) {
			t.Errorf("未渲染 %q", want)
		}
	}
	if tt.FramePending() {
		t.Error("原生包装控件持续请求帧")
	}
}

func TestPopoverOpens(t *testing.T) {
	open := false
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Popover(c, func(c *ui.Context) {
			widgets.Button(c, "详情", widgets.OnClick(func() { open = true }))
		}, &open, func() {
			widgets.Paragraph(c, "气泡内容")
		})
	}
	tt := ui.NewTester(view, 420, 300)
	tt.Frame()
	if err := tt.Click("详情"); err != nil {
		t.Fatalf("点击触发器失败: %v", err)
	}
	tt.Frame()
	if !tt.HasText("气泡内容") {
		t.Error("Popover 内容未显示")
	}
}

func TestContextMenuChooses(t *testing.T) {
	chosen := ""
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.ContextMenu(c, func(c *ui.Context) {
			widgets.Paragraph(c, "右键这一行")
		}, []widgets.ContextMenuItem{
			{ID: "rename", Label: "重命名"},
			{ID: "delete", Label: "删除", Danger: true},
		}, func(id string) { chosen = id })
	}
	tt := ui.NewTester(view, 480, 320)
	tt.Frame()
	if err := tt.RightClick("右键这一行"); err != nil {
		t.Fatalf("右键失败: %v", err)
	}
	tt.Frame()
	if !tt.HasText("重命名") {
		t.Fatal("右键菜单未弹出")
	}
	if err := tt.Click("删除"); err != nil {
		t.Fatalf("点击菜单项失败: %v", err)
	}
	tt.Frame()
	if chosen != "delete" {
		t.Errorf("选择 = %q, want delete", chosen)
	}
}
