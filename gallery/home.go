package main

import (
	"github.com/egoist/mygo/ui"
	"github.com/lishuangquan1987/mygo-antdesign/icon"
	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

// catIcon 分类图标（图标库为 antd outlined 子集；未收录的组件图标
// 在卡片上回退到分类图标）。
var catIcon = map[samples.Category]string{
	samples.CatGeneral:  icon.Appstore,
	samples.CatLayout:   icon.Menu,
	samples.CatNav:      icon.Export,
	samples.CatInput:    icon.Edit,
	samples.CatDisplay:  icon.Picture,
	samples.CatFeedback: icon.Bell,
	samples.CatOther:    icon.More,
}

// pageIcon 组件页图标（对齐 antd 组件总览「每组件一图标」的形制；
// 按可用图标近似映射，未列出的组件回退分类图标）。
var pageIcon = map[string]string{
	"Button": icon.Appstore, "FloatButton": icon.Plus, "Title": icon.File,
	"Paragraph": icon.File, "Code": icon.File,
	"Hyperlink": icon.Export, "Divider": icon.Minus,
	"Splitter": icon.Menu, "SplitterPercent": icon.Menu,
	"SplitterVerticalPercent": icon.Menu, "SplitterVertical": icon.Menu,
	"Affix": icon.Lock, "BorderBeam": icon.Star, "IDE 风格": icon.Setting,
	"Anchor": icon.Export, "Breadcrumb": icon.Right, "Dropdown": icon.Down,
	"Pagination": icon.Left, "Steps": icon.Right, "Tabs": icon.Folder,
	"Sidebar": icon.Menu,
	"Input":   icon.Edit, "TextArea": icon.Edit, "SearchInput": icon.Search,
	"SearchField": icon.Search, "NumberInput": icon.Edit,
	"DecimalInput": icon.Edit, "Stepper": icon.Up, "DateInput": icon.Calendar,
	"TimeInput": icon.Calendar, "Calendar": icon.Calendar, "Form": icon.Edit,
	"Autocomplete": icon.Search, "Mentions": icon.User, "Checkbox": icon.Check,
	"CheckboxGroup": icon.Check, "Radio": icon.Check, "Switch": icon.Check,
	"Toggle": icon.Check, "ToggleGroup": icon.Check, "Slider": icon.Filter,
	"RangeSlider": icon.Filter, "Rating": icon.Star, "Segmented": icon.Filter,
	"Select": icon.Down, "TokenField": icon.More, "Cascader": icon.Folder,
	"TreeSelect": icon.Folder, "Tree": icon.Folder, "Upload": icon.Upload,
	"Transfer": icon.Export, "ColorPicker": icon.Eye, "ColorWell": icon.Eye,
	"Avatar": icon.User, "Badge": icon.InfoCircle, "Card": icon.Appstore,
	"Carousel": icon.Picture, "Collapse": icon.Down, "Collapsible": icon.Down,
	"Accordion": icon.Down, "Descriptions": icon.File, "Empty": icon.InfoCircle,
	"Image": icon.Picture, "RemoteImage": icon.Picture, "List": icon.Menu,
	"Popover": icon.InfoCircle, "QRCode": icon.Appstore, "Statistic": icon.File,
	"Table": icon.Filter, "Tag": icon.More, "Timeline": icon.Menu,
	"Tooltip": icon.InfoCircle, "Tour": icon.QuestionCircle, "Meter": icon.Filter,
	"Markdown": icon.File, "CodeEditor": icon.File, "ChartPlot": icon.Picture,
	"RichTextEditor": icon.Edit, "GridView": icon.Appstore,
	"Alert": icon.InfoCircle, "Drawer": icon.Menu, "Message": icon.Mail,
	"Modal": icon.Appstore, "Notification": icon.Notification,
	"Popconfirm": icon.QuestionCircle, "ContextMenu": icon.More,
	"Confirm": icon.QuestionCircle, "Progress": icon.Loading, "Spin": icon.Loading,
	"Result": icon.Check, "Skeleton": icon.Picture, "Watermark": icon.Picture,
	"EditableText": icon.Edit,
}

// homePage 组件总览（对齐 ant.design/components）：标题 + 说明 +
// 按类别分区的组件卡片栅格（图标 + 中文名 + 英文名，点击进入控件页）。
// 本项目有而 antd 没有的控件卡片带「MyGo」角标（对齐 antd 给新组件
// 标版本号的形制）。
func (g *Gallery) homePage(c *ui.Context) {
	t := c.Theme()

	ui.Scroll(c).Fill().Children(func() {
		col := ui.Column(c).FillWidth().Padding(t.Space(8), t.Space(10)).Gap(t.Space(6))

		// 页头：标题 + 说明。
		col.Children(func() {
			head := ui.Column(c).FillWidth().Gap(t.Space(2))
			head.Children(func() {
				ui.Text(c, "组件总览").FontSize(t.Rem(1.71)).FontWeight(600).TextColor(t.Text)
				ui.Text(c, "mygo-antdesign 为 MyGo 原生 UI 提供 Ant Design 风格控件：换肤 MyGo 原生控件、补齐 antd 有而 MyGo 无的控件；全部控件亮暗双主题、标准/紧凑双密度，配真实代码示例。").
					TextColor(t.TextMuted).FontSize(t.Rem(0.9))
			})
		})

		// 按类别分区：类别名 + 数量 + 卡片栅格（antd 组件总览形制）。
		pages := g.filtered()
		for _, cat := range samples.CategoryOrder {
			var group []samples.ControlPage
			for _, p := range pages {
				if p.Category == cat {
					group = append(group, p)
				}
			}
			if len(group) == 0 {
				continue
			}
			cat := cat
			col.Children(func() {
				secHead := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(2))
				secHead.Children(func() {
					ui.Text(c, string(cat)).FontSize(t.Rem(1.29)).FontWeight(600).TextColor(t.Text)
					ui.Textf(c, "%d", len(group)).TextColor(t.TextMuted).FontSize(t.Rem(0.78))
				})

				row := ui.Row(c).FillWidth().Gap(t.Space(3)).Wrap()
				row.Children(func() {
					for _, p := range group {
						g.overviewCard(c, p, cat)
					}
				})
			})
		}
	})
}

