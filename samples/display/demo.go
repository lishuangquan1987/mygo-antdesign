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
		Summary:  "Tag 多彩标签、Badge 徽标、Avatar 头像、Card 卡片。",
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
