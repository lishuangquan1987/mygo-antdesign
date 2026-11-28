package ant

import (
	"strconv"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/icon"
)

// ---------- Empty ----------

// EmptyOption 定制 Empty。
type EmptyOption func(*emptyConfig)

type emptyConfig struct {
	description string
	image       func(*ui.Context)
	footer      func(*ui.Context)
}

// EmptyDescription 覆盖默认描述文字（antd description）。
func EmptyDescription(s string) EmptyOption { return func(c *emptyConfig) { c.description = s } }

// EmptyImage 自定义插画（替换默认的空状态图形）。
func EmptyImage(fn func(*ui.Context)) EmptyOption { return func(c *emptyConfig) { c.image = fn } }

// EmptyFooter 底部操作区（如「新建」按钮）。
func EmptyFooter(fn func(*ui.Context)) EmptyOption { return func(c *emptyConfig) { c.footer = fn } }

// Empty 空状态（antd Empty）：居中的插画 + 描述 + 可选操作区。
//
//	ant.Empty(c, ant.EmptyDescription("暂无数据"))
//	ant.Empty(c, ant.EmptyFooter(func() { ant.Button(c, "新建", ant.Primary()) }))
func Empty(c *ui.Context, opts ...EmptyOption) ui.Element {
	cfg := &emptyConfig{description: "暂无数据"}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	el := ui.Column(c).FillWidth().Padding(t.Space(6), t.Space(4)).
		AlignItems(ui.Center).Gap(t.Space(2))
	el.Children(func() {
		if cfg.image != nil {
			cfg.image(c)
		} else {
			// 默认插画：虚线框 + 文件图标，全部取主题色（antd 空状态 64px 图）。
			frame := ui.Box(c).Size(t.Space(16), t.Space(11)).
				Radius(t.Radius*1.34).
				Border(1, t.Border).BorderStyle(ui.BorderDashed).
				Background(t.SurfaceHover).
				Center()
			frame.Children(func() {
				icon.Icon(c, icon.File, t.Space(5))
			})
		}
		ui.Text(c, cfg.description).TextColor(t.TextMuted).FontSize(t.FontSize)
		if cfg.footer != nil {
			row := ui.Row(c).Gap(t.Space(2)).Padding(t.Space(2), 0, 0, 0)
			row.Children(func() { cfg.footer(c) })
		}
	})
	return el
}

// ---------- Result ----------

// ResultStatus 结果页状态（antd Result status）。
type ResultStatus string

const (
	ResultSuccess ResultStatus = "success"
	ResultInfo    ResultStatus = "info"
	ResultWarning ResultStatus = "warning"
	ResultError   ResultStatus = "error"
)

// ResultOption 定制 Result。
type ResultOption func(*resultConfig)

type resultConfig struct {
	subTitle string
	extra    func(*ui.Context)
}

// ResultSubTitle 副标题（antd subTitle）。
func ResultSubTitle(s string) ResultOption { return func(c *resultConfig) { c.subTitle = s } }

// ResultExtra 操作区（如「返回首页」按钮）。
func ResultExtra(fn func(*ui.Context)) ResultOption { return func(c *resultConfig) { c.extra = fn } }

// Result 结果页（antd Result）：状态图标 + 标题 + 副标题 + 操作区。
//
//	ant.Result(c, ant.ResultSuccess, "操作成功", ant.ResultSubTitle("订单已提交"))
//	ant.Result(c, ant.ResultError, "提交失败", ant.ResultExtra(func() {
//	    ant.Button(c, "重试", ant.Primary())
//	}))
func Result(c *ui.Context, status ResultStatus, title string, opts ...ResultOption) ui.Element {
	cfg := &resultConfig{}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	var (
		tip   = icon.InfoCircle
		color = t.Accent
	)
	switch status {
	case ResultSuccess:
		tip, color = icon.Check, t.Success
	case ResultError:
		tip, color = icon.Close, t.Danger
	case ResultWarning:
		tip, color = icon.InfoCircle, t.Warning
	}

	el := ui.Column(c).FillWidth().Padding(t.Space(8), t.Space(4)).
		AlignItems(ui.Center).Gap(t.Space(2))
	el.Children(func() {
		holder := ui.Box(c).Size(t.Space(18), t.Space(18)).
			Radius(t.Space(9)).
			Background(color.Mix(t.Background, 0.88)).
			Center()
		holder.Children(func() {
			icon.Icon(c, tip, t.Space(9), icon.WithColor(color))
		})
		ui.Text(c, title).FontSize(t.Rem(1.71)).FontWeight(600).TextColor(t.Text).SingleLine()
		if cfg.subTitle != "" {
			ui.Text(c, cfg.subTitle).TextColor(t.TextMuted).FontSize(t.FontSize).SingleLine()
		}
		if cfg.extra != nil {
			row := ui.Row(c).Gap(t.Space(2)).Padding(t.Space(4), 0, 0, 0)
			row.Children(func() { cfg.extra(c) })
		}
	})
	return el
}

