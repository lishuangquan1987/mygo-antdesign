// Package layout 提供 Space / Flex / Grid 布局控件的示例。
package layout

import (
	_ "embed"
	"strconv"

	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

//go:embed demo.go
var demoSrc string

// Demo 返回布局控件的示例集合。
func Demo() samples.WidgetDemo {
	return samples.WidgetDemo{
		Name:     "布局",
		Package:  "widgets",
		Category: samples.CatGeneral,
		Summary:  "布局：Space 间距容器、Flex 弹性布局、Grid 网格。",
		Samples: []samples.Sample{
			{
				Title:       "Space 间距容器",
				Description: "给一组子元素统一加间距，支持三档尺寸与纵向排列。",
				Code:        samples.Extract(demoSrc, "space"),
				View:        SpaceDemo,
			},
			{
				Title:       "Flex 弹性布局",
				Description: "主轴对齐、交叉轴对齐、间距与换行。",
				Code:        samples.Extract(demoSrc, "flex"),
				View:        FlexDemo,
			},
			{
				Title:       "Grid 网格",
				Description: "等宽网格，间距随密度缩放。",
				Code:        samples.Extract(demoSrc, "grid"),
				View:        GridDemo,
			},
			{
				Title:       "分隔面板",
				Description: "Splitter 左右两栏，中间分隔条可拖拽调整宽度。",
				Code:        samples.Extract(demoSrc, "splitter"),
				View:        SplitterDemo,
			},
			{
				Title:       "悬浮按钮",
				Description: "FloatButton 固定在窗口右下角的悬浮按钮。",
				Code:        samples.Extract(demoSrc, "floatbutton"),
				View:        FloatButtonDemo,
			},
			{
				Title:       "固钉",
				Description: "Affix 把内容固定在滚动容器顶部。",
				Code:        samples.Extract(demoSrc, "affix"),
				View:        AffixDemo,
			},
			{
				Title:       "百分比与垂直分隔",
				Description: "SplitterPercent 百分比分栏；SplitterVerticalPercent 上下按比例分栏。",
				Code:        samples.Extract(demoSrc, "splitterpercent"),
				View:        SplitterDemo2,
			},
			{
				Title:       "垂直分隔",
				Description: "SplitterVertical 上下分栏，尺寸由调用方持有。",
				Code:        samples.Extract(demoSrc, "verticalsplit"),
				View:        VerticalSplitDemo,
			},
		},
	}
}

// sample:space
func SpaceDemo(c *ui.Context) {
	col := ui.Column(c).FillWidth().Gap(12)
	col.Children(func() {
		widgets.Space(c).Children(func() {
			widgets.Button(c, "取消")
			widgets.Button(c, "默认")
			widgets.Button(c, "确定", widgets.Primary())
		})
		small := widgets.Space(c, widgets.Gap(widgets.SpaceSmall), widgets.SpaceVertical(),
			widgets.SpaceAlign(ui.Start), widgets.SpaceFillWidth())
		small.Children(func() {
			widgets.Paragraph(c, "纵向小间距一")
			widgets.Paragraph(c, "纵向小间距二")
		})
	})
}

// /sample

// sample:flex
func FlexDemo(c *ui.Context) {
	col := ui.Column(c).FillWidth().Gap(12)
	col.Children(func() {
		head := widgets.Flex(c, widgets.FlexJustify(ui.SpaceBetween), widgets.FlexAlign(ui.Center), widgets.FlexFill())
		head.Children(func() {
			widgets.Title(c, "用户管理", widgets.H4)
			widgets.Space(c, widgets.GapUnits(1)).Children(func() {
				widgets.Button(c, "刷新", widgets.Text())
				widgets.Button(c, "新建", widgets.Primary())
			})
		})
		wrap := widgets.Flex(c, widgets.FlexWrap(), widgets.FlexGap(2), widgets.FlexFill())
		wrap.Children(func() {
			widgets.Button(c, "筛选条件一", widgets.Text())
			widgets.Button(c, "筛选条件二", widgets.Text())
			widgets.Button(c, "筛选条件三", widgets.Text())
		})
	})
}

// /sample

// sample:grid
func GridDemo(c *ui.Context) {
	grid := widgets.Grid(c, 3)
	grid.Children(func() {
		for _, name := range []string{"概览", "明细", "统计"} {
			widgets.Title(c, name, widgets.H5)
		}
	})
}

// /sample

// sample:splitter
func SplitterDemo(c *ui.Context) {
	host := ui.Column(c).FillWidth().Height(c.Theme().Space(32))
	size := ui.Local(host, "ant.sample.layout.split", func() float32 { return 200 })
	host.Children(func() {
		widgets.Splitter(c, size, func() {
			widgets.Paragraph(c, "左侧面板")
		}, func() {
			widgets.Paragraph(c, "右侧面板")
		}).Fill()
	})
}

// /sample

// sample:floatbutton
func FloatButtonDemo(c *ui.Context) {
	host := ui.Column(c).FillWidth()
	count := ui.Local(host, "ant.sample.layout.float", func() int { return 0 })
	host.Children(func() {
		widgets.Paragraph(c, "右下角悬浮按钮固定在窗口上。")
		widgets.Paragraph(c, "点击次数："+strconv.Itoa(*count), widgets.Secondary())
		widgets.FloatButton(c, "↑", func() { *count++ })
	})
}

// /sample

// sample:affix
func AffixDemo(c *ui.Context) {
	host := ui.Column(c).FillWidth()
	host.Children(func() {
		widgets.Affix(c, 0, func() {
			box := ui.Box(c).FillWidth().Padding(c.Theme().Space(2)).
				Background(c.Theme().Surface).Border(1, c.Theme().Border)
			box.Children(func() { ui.Text(c, "吸顶工具条（滚动时保持顶部）") })
		})
		for i := 0; i < 12; i++ {
			widgets.Paragraph(c, "可滚动内容行")
		}
	})
}

// /sample

// sample:splitterpercent
func SplitterDemo2(c *ui.Context) {
	t := c.Theme()
	host := ui.Column(c).FillWidth().Height(c.Theme().Space(40)).Gap(c.Theme().Space(2))
	percent := ui.Local(host, "ant.sample.layout.splitpercent", func() float32 { return 30 })
	verticalPercent := ui.Local(host, "ant.sample.layout.splitverticalpercent", func() float32 { return 50 })
	host.Children(func() {
		// 横向分栏占上部（Grow 撑满剩余高度），竖向分栏固定高度，
		// 避免 Fill 百分比在无界方向上解析异常。
		row := ui.Row(c).FillWidth().Grow(1).AlignItems(ui.Stretch)
		row.Children(func() {
			widgets.SplitterPercent(c, percent, func() {
				widgets.Paragraph(c, "左栏 30%")
			}, func() {
				widgets.Paragraph(c, "右栏")
			}).Fill()
		})
		widgets.SplitterVerticalPercent(c, verticalPercent, func() {
			widgets.Paragraph(c, "上栏（50%）")
		}, func() {
			widgets.Paragraph(c, "下栏")
		}).FillWidth().Height(t.Space(20))
	})
}

// sample:verticalsplit
func VerticalSplitDemo(c *ui.Context) {
	host := ui.Column(c).FillWidth().Height(c.Theme().Space(30))
	size := ui.Local(host, "ant.sample.layout.splitvertical", func() float32 { return 100 })
	host.Children(func() {
		widgets.SplitterVertical(c, size, func() {
			widgets.Paragraph(c, "上栏（固定 100 DIP）")
		}, func() {
			widgets.Paragraph(c, "下栏")
		})
	})
}

// /sample
