package extra

import (
	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/icon"
	"github.com/lishuangquan1987/mygo-antdesign/internal"
)

// triggerBox Cascader/TreeSelect 共用的下拉触发器，对齐 antd Select
// 触发器样式：32px 高、11px 水平内边距、6px 圆角、右端下拉箭头，
// 占位文字弱化，悬停时描边转品牌色。
func triggerBox(c *ui.Context, t *ui.Theme, text string, isPlaceholder bool, disabled bool) ui.Element {
	el := ui.Row(c).
		Height(internal.ControlHeight(t, internal.SizeMiddle)).
		Padding(0, internal.ControlPaddingX(t, internal.SizeMiddle)).
		Radius(t.Radius).
		MinWidth(t.Space(40)).
		Gap(t.Space(2)).
		Justify(ui.SpaceBetween).
		AlignItems(ui.Center)
	if disabled {
		el.Background(t.Background).Border(1, t.Border)
	} else {
		el.Background(t.Surface).Border(1, t.Border)
	}
	if !disabled && el.Hovered() {
		el.Border(1, t.Accent)
	}
	el.Children(func() {
		color := t.Text
		if isPlaceholder {
			color = t.TextMuted
		}
		ui.Text(c, text).TextColor(color).
			FontSize(internal.ControlFont(t, internal.SizeMiddle)).SingleLine()
		// 下拉箭头使用 antd 的 chevron-down 图标。
		icon.Icon(c, icon.Down, t.Space(3), icon.WithColor(t.TextMuted))
	})
	return el
}