// ---------- Skeleton ----------

// SkeletonOption 定制 Skeleton。
type SkeletonOption func(*skeletonConfig)

type skeletonConfig struct {
	rows   int
	avatar bool
	title  bool
}

// SkeletonRows 段落行数（默认 3）。
func SkeletonRows(n int) SkeletonOption { return func(c *skeletonConfig) { c.rows = n } }

// SkeletonAvatar 左侧显示圆形头像占位（antd avatar）。
func SkeletonAvatar(b bool) SkeletonOption { return func(c *skeletonConfig) { c.avatar = b } }

// SkeletonTitle 顶部显示标题条（antd title，默认 true）。
func SkeletonTitle(b bool) SkeletonOption { return func(c *skeletonConfig) { c.title = b } }

// Skeleton 骨架屏（antd Skeleton）。当前为静态骨架，不含呼吸动画。
//
//	ant.Skeleton(c, ant.SkeletonAvatar(true), ant.SkeletonRows(4))
func Skeleton(c *ui.Context, opts ...SkeletonOption) ui.Element {
	cfg := &skeletonConfig{rows: 3, title: true}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()
	fill := t.SurfaceHover

	bar := func(w float32) {
		b := ui.Box(c).Height(t.Space(3)).Width(w).
			Radius(t.Radius * 0.5).Background(fill)
		b.Children(func() {})
	}

	el := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Start)
	el.Children(func() {
		if cfg.avatar {
			av := ui.Box(c).Size(t.Space(8), t.Space(8)).
				Radius(t.Space(4)).Background(fill)
			av.Children(func() {})
		}
		body := ui.Column(c).Gap(t.Space(2)).Grow(1).AlignItems(ui.Start)
		body.Children(func() {
			if cfg.title {
				bar(t.Space(30))
			}
			for i := 0; i < cfg.rows; i++ {
				if i == cfg.rows-1 {
					bar(t.Space(24))
				} else {
					line := ui.Box(c).Height(t.Space(3)).FillWidth().
						Radius(t.Radius * 0.5).Background(fill)
					line.Children(func() {})
				}
			}
		})
	})
	return el
}

// ---------- Statistic ----------

// StatisticOption 定制 Statistic。
type StatisticOption func(*statisticConfig)

type statisticConfig struct {
	prefix    string
	suffix    string
	trend     int // 1 上升 / -1 下降 / 0 无
	precision int
	formatter func(float64) string
}

// StatisticPrefix 数值前缀（如 "¥"）。
func StatisticPrefix(s string) StatisticOption { return func(c *statisticConfig) { c.prefix = s } }

// StatisticSuffix 数值后缀（如 "%"、"次"）。
func StatisticSuffix(s string) StatisticOption { return func(c *statisticConfig) { c.suffix = s } }

// StatisticPrecision 小数位数（antd precision）。
func StatisticPrecision(n int) StatisticOption { return func(c *statisticConfig) { c.precision = n } }

// StatisticTrendUp 上升趋势（中文习惯：红涨）。
func StatisticTrendUp() StatisticOption { return func(c *statisticConfig) { c.trend = 1 } }

// StatisticTrendDown 下降趋势（中文习惯：绿跌）。
func StatisticTrendDown() StatisticOption { return func(c *statisticConfig) { c.trend = -1 } }

// StatisticFormatter 自定义数值格式（antd formatter）；未设置时按 precision
// 保留小数并加千分位。
func StatisticFormatter(fn func(float64) string) StatisticOption {
	return func(c *statisticConfig) { c.formatter = fn }
}