// overviewCard 组件总览的一张卡片：图标 + 中文名 + 英文名（+ MyGo 角标）。
// hover 浮起（主色描边 + 阴影），点击进入控件页。
func (g *Gallery) overviewCard(c *ui.Context, p samples.ControlPage, cat samples.Category) {
	t := c.Theme()
	card := ui.ButtonBase(c).Width(t.Space(58)).
		Height(t.Space(30)).
		Padding(t.Space(4)).
		Radius(t.Radius+t.Space(0.5)).
		Background(t.Surface).
		Border(1, t.Border).
		AlignItems(ui.Center).
		Gap(t.Space(3))
	if p.Name == g.selected {
		card.Border(1, t.Accent)
		card.Shadow(0, 2, 8, 0, t.Accent.Alpha(0.15))
	}
	if card.Hovered() {
		card.Border(1, t.Accent.Alpha(0.6))
		card.Shadow(0, 2, 8, 0, t.Text.Alpha(0.10))
	}
	card.Children(func() {
		ic := pageIcon[p.Name]
		if ic == "" {
			ic = catIcon[cat]
		}
		icon.Icon(c, ic, t.Space(6), icon.WithColor(t.Accent))
		name := ui.Column(c).Gap(0).Grow(1)
		name.Children(func() {
			title := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(1))
			title.Children(func() {
				ui.Text(c, samples.PageZh(p.Name)).
					FontSize(t.Rem(0.9)).FontWeight(600).TextColor(t.Text).SingleLine()
				if samples.IsMyGoOnly(p.Name) {
					widgets.Tag(c, "MyGo", widgets.TagPrimary())
				}
			})
			ui.Text(c, p.Name).TextColor(t.TextMuted).FontSize(t.Rem(0.72)).SingleLine()
		})
	})
	if card.Clicked() {
		g.selected = p.Name
		g.codeTab = 0
	}
}
