package ant_test

import (
	"regexp"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

// TestSelectMultiToggle 多选下拉：打开浮层后点选项切换选中，回调触发一次。
func TestSelectMultiToggle(t *testing.T) {
	tags := []string{"后端"}
	calls := 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.SelectMulti(c, &tags, []string{"前端", "后端", "设计"},
			widgets.OnChange(func() { calls++ }))
	}
	tt := ui.NewTester(view, 360, 320)
	tt.Frame()
	if len(tags) != 1 || tags[0] != "后端" {
		t.Fatalf("初值 = %v, want [后端]", tags)
	}
	if err := tt.Click("多选下拉"); err != nil {
		t.Fatalf("打开浮层失败: %v", err)
	}
	tt.Frame()
	if err := tt.Click("前端"); err != nil {
		t.Fatalf("点选项失败: %v", err)
	}
	tt.Frame()
	if len(tags) != 2 {
		t.Fatalf("选中后 tags = %v, want 2 项", tags)
	}
	if calls != 1 {
		t.Errorf("OnChange 触发 %d 次, want 1", calls)
	}
}

// TestSelectSearchFilters 可搜索下拉：键入关键字后只显示匹配候选。
func TestSelectSearchFilters(t *testing.T) {
	city := ""
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.SelectSearch(c, &city, []string{"北京", "上海", "广州"})
	}
	tt := ui.NewTester(view, 320, 280)
	tt.Frame()
	if err := tt.Click("城市选择"); err != nil {
		t.Fatalf("聚焦搜索框失败: %v", err)
	}
	tt.Type("上")
	tt.Frame()
	if !tt.HasText("上海") {
		t.Error("过滤后应保留「上海」")
	}
	if tt.HasText("北京") {
		t.Error("过滤后不应显示「北京」")
	}
}

// TestFormRulesLive 即时校验：必填为空时显示错误，填入后错误消失。
func TestFormRulesLive(t *testing.T) {
	name := ""
	re := regexp.MustCompile(`^[a-z]+$`)
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Form(c, func() {
			widgets.FormItem(c, "用户名", func() { widgets.Input(c, &name) },
				widgets.FormItemValue(func() string { return name }),
				widgets.FormItemRules(
					widgets.FormRule{Required: true},
					widgets.FormRule{Pattern: re, PatternHint: "只允许小写字母"},
				))
		})
	}
	tt := ui.NewTester(view, 420, 200)
	tt.Frame()
	if !tt.HasText("此项为必填项") {
		t.Fatal("空值应显示必填错误")
	}
	name = "abc123"
	tt.Frame()
	if !tt.HasText("只允许小写字母") {
		t.Fatal("不合格式应显示正则错误")
	}
	name = "abc"
	tt.Frame()
	if tt.HasText("只允许小写字母") || tt.HasText("此项为必填项") {
		t.Fatal("合法值不应显示错误")
	}
}

// TestTableRenderPagination 自定义单元格 + 分页。
func TestTableRenderPagination(t *testing.T) {
	rows := [][]string{
		{"A1", "v1"}, {"A2", "v2"}, {"A3", "v3"}, {"A4", "v4"}, {"A5", "v5"},
	}
	page := 1
	cols := []widgets.TableColumn{{Title: "名称"}, {Title: "值"}}
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Table(c, cols, rows,
			widgets.TableRender(func(c *ui.Context, row, col int) ui.Element {
				if col == 0 {
					return widgets.Tag(c, "R"+rows[row][0], widgets.TagPrimary())
				}
				return ui.Text(c, rows[row][col])
			}),
			widgets.TablePagination(2, &page, func(p int) {}))
	}
	tt := ui.NewTester(view, 640, 520)
	tt.Frame()
	if !tt.HasText("RA1") || !tt.HasText("RA2") {
		t.Fatal("第 1 页应渲染自定义单元格 RA1/RA2")
	}
	if tt.HasText("RA3") {
		t.Fatal("第 1 页不应出现第 2 页内容")
	}
	if err := tt.Click("2"); err != nil {
		t.Fatalf("点第 2 页失败: %v", err)
	}
	tt.Frame()
	if !tt.HasText("RA3") || !tt.HasText("RA4") {
		t.Fatal("第 2 页应渲染 RA3/RA4")
	}
}

