package extra

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/icon"
	"github.com/lishuangquan1987/mygo-antdesign/internal"
)

// CascaderNode 级联选项节点（antd Cascader options 项）。叶子为
// Children 为空的节点。
type CascaderNode struct {
	// Label 选项文字（也是选中路径 path 中的元素）。
	Label string
	// Children 下一级选项；为空时本节点是叶子。
	Children []CascaderNode
}

// CascaderOption 定制 Cascader。
type CascaderOption func(*cascaderConfig)

type cascaderConfig struct {
	disabled    bool
	placeholder string
	onChange    func(path []string)
}

// CascaderDisabled 禁用触发器。
func CascaderDisabled(b bool) CascaderOption {
	return func(c *cascaderConfig) { c.disabled = b }
}

// CascaderPlaceholder 未选中时触发器显示的占位文字。
func CascaderPlaceholder(s string) CascaderOption {
	return func(c *cascaderConfig) { c.placeholder = s }
}

// CascaderOnChange 选中叶子（完成一次级联选择）后的回调，path 为
// 从根到叶子的 Label 路径。
func CascaderOnChange(fn func(path []string)) CascaderOption {
	return func(c *cascaderConfig) { c.onChange = fn }
}

// Cascader 级联选择（antd Cascader）。点击触发器展开浮层，浮层按
// 下钻深度逐级显示选项列，点击叶子项完成选择并把路径写入 path：
//
//	ant.Cascader(c, &app.region, []ant.CascaderNode{
//	    {Label: "浙江", Children: []ant.CascaderNode{
//	        {Label: "杭州"},
//	        {Label: "宁波"},
//	    }},
//	    {Label: "广东", Children: []ant.CascaderNode{{Label: "广州"}}},
//	}, ant.CascaderOnChange(func(path []string) { app.load(path) }))
func Cascader(c *ui.Context, path *[]string, nodes []CascaderNode, opts ...CascaderOption) ui.Element {
	cfg := &cascaderConfig{placeholder: "请选择"}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	host := ui.Row(c)
	open := ui.Local(host, "ant.extra.cascader.open", func() bool { return false })
	// 下钻下标路径：跨帧持有，让浮层在多次开关间保持浏览进度。
	idx := ui.Local(host, "ant.extra.cascader.idx", func() []int { return nil })

	// 触发器本身作为浮层的 anchor（保证浮层紧贴触发器下方）。
	trigger := triggerBox(c, t, triggerPathText(*path, cfg.placeholder), len(*path) == 0, cfg.disabled)
	if trigger.Clicked() && !cfg.disabled {
		*open = true
	}
	ui.Popover(c, trigger, open, func() {
		panel := ui.Row(c).
			Background(t.Surface).
			Border(1, t.Border).
			Radius(t.Radius)
		panel.Children(func() {
			renderCascaderColumns(c, t, nodes, idx, open, path, cfg)
		})
	})
	return host
}

// renderCascaderColumns 按 idx 下钻深度逐级渲染选项列。
func renderCascaderColumns(c *ui.Context, t *ui.Theme, nodes []CascaderNode, idx *[]int, open *bool, path *[]string, cfg *cascaderConfig) {
	cur := nodes
	for level := 0; ; level++ {
		column := ui.Column(c).MinHeight(t.Space(16)).MinWidth(t.Space(24))
		level := level
		column.Children(func() {
			for i, n := range cur {
				i := i
				hasChild := len(n.Children) > 0
				// ButtonBase 让每个选项可聚焦：Tab 移动、Enter/Space 选择
				// （纯 Row 只有鼠标命中区，键盘无法操作）。
				item := ui.ButtonBase(c).
					FillWidth().
					Padding(t.Space(1), t.Space(2)).
					Gap(t.Space(2)).
					AlignItems(ui.Center).
					Label(n.Label)
				item.Children(func() {
					ui.Text(c, n.Label).TextColor(t.Text).
						FontSize(internal.ControlFont(t, internal.SizeMiddle)).Grow(1)
					if hasChild {
						icon.Icon(c, icon.Right, t.Space(2.5), icon.WithColor(t.TextMuted))
					}
				})
				// 悬停高亮：antd menu-item-hover。
				if item.Hovered() && !cfg.disabled {
					item.Background(t.SurfaceHover)
				}
				if item.Clicked() && !cfg.disabled {
					// 更新下钻路径到本层（截断更深层的旧选择）。
					next := append(append([]int{}, (*idx)[:min(level, len(*idx))]...), i)
					*idx = next
					if !hasChild {
						// 叶子：写入选中路径并收起浮层。
						labels := cascadeLabels(nodes, next)
						*path = labels
						*idx = nil
						*open = false
						if cfg.onChange != nil {
							cfg.onChange(labels)
						}
					}
				}
			}
		})
		// 下钻到下一级：cur 为选中节点的 children；无则结束。
		if level >= len(*idx) {
			break
		}
		next := (*idx)[level]
		if next >= len(cur) || len(cur[next].Children) == 0 {
			break
		}
		cur = cur[next].Children
	}
}

// cascadeLabels 按下标路径取各级 Label（路径必然以叶子结尾）。
func cascadeLabels(nodes []CascaderNode, idx []int) []string {
	labels := make([]string, 0, len(idx))
	cur := nodes
	for _, i := range idx {
		if i >= len(cur) {
			break
		}
		labels = append(labels, cur[i].Label)
		cur = cur[i].Children
	}
	return labels
}

// triggerPathText 触发器文字：选中路径用 " / " 连接，未选中用占位。
func triggerPathText(path []string, placeholder string) string {
	if len(path) == 0 {
		return placeholder
	}
	return strings.Join(path, " / ")
}
