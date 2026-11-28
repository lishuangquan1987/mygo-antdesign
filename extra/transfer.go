package extra

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/icon"
)

// TransferItem 穿梭框的一项。
type TransferItem struct {
	// Key 唯一标识。
	Key string
	// Label 显示文字。
	Label string
}

// TransferOption 定制 Transfer。
type TransferOption func(*transferConfig)

type transferConfig struct {
	leftTitle  string
	rightTitle string
	onChange   func(right []string)
	searchable bool
}

// TransferTitles 两列标题（默认「未选」/「已选」）。
func TransferTitles(left, right string) TransferOption {
	return func(c *transferConfig) { c.leftTitle, c.rightTitle = left, right }
}

// TransferSearchable 两列都显示搜索框（antd showSearch），按 Label 过滤。
func TransferSearchable(b bool) TransferOption {
	return func(c *transferConfig) { c.searchable = b }
}

// TransferOnChange 右列（已选 Key 列表）变化回调。
func TransferOnChange(fn func(right []string)) TransferOption {
	return func(c *transferConfig) { c.onChange = fn }
}

// Transfer 穿梭框（antd Transfer）：左列候选、右列已选，中间按钮移动。
// 交互：点击行选中（可多选），▶ 移到右列、◀ 移回左列。
// right 由调用方持有（已选 Key 列表，保持 items 顺序）：
//
//	ant.Transfer(c, items, &right, ant.TransferOnChange(app.save))
func Transfer(c *ui.Context, items []TransferItem, right *[]string, opts ...TransferOption) ui.Element {
	cfg := &transferConfig{leftTitle: "未选", rightTitle: "已选"}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	host := ui.Row(c).Gap(t.Space(3)).AlignItems(ui.Start)
	chosen := ui.Local(host, "ant.extra.transfer.chosen", func() map[string]bool {
		return map[string]bool{}
	})
	picked := ui.Local(host, "ant.extra.transfer.picked", func() map[string]bool {
		return map[string]bool{}
	})
	initialized := ui.Local(host, "ant.extra.transfer.init", func() bool { return false })
	queryLeft := ui.Local(host, "ant.extra.transfer.q.left", func() string { return "" })
	queryRight := ui.Local(host, "ant.extra.transfer.q.right", func() string { return "" })

	// 首帧：以调用方传入的 right 初始化已选集合。
	if !*initialized {
		for _, k := range *right {
			(*chosen)[k] = true
		}
		*initialized = true
	}

	sync := func() {
		out := make([]string, 0, len(*chosen))
		for _, it := range items {
			if (*chosen)[it.Key] {
				out = append(out, it.Key)
			}
		}
		*right = out
		if cfg.onChange != nil {
			cfg.onChange(out)
		}
	}

	move := func(toRight bool) {
		changed := false
		for _, it := range items {
			if !(*picked)[it.Key] {
				continue
			}
			if toRight != (*chosen)[it.Key] {
				(*chosen)[it.Key] = toRight
				changed = true
			}
			delete(*picked, it.Key)
		}
		if changed {
			sync()
		}
	}

	// 左右两列内容。
	var leftItems, rightItems []TransferItem
	for _, it := range items {
		if (*chosen)[it.Key] {
			rightItems = append(rightItems, it)
		} else {
			leftItems = append(leftItems, it)
		}
	}

	panel := func(title string, list []TransferItem, query *string) ui.Element {
		col := ui.Column(c).Width(t.Space(44)).
			Radius(t.Radius).
			Background(t.Surface).
			Border(1, t.Border)
		col.Children(func() {
			// antd Transfer 头部：40px 高、浅灰底、14px 标题 + 12px 计数。
			head := ui.Row(c).FillWidth().
				Height(t.Space(10)).
				Padding(0, t.Space(4)).
				Background(t.Background.Mix(t.Surface, 0.5)).
				Justify(ui.SpaceBetween).AlignItems(ui.Center)
			head.Border(1, t.Border)
			head.BorderWidth(0, 0, 1, 0)
			head.Children(func() {
				ui.Text(c, title).FontWeight(600).FontSize(t.FontSize).TextColor(t.Text)
				ui.Textf(c, "%d 项", len(list)).TextColor(t.TextMuted).FontSize(t.Rem(0.86))
			})
			if cfg.searchable {
				ui.TextInput(c, query).Placeholder("搜索").Search().
					FillWidth().
					Height(t.Space(7)).
					Margin(t.Space(2)).
					Padding(0, t.Space(2)).
					Radius(t.Radius).
					Border(1, t.Border).
					FontSize(t.FontSize)
			}
			q := strings.ToLower(strings.TrimSpace(*query))
			shown := list
			if cfg.searchable && q != "" {
				shown = nil
				for _, it := range list {
					if strings.Contains(strings.ToLower(it.Label), q) {
						shown = append(shown, it)
					}
				}
			}
			for _, it := range shown {
				key, label := it.Key, it.Label
				row := ui.ButtonBase(c).FillWidth().Height(t.Space(8)).
					Padding(0, t.Space(4)).
					AlignItems(ui.Center).Gap(t.Space(2))
				row.Children(func() {
					// antd Transfer 行首是方形复选框，选中为主色实心 + 勾。
					box := ui.Box(c).Size(t.Space(4), t.Space(4)).
						Radius(t.Space(0.5)).
						Border(1, t.Border).
						Center()
					if (*picked)[key] {
						box.Background(t.Accent)
						box.Border(1, t.Accent)
						box.Children(func() {
							icon.Icon(c, icon.Check, t.Space(2.5), icon.WithColor(t.AccentText))
						})
					} else {
						box.Children(func() {})
					}
					ui.Text(c, label).TextColor(t.Text).FontSize(t.FontSize)
				})
				if (*picked)[key] {
					row.Background(t.Selection)
				} else if row.Hovered() {
					row.Background(t.SurfaceHover)
				}
				if row.Clicked() {
					if (*picked)[key] {
						delete(*picked, key)
					} else {
						(*picked)[key] = true
					}
				}
			}
			if len(shown) == 0 {
				empty := ui.Row(c).FillWidth().Height(t.Space(12)).Center()
				empty.Children(func() {
					ui.Text(c, "暂无数据").TextColor(t.TextMuted).FontSize(t.Rem(0.86))
				})
			}
		})
		return col
	}

	host.Children(func() {
		panel(cfg.leftTitle, leftItems, queryLeft)
		mid := ui.Column(c).Gap(t.Space(2)).AlignItems(ui.Center)
		mid.Children(func() {
			// antd Transfer 中间操作：32×32 描边按钮 + 箭头。
			arrow := func(svg string, toRight bool) {
				btn := ui.ButtonBase(c).Size(t.Space(8), t.Space(8)).
					Radius(t.Radius).Border(1, t.Border).
					Background(t.Surface).Center()
				if toRight {
					btn.Label("移到右侧")
				} else {
					btn.Label("移回左侧")
				}
				if btn.Hovered() {
					btn.Border(1, t.Accent)
				}
				btn.Children(func() {
					color := t.TextMuted
					if btn.Hovered() {
						color = t.Accent
					}
					icon.Icon(c, svg, t.Space(3), icon.WithColor(color))
				})
				if btn.Clicked() {
					move(toRight)
				}
			}
			arrow(icon.Right, true)
			arrow(icon.Left, false)
		})
		panel(cfg.rightTitle, rightItems, queryRight)
	})
	return host
}
