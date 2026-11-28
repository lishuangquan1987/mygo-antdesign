// Package table 提供表格控件的示例。
package table

import (
	_ "embed"

	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

//go:embed demo.go
var demoSrc string

// Demo 返回表格控件的示例集合。
func Demo() samples.WidgetDemo {
	return samples.WidgetDemo{
		Name:     "表格",
		Package:  "widgets",
		Category: samples.CatDisplay,
		Summary:  "Table：列定义、表头排序、限高滚动、行多选与自定义单元格。",
		Samples: []samples.Sample{
			{
				Title:       "基础表格",
				Description: "可排序列表头；TableVirtual 限高，滚动看全部行。",
				Code:        samples.Extract(demoSrc, "basic"),
				View:        Basic,
				Page:        "Table",
			},
			{
				Title:       "行多选与自定义单元格",
				Description: "TableSelection 行多选，TableRender 自绘单元格。",
				Code:        samples.Extract(demoSrc, "render"),
				View:        Render,
				Page:        "Table",
			},
		},
	}
}

// sample:basic
func Basic(c *ui.Context) {
	col := ui.Column(c).FillWidth()
	col.Children(func() {
		rows := [][]string{
			{"张三", "研发部", "34"},
			{"李四", "产品部", "28"},
			{"王五", "运营部", "42"},
			{"赵六", "研发部", "25"},
		}
		widgets.Table(c, []widgets.TableColumn{
			{Title: "姓名", Sortable: true},
			{Title: "部门", Sortable: true},
			{Title: "年龄", Sortable: true},
		}, rows, widgets.TableVirtual(c.Theme().Space(30)))
	})
}

// /sample

// sample:render
func Render(c *ui.Context) {
	col := ui.Column(c).FillWidth()
	chosen := ui.Local(col, "ant.sample.table.render.chosen", func() ui.Selection[int] {
		var sel ui.Selection[int]
		sel.Add(0) // 初始选中首行，示例一进来就能看到选中态
		return sel
	})
	col.Children(func() {
		rows := [][]string{
			{"张三", "研发部", "在职"},
			{"李四", "产品部", "试用"},
			{"王五", "运营部", "离职"},
		}
		widgets.Table(c, []widgets.TableColumn{
			{Title: "姓名", Sortable: true},
			{Title: "部门"},
			{Title: "状态"},
		}, rows,
			widgets.TableSelection(chosen),
			widgets.TableRender(func(c *ui.Context, row, col int) ui.Element {
				switch col {
				case 1:
					return widgets.Tag(c, rows[row][col], widgets.TagPrimary())
				case 2:
					if rows[row][col] == "离职" {
						return widgets.Tag(c, rows[row][col], widgets.TagError())
					}
					return widgets.Tag(c, rows[row][col], widgets.TagSuccess())
				default:
					return ui.Text(c, rows[row][col])
				}
			}),
		)
	})
}

// /sample
