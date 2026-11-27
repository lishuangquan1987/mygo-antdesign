package ant

import (
	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/icon"
)

// CollapseItem 折叠面板的一项（antd Collapse.Panel）。
type CollapseItem struct {
	// Key 唯一标识（用于记住展开状态），留空时用 Header。
	Key string
	// Header 头部标题。
	Header string
	// Content 展开后构建的内容。
	Content func()
}

func (it CollapseItem) key() string {
	if it.Key != "" {
		return it.Key
	}
	return it.Header
}

// CollapseOption 定制 Collapse。
type CollapseOption func(*collapseConfig)

type collapseConfig struct {
	accordion   bool
	defaultOpen map[string]bool
}

// CollapseAccordion 手风琴模式：同一时间只展开一项（antd accordion）。
func CollapseAccordion(b bool) CollapseOption {
	return func(c *collapseConfig) { c.accordion = b }
}

// CollapseDefaultOpen 初始展开的项（按 Key 或 Header 匹配）。
func CollapseDefaultOpen(keys ...string) CollapseOption {
	return func(c *collapseConfig) {
		if c.defaultOpen == nil {
			c.defaultOpen = map[string]bool{}
		}
		for _, k := range keys {
			c.defaultOpen[k] = true
		}
	}
}

// Collapse 折叠面板（antd Collapse）：带边框的容器，头部点击展开/收起，
// 展开状态按 Key 持久化：
//
//	ant.Collapse(c, []ant.CollapseItem{
//	    {Header: "基础设置", Content: func() { ... }},
//	    {Header: "高级设置", Content: func() { ... }},
//	}, ant.CollapseAccordion(true))
func Collapse(c *ui.Context, items []CollapseItem, opts ...CollapseOption) ui.Element {
	cfg := &collapseConfig{}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	host := ui.Column(c).FillWidth().
		Border(1, t.Border).Radius(t.Radius).Background(t.Surface)
	open := ui.Local(host, "ant.collapse.open", func() map[string]bool {
		m := make(map[string]bool, len(items))
		for _, it := range items {
			m[it.key()] = cfg.defaultOpen[it.key()]
		}
		return m
	})
	for _, it := range items {
		if _, ok := (*open)[it.key()]; !ok {
			(*open)[it.key()] = cfg.defaultOpen[it.key()]
		}
	}

	host.Children(func() {
		for i, it := range items {
			it := it
			if i > 0 {
				// 项之间的分隔线（1px，从主题单位派生）。
				line := ui.Box(c).FillWidth().Height(t.Space(1) / 4).Background(t.Border)
				line.Children(func() {})
			}
			opened := (*open)[it.key()]
			head := ui.ButtonBase(c).FillWidth().
				Height(t.Space(8)).
				Padding(0, t.Space(3)).
				Background(t.Surface)
			head.Children(func() {
				row := ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(t.Space(2))
				row.Children(func() {
					arrow := icon.Right
					if opened {
						arrow = icon.Down
					}
					icon.Icon(c, arrow, t.Space(4), icon.WithColor(t.TextMuted))
					ui.Text(c, it.Header).TextColor(t.Text).FontSize(t.FontSize).FontWeight(500)
				})
			})
			if head.Clicked() {
				(*open)[it.key()] = !opened
				if !opened && cfg.accordion {
					for _, other := range items {
						if other.key() != it.key() {
							(*open)[other.key()] = false
						}
					}
				}
			}
			if opened && it.Content != nil {
				body := ui.Column(c).FillWidth().Padding(t.Space(3)).
					Background(t.Surface)
				body.Children(it.Content)
			}
		}
	})
	return host
}
