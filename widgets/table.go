package ant

import (
	"slices"
	"sort"
	"strconv"
	"strings"

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
// 降序切换，数据顺序自动重排。
//
// 样式对齐 antd v5 的默认表格：8px 圆角外框、#fafafa 表头（深色
// 主题自动派生）、16px 单元格内边距、1px 行分隔线；启用 TableSelection
// 时行可点击并显示选中/悬停底色。
func Table(c *ui.Context, columns []TableColumn, rows [][]string, opts ...TableOption) ui.Element {
	cfg := &tableConfig{}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	// 跨帧状态：排序方向与排序后的行下标，绑定在宿主元素上持久化。
	host := ui.Column(c).FillWidth()
	sortOrder := ui.Local(host, "ant.table.sort", func() ui.SortOrder { return ui.SortOrder{} })
	sortedFor := ui.Local(host, "ant.table.sortedFor", func() ui.SortOrder {
		return ui.SortOrder{Column: "\x00initial\x00"}
	})
	order := ui.Local(host, "ant.table.order", func() []int { return nil })

	if *sortOrder != *sortedFor || len(*order) != len(rows) {
		*order = sortIndexes(columns, rows, sortOrder)
		*sortedFor = *sortOrder
		if cfg.onSort != nil {
			cfg.onSort()
		}
	}

	host.Children(func() {
		// 外框：antd 表格默认圆角 8（主题 Radius 6 + 2），1px 描边。
		table := ui.Column(c).FillWidth().
			Background(t.Surface).
			Radius(t.Radius+t.Space(0.5)).
			Border(1, t.Border)
		table.Children(func() {
			// 表头：浅灰底 + 底部分隔线 + 600 字重。
			header := ui.Row(c).FillWidth().
				Background(t.Background.Mix(t.Surface, 0.5)).
				Border(1, t.Border).
				BorderWidth(0, 0, 1, 0)
			header.Children(func() {
				for _, col := range columns {
					col := col
					id := col.ID
					if id == "" {
						id = col.Title
					}
					cell := ui.ButtonBase(c).
						Padding(t.Space(4), t.Space(4)).
						AlignItems(ui.Center).
						Gap(t.Space(1))
					if col.Width > 0 {
						cell.Width(col.Width)
					} else {
						cell.Grow(1)
					}
					applyCellAlign(cell, col.Align)
					cell.Children(func() {
						ui.Text(c, col.Title).
							TextColor(t.Text).
							FontWeight(600).
							FontSize(t.FontSize)
						if col.Sortable {
							caret := "▲"
							color := t.TextMuted
							if sortOrder.Column == id {
								if sortOrder.Descending {
									caret = "▼"
								}
								color = t.Accent
							}
							ui.Text(c, caret).TextColor(color).FontSize(t.Rem(0.6))
						}
					})
					if col.Sortable && cell.Clicked() {
						if sortOrder.Column == id {
							sortOrder.Descending = !sortOrder.Descending
						} else {
							*sortOrder = ui.SortOrder{Column: id, Descending: false}
						}
					}
				}
			})

			// 数据行：1px 底部分隔线（最后一行不画），可选选中/悬停底色。
			for r := range rows {
				idx := r
				if r < len(*order) {
					idx = (*order)[r]
				}
				last := r == len(rows)-1
				var rowEl ui.Element
				if cfg.selection != nil {
					rowEl = ui.ButtonBase(c).FillWidth()
				} else {
					rowEl = ui.Row(c).FillWidth()
				}
				rowEl.Border(1, t.Border).BorderWidth(0, 0, 1, 0)
				selected := cfg.selection != nil && cfg.selection.Has(idx)
				switch {
				case selected:
					rowEl.Background(t.Selection)
				case cfg.selection != nil && rowEl.Hovered():
					rowEl.Background(t.Background.Mix(t.Surface, 0.5))
				}
				rowEl.Children(func() {
					for col := range columns {
						col := col
						value := ""
						if col < len(rows[idx]) {
							value = rows[idx][col]
						}
						cell := ui.Row(c).
							Padding(t.Space(4), t.Space(4)).
							AlignItems(ui.Center)
						if col < len(columns) && columns[col].Width > 0 {
							cell.Width(columns[col].Width)
						} else {
							cell.Grow(1)
						}
						applyCellAlign(cell, columns[col].Align)
						cell.Children(func() {
							ui.Text(c, value).TextColor(t.Text).FontSize(t.FontSize)
						})
					}
				})
				if cfg.selection != nil && rowEl.Clicked() {
					// antd rowSelection：普通单击替换选择；Cmd/Shift 单击切换
					// 该行（web 网格习惯用 Shift，macOS 用 Cmd）。
					if rowEl.ClickModifiers()&(ui.Cmd|ui.Shift) != 0 {
						if cfg.selection.Has(idx) {
							cfg.selection.Remove(idx)
						} else {
							cfg.selection.Add(idx)
						}
					} else {
						cfg.selection.Clear()
						cfg.selection.Add(idx)
					}
				}
				if last {
					rowEl.BorderWidth(0, 0, 0, 0)
				}
			}
			if len(rows) == 0 {
				ui.Text(c, "暂无数据").
					TextColor(t.TextMuted).
					FontSize(t.FontSize).
					Padding(t.Space(6)).
					Center()
			}
		})
	})
	return host
}

// applyCellAlign 设置单元格内容对齐（antd column.align）。
func applyCellAlign(el ui.Element, align ui.Align) {
	switch align {
	case ui.Center:
		el.Justify(ui.Center)
	case ui.End:
		el.Justify(ui.End)
	default:
		el.Justify(ui.Start)
	}
}

// colAlign 返回第 col 列的对齐方式。
func colAlign(columns []TableColumn, col int) ui.Align {
	for i := range columns {
		if i == col {
			return columns[i].Align
		}
	}
	return ui.Start
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
		// 数字列按数值比较（「9」排在「10」前），其余按字符串比较。
		if na, okA := parseSortNumber(cellA); okA {
			if nb, okB := parseSortNumber(cellB); okB {
				if sortOrder.Descending {
					return nb < na
				}
				return na < nb
			}
		}
		if sortOrder.Descending {
			return cellB < cellA
		}
		return cellA < cellB
	})
	return order
}

// parseSortNumber 解析去掉千分位逗号与 % 后缀的数值；非数字返回 false。
func parseSortNumber(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "%")
	s = strings.ReplaceAll(s, ",", "")
	if s == "" {
		return 0, false
	}
	n, err := strconv.ParseFloat(s, 64)
	return n, err == nil
}

func cellAt(rows [][]string, row, col int) string {
	if row < len(rows) && col < len(rows[row]) {
		return rows[row][col]
	}
	return ""
}
