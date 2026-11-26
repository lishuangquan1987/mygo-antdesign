package ant

import (
	"slices"
	"sort"

	"github.com/egoist/mygo/ui"
)

// TableColumn 表格列定义（antd Table 的 columns）。
type TableColumn struct {
	// Title 表头文字；ID 标识列（排序/布局），留空时用 Title。
	Title string
	ID    string
	// Width 列宽（DIP），0 表示与其他列分摊剩余空间。
	Width float32
	// Sortable 允许点击表头排序（antd sorter）。
	Sortable bool
	// Align 单元格内容对齐。
	Align ui.Align
}

// TableOption 定制 Table。
type TableOption func(*tableConfig)

type tableConfig struct {
	selection *ui.Selection[int]
	onSort    func()
}

// TableSelection 启用行多选（点击行选中/再点取消），选中集由调用方持有，
// 用 ui.Selection[int].Has(row) 查询（antd rowSelection）：
//
//	var chosen ui.Selection[int]
//	ant.Table(c, cols, rows, ant.TableSelection(&chosen))
func TableSelection(sel *ui.Selection[int]) TableOption {
	return func(c *tableConfig) { c.selection = sel }
}

// TableOnSort 排序变化回调（点击表头触发）。
func TableOnSort(fn func()) TableOption { return func(c *tableConfig) { c.onSort = fn } }

// Table 数据表格（antd Table）。rows 为字符串单元格（行为主序），
// 排序按字符串比较（数字列请先格式化为定宽文本，或后续版本扩展
// 自定义比较器）。排序状态由库内持久化，点击 Sortable 列表头升序/
// 降序切换，数据顺序自动重排：
//
//	ant.Table(c, []ant.TableColumn{{Title: "名称", Sortable: true}}, rows)
func Table(c *ui.Context, columns []TableColumn, rows [][]string, opts ...TableOption) ui.Element {
	cfg := &tableConfig{}
	for _, o := range opts {
		o(cfg)
	}

	// 跨帧状态：ListState（原语要求持久）、排序方向与排序后的行下标，
	// 绑定在宿主元素上持久化。
	host := ui.Column(c).FillWidth()
	state := ui.Local(host, "ant.table.state", func() ui.ListState { return ui.ListState{} })
	sortOrder := ui.Local(host, "ant.table.sort", func() ui.SortOrder { return ui.SortOrder{} })
	sortedFor := ui.Local(host, "ant.table.sortedFor", func() ui.SortOrder {
		return ui.SortOrder{Column: "\x00initial\x00"}
	})
	order := ui.Local(host, "ant.table.order", func() []int { return nil })

	// 原语在表头点击时更新 *sortOrder；本帧按其重排行下标。
	state.Sort = sortOrder
	if cfg.selection != nil {
		state.Selection = cfg.selection
	}
	if *sortOrder != *sortedFor || len(*order) != len(rows) {
		*order = sortIndexes(columns, rows, sortOrder)
		*sortedFor = *sortOrder
		if cfg.onSort != nil {
			cfg.onSort()
		}
	}

	cols := make([]ui.TableColumn, len(columns))
	for i, col := range columns {
		id := col.ID
		if id == "" {
			id = col.Title
		}
		cols[i] = ui.TableColumn{
			Title:    col.Title,
			ID:       id,
			Width:    col.Width,
			Align:    col.Align,
			Sortable: col.Sortable,
		}
	}

	host.Children(func() {
		ui.Table(c, state, cols, len(rows), func(row, col int) {
			i := row
			if row < len(*order) {
				i = (*order)[row]
			}
			if col < len(rows[i]) {
				ui.Text(c, rows[i][col])
			}
		}).FillWidth()
	})
	return host
}

// sortIndexes 返回按 sortOrder 排好的行下标；无有效排序列时保持原序。
func sortIndexes(columns []TableColumn, rows [][]string, sortOrder *ui.SortOrder) []int {
	order := make([]int, len(rows))
	for i := range order {
		order[i] = i
	}
	if sortOrder == nil || sortOrder.Column == "" || len(rows) == 0 {
		return order
	}
	col := slices.IndexFunc(columns, func(col TableColumn) bool {
		id := col.ID
		if id == "" {
			id = col.Title
		}
		return id == sortOrder.Column
	})
	if col < 0 {
		return order
	}
	sort.SliceStable(order, func(a, b int) bool {
		ra, rb := order[a], order[b]
		cellA := cellAt(rows, ra, col)
		cellB := cellAt(rows, rb, col)
		if sortOrder.Descending {
			return cellB < cellA
		}
		return cellA < cellB
	})
	return order
}

func cellAt(rows [][]string, row, col int) string {
	if row < len(rows) && col < len(rows[row]) {
		return rows[row][col]
	}
	return ""
}
