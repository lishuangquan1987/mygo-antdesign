// Package feedback 提供反馈类控件的示例。
package feedback

import (
	_ "embed"

	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

//go:embed demo.go
var demoSrc string

// Demo 返回反馈控件的示例集合。
func Demo() samples.WidgetDemo {
	return samples.WidgetDemo{
		Name:     "反馈",
		Package:  "widgets",
		Category: samples.CatFeedback,
		Summary:  "Alert 提示条、Modal 模态框、Message 轻提示、Progress 进度条、Spin 加载中。",
		Samples: []samples.Sample{
			{
				Title:       "提示条",
				Description: "四种级别（信息/成功/警告/错误），淡底配色随主题派生。",
				Code:        samples.Extract(demoSrc, "alert"),
				View:        AlertDemo,
			},
			{
				Title:       "模态对话框",
				Description: "点击「打开」弹出确认框，确定 / 取消均有关闭反馈。",
				Code:        samples.Extract(demoSrc, "modal"),
				View:        ModalDemo,
			},
			{
				Title:       "轻提示与通知",
				Description: "Message 短消息与带操作按钮的 Notify 通知。",
				Code:        samples.Extract(demoSrc, "message"),
				View:        MessageDemo,
			},
			{
				Title:       "进度条",
				Description: "线形进度条，百分比文本，颜色随品牌色。",
				Code:        samples.Extract(demoSrc, "progress"),
				View:        ProgressDemo,
			},
			{
				Title:       "加载中",
				Description: "Spin 旋转指示器 + 提示文字，可选容器内居中。",
				Code:        samples.Extract(demoSrc, "spin"),
				View:        SpinDemo,
			},
			{
				Title:       "文字提示与气泡确认",
				Description: "Tooltip 悬停显示反色提示；Popconfirm 点击弹出确认气泡。",
				Code:        samples.Extract(demoSrc, "overlay"),
				View:        OverlayDemo,
			},
			{
				Title:       "抽屉",
				Description: "Drawer 从右侧滑入的面板，带遮罩、标题栏与滚动内容区。",
				Code:        samples.Extract(demoSrc, "drawer"),
				View:        DrawerDemo,
			},
			{
				Title:       "气泡与右键菜单",
				Description: "Popover 锚定气泡；ContextMenu 右键/菜单键弹出菜单。",
				Code:        samples.Extract(demoSrc, "popover"),
				View:        PopoverDemo,
			},
		},
	}
}

// sample:alert
func AlertDemo(c *ui.Context) {
	col := ui.Column(c).FillWidth().Gap(8)
	col.Children(func() {
		widgets.Alert(c, "信息提示", widgets.AlertTypeOpt(widgets.AlertInfo),
			widgets.AlertDescription("这是一条补充说明文字。"))
		widgets.Alert(c, "保存成功", widgets.AlertTypeOpt(widgets.AlertSuccess))
		widgets.Alert(c, "磁盘空间不足", widgets.AlertTypeOpt(widgets.AlertWarning))
		widgets.Alert(c, "网络连接失败", widgets.AlertTypeOpt(widgets.AlertError))
	})
}

// /sample

// sample:modal
func ModalDemo(c *ui.Context) {
	row := ui.Row(c).Gap(8).AlignItems(ui.Center)
	open := ui.Local(row, "ant.sample.modal.open", func() bool { return false })
	result := ui.Local(row, "ant.sample.modal.result", func() string { return "" })
	row.Children(func() {
		widgets.Button(c, "打开对话框", widgets.Primary(), widgets.OnClick(func() { *open = true }))
		if *result != "" {
			widgets.Tag(c, *result)
		}
		widgets.Modal(c, open, "确认操作",
			func() {
				widgets.Paragraph(c, "此操作不可撤销，是否继续？")
			},
			widgets.OnOk(func() { *result = "已确认" }),
			widgets.OnCancel(func() { *result = "已取消" }))
	})
}

// /sample

// sample:message
func MessageDemo(c *ui.Context) {
	row := ui.Row(c).Gap(8).Wrap()
	row.Children(func() {
		widgets.Button(c, "普通消息", widgets.OnClick(func() {
			widgets.Message(c, "已保存")
		}))
		widgets.Button(c, "成功消息", widgets.OnClick(func() {
			widgets.MessageTypeTyped(c, "success", "操作成功")
		}))
		widgets.Button(c, "错误消息", widgets.OnClick(func() {
			widgets.MessageTypeTyped(c, "error", "操作失败")
		}))
		widgets.Button(c, "通知", widgets.OnClick(func() {
			widgets.Notify(c, "部署完成", "v0.3.1 已上线", "查看日志", func() {})
		}))
	})
}

// /sample

// sample:progress
func ProgressDemo(c *ui.Context) {
	col := ui.Column(c).FillWidth().Gap(8)
	value := ui.Local(col, "ant.sample.progress.value", func() float64 { return 68 })
	col.Children(func() {
		widgets.ProgressLine(c, *value).FillWidth()
		widgets.Slider(c, value, 0, 100).FillWidth()
		widgets.ProgressLine(c, 100, widgets.ProgressHideText()).FillWidth()
		widgets.ProgressLine(c, 0).FillWidth()
		row := ui.Row(c).Gap(c.Theme().Space(4)).AlignItems(ui.Center)
		row.Children(func() {
			widgets.ProgressCircle(c, *value)
			widgets.ProgressCircle(c, 100, widgets.ProgressStatusOpt(widgets.ProgressSuccess))
			widgets.ProgressCircle(c, 30, widgets.ProgressStatusOpt(widgets.ProgressError))
		})
	})
}

// /sample

// sample:spin
func SpinDemo(c *ui.Context) {
	t := c.Theme()
	col := ui.Column(c).Gap(t.Space(4))
	col.Children(func() {
		widgets.Spin(c, widgets.SpinTip("加载中…"))
		box := ui.Box(c).FillWidth().Height(t.Space(20)).
			Radius(t.Radius).Background(t.SurfaceHover)
		box.Children(func() {
			widgets.Spin(c, widgets.SpinCenter(), widgets.SpinTip("正在保存"))
		})
	})
}

// /sample

// sample:overlay
func OverlayDemo(c *ui.Context) {
	t := c.Theme()
	row := ui.Row(c).Gap(t.Space(3)).AlignItems(ui.Center)
	removed := ui.Local(row, "ant.sample.feedback.removed", func() bool { return false })
	row.Children(func() {
		widgets.Tooltip(c, "删除后不可恢复", func(c *ui.Context) {
			widgets.Button(c, "悬停提示")
		})
		widgets.Popconfirm(c, func(c *ui.Context) ui.Element {
			return widgets.Button(c, "删除", widgets.Danger())
		}, "确定删除这一项？", func() { *removed = true },
			widgets.PopconfirmTitle("删除确认"))
		if *removed {
			widgets.Tag(c, "已删除", widgets.TagError())
		}
	})
}

// /sample

// sample:drawer
func DrawerDemo(c *ui.Context) {
	host := ui.Column(c).FillWidth()
	open := ui.Local(host, "ant.sample.feedback.drawer", func() bool { return false })
	host.Children(func() {
		widgets.Button(c, "打开抽屉", widgets.Primary(), widgets.OnClick(func() { *open = true }))
		widgets.Drawer(c, open, "设置", func() {
			widgets.Paragraph(c, "抽屉内容：可以放任意控件。")
			widgets.Paragraph(c, "点遮罩、Esc 或右上角 ✕ 关闭。", widgets.Secondary())
		}, widgets.DrawerWidth(80))
	})
}

// /sample

// sample:popover
func PopoverDemo(c *ui.Context) {
	host := ui.Column(c).FillWidth().Gap(c.Theme().Space(3))
	open := ui.Local(host, "ant.sample.feedback.popover", func() bool { return false })
	chosen := ui.Local(host, "ant.sample.feedback.menu", func() string { return "" })
	host.Children(func() {
		widgets.Popover(c, func(c *ui.Context) {
			if widgets.Button(c, "显示气泡").Clicked() {
				*open = true
			}
		}, open, func() {
			widgets.Paragraph(c, "这是 Popover 的气泡内容。")
		})
		widgets.ContextMenu(c, func(c *ui.Context) {
			box := ui.Box(c).Padding(c.Theme().Space(3)).Radius(c.Theme().Radius).
				Background(c.Theme().SurfaceHover)
			box.Children(func() { ui.Text(c, "右键这一行打开菜单") })
		}, []widgets.ContextMenuItem{
			{ID: "rename", Label: "重命名"},
			{ID: "delete", Label: "删除", Danger: true},
		}, func(id string) { *chosen = id })
		if *chosen != "" {
			widgets.Tag(c, "菜单选择："+*chosen, widgets.TagPrimary())
		}
	})
}

// /sample