// TestTabsOnChange 页签切换回调。
func TestTabsOnChange(t *testing.T) {
	tab := 0
	got := -1
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Tabs(c, &tab, []string{"一", "二"}, widgets.TabsOnChange(func(i int) { got = i }))
	}
	tt := ui.NewTester(view, 320, 120)
	tt.Frame()
	if err := tt.Click("二"); err != nil {
		t.Fatalf("点页签失败: %v", err)
	}
	tt.Frame()
	if tab != 1 || got != 1 {
		t.Fatalf("tab=%d onChange=%d, want 1/1", tab, got)
	}
}

// TestPaginationOptions 禁用、总数文案与每页条数。
func TestPaginationOptions(t *testing.T) {
	page := 1
	size := 10
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Pagination(c, &page, 45,
			widgets.PaginationShowTotal(), widgets.PaginationPageSize(&size))
	}
	tt := ui.NewTester(view, 720, 120)
	tt.Frame()
	if !tt.HasText("共 45 条") {
		t.Error("应显示总数「共 45 条」")
	}
	if !tt.HasText("10 条/页") {
		t.Error("应显示每页条数选择")
	}

	// 禁用态：点下一页不变。
	page = 1
	view2 := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Pagination(c, &page, 3, widgets.PaginationDisabled(true))
	}
	tt2 := ui.NewTester(view2, 480, 120)
	tt2.Frame()
	// 点第 2 页（禁用时按钮仍显示，但点击不生效）。
	_ = tt2.Click("2")
	tt2.Frame()
	if page != 1 {
		t.Errorf("禁用分页器被点击后 page=%d, want 1", page)
	}
}

// TestBadgeStatusAndTimelinePending 状态徽标与时间轴待办节点。
func TestBadgeStatusAndTimelinePending(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.BadgeStatus(c, widgets.BadgeStatusSuccess, "运行中")
		widgets.Timeline(c, []widgets.TimelineItem{{Title: "已创建"}},
			widgets.TimelinePending("等待审核"))
	}
	tt := ui.NewTester(view, 400, 240)
	tt.Frame()
	for _, want := range []string{"运行中", "已创建", "等待审核"} {
		if !tt.HasText(want) {
			t.Errorf("未渲染 %q", want)
		}
	}
}

// TestStatisticFormatter 自定义格式化覆盖默认千分位。
func TestStatisticFormatter(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Statistic(c, "转化", 1234.5, widgets.StatisticFormatter(func(v float64) string {
			return "约 1.2k"
		}))
	}
	tt := ui.NewTester(view, 320, 120)
	tt.Frame()
	if !tt.HasText("约 1.2k") {
		t.Fatal("自定义 formatter 未生效")
	}
}

// TestNotificationTypeRenders 带类型图标的通知渲染标题与描述。
func TestNotificationTypeRenders(t *testing.T) {
	open := true
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Notification(c, &open, "部署完成", "v0.6.0 已上线",
			widgets.NotificationType(widgets.AlertSuccess))
	}
	tt := ui.NewTester(view, 520, 260)
	tt.Frame()
	if !tt.HasText("部署完成") || !tt.HasText("v0.6.0 已上线") {
		t.Fatal("通知标题/描述未渲染")
	}
}

// TestFormItemRequiredValidates 守护「星号即校验」：只设 FormItemRequired
// （不写 FormRule.Required）时，空值也要出必填错误——否则标了必填却放行。
func TestFormItemRequiredValidates(t *testing.T) {
	name := ""
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Form(c, func() {
			widgets.FormItem(c, "用户名", func() { widgets.Input(c, &name) },
				widgets.FormItemRequired(true),
				widgets.FormItemValue(func() string { return name }),
				widgets.FormItemRules(widgets.FormRule{MinLen: 3, Message: "至少 3 个字符"}))
		})
	}
	tt := ui.NewTester(view, 420, 200)
	tt.Frame()
	if !tt.HasText("此项为必填项") {
		t.Fatal("标了必填但没有必填错误")
	}
	name = "ab"
	tt.Frame()
	if !tt.HasText("至少 3 个字符") {
		t.Fatal("长度不足没有提示")
	}
}
