package ant_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func TestBreadcrumbAndTabs(t *testing.T) {
	// Breadcrumbs 语义：chosen=-1 表示未选择；最后一项是当前页不可点。
	level, tab := -1, 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		col := ui.Column(c).Gap(8)
		col.Children(func() {
			widgets.Breadcrumb(c, []string{"首页", "列表", "详情"}, &level)
			widgets.Tabs(c, &tab, "常规", "外观", "高级")
		})
	}
	tt := ui.NewTester(view, 400, 160)
	tt.Frame()
	if err := tt.Click("列表"); err != nil {
		t.Fatalf("点击面包屑失败: %v", err)
	}
	tt.Frame()
	if level != 1 {
		t.Errorf("面包屑层级 = %d, want 1", level)
	}
	if err := tt.Click("外观"); err != nil {
		t.Fatalf("点击页签失败: %v", err)
	}
	tt.Frame()
	if tab != 1 {
		t.Errorf("页签 = %d, want 1", tab)
	}
}

func TestDropdownMenu(t *testing.T) {
	selected := ""
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Dropdown(c, "操作", []widgets.MenuItemData{
			{Label: "编辑", OnSelect: func() { selected = "编辑" }},
			{Separator: true},
			{Label: "删除", Danger: true, OnSelect: func() { selected = "删除" }},
		})
	}
	tt := ui.NewTester(view, 320, 200)
	tt.Frame()
	if err := tt.Click("操作"); err != nil {
		t.Fatalf("点击菜单按钮失败: %v", err)
	}
	tt.Frame()
	// 菜单在独立弹层渲染（HasText 扫不到），直接用 ChooseMenuItem 选择。
	if err := tt.ChooseMenuItem("编辑"); err != nil {
		t.Fatalf("选择菜单项失败: %v", err)
	}
	tt.Frame()
	if selected != "编辑" {
		t.Errorf("选择 = %q, want 编辑", selected)
	}
}

func TestPagination(t *testing.T) {
	page, changed := 1, 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Pagination(c, &page, 42, widgets.OnPageChange(func(p int) { changed++; _ = p }))
	}
	tt := ui.NewTester(view, 480, 120)
	tt.Frame()
	// 总数 42 > 7，窗口为 1 2 3 4 … 42。
	if !tt.HasText("…") {
		t.Error("省略号未渲染")
	}
	if err := tt.Click("42"); err != nil {
		t.Fatalf("点击尾页失败: %v", err)
	}
	tt.Frame()
	if page != 42 {
		t.Errorf("页码 = %d, want 42", page)
	}
	if changed != 1 {
		t.Errorf("OnPageChange 触发 %d 次, want 1", changed)
	}
	// 跳到尾页后窗口翻转。
	if !tt.HasText("39") {
		t.Error("尾页窗口未渲染 39")
	}
}

func TestStepsStates(t *testing.T) {
	step := 1
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Steps(c, &step, []string{"填写信息", "确认订单", "支付"})
	}
	tt := ui.NewTester(view, 480, 140)
	tt.Frame()
	for _, want := range []string{"填写信息", "确认订单", "支付"} {
		if !tt.HasText(want) {
			t.Errorf("步骤 %q 未渲染", want)
		}
	}
	if !tt.HasText("✓") {
		t.Error("完成步骤的对勾未渲染")
	}
}

func TestNavFollowsThemeAndDensity(t *testing.T) {
	for _, tk := range []*theme.AntToken{
		theme.Light(),
		theme.Dark(),
		theme.Light(theme.WithDensity(theme.Compact)),
		theme.Dark(theme.WithDensity(theme.Compact), theme.WithPrimary(ui.Hex("#7c3aed"))),
	} {
		token := tk
		level, tab, page, step := 0, 0, 5, 1
		view := func(c *ui.Context) {
			theme.Apply(c, token)
			col := ui.Column(c).Gap(8).FillWidth()
			col.Children(func() {
				widgets.Breadcrumb(c, []string{"首页", "设置"}, &level)
				widgets.Tabs(c, &tab, "一", "二")
				widgets.Pagination(c, &page, 42)
				widgets.Steps(c, &step, []string{"甲", "乙", "丙"})
				widgets.Dropdown(c, "菜单", []widgets.MenuItemData{{Label: "项"}})
			})
		}
		tt := ui.NewTester(view, 480, 320)
		tt.Frame()
		if tt.FramePending() {
			t.Errorf("密度 %v 下导航控件持续请求帧", token.Density)
		}
	}
}
