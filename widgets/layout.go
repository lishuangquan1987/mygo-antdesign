package ant

import (
	"github.com/egoist/mygo/ui"
)

// SpaceSize Space 组件的间距档（antd size：small 8 / middle 16 / large 24px）。
type SpaceSize int

const (
	SpaceSmall  SpaceSize = iota // 2 单位（标准密度 8px）
	SpaceMiddle                  // 4 单位（标准密度 16px，默认）
	SpaceLarge                   // 6 单位（标准密度 24px）
)

func (s SpaceSize) units() float32 {
	switch s {
	case SpaceSmall:
		return 2
	case SpaceLarge:
		return 6
	default:
		return 4
	}
}

// SpaceOption 定制 Space。
type SpaceOption func(*spaceConfig)

type spaceConfig struct {
	size      SpaceSize
	custom    bool
	gapUnits  float32
	vertical  bool
	align     ui.Align
	fillWidth bool
}

// Gap 设置间距档（antd size）。
func Gap(size SpaceSize) SpaceOption { return func(c *spaceConfig) { c.size = size } }

// GapUnits 覆盖间距（单位数，经 t.Space 换算，随密度缩放）。
func GapUnits(units float32) SpaceOption {
	return func(c *spaceConfig) { c.custom = true; c.gapUnits = units }
}

// SpaceVertical 纵向排列（antd direction="vertical"）。
func SpaceVertical() SpaceOption { return func(c *spaceConfig) { c.vertical = true } }

// SpaceAlign 对齐方式（默认 ui.Center，对齐 antd align="center"）。
func SpaceAlign(a ui.Align) SpaceOption { return func(c *spaceConfig) { c.align = a } }

// SpaceFillWidth 宽度撑满。
func SpaceFillWidth() SpaceOption { return func(c *spaceConfig) { c.fillWidth = true } }

// Space 间距容器（antd Space）：给一组子元素统一加间距，子元素通过
// 返回值的 .Children(func() { ... }) 提供：
//
//	ant.Space(c, ant.SpaceAlign(ui.Center)).Children(func() {
//	    ant.Button(c, "取消")
//	    ant.Button(c, "确定", ant.Primary())
//	})
func Space(c *ui.Context, opts ...SpaceOption) ui.Element {
	cfg := &spaceConfig{size: SpaceMiddle, align: ui.Center}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()
	gap := cfg.size.units()
	if cfg.custom {
		gap = cfg.gapUnits
	}

	if cfg.vertical {
		el := ui.Column(c).Gap(t.Space(gap)).AlignItems(cfg.align)
		if cfg.fillWidth {
			el.FillWidth()
		}
		return el
	}
	el := ui.Row(c).Gap(t.Space(gap)).AlignItems(cfg.align)
	if cfg.fillWidth {
		el.FillWidth()
	}
	return el
}

// FlexOption 定制 Flex。
type FlexOption func(*flexConfig)

type flexConfig struct {
	vertical bool
	justify  ui.Align
	align    ui.Align
	gapUnits float32
	wrap     bool
	grow     float32
	hasGrow  bool
	fill     bool
}

// FlexColumn 主轴改为纵向（antd vertical）。
func FlexColumn() FlexOption { return func(c *flexConfig) { c.vertical = true } }

// FlexJustify 主轴对齐（Start/Center/End/SpaceBetween…）。
func FlexJustify(a ui.Align) FlexOption { return func(c *flexConfig) { c.justify = a } }

// FlexAlign 交叉轴对齐。
func FlexAlign(a ui.Align) FlexOption { return func(c *flexConfig) { c.align = a } }

// FlexGap 间距（单位数，经 t.Space 换算）。
func FlexGap(units float32) FlexOption { return func(c *flexConfig) { c.gapUnits = units } }

// FlexWrap 允许换行（antd wrap）。
func FlexWrap() FlexOption { return func(c *flexConfig) { c.wrap = true } }

// FlexGrow 拉伸比例。
func FlexGrow(f float32) FlexOption { return func(c *flexConfig) { c.grow, c.hasGrow = f, true } }

// FlexFill 撑满父容器（宽 + 高）。
func FlexFill() FlexOption { return func(c *flexConfig) { c.fill = true } }

// Flex 弹性布局容器（antd Flex）：统一封装主轴方向、对齐与间距：
//
//	ant.Flex(c, ant.FlexJustify(ui.SpaceBetween), ant.FlexAlign(ui.Center)).Children(func() {
//	    ant.Title(c, "标题", ant.H4)
//	    ant.Button(c, "新建", ant.Primary())
//	})
func Flex(c *ui.Context, opts ...FlexOption) ui.Element {
	cfg := &flexConfig{align: ui.Stretch}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	var el ui.Element
	if cfg.vertical {
		el = ui.Column(c)
	} else {
		el = ui.Row(c)
	}
	el.Justify(cfg.justify).AlignItems(cfg.align)
	if cfg.gapUnits > 0 {
		el.Gap(t.Space(cfg.gapUnits))
	}
	if cfg.wrap && !cfg.vertical {
		el.Wrap()
	}
	if cfg.hasGrow {
		el.Grow(cfg.grow)
	}
	if cfg.fill {
		el.Fill()
	}
	return el
}

// Grid 网格布局（antd Row/Col 的栅格简化版）：cols 列，间距
// t.Space(3)（标准密度 12px）：
//
//	ant.Grid(c, 3).Children(func() {
//	    ant.Card(c, "卡片一", ...)
//	    ant.Card(c, "卡片二", ...)
//	})
func Grid(c *ui.Context, cols int) ui.Element {
	t := c.Theme()
	return ui.Grid(c).Columns(cols).Gap(t.Space(3)).FillWidth()
}
