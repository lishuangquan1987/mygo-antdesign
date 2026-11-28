package main

// 本文件由 README「快速开始」代码块同步而来（gallery/quickstart_test.go
// 断言两者逐字一致，改 README 时同步这里，避免示例与文档漂移）。

const quickStartCode = `package main

import (
	"fmt"
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func main() {
	name := ""
	remember := false

	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title: "Hello mygo-antdesign",
			Content: ui.View(func(c *ui.Context) {
				// 一行接入主题：亮色预设 + 自定义品牌色。
				theme.Apply(c, theme.Light(theme.WithPrimary(ui.RGB(0x16, 0x77, 0xff))))

				// 间距一律走主题度量，切换密度时自动缩放。
				t := c.Theme()
				col := ui.Column(c).Gap(t.Space(3)).Padding(t.Space(6))
				col.Children(func() {
					widgets.Title(c, "登录", widgets.H2)
					widgets.Input(c, &name, widgets.PlaceholderText("用户名"))
					widgets.Checkbox(c, &remember, "记住我")
					widgets.Button(c, "登录", widgets.Primary(), widgets.OnClick(func() {
						fmt.Println("登录：", name, remember)
					}))
				})
			}),
		})
	})

	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}`
