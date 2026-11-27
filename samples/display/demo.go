// Package display 提供数据展示类控件的示例。
package display

import (
	_ "embed"

	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

//go:embed demo.go
var demoSrc string

// Demo 返回数据展示控件的示例集合。
func Demo() samples.WidgetDemo {
	return samples.WidgetDemo{
		Name:     "数据展示",
		Package:  "widgets",
		Category: samples.CatDisplay,
		Summary:  "Tag/Badge/Avatar/Card、Empty/Result/Skeleton、Statistic/Descriptions/Timeline 等数据展示控件。",
		Samples: []samples.Sample{
			{
				Title:       "标签",
				Description: "七种预设色系（从主题语义色派生），支持关闭。",
				Code:        samples.Extract(demoSrc, "tags"),
				View:        Tags,
			},
			{
				Title:       "徽标",
				Description: "数字徽标（超过 99 显示 99+）与红点徽标。",
				Code:        samples.Extract(demoSrc, "badge"),
				View:        BadgeDemo,
			},
			{
				Title:       "头像",
				Description: "文字头像，尺寸随密度缩放。",
				Code:        samples.Extract(demoSrc, "avatar"),
				View:        AvatarDemo,
			},
			{
				Title:       "卡片",
				Description: "标题 + 内容 + 头部附加 + 底部操作区。",
				Code:        samples.Extract(demoSrc, "card"),
				View:        CardDemo,
			},
			{
				Title:       "空状态与结果页",
				Description: "Empty 插画 + 描述 + 操作区；Result 四种状态图标 + 标题 + 副标题。",
				Code:        samples.Extract(demoSrc, "emptyresult"),
				View:        EmptyResultDemo,
			},
			{
				Title:       "骨架屏与统计数值",
				Description: "Skeleton 静态占位（头像/标题/段落）；Statistic 大字号数值 + 前缀后缀 + 涨跌趋势。",
				Code:        samples.Extract(demoSrc, "skeletonstat"),
				View:        SkeletonStatDemo,
			},
			{
				Title:       "描述列表与时间轴",
				Description: "Descriptions 标签/值成对布局（可带外框）；Timeline 节点与连接线。",
				Code:        samples.Extract(demoSrc, "desctimeline"),
				View:        DescTimelineDemo,
			},
		},
	}
}

// sample:tags
func Tags(c *ui.Context) {
	row := ui.Row(c).Gap(8).Wrap()
	row.Children(func() {
		widgets.Tag(c, "默认")
		widgets.Tag(c, "成功", widgets.TagSuccess())
		widgets.Tag(c, "警告", widgets.TagWarning())
		widgets.Tag(c, "错误", widgets.TagError())
		widgets.Tag(c, "品牌", widgets.TagPrimary())
		widgets.Tag(c, "洋红", widgets.TagMagenta())
		widgets.Tag(c, "青色", widgets.TagCyan())
		widgets.Tag(c, "金色", widgets.TagGold())
	})
}

// /sample

// sample:badge
func BadgeDemo(c *ui.Context) {
	row := ui.Row(c).Gap(24).AlignItems(ui.Center)
	row.Children(func() {
		widgets.Badge(c, 5, func() { widgets.Button(c, "消息") })
		widgets.Badge(c, 120, func() { widgets.Button(c, "通知") })
		widgets.BadgeDot(c, func() { widgets.AvatarEl(c, "李") })
	})
}

// /sample

// sample:avatar
func AvatarDemo(c *ui.Context) {
	row := ui.Row(c).Gap(12).AlignItems(ui.Center)
	row.Children(func() {
		widgets.AvatarEl(c, "李")
		widgets.AvatarEl(c, "王")
		widgets.AvatarEl(c, "Chen")
	})
}

// /sample

// sample:card
func CardDemo(c *ui.Context) {
	widgets.Card(c, "卡片标题",
		func() {
			widgets.Paragraph(c, "这是卡片的内容区，可以放任意控件。")
			widgets.Paragraph(c, "支持头部附加与底部操作区。", widgets.Secondary())
		},
		widgets.CardExtra(func(c *ui.Context) {
			widgets.Button(c, "更多", widgets.Text())
		}),
		widgets.CardActions(func(c *ui.Context) {
			widgets.Button(c, "取消")
			widgets.Button(c, "确定", widgets.Primary())
		}))
}

// /sample

// sample:emptyresult
func EmptyResultDemo(c *ui.Context) {
	t := c.Theme()
	row := ui.Row(c).Gap(t.Space(6)).AlignItems(ui.Start)
	row.Children(func() {
		widgets.Empty(c, widgets.EmptyFooter(func(c *ui.Context) {
			widgets.Button(c, "新建", widgets.Primary())
		}))
		widgets.Result(c, widgets.ResultSuccess, "操作成功",
			widgets.ResultSubTitle("订单已提交"))
		widgets.Result(c, widgets.ResultError, "提交失败",
			widgets.ResultExtra(func(c *ui.Context) {
				widgets.Button(c, "重试", widgets.Primary())
			}))
	})
}

// /sample

// sample:skeletonstat
func SkeletonStatDemo(c *ui.Context) {
	t := c.Theme()
	col := ui.Column(c).Gap(t.Space(4))
	col.Children(func() {
		widgets.Skeleton(c, widgets.SkeletonAvatar(true), widgets.SkeletonRows(3))
		row := ui.Row(c).Gap(t.Space(6))
		row.Children(func() {
			widgets.Statistic(c, "销售额", 126560,
				widgets.StatisticPrefix("¥"), widgets.StatisticTrendUp())
			widgets.Statistic(c, "转化率", 12.3,
				widgets.StatisticSuffix("%"), widgets.StatisticPrecision(1),
				widgets.StatisticTrendDown())
		})
	})
}

// /sample

// sample:desctimeline
func DescTimelineDemo(c *ui.Context) {
	t := c.Theme()
	col := ui.Column(c).Gap(t.Space(4))
	col.Children(func() {
		widgets.Descriptions(c, []widgets.DescriptionItem{
			{Label: "用户名", Value: "张三"},
			{Label: "手机号", Value: "13800000000"},
			{Label: "状态", Value: "已启用"},
		}, widgets.DescriptionsColumns(2), widgets.DescriptionsTitle("账号信息"),
			widgets.DescriptionsBordered(true))
		widgets.Timeline(c, []widgets.TimelineItem{
			{Title: "创建订单", Description: "2026-11-27 09:00"},
			{Title: "已支付", Description: "2026-11-27 09:05", Color: widgets.TimelineGreen},
			{Title: "已发货", Description: "2026-11-27 12:30", Color: widgets.TimelineRed},
			{Title: "已完成", Color: widgets.TimelineGray},
		})
	})
}

// /sample
