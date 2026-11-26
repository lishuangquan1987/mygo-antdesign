package ant_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

var tableCols = []widgets.TableColumn{
	{Title: "名称", Sortable: true},
	{Title: "年龄", Sortable: true},
}

var tableRows = [][]string{
	{"张三", "34"},
	{"李四", "28"},
	{"王五", "42"},
}

func firstRowTexts(tt *ui.Tester) []string {
	// Texts 顺序即渲染顺序；截取单元格文本（表头之后）。
	return tt.Texts()
}

func TestTableRendersHeadersAndCells(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Table(c, tableCols, tableRows)
	}
	tt := ui.NewTester(view, 400, 200)
	tt.Frame()
	for _, want := range []string{"名称", "年龄", "张三", "李四", "王五"} {
		if !tt.HasText(want) {
			t.Errorf("表格 %q 未渲染", want)
		}
	}
}

func TestTableSortByHeader(t *testing.T) {
	rows := [][]string{{"3", "c"}, {"1", "a"}, {"2", "b"}}
	cols := []widgets.TableColumn{{Title: "编号", Sortable: true}, {Title: "字母"}}
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Table(c, cols, rows)
	}
	tt := ui.NewTester(view, 320, 220)
	tt.Frame()
	// 初始顺序 3 1 2；点击表头按「编号」升序 → 1 2 3。
	if err := tt.Click("编号"); err != nil {
		t.Fatalf("点击表头失败: %v", err)
	}
	tt.Frame()
	texts := tt.Texts()
	pos := func(s string) int {
		for i, x := range texts {
			if x == s {
				return i
			}
		}
		return -1
	}
	if pos("1") > pos("2") || pos("2") > pos("3") {
		t.Errorf("升序排序错误，texts=%v", texts)
	}
	// 再次点击切换降序 → 3 2 1。
	if err := tt.Click("编号"); err != nil {
		t.Fatalf("再次点击表头失败: %v", err)
	}
	tt.Frame()
	texts = tt.Texts()
	if pos("3") > pos("2") || pos("2") > pos("1") {
		t.Errorf("降序排序错误，texts=%v", texts)
	}
}

func TestTableSelection(t *testing.T) {
	var chosen ui.Selection[int]
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Table(c, tableCols, tableRows, widgets.TableSelection(&chosen))
	}
	tt := ui.NewTester(view, 400, 200)
	tt.Frame()
	if err := tt.Click("李四"); err != nil {
		t.Fatalf("点击行失败: %v", err)
	}
	tt.Frame()
	if !chosen.Has(1) {
		t.Error("点击行后未选中 row 1")
	}
	// 单击替换选择；Cmd-click（Linux 为 Ctrl-click）添加/移除。
	if err := tt.Click("王五"); err != nil {
		t.Fatalf("点击第二行失败: %v", err)
	}
	tt.Frame()
	if !chosen.Has(2) || chosen.Has(1) {
		t.Error("单击应替换选择为 row 2")
	}
	if err := tt.ClickWith(ui.Cmd, "李四"); err != nil {
		t.Fatalf("Cmd 点击失败: %v", err)
	}
	tt.Frame()
	if !chosen.Has(1) || !chosen.Has(2) || chosen.Len() != 2 {
		t.Errorf("Cmd-click 未添加 row 1，len=%d", chosen.Len())
	}
	if err := tt.ClickWith(ui.Cmd, "李四"); err != nil {
		t.Fatalf("Cmd 再次点击失败: %v", err)
	}
	tt.Frame()
	if chosen.Has(1) || chosen.Len() != 1 {
		t.Errorf("Cmd-click 未取消 row 1，len=%d", chosen.Len())
	}
}

func TestTableFollowsThemeAndDensity(t *testing.T) {
	for _, tk := range []*theme.AntToken{
		theme.Light(),
		theme.Dark(),
		theme.Light(theme.WithDensity(theme.Compact)),
		theme.Dark(theme.WithDensity(theme.Compact), theme.WithPrimary(ui.Hex("#7c3aed"))),
	} {
		token := tk
		var chosen ui.Selection[int]
		view := func(c *ui.Context) {
			theme.Apply(c, token)
			widgets.Table(c, tableCols, tableRows, widgets.TableSelection(&chosen))
		}
		tt := ui.NewTester(view, 360, 200)
		tt.Frame()
		if tt.FramePending() {
			t.Errorf("密度 %v 下表格持续请求帧", token.Density)
		}
		if !tt.HasText("张三") {
			t.Errorf("密度 %v 下表格未渲染完整", token.Density)
		}
	}
}
