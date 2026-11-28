package ant_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func TestTooltipRendersAnchor(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Tooltip(c, "删除后不可恢复", func(c *ui.Context) {
			widgets.Button(c, "删除", widgets.Danger())
		})
	}
	tt := ui.NewTester(view, 360, 160)
	tt.Frame()
	if !tt.HasText("删除") {
		t.Error("Tooltip 锚点未渲染")
	}
}

func TestPopconfirmConfirmAndCancel(t *testing.T) {
	confirmed := 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Popconfirm(c, func(c *ui.Context) ui.Element {
			return widgets.Button(c, "删除")
		}, "确定删除这一项？", func() { confirmed++ },
			widgets.PopconfirmTitle("删除确认"))
	}
	tt := ui.NewTester(view, 480, 260)
	tt.Frame()
	if err := tt.Click("删除"); err != nil {
		t.Fatalf("点击触发器失败: %v", err)
	}
	tt.Frame()
	if !tt.HasText("确定删除这一项？") || !tt.HasText("删除确认") {
		t.Fatal("Popconfirm 气泡未弹出")
	}
	if err := tt.Click("确定"); err != nil {
		t.Fatalf("点击确定失败: %v", err)
	}
	tt.Frame()
	if confirmed != 1 {
		t.Errorf("onConfirm 触发 %d 次, want 1", confirmed)
	}

	// 再次打开并取消：不应触发回调。
	if err := tt.Click("删除"); err != nil {
		t.Fatalf("再次点击触发器失败: %v", err)
	}
	tt.Frame()
	if err := tt.Click("取消"); err != nil {
		t.Fatalf("点击取消失败: %v", err)
	}
	tt.Frame()
	if confirmed != 1 {
		t.Errorf("取消后 onConfirm 触发 %d 次, want 1", confirmed)
	}
}

func TestCollapseToggle(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Collapse(c, []widgets.CollapseItem{
			{Header: "基础设置", Content: func() { widgets.Paragraph(c, "基础内容") }},
			{Header: "高级设置", Content: func() { widgets.Paragraph(c, "高级内容") }},
		}, widgets.CollapseDefaultOpen("基础设置"))
	}
	tt := ui.NewTester(view, 480, 320)
	tt.Frame()
	if !tt.HasText("基础内容") {
		t.Fatal("默认展开项的内容未渲染")
	}
	if tt.HasText("高级内容") {
		t.Error("未展开项的内容不应渲染")
	}
	if err := tt.Click("高级设置"); err != nil {
		t.Fatalf("点击高级设置失败: %v", err)
	}
	tt.Frame()
	if !tt.HasText("高级内容") {
		t.Error("点击后高级设置未展开")
	}
}

func TestCollapseAccordion(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Collapse(c, []widgets.CollapseItem{
			{Header: "第一项", Content: func() { widgets.Paragraph(c, "内容一") }},
			{Header: "第二项", Content: func() { widgets.Paragraph(c, "内容二") }},
		}, widgets.CollapseAccordion(true), widgets.CollapseDefaultOpen("第一项"))
	}
	tt := ui.NewTester(view, 480, 320)
	tt.Frame()
	if !tt.HasText("内容一") {
		t.Fatal("第一项未默认展开")
	}
	if err := tt.Click("第二项"); err != nil {
		t.Fatalf("点击第二项失败: %v", err)
	}
	// 手风琴的关闭在展开帧之后生效，折叠动画需要多帧；推进到稳定。
	for i := 0; i < 40; i++ {
		tt.Frame()
		if !tt.FramePending() {
			break
		}
	}
	if !tt.HasText("内容二") {
		t.Error("第二项未展开")
	}
	if tt.HasText("内容一") {
		t.Error("手风琴模式下其余项应关闭")
	}
}

func TestFormItems(t *testing.T) {
	name := ""
	email := ""
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Form(c, func() {
			widgets.FormItem(c, "用户名", func() { widgets.Input(c, &name) },
				widgets.FormItemRequired(true))
			widgets.FormItem(c, "邮箱", func() { widgets.Input(c, &email) },
				widgets.FormItemHelp("用于接收通知"),
				widgets.FormItemError("邮箱格式不正确"))
			widgets.FormGroup(c, "其他", func() {
				widgets.FormItem(c, "备注", func() { widgets.Input(c, new(string)) })
			})
		})
	}
	tt := ui.NewTester(view, 560, 320)
	tt.Frame()
	for _, want := range []string{"用户名", "邮箱", "邮箱格式不正确", "其他", "备注"} {
		if !tt.HasText(want) {
			t.Errorf("表单元素 %q 未渲染", want)
		}
	}
	// antd：进入错误态后帮助文字让位给错误文案。
	if tt.HasText("用于接收通知") {
		t.Error("错误态不应再显示帮助文字")
	}
}

func TestOverlayFollowsThemeAndDensity(t *testing.T) {
	for _, tk := range []*theme.AntToken{
		theme.Light(),
		theme.Dark(),
		theme.Light(theme.WithDensity(theme.Compact)),
		theme.Dark(theme.WithDensity(theme.Compact), theme.WithPrimary(ui.Hex("#7c3aed"))),
	} {
		token := tk
		view := func(c *ui.Context) {
			theme.Apply(c, token)
			col := ui.Column(c).Gap(8)
			col.Children(func() {
				widgets.Tooltip(c, "提示", func(c *ui.Context) { widgets.Button(c, "悬停") })
				widgets.Popconfirm(c, func(c *ui.Context) ui.Element {
					return widgets.Button(c, "删除")
				}, "确定？", func() {})
				widgets.Collapse(c, []widgets.CollapseItem{
					{Header: "分组", Content: func() { widgets.Paragraph(c, "内容") }},
				}, widgets.CollapseDefaultOpen("分组"))
				widgets.Form(c, func() {
					widgets.FormItem(c, "字段", func() { widgets.Input(c, new(string)) })
				})
			})
		}
		tt := ui.NewTester(view, 480, 420)
		tt.Frame()
		if tt.FramePending() {
			t.Errorf("密度 %v 下浮层/表控件持续请求帧", token.Density)
		}
	}
}

func TestTooltipShowsOnHover(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		t := c.Theme()
		ui.Column(c).Padding(t.Space(16)).Children(func() {
			widgets.Tooltip(c, "提示内容", func(c *ui.Context) {
				widgets.Button(c, "悬停我")
			})
		})
	}
	tt := ui.NewTester(view, 360, 200)
	tt.Frame()
	r, ok := tt.Find("悬停我")
	if !ok {
		t.Fatal("找不到锚点")
	}
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Advance(2 * time.Second)
	tt.Frame()
	if !tt.HasText("提示内容") {
		t.Fatalf("悬停后未显示提示，texts=%v", tt.Texts())
	}
}
