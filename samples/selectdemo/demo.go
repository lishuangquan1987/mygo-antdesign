// Package select 提供选择类控件的示例。
package selectdemo

import (
	_ "embed"

	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

//go:embed demo.go
var demoSrc string

// Demo 返回选择类控件的示例集合。
func Demo() samples.WidgetDemo {
	return samples.WidgetDemo{
		Name:     "选择控件",
		Package:  "widgets",
		Category: samples.CatInput,
		Summary:  "Checkbox / Radio / Switch / Slider / Rate / Segmented / Select / Combobox：选中色随品牌色。",
		Samples: []samples.Sample{
			{
				Title:       "复选框",
				Description: "Checkbox 与三态全选框（全选 / 半选 / 未选）。",
				Code:        samples.Extract(demoSrc, "checkbox"),
				View:        CheckboxDemo,
			},
			{
				Title:       "单选组",
				Description: "RadioGroup 内一组 Radio，选中值受控。",
				Code:        samples.Extract(demoSrc, "radio"),
				View:        RadioDemo,
			},
			{
				Title:       "开关",
				Description: "Switch 无文字开关与带文字开关。",
				Code:        samples.Extract(demoSrc, "switch"),
				View:        SwitchDemo,
			},
			{
				Title:       "滑动条与评分",
				Description: "Slider 数值滑轨、Rate 星级评分。",
				Code:        samples.Extract(demoSrc, "slider"),
				View:        SliderDemo,
			},
			{
				Title:       "分段控制器",
				Description: "Segmented 切换互斥选项。",
				Code:        samples.Extract(demoSrc, "segmented"),
				View:        SegmentedDemo,
			},
			{
				Title:       "下拉选择",
				Description: "Select 下拉单选与 Combobox 输入过滤。",
				Code:        samples.Extract(demoSrc, "selectdemo"),
				View:        SelectDemo,
			},
			{
				Title:       "开关按钮",
				Description: "Toggle 保持按下的按钮，适合编辑器工具条。",
				Code:        samples.Extract(demoSrc, "toggle"),
				View:        ToggleDemo,
			},
		},
	}
}

// sample:checkbox
func CheckboxDemo(c *ui.Context) {
	col := ui.Column(c).Gap(6)
	agree := ui.Local(col, "ant.sample.select.agree", func() bool { return false })
	all := ui.Local(col, "ant.sample.select.all", func() ui.CheckState { return ui.CheckOff })
	col.Children(func() {
		widgets.Checkbox(c, agree, "同意用户协议")
		widgets.Checkbox(c, agree, "同意隐私政策", widgets.ControlDisabled(true))
		widgets.CheckboxAll(c, all, "全选（三态）")
	})
}

// /sample

// sample:radio
func RadioDemo(c *ui.Context) {
	col := ui.Column(c)
	size := ui.Local(col, "ant.sample.select.size", func() string { return "middle" })
	col.Children(func() {
		widgets.RadioGroup(c, func() {
			widgets.Radio(c, size, "small", "小")
			widgets.Radio(c, size, "middle", "中")
			widgets.Radio(c, size, "large", "大")
		})
	})
}

// /sample

// sample:switch
func SwitchDemo(c *ui.Context) {
	row := ui.Row(c).Gap(16).AlignItems(ui.Center)
	on := ui.Local(row, "ant.sample.switch.on", func() bool { return true })
	row.Children(func() {
		widgets.Switch(c, on)
		widgets.SwitchWithLabel(c, on, "自动保存")
	})
}

// /sample

// sample:slider
func SliderDemo(c *ui.Context) {
	col := ui.Column(c).FillWidth().Gap(12)
	vol := ui.Local(col, "ant.sample.slider.vol", func() float64 { return 30 })
	score := ui.Local(col, "ant.sample.slider.score", func() int { return 3 })
	col.Children(func() {
		widgets.Slider(c, vol, 0, 100).FillWidth()
		widgets.Rate(c, score, 5)
	})
}

// /sample

// sample:segmented
func SegmentedDemo(c *ui.Context) {
	row := ui.Row(c)
	tab := ui.Local(row, "ant.sample.segmented.tab", func() int { return 0 })
	row.Children(func() {
		widgets.Segmented(c, tab, []string{"全部", "进行中", "已完成"})
	})
}

// /sample

// sample:selectdemo
func SelectDemo(c *ui.Context) {
	col := ui.Column(c).FillWidth().Gap(8)
	city := ui.Local(col, "ant.sample.select.city", func() string { return "上海" })
	cmd := ui.Local(col, "ant.sample.select.cmd", func() string { return "" })
	col.Children(func() {
		widgets.Select(c, city, []string{"北京", "上海", "广州", "深圳"}, widgets.OnChange(func() {}))
		widgets.Combobox(c, cmd, []string{"重启服务", "重启网关", "重启集群"}, widgets.OnChange(func() {}))
	})
}

// /sample

// sample:toggle
func ToggleDemo(c *ui.Context) {
	host := ui.Column(c).FillWidth().Gap(c.Theme().Space(2))
	bold := ui.Local(host, "ant.sample.select.bold", func() bool { return true })
	italic := ui.Local(host, "ant.sample.select.italic", func() bool { return false })
	host.Children(func() {
		row := ui.Row(c).Gap(c.Theme().Space(2)).AlignItems(ui.Center)
		row.Children(func() {
			widgets.Toggle(c, bold, "B")
			widgets.Toggle(c, italic, "I")
			widgets.Paragraph(c, "Toggle 保持按下的按钮，适合编辑器工具条。", widgets.Secondary())
		})
	})
}

// /sample
