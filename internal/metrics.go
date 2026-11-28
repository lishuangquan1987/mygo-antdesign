// Package internal 存放换肤控件共享的状态映射与度量 helper。
// 样式铁律要求所有控件的颜色/间距经由此处或 theme 包取值（AGENTS.md 2.1）。
package internal

import "github.com/egoist/mygo/ui"

// SizeKind 控件尺寸档（对齐 antd size: large / middle / small）。
type SizeKind int

const (
	SizeSmall SizeKind = iota
	SizeMiddle
	SizeLarge
)

// ControlHeight 控件标准高度：large 40 / middle 32 / small 24（antd
// controlHeightLG/HS）。数值以主题 Spacing 为单位表达（4×10 / 4×8 / 4×6），
// 紧凑密度（Spacing=3）下自动缩为 30 / 24 / 18。
func ControlHeight(t *ui.Theme, s SizeKind) float32 {
	switch s {
	case SizeLarge:
		return t.Space(10)
	case SizeSmall:
		return t.Space(6)
	default:
		return t.Space(8)
	}
}

// ControlPaddingX 控件水平内边距（antd paddingContentHorizontal）。
// 小档 12 / 中大档 16：两档只差 4px（antd 为 7/11 的比例），同一列里
// 不同尺寸控件的文字起始位置肉眼基本对齐。
func ControlPaddingX(t *ui.Theme, s SizeKind) float32 {
	switch s {
	case SizeSmall:
		return t.Space(3)
	default:
		return t.Space(4)
	}
}

// ControlFont 控件字号：large 约为基准字号的 1.15 倍（antd 16px），
// middle / small 用基准字号。
func ControlFont(t *ui.Theme, s SizeKind) float32 {
	if s == SizeLarge {
		return t.Rem(1.15)
	}
	return t.FontSize
}