// formatNumber 千分位格式化（antd Statistic 默认分组），支持负数与小数。
func formatNumber(v float64, precision int) string {
	s := strconv.FormatFloat(v, 'f', precision, 64)
	neg := ""
	if len(s) > 0 && s[0] == '-' {
		neg, s = "-", s[1:]
	}
	intPart, frac := s, ""
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			intPart, frac = s[:i], s[i:]
			break
		}
	}
	var b []byte
	for i := 0; i < len(intPart); i++ {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			b = append(b, ',')
		}
		b = append(b, intPart[i])
	}
	return neg + string(b) + frac
}

// Statistic 统计数值（antd Statistic）：标题 + 大字号数值 + 可选前缀/后缀/趋势。
//
//	ant.Statistic(c, "销售额", 126560, ant.StatisticPrefix("¥"), ant.StatisticTrendUp())
//	ant.Statistic(c, "转化率", 12.3, ant.StatisticSuffix("%"), ant.StatisticPrecision(1))
func Statistic(c *ui.Context, title string, value float64, opts ...StatisticOption) ui.Element {
	cfg := &statisticConfig{}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	el := ui.Column(c).Gap(t.Space(1))
	el.Children(func() {
		ui.Text(c, title).TextColor(t.TextMuted).FontSize(t.Rem(0.9))
		row := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(1))
		row.Children(func() {
			if cfg.prefix != "" {
				ui.Text(c, cfg.prefix).TextColor(t.Text).FontSize(t.Rem(1.14))
			}
			text := ""
			if cfg.formatter != nil {
				text = cfg.formatter(value)
			} else {
				text = formatNumber(value, cfg.precision)
			}
			ui.Text(c, text).
				TextColor(t.Text).FontSize(t.Rem(1.86)).FontWeight(600)
			if cfg.suffix != "" {
				ui.Text(c, cfg.suffix).TextColor(t.TextMuted).FontSize(t.Rem(0.9))
			}
			switch cfg.trend {
			case 1:
				icon.Icon(c, icon.Up, t.Rem(0.8), icon.WithColor(t.Danger))
			case -1:
				icon.Icon(c, icon.Down, t.Rem(0.8), icon.WithColor(t.Success))
			}
		})
	})
	return el
}

// ---------- Descriptions ----------

// DescriptionItem 描述列表的一项（antd Descriptions.Item）。
type DescriptionItem struct {
	Label string
	Value string
}

// DescriptionsOption 定制 Descriptions。
type DescriptionsOption func(*descriptionsConfig)

type descriptionsConfig struct {
	columns  int
	title    string
	bordered bool
}

// DescriptionsColumns 每行列数（默认 3）。
func DescriptionsColumns(n int) DescriptionsOption {
	return func(c *descriptionsConfig) { c.columns = n }
}

// DescriptionsTitle 列表标题。
func DescriptionsTitle(s string) DescriptionsOption {
	return func(c *descriptionsConfig) { c.title = s }
}

// DescriptionsBordered 带外框（antd bordered；简化为外框 + 内边距，无内部网格线）。
func DescriptionsBordered(b bool) DescriptionsOption {
	return func(c *descriptionsConfig) { c.bordered = b }
}

// Descriptions 描述列表（antd Descriptions）：标签 / 值成对布局，支持列数与标题。
//
//	ant.Descriptions(c, []ant.DescriptionItem{
//	    {Label: "用户名", Value: "张三"},
//	    {Label: "状态", Value: "已启用"},
//	}, ant.DescriptionsColumns(2), ant.DescriptionsTitle("账号信息"))
func Descriptions(c *ui.Context, items []DescriptionItem, opts ...DescriptionsOption) ui.Element {
	cfg := &descriptionsConfig{columns: 3}
	for _, o := range opts {
		o(cfg)
	}
	if cfg.columns < 1 {
		cfg.columns = 1
	}
	t := c.Theme()

	el := ui.Column(c).FillWidth().Gap(t.Space(2))
	el.Children(func() {
		if cfg.title != "" {
			ui.Text(c, cfg.title).FontSize(t.Rem(1.14)).FontWeight(600).TextColor(t.Text)
		}
		grid := ui.Grid(c).Columns(cfg.columns).Gap(t.Space(2)).FillWidth()
		if cfg.bordered {
			// antd 带边框描述列表：8px 圆角、16px 内边距、浅色标签。
			grid = ui.Grid(c).Columns(cfg.columns).Gap(t.Space(4)).FillWidth().
				Padding(t.Space(4)).
				Radius(t.Radius+t.Space(0.5)).
				Border(1, t.Border).
				Background(t.Surface)
		}
		grid.Children(func() {
			for _, it := range items {
				cell := ui.Column(c).Gap(t.Space(1))
				cell.Children(func() {
					ui.Text(c, it.Label).TextColor(t.TextMuted).FontSize(t.FontSize)
					ui.Text(c, it.Value).TextColor(t.Text).FontSize(t.FontSize).FontWeight(500)
				})
			}
		})
	})
	return el
}

