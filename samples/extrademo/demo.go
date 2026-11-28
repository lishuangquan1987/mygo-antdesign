// Package extrademo 提供缺口控件（extra 包）的示例。
package extrademo

import (
	_ "embed"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/extra"
	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

//go:embed demo.go
var demoSrc string

// Demo 返回缺口控件的示例集合。
func Demo() samples.WidgetDemo {
	return samples.WidgetDemo{
		Name:     "缺口控件",
		Package:  "extra",
		Category: samples.CatDisplay,
		Summary:  "Tree / Cascader / TreeSelect / Upload / Transfer / Carousel / QRCode：基于底座组合实现。",
		Samples: []samples.Sample{
			{
				Title:       "树形控件",
				Description: "Tree 分支展开/折叠，叶子与分支点击选中。",
				Code:        samples.Extract(demoSrc, "tree"),
				View:        TreeDemo,
			},
			{
				Title:       "级联选择",
				Description: "Cascader 逐级下钻的多列浮层，叶子完成选择。",
				Code:        samples.Extract(demoSrc, "cascader"),
				View:        CascaderDemo,
			},
			{
				Title:       "树形下拉选择",
				Description: "TreeSelect 触发器展开树形浮层，选中节点后收起。",
				Code:        samples.Extract(demoSrc, "treeselect"),
				View:        TreeSelectDemo,
			},
			{
				Title:       "上传",
				Description: "Upload 点击选择文件（系统对话框）或拖入，列表可移除。",
				Code:        samples.Extract(demoSrc, "upload"),
				View:        UploadDemo,
			},
			{
				Title:       "穿梭框",
				Description: "Transfer 点击行选中，▶ 移入右列、◀ 移回。",
				Code:        samples.Extract(demoSrc, "transfer"),
				View:        TransferDemo,
			},
			{
				Title:       "走马灯与二维码",
				Description: "Carousel 指示点切换页面；QRCode 码色随主题。",
				Code:        samples.Extract(demoSrc, "carousel"),
				View:        CarouselDemo,
			},
		},
	}
}

// sample:tree
func TreeDemo(c *ui.Context) {
	col := ui.Column(c).FillWidth()
	picked := ui.Local(col, "ant.sample.extra.tree", func() string { return "" })
	col.Children(func() {
		nodes := []extra.TreeNode{
			{Label: "src", Children: []extra.TreeNode{
				{Label: "main.go", Key: "main"},
				{Label: "view.go", Key: "view"},
				{Label: "widgets", Children: []extra.TreeNode{
					{Label: "button.go", Key: "btn"},
					{Label: "input.go", Key: "input"},
				}},
			}},
			{Label: "theme", Children: []extra.TreeNode{
				{Label: "token.go", Key: "token"},
			}},
			{Label: "go.mod", Key: "gomod"},
		}
		extra.Tree(c, nodes, extra.TreeDefaultExpand(),
			extra.OnTreeNodeSelect(func(key string) { *picked = key }))
		if *picked != "" {
			widgets.Tag(c, "已选 "+*picked, widgets.TagSuccess())
		}
	})
}

// /sample

// sample:cascader
func CascaderDemo(c *ui.Context) {
	col := ui.Column(c).FillWidth()
	region := ui.Local(col, "ant.sample.extra.cascader", func() []string { return nil })
	col.Children(func() {
		extra.Cascader(c, region, []extra.CascaderNode{
			{Label: "浙江", Children: []extra.CascaderNode{
				{Label: "杭州"},
				{Label: "宁波"},
			}},
			{Label: "广东", Children: []extra.CascaderNode{{Label: "广州"}}},
		}, extra.CascaderPlaceholder("请选择地区"))
		if len(*region) > 0 {
			widgets.Tag(c, strings.Join(*region, " / "), widgets.TagSuccess())
		}
	})
}

// /sample

// sample:treeselect
func TreeSelectDemo(c *ui.Context) {
	col := ui.Column(c).FillWidth()
	org := ui.Local(col, "ant.sample.extra.treeselect", func() string { return "" })
	col.Children(func() {
		extra.TreeSelect(c, org, []extra.TreeNode{
			{Label: "研发部", Key: "rd", Children: []extra.TreeNode{
				{Label: "前端组", Key: "fe"},
				{Label: "后端组", Key: "be"},
			}},
			{Label: "产品部", Key: "pm"},
		}, extra.TreeSelectPlaceholder("选择部门"))
	})
}

// /sample

// sample:upload
func UploadDemo(c *ui.Context) {
	col := ui.Column(c).FillWidth()
	count := ui.Local(col, "ant.sample.extra.upload", func() int { return 0 })
	col.Children(func() {
		// 只接受图片与 PDF，最多 3 个（超出部分被忽略）。
		extra.Upload(c, extra.UploadMultiple(true),
			extra.UploadAccept("PNG / JPG / PDF"),
			extra.UploadExtensions("png", "jpg", "pdf"),
			extra.UploadMaxCount(3),
			extra.UploadOnChange(func(files []extra.UploadFile) { *count = len(files) }))
		if *count > 0 {
			widgets.TextMuted(c, "已选择文件（见列表）")
		}
	})
}

// /sample

// sample:transfer
func TransferDemo(c *ui.Context) {
	col := ui.Column(c).FillWidth()
	right := ui.Local(col, "ant.sample.extra.transfer", func() []string { return []string{"b"} })
	col.Children(func() {
		items := []extra.TransferItem{
			{Key: "a", Label: "列表项 A"},
			{Key: "b", Label: "列表项 B"},
			{Key: "c", Label: "列表项 C"},
		}
		extra.Transfer(c, items, right,
			extra.TransferTitles("候选", "已选"),
			extra.TransferSearchable(true))
	})
}

// /sample

// sample:carousel
func CarouselDemo(c *ui.Context) {
	t := c.Theme()
	col := ui.Column(c).FillWidth().Gap(t.Space(5))
	col.Children(func() {
		// 箭头 + 圆点 + 3 秒自动播放（系统偏好「减少动态」时自动暂停）。
		extra.Carousel(c, []func(*ui.Context){
			func(c *ui.Context) { widgets.Title(c, "欢迎使用", widgets.H3) },
			func(c *ui.Context) { widgets.Title(c, "第二页", widgets.H3) },
			func(c *ui.Context) { widgets.Title(c, "第三页", widgets.H3) },
		}, extra.CarouselHeight(96),
			extra.CarouselArrows(true),
			extra.CarouselDots(true),
			extra.CarouselAutoplay(3*time.Second))
		row := ui.Row(c).Gap(t.Space(5)).AlignItems(ui.Center)
		row.Children(func() {
			extra.QRCode(c, "https://github.com/lishuangquan1987/mygo-antdesign", extra.QRCodeSize(88))
			info := ui.Column(c).Gap(t.Space(1))
			info.Children(func() {
				widgets.Title(c, "扫码访问项目仓库", widgets.H5)
				widgets.Paragraph(c, "github.com/lishuangquan1987/mygo-antdesign", widgets.Secondary())
			})
		})
	})
}

// /sample
