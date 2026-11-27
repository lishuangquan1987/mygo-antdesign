package ant

import (
	"github.com/egoist/mygo/ui"
)

// ListOption 定制 List。
type ListOption func(*listConfig)

type listConfig struct {
	selected  *int
	selection *ui.Selection[int]
}

// ListSelected 绑定选中行（antd List 的选中态；-1 表示未选）。
func ListSelected(p *int) ListOption { return func(c *listConfig) { c.selected = p } }

// ListSelection 启用多选（按行下标记录；Ctrl/Cmd-click 加选，Shift-click 连选）。
func ListSelection(sel *ui.Selection[int]) ListOption {
	return func(c *listConfig) { c.selection = sel }
}

// List 列表（antd List）：虚拟滚动由 MyGo 的 ui.List 提供，行内容由
// row 构建，统一加 antd 风格的横向内边距：
//
//	ant.List(c, &app.list, len(rows), func(i int) {
//	    ant.Paragraph(c, rows[i])
//	}, ant.ListSelected(&app.chosen))
func List(c *ui.Context, state *ui.ListState, n int, row func(int), opts ...ListOption) ui.Element {
	cfg := &listConfig{}
	for _, o := range opts {
		o(cfg)
	}
	if cfg.selected != nil {
		state.Selected = cfg.selected
	}
	if cfg.selection != nil {
		state.Selection = cfg.selection
	}
	t := c.Theme()

	return ui.List(c, state, n, func(i int) {
		el := ui.Row(c).FillWidth().Padding(t.Space(3), t.Space(4)).AlignItems(ui.Center).
			Border(1, t.Border)
		el.BorderWidth(0, 0, 1, 0) // 行分隔线
		// antd List.Item：选中/悬停时给整行浅底。
		if cfg.selection != nil && cfg.selection.Has(i) {
			el.Background(t.Selection)
		} else if el.Hovered() {
			el.Background(t.Background.Mix(t.Surface, 0.5))
		}
		el.Children(func() { row(i) })
	}).FillWidth()
}
