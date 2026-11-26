package extra

import (
	"github.com/egoist/mygo/ui"
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
}

// TransferTitles 两列标题（默认「未选」/「已选」）。
func TransferTitles(left, right string) TransferOption {
	return func(c *transferConfig) { c.leftTitle, c.rightTitle = left, right }
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

	panel := func(title string, list []TransferItem) ui.Element {
		col := ui.Column(c).Width(t.Space(44)).
			Radius(t.Radius).
			Background(t.Surface).
			Border(1, t.Border)
		col.Children(func() {
			head := ui.Row(c).FillWidth().Padding(t.Space(2), t.Space(3)).
				Justify(ui.SpaceBetween).AlignItems(ui.Center)
			head.Border(1, t.Border)
			head.BorderWidth(0, 0, 1, 0)
			head.Children(func() {
				ui.Text(c, title).FontWeight(600).TextColor(t.Text)
				ui.Textf(c, "%d 项", len(list)).TextColor(t.TextMuted).FontSize(t.Rem(0.86))
			})
			for _, it := range list {
				key, label := it.Key, it.Label
				row := ui.ButtonBase(c).FillWidth().Height(t.Space(6)).
					Padding(0, t.Space(3))
				row.Children(func() {
					mark := "○"
					if (*picked)[key] {
						mark = "●"
					}
					ui.Text(c, mark).TextColor(t.Accent)
					ui.Text(c, label).TextColor(t.Text)
				})
				if (*picked)[key] {
					row.Background(t.Selection)
				}
				if row.Clicked() {
					if (*picked)[key] {
						delete(*picked, key)
					} else {
						(*picked)[key] = true
					}
				}
			}
			if len(list) == 0 {
				empty := ui.Row(c).FillWidth().Height(t.Space(6)).Center()
				empty.Children(func() {
					ui.Text(c, "空").TextColor(t.TextMuted).FontSize(t.Rem(0.86))
				})
			}
		})
		return col
	}

	host.Children(func() {
		panel(cfg.leftTitle, leftItems)
		mid := ui.Column(c).Gap(t.Space(2)).AlignItems(ui.Center)
		mid.Children(func() {
			if ui.Button(c, "▶").Clicked() {
				move(true)
			}
			if ui.Button(c, "◀").Clicked() {
				move(false)
			}
		})
		panel(cfg.rightTitle, rightItems)
	})
	return host
}
