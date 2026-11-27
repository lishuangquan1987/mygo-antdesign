package extra

import (
	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/internal"
)

// triggerBox Cascader/TreeSelect 共用的下拉触发器，对齐 antd Select
// 触发器样式：输入框高度、边框、右端下拉箭头，占位文字弱化。
func triggerBox(c *ui.Context, t *ui.Theme, text string, isPlaceholder bool, disabled bool) ui.Element {
	el := ui.Row(c).
		Height(internal.ControlHeight(t, internal.SizeMiddle)).
		Padding(0, internal.ControlPaddingX(t, internal.SizeMiddle)).
		Radius(t.Radius).
		MinWidth(t.Space(40)).
		Gap(t.Space(2)).
		AlignItems(ui.Center)
	if disabled {
		el.Background(t.Background).Border(1, t.Border)
	} else {
		el.Background(t.Surface).Border(1, t.Border)
	}
	el.Children(func() {
		color := t.Text
		if isPlaceholder {
			color = t.TextMuted
		}
		ui.Text(c, text).TextColor(color).
			FontSize(internal.ControlFont(t, internal.SizeMiddle))
		// 下拉箭头。
		ui.Text(c, "▾").TextColor(t.TextMuted)
	})
	return el
}
