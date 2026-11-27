package extra

import (
	"github.com/egoist/mygo/ui"
)

// TreeSelectOption 定制 TreeSelect。
type TreeSelectOption func(*treeSelectConfig)

type treeSelectConfig struct {
	disabled    bool
	placeholder string
	onSelect    func(key string)
}

// TreeSelectDisabled 禁用触发器。
func TreeSelectDisabled(b bool) TreeSelectOption {
	return func(c *treeSelectConfig) { c.disabled = b }
}

// TreeSelectPlaceholder 未选中时触发器显示的占位文字。
func TreeSelectPlaceholder(s string) TreeSelectOption {
	return func(c *treeSelectConfig) { c.placeholder = s }
}

// TreeSelectOnSelect 选中树节点（叶子或分支）后的回调。
func TreeSelectOnSelect(fn func(key string)) TreeSelectOption {
	return func(c *treeSelectConfig) { c.onSelect = fn }
}

// TreeSelect 树形下拉选择（antd TreeSelect）：触发器显示当前选中
// 节点，点击展开含树形结构的浮层，选中节点后浮层收起：
//
//	ant.TreeSelect(c, &app.org, []ant.TreeNode{
//	    {Label: "研发部", Key: "rd", Children: []ant.TreeNode{
//	        {Label: "前端组", Key: "fe"},
//	    }},
//	}, ant.TreeSelectOnSelect(func(key string) { app.setOrg(key) }))
func TreeSelect(c *ui.Context, selected *string, nodes []TreeNode, opts ...TreeSelectOption) ui.Element {
	cfg := &treeSelectConfig{placeholder: "请选择"}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	host := ui.Row(c)
	open := ui.Local(host, "ant.extra.treeselect.open", func() bool { return false })

	// 触发器本身作为浮层的 anchor（保证浮层紧贴触发器下方）。
	trigger := triggerBox(c, t, triggerSelectedText(*selected, nodes, cfg.placeholder), *selected == "", cfg.disabled)
	if trigger.Clicked() && !cfg.disabled {
		*open = true
	}
	ui.Popover(c, trigger, open, func() {
		panel := ui.Column(c).
			Width(t.Space(56)).
			MaxHeight(t.Space(40)).
			Padding(t.Space(1), 0).
			Background(t.Surface).
			Border(1, t.Border).
			Radius(t.Radius)
		panel.Children(func() {
			Tree(c, nodes, OnTreeNodeSelect(func(key string) {
				*selected = key
				*open = false
				if cfg.onSelect != nil {
					cfg.onSelect(key)
				}
			}))
		})
	})
	return host
}

// triggerSelectedText 触发器文字：selected 为 key，翻译成节点 Label。
func triggerSelectedText(selected string, nodes []TreeNode, placeholder string) string {
	if selected == "" {
		return placeholder
	}
	if label, ok := findTreeNodeLabel(nodes, selected); ok {
		return label
	}
	return selected
}

// findTreeNodeLabel 递归查 key 对应的 Label。
func findTreeNodeLabel(nodes []TreeNode, key string) (string, bool) {
	for _, n := range nodes {
		k := n.Key
		if k == "" {
			k = n.Label
		}
		if k == key {
			return n.Label, true
		}
		if len(n.Children) > 0 {
			if label, ok := findTreeNodeLabel(n.Children, key); ok {
				return label, true
			}
		}
	}
	return "", false
}
