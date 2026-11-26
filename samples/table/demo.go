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
		Summary:  "Table：列定义、点击表头排序、行多选。",
		Samples: []samples.Sample{
			{
				Title:       "基础表格",
				Description: "可排序列与行多选；单击选中、Ctrl-click 加选或移除。",
				Code:        samples.Extract(demoSrc, "basic"),
				View:        Basic,
			},
		},
	}
}

// sample:basic
func Basic(c *ui.Context) {
	col := ui.Column(c).FillWidth()
	chosen := ui.Local(col, "ant.sample.table.chosen", func() ui.Selection[int] { return ui.Selection[int]{} })
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
		}, rows, widgets.TableSelection(chosen))
	})
}

// /sample