// ---------- Timeline ----------

// TimelineColor 时间轴节点颜色（antd Timeline color）。
type TimelineColor string

const (
	TimelineDefault TimelineColor = "default"
	TimelineBlue    TimelineColor = "blue"
	TimelineGreen   TimelineColor = "green"
	TimelineRed     TimelineColor = "red"
	TimelineGray    TimelineColor = "gray"
)

// TimelineItem 时间轴的一项。
type TimelineItem struct {
	// Title 节点标题。
	Title string
	// Description 节点描述（可选）。
	Description string
	// Color 节点颜色（默认 TimelineDefault = 主色）。
	Color TimelineColor
}

// TimelineOption 定制 Timeline。
type TimelineOption func(*timelineConfig)

type timelineConfig struct{ pending string }

// TimelinePending 在末尾追加一个待办节点（antd Timeline pending）：
// 空心灰点 + 次要色文字，表示尚未发生的步骤。
func TimelinePending(text string) TimelineOption {
	return func(c *timelineConfig) { c.pending = text }
}

func timelineDotColor(t *ui.Theme, c TimelineColor) ui.Color {
	switch c {
	case TimelineGreen:
		return t.Success
	case TimelineRed:
		return t.Danger
	case TimelineGray:
		return t.TextMuted
	case TimelineDefault, TimelineBlue:
		return t.Accent
	default:
		return t.Accent
	}
}

// Timeline 时间轴（antd Timeline）：竖向节点 + 连接线。
//
//	ant.Timeline(c, []ant.TimelineItem{
//	    {Title: "创建订单", Description: "09:00"},
//	    {Title: "已发货", Description: "12:30", Color: ant.TimelineGreen},
//	})
func Timeline(c *ui.Context, items []TimelineItem, opts ...TimelineOption) ui.Element {
	cfg := &timelineConfig{}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()
	// 时间轴连接线的 1px 基准：从主题单位派生，避免裸像素（标准密度下即 1px）。
	hairline := t.Space(0.5)

	el := ui.Column(c).FillWidth()
	el.Children(func() {
		for i, it := range items {
			row := ui.Row(c).FillWidth().AlignItems(ui.Stretch).Gap(t.Space(2))
			row.Children(func() {
				// 左列：节点圆点 + 向下连接线（最后一项不画线）。
				rail := ui.Column(c).Width(t.Space(3)).AlignItems(ui.Center)
				rail.Children(func() {
					dot := ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).
						Radius(t.Space(1.25)).
						Background(timelineDotColor(t, it.Color))
					dot.Children(func() {})
					if i < len(items)-1 {
						line := ui.Box(c).Width(hairline).Grow(1).Background(t.Border)
						line.Children(func() {})
					}
				})
				// 右列：标题与描述。
				body := ui.Column(c).Gap(t.Space(1)).Grow(1).
					Padding(0, 0, t.Space(5), t.Space(3))
				body.Children(func() {
					ui.Text(c, it.Title).TextColor(t.Text).FontSize(t.FontSize)
					if it.Description != "" {
						ui.Text(c, it.Description).TextColor(t.TextMuted).FontSize(t.FontSize)
					}
				})
			})
		}
		// 待办节点：空心灰点 + 次要色文字（antd pending）。
		if cfg.pending != "" {
			row := ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(t.Space(2))
			row.Children(func() {
				rail := ui.Column(c).Width(t.Space(3)).AlignItems(ui.Center)
				rail.Children(func() {
					dot := ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).
						Radius(t.Space(1.25)).
						Border(1, t.Border).
						Background(t.Surface)
					dot.Children(func() {})
				})
				ui.Text(c, cfg.pending).TextColor(t.TextMuted).FontSize(t.FontSize)
			})
		}
	})
	return el
}
