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
			{
				Title:       "折叠面板",
				Description: "Collapse 头部点击展开/收起，支持手风琴模式与默认展开。",
				Code:        samples.Extract(demoSrc, "collapse"),
				View:        CollapseDemo,
			},
			{
				Title:       "列表",
				Description: "List 虚拟滚动列表，行选中态由调用方持有。",
				Code:        samples.Extract(demoSrc, "list"),
				View:        ListDemo,
			},
			{
				Title:       "图片",
				Description: "ImageEl 可指定宽度，点击后在模态框中放大预览。",
				Code:        samples.Extract(demoSrc, "image"),
				View:        ImageDemo,
			},
			{
				Title:       "折叠块与电平条",
				Description: "Collapsible / Accordion 原生折叠动画；Meter 按阈值变色。",
				Code:        samples.Extract(demoSrc, "accordion"),
				View:        AccordionDemo,
			},
			{
				Title:       "锚点、瀑布流与边框流光",
				Description: "Anchor 目录高亮、Masonry 多列瀑布流、BorderBeam 渐变流光边框。",
				Code:        samples.Extract(demoSrc, "masonry"),
				View:        MasonryDemo,
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
	// 纵向堆叠：Row 内百分比宽度子元素会被压缩，标题被迫截断；
	// 堆叠后每个 Result 占满整行宽度。
	col := ui.Column(c).FillWidth().Gap(t.Space(4))
	col.Children(func() {
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
	// 紧凑间距：Descriptions + Timeline 合计高度需在预览面板内完整
	// 展示（含末项），不留被底边截断的半行。描述 3 列压成单行，
	// 时间轴 3 项，保证 800px 视口下完整收尾。
	col := ui.Column(c).Gap(t.Space(1))
	col.Children(func() {
		widgets.Descriptions(c, []widgets.DescriptionItem{
			{Label: "用户名", Value: "张三"},
			{Label: "手机号", Value: "13800000000"},
			{Label: "状态", Value: "已启用"},
		}, widgets.DescriptionsColumns(3),
			widgets.DescriptionsBordered(true))
		widgets.Timeline(c, []widgets.TimelineItem{
			{Title: "创建订单", Description: "2026-11-27 09:00"},
			{Title: "已支付", Description: "2026-11-27 09:05", Color: widgets.TimelineGreen},
			{Title: "已完成", Color: widgets.TimelineGray},
		})
	})
}

// /sample

// sample:collapse
func CollapseDemo(c *ui.Context) {
	widgets.Collapse(c, []widgets.CollapseItem{
		{Header: "基础设置", Content: func() {
			widgets.Paragraph(c, "这里是基础设置的内容。")
		}},
		{Header: "高级设置", Content: func() {
			widgets.Paragraph(c, "这里是高级设置的内容。")
		}},
	}, widgets.CollapseAccordion(true), widgets.CollapseDefaultOpen("基础设置"))
}

// /sample

// sample:list
func ListDemo(c *ui.Context) {
	items := []string{"苹果", "香蕉", "橙子", "葡萄"}
	host := ui.Column(c).FillWidth()
	state := ui.Local(host, "ant.sample.display.list", func() ui.ListState { return ui.ListState{} })
	chosen := ui.Local(host, "ant.sample.display.chosen", func() int { return -1 })
	host.Children(func() {
		widgets.List(c, state, len(items), func(i int) {
			widgets.Paragraph(c, items[i])
		}, widgets.ListSelected(chosen))
	})
}

// /sample

// pictureSVG 是 Image 示例用的风景占位图（蓝天、太阳与山丘）。
const pictureSVG = `<?xml version="1.0" encoding="utf-8"?><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 160 120"><rect width="160" height="120" fill="#e6f4ff"/><circle cx="118" cy="34" r="16" fill="#ffd666"/><path d="M0 96 L44 50 L84 96 Z" fill="#95de64"/><path d="M50 96 L100 38 L160 96 Z" fill="#52c41a"/><rect y="96" width="160" height="24" fill="#d9f7be"/></svg>`

// sample:image
func ImageDemo(c *ui.Context) {
	svg, err := ui.ParseSVG([]byte(pictureSVG))
	if err != nil {
		return
	}
	row := ui.Row(c).Gap(c.Theme().Space(4)).AlignItems(ui.Center)
	row.Children(func() {
		widgets.ImageEl(c, svg, widgets.ImageWidth(14))
		widgets.ImageEl(c, svg, widgets.ImageWidth(22), widgets.ImagePreview(false))
	})
}

// /sample

// sample:accordion
func AccordionDemo(c *ui.Context) {
	host := ui.Column(c).FillWidth().Gap(c.Theme().Space(3))
	open := ui.Local(host, "ant.sample.display.collapsible", func() bool { return true })
	general := ui.Local(host, "ant.sample.display.accordion", func() bool { return false })
	host.Children(func() {
		widgets.Collapsible(c, "高级设置（单个折叠块）", open, func() {
			widgets.Paragraph(c, "带展开/收起动画，状态由调用方持有。")
		})
		widgets.Accordion(c, func() {
			widgets.AccordionItem(c, "常规", general, func() {
				widgets.Paragraph(c, "常规设置内容")
			})
			widgets.AccordionItem(c, "隐私", new(bool), func() {
				widgets.Paragraph(c, "隐私设置内容")
			})
		})
		widgets.Meter(c, 72, 0, 100, &ui.MeterLevels{Warning: 60, Critical: 85})
	})
}

// /sample

// sample:masonry
func MasonryDemo(c *ui.Context) {
	host := ui.Column(c).FillWidth().Gap(c.Theme().Space(3))
	active := ui.Local(host, "ant.sample.display.anchor", func() int { return 0 })
	host.Children(func() {
		row := ui.Row(c).FillWidth().Gap(c.Theme().Space(4)).AlignItems(ui.Start)
		row.Children(func() {
			widgets.Anchor(c, []widgets.AnchorItem{
				{ID: "intro", Label: "介绍"},
				{ID: "usage", Label: "用法"},
				{ID: "api", Label: "API"},
			}, active, func(id string) {})
			widgets.BorderBeam(c, func() {
				widgets.Masonry(c, 2, []func(*ui.Context){
					func(c *ui.Context) { widgets.Card(c, "卡片一", func() { widgets.Paragraph(c, "瀑布流内容") }) },
					func(c *ui.Context) { widgets.Card(c, "卡片二", func() { widgets.Paragraph(c, "瀑布流内容") }) },
					func(c *ui.Context) { widgets.Card(c, "卡片三", func() { widgets.Paragraph(c, "瀑布流内容") }) },
				})
			}).Grow(1)
		})
	})
}

// /sample
