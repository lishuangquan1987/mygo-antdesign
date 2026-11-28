package ant_test

import (
	"strconv"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

// BenchmarkTable1k 测千行表格整帧构建成本；虚拟滚动（TableVirtual）应
// 显著低于全量渲染。运行：
//
//	go test ./widgets -bench BenchmarkTable1k -benchmem
func BenchmarkTable1k(b *testing.B) {
	rows := make([][]string, 1000)
	for i := range rows {
		rows[i] = []string{
			"订单 " + strconv.Itoa(i+1),
			"¥" + strconv.Itoa(100+i),
			"已完成",
		}
	}
	cols := []widgets.TableColumn{{Title: "名称"}, {Title: "金额"}, {Title: "状态"}}
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Table(c, cols, rows)
	}
	tt := ui.NewTester(view, 900, 600)
	tt.Frame()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tt.Frame()
	}
}

// BenchmarkTable1kVirtual 测启用虚拟滚动的同规模表格。
func BenchmarkTable1kVirtual(b *testing.B) {
	rows := make([][]string, 1000)
	for i := range rows {
		rows[i] = []string{
			"订单 " + strconv.Itoa(i+1),
			"¥" + strconv.Itoa(100+i),
			"已完成",
		}
	}
	cols := []widgets.TableColumn{{Title: "名称"}, {Title: "金额"}, {Title: "状态"}}
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Table(c, cols, rows, widgets.TableVirtual(420))
	}
	tt := ui.NewTester(view, 900, 600)
	tt.Frame()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tt.Frame()
	}
}
