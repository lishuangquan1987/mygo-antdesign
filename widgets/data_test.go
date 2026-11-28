package ant_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func TestEmptyAndResult(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		col := ui.Column(c).Gap(8)
		col.Children(func() {
			widgets.Empty(c)
			widgets.Empty(c, widgets.EmptyDescription("没有匹配的记录"),
				widgets.EmptyFooter(func(c *ui.Context) {
					widgets.Button(c, "新建", widgets.Primary())
				}))
			widgets.Result(c, widgets.ResultSuccess, "操作成功",
				widgets.ResultSubTitle("订单已提交"))
			widgets.Result(c, widgets.ResultError, "提交失败",
				widgets.ResultExtra(func(c *ui.Context) {
					widgets.Button(c, "重试", widgets.Primary())
				}))
			widgets.Result(c, widgets.ResultWarning, "存在风险")
			widgets.Result(c, widgets.ResultInfo, "处理中")
		})
	}
	tt := ui.NewTester(view, 560, 900)
	tt.Frame()
	for _, want := range []string{
		"暂无数据", "没有匹配的记录", "新建",
		"操作成功", "订单已提交", "提交失败", "重试", "存在风险", "处理中",
	} {
		if !tt.HasText(want) {
			t.Errorf("未渲染 %q", want)
		}
	}
}

func TestSkeletonAndStatistic(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		col := ui.Column(c).Gap(12)
		col.Children(func() {
			widgets.Skeleton(c)
			widgets.Skeleton(c, widgets.SkeletonAvatar(true), widgets.SkeletonRows(2))
			row := ui.Row(c).Gap(24)
			row.Children(func() {
				widgets.Statistic(c, "销售额", 126560,
					widgets.StatisticPrefix("¥"), widgets.StatisticTrendUp())
				widgets.Statistic(c, "转化率", 12.3,
					widgets.StatisticSuffix("%"), widgets.StatisticPrecision(1),
					widgets.StatisticTrendDown())
				widgets.Statistic(c, "订单数", 88)
			})
		})
	}
	tt := ui.NewTester(view, 640, 400)
	tt.Frame()
	for _, want := range []string{"销售额", "¥", "126,560", "转化率", "12.3", "订单数", "88"} {
		if !tt.HasText(want) {
			t.Errorf("未渲染 %q", want)
		}
	}
}

func TestDescriptionsAndTimeline(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		col := ui.Column(c).Gap(16)
		col.Children(func() {
			widgets.Descriptions(c, []widgets.DescriptionItem{
				{Label: "用户名", Value: "张三"},
				{Label: "手机号", Value: "13800000000"},
				{Label: "状态", Value: "已启用"},
			}, widgets.DescriptionsColumns(2), widgets.DescriptionsTitle("账号信息"),
				widgets.DescriptionsBordered(true))
			widgets.Timeline(c, []widgets.TimelineItem{
				{Title: "创建订单", Description: "09:00"},
				{Title: "已支付", Description: "09:05", Color: widgets.TimelineGreen},
				{Title: "已发货", Description: "12:30", Color: widgets.TimelineRed},
				{Title: "已完成", Color: widgets.TimelineGray},
			})
		})
	}
	tt := ui.NewTester(view, 560, 520)
	tt.Frame()
	for _, want := range []string{
		"账号信息", "用户名", "张三", "手机号", "状态", "已启用",
		"创建订单", "已支付", "已发货", "已完成", "09:05",
	} {
		if !tt.HasText(want) {
			t.Errorf("未渲染 %q", want)
		}
	}
}

func TestDataGapsFollowThemeAndDensity(t *testing.T) {
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
				widgets.Empty(c)
				widgets.Result(c, widgets.ResultSuccess, "成功")
				widgets.Skeleton(c, widgets.SkeletonAvatar(true))
				widgets.Statistic(c, "数量", 42, widgets.StatisticTrendUp())
				widgets.Descriptions(c, []widgets.DescriptionItem{
					{Label: "标签", Value: "值"},
				}, widgets.DescriptionsColumns(1), widgets.DescriptionsBordered(true))
				widgets.Timeline(c, []widgets.TimelineItem{
					{Title: "节点一"},
					{Title: "节点二", Color: widgets.TimelineGreen},
				})
			})
		}
		tt := ui.NewTester(view, 480, 760)
		tt.Frame()
		if tt.FramePending() {
			t.Errorf("密度 %v 下数据展示缺口控件持续请求帧", token.Density)
		}
	}
}
