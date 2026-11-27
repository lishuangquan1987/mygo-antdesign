package ant_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

// TestMygoNamedAliases 验证按 MyGo 命名补齐的入口都能渲染（docs/02 覆盖矩阵）。
func TestMygoNamedAliases(t *testing.T) {
	chosen := 1
	rating := 3
	progress := 42.0
	low, high := 20.0, 80.0
	number := 12.0
	name := "李"
	query := ""
	findOpen := false

	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		t := c.Theme()
		ui.Column(c).Padding(t.Space(8)).Gap(t.Space(3)).Children(func() {
			widgets.Breadcrumbs(c, []string{"首页", "列表", "详情"}, &chosen)
			widgets.Avatar(c, name)
			widgets.NumberInput(c, &number, 0, 100, 1)
			widgets.Progress(c, progress)
			widgets.RangeSlider(c, &low, &high, 0, 100, 1)
			widgets.Rating(c, &rating, 5)
			widgets.Spinner(c)
			widgets.RichText(c, ui.Span{Text: "富文本"}, ui.Span{Text: "片段", Weight: 600})
			widgets.PrimaryButton(c, "保存")
			widgets.Toolbar(c, func() {
				widgets.Button(c, "加粗")
				widgets.Button(c, "斜体")
			})
			widgets.FindBar(c, &findOpen, &query, 0, new(int))
		})
	}
	tt := ui.NewTester(view, 640, 520)
	tt.Frame()
	for _, want := range []string{"首页", "保存", "富文本"} {
		if !tt.HasText(want) {
			t.Errorf("缺少文本 %q", want)
		}
	}
}

// TestMygoNamedSplitButton 验证分裂按钮按 MyGo 命名可用。
func TestMygoNamedSplitButton(t *testing.T) {
	clicked := 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		parts := widgets.SplitButton(c, "保存")
		if parts.Button.Clicked() {
			clicked++
		}
		parts.Arrow.Menu(func(m *ui.Menu) {
			if m.Item("另存为").Chosen() {
				clicked++
			}
		})
	}
	tt := ui.NewTester(view, 360, 160)
	tt.Frame()
	if !tt.HasText("保存") {
		t.Fatal("分裂按钮未渲染主按钮")
	}
	if err := tt.Click("保存"); err != nil {
		t.Fatalf("点击分裂按钮失败: %v", err)
	}
	tt.Frame()
	if clicked != 1 {
		t.Errorf("主按钮点击 %d 次, want 1", clicked)
	}
}

// TestMygoNamedGridView 验证网格视图入口可渲染惰性项。
func TestMygoNamedGridView(t *testing.T) {
	var st ui.GridState
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		ui.Column(c).Padding(24).Width(560).Children(func() {
			widgets.GridView(c, &st, 8, 120, 80, func(i int) {
				ui.Textf(c, "项目 %d", i)
			})
		})
	}
	tt := ui.NewTester(view, 600, 320)
	tt.Frame()
	if !tt.HasText("项目 0") {
		t.Error("网格视图未渲染首项")
	}
}

// TestMygoNamedCheckboxGroup 验证复选组入口渲染标签。
func TestMygoNamedCheckboxGroup(t *testing.T) {
	var mail, chat bool
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.CheckboxGroup(c, "通知", func() {
			widgets.Checkbox(c, &mail, "邮件")
			widgets.Checkbox(c, &chat, "聊天")
		})
	}
	tt := ui.NewTester(view, 360, 200)
	tt.Frame()
	if !tt.HasText("通知") || !tt.HasText("邮件") {
		t.Error("复选组未渲染标签或子项")
	}
}
