// Package extra 提供 Ant Design 有而 MyGo 原语没有的缺口控件
// （基于 ui.Tree / ui.ButtonBase / ui.Row / ui.Overlay 等底座组合，
// 见 docs/03）。
//
// 样式铁律与 widgets 一致：颜色只读 c.Theme()，间距只用 t.Space / t.Rem。
package extra

import (
	"github.com/egoist/mygo/ui"
)

// TreeNode 树节点（antd Tree 的 treeData 项）。
type TreeNode struct {
	// Label 节点文字。
	Label string
	// Key 节点标识（OnSelect 回调参数），留空时用 Label。
	Key string
	// Children 子节点；非空时渲染为可展开分支，空时为叶子。
	Children []TreeNode
}

// TreeOption 定制 Tree。
type TreeOption func(*treeConfig)

type treeConfig struct {
	onSelect      func(key string)
	defaultExpand bool
}

// OnTreeNodeSelect 节点选中回调（叶子或分支被点击时触发）。
func OnTreeNodeSelect(fn func(key string)) TreeOption {
	return func(c *treeConfig) { c.onSelect = fn }
}

// TreeDefaultExpand 分支默认展开（antd Tree 的 defaultExpandAll）。
// 首帧即展示层级关系；用户折叠后仍以用户操作为准。
func TreeDefaultExpand() TreeOption {
	return func(c *treeConfig) { c.defaultExpand = true }
}

// Tree 树形控件（antd Tree）。基于 MyGo Tree 原语：分支自带展开/折叠
// 箭头与缩进，叶子与分支点击触发 OnTreeNodeSelect：
//
//	ant.Tree(c, []ant.TreeNode{
//	    {Label: "src", Children: []ant.TreeNode{{Label: "main.go"}}},
//	}, ant.OnTreeNodeSelect(func(key string) { app.open(key) }))
func Tree(c *ui.Context, nodes []TreeNode, opts ...TreeOption) ui.Element {
	cfg := &treeConfig{}
	for _, o := range opts {
		o(cfg)
	}
	host := ui.Column(c).FillWidth()
	// 分支展开状态跨帧持久（Key → 展开指针）。
	opens := ui.Local(host, "ant.extra.tree.opens", func() map[string]*bool {
		return map[string]*bool{}
	})
	host.Children(func() {
		ui.Tree(c, func() {
			buildNodes(c, nodes, cfg, *opens)
		})
	})
	return host
}

// buildNodes 递归构建节点：分支展开状态由 opens 跨帧持有。
func buildNodes(c *ui.Context, nodes []TreeNode, cfg *treeConfig, opens map[string]*bool) {
	for _, n := range nodes {
		key := n.Key
		if key == "" {
			key = n.Label
		}
		if len(n.Children) == 0 {
			// 叶子：children 与 open 均为 nil。
			item := ui.TreeItem(c, n.Label, nil, nil)
			if item.Clicked() && cfg.onSelect != nil {
				cfg.onSelect(key)
			}
			continue
		}
		// 分支：展开状态跨帧持有（初始值由 defaultExpand 决定）。
		item := ui.TreeItem(c, n.Label, branchOpen(opens, key, cfg.defaultExpand), func() {
			buildNodes(c, n.Children, cfg, opens)
		})
		if item.Clicked() && cfg.onSelect != nil {
			cfg.onSelect(key)
		}
	}
}

// branchOpen 取分支的跨帧展开状态指针。
func branchOpen(opens map[string]*bool, key string, dft bool) *bool {
	b, ok := opens[key]
	if !ok {
		b = new(bool)
		*b = dft
		opens[key] = b
	}
	return b
}
