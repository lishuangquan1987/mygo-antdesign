package ant

import (
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/icon"
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
	selection    *ui.Selection[int]
	onSort       func()
	render       func(c *ui.Context, row, col int) ui.Element
	pageSize     int
	current      *int
	onPageChange func(page int)
	virtualH     float32
}

// TableRender 自定义单元格内容（antd column.render）：传入原始行下标与
// 列下标，调用方在 cell 内构建任意元素（图标、标签、按钮等）：
//
//	ant.Table(c, cols, rows, ant.TableRender(func(c *ui.Context, row, col int) ui.Element {
//	    if col == 1 {
//	        return ant.Tag(c, rows[row][col], ant.TagSuccess())
//	    }
//	    return ant.Text(c, rows[row][col])
//	}))
func TableRender(fn func(c *ui.Context, row, col int) ui.Element) TableOption {
	return func(c *tableConfig) { c.render = fn }
}

// TablePagination 分页展示（antd pagination）：pageSize 每页行数，current
// 为当前页（1 起，由调用方持有），切换时回调。
func TablePagination(pageSize int, current *int, onPageChange func(page int)) TableOption {
	return func(c *tableConfig) {
		c.pageSize, c.current, c.onPageChange = pageSize, current, onPageChange
	}
}

// TableVirtual 虚拟滚动（行数多时只构建可见行）：height 为表格主体最大
// 高度，超出后滚动。启用行多选时自动退回全量渲染（多选需要全部行的
// 命中区）。
func TableVirtual(height float32) TableOption {
	return func(c *tableConfig) { c.virtualH = height }
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
							caret := icon.Up
							color := t.TextMuted
							if sortOrder.Column == id {
								if sortOrder.Descending {
									caret = icon.Down
								}
								color = t.Accent
							}
							icon.Icon(c, caret, t.Rem(0.6), icon.WithColor(color))
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

			// 可见行：排序后的原始行下标，按分页切片。
			view := make([]int, 0, len(rows))
			for r := range rows {
				idx := r
				if r < len(*order) {
					idx = (*order)[r]
				}
				view = append(view, idx)
			}
			if cfg.pageSize > 0 {
				pages := (len(rows) + cfg.pageSize - 1) / cfg.pageSize
				if pages < 1 {
					pages = 1
				}
				if cfg.current != nil {
					if *cfg.current < 1 {
						*cfg.current = 1
					}
					if *cfg.current > pages {
						*cfg.current = pages
					}
				}
				page := 1
				if cfg.current != nil {
					page = *cfg.current
				}
				start := (page - 1) * cfg.pageSize
				end := start + cfg.pageSize
				if start > len(view) {
					start = len(view)
				}
				if end > len(view) {
					end = len(view)
				}
				view = view[start:end]
			}

			// 行内容：普通行与虚拟滚动行共用。
			buildCells := func(idx int, rowEl ui.Element) {
				rowEl.Children(func() {
					for col := range columns {
						col := col
						if cfg.render != nil {
							cell := ui.Row(c).
								Padding(t.Space(4), t.Space(4)).
								AlignItems(ui.Center)
							if columns[col].Width > 0 {
								cell.Width(columns[col].Width)
							} else {
								cell.Grow(1)
							}
							applyCellAlign(cell, columns[col].Align)
							cell.Children(func() { cfg.render(c, idx, col) })
							continue
						}
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
			}

			// 限高滚动：设置了 TableVirtual 且未启用多选时只构建可见行
			// （antd scroll.y）。行数很少时同样走这条路——此时等价于给
			// 表格加一个固定高度的滚动区，不会带来额外开销。
			if cfg.virtualH > 0 && cfg.selection == nil {
				lst := ui.List(c, nil, len(view), func(i int) {
					idx := view[i]
					rowEl := ui.Row(c).FillWidth().
						Border(1, t.Border).BorderWidth(0, 0, 1, 0)
					if i == len(view)-1 {
						rowEl.BorderWidth(0, 0, 0, 0)
					}
					if rowEl.Hovered() {
						rowEl.Background(t.Background.Mix(t.Surface, 0.5))
					}
					buildCells(idx, rowEl)
				}).Height(cfg.virtualH)
				lst.FillWidth()
			} else {
				// 数据行：1px 底部分隔线（最后一行不画），可选选中/悬停底色。
				for r, idx := range view {
					last := r == len(view)-1
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
					buildCells(idx, rowEl)
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
			}
			if len(rows) == 0 {
				ui.Text(c, "暂无数据").
					TextColor(t.TextMuted).
					FontSize(t.FontSize).
					Padding(t.Space(6)).
					Center()
			}
		})

		// 分页器：表格外框下方（antd Table pagination 位置）。
		if cfg.pageSize > 0 && cfg.current != nil && len(rows) > cfg.pageSize {
			pages := (len(rows) + cfg.pageSize - 1) / cfg.pageSize
			ui.Row(c).FillWidth().Justify(ui.End).
				Padding(t.Space(3), 0, 0, 0).
				Children(func() {
					Pagination(c, cfg.current, pages, OnPageChange(func(page int) {
						if cfg.onPageChange != nil {
							cfg.onPageChange(page)
						}
					}))
				})
		}
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
