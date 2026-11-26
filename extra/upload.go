package extra

import (
	"path/filepath"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/transfer"
	"github.com/egoist/mygo/ui"
)

// UploadFile 已加入上传列表的一个文件。
type UploadFile struct {
	Name string // 文件名（不含目录）
	Path string // 完整路径
}

// UploadOption 定制 Upload。
type UploadOption func(*uploadConfig)

type uploadConfig struct {
	accept   string // 提示的文件类型说明
	multiple bool
	onChange func(files []UploadFile)
	hint     string
}

// UploadMultiple 允许一次选择多个文件。
func UploadMultiple(b bool) UploadOption { return func(c *uploadConfig) { c.multiple = b } }

// UploadAccept 提示的文件类型说明（如 "图片或 PDF"）。
func UploadAccept(desc string) UploadOption { return func(c *uploadConfig) { c.accept = desc } }

// UploadHint 拖放区的提示文字（默认「点击选择或拖拽文件到此处」）。
func UploadHint(s string) UploadOption { return func(c *uploadConfig) { c.hint = s } }

// UploadOnChange 文件列表变化回调。
func UploadOnChange(fn func(files []UploadFile)) UploadOption {
	return func(c *uploadConfig) { c.onChange = fn }
}

// Upload 上传控件（antd Upload）。点击拖放区打开系统文件对话框；
// 接受从系统拖入的文件；文件列表由库内跨帧持有，每项可移除：
//
//	ant.Upload(c, ant.UploadMultiple(true), ant.UploadOnChange(app.upload))
func Upload(c *ui.Context, opts ...UploadOption) ui.Element {
	cfg := &uploadConfig{hint: "点击选择或拖拽文件到此处"}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	host := ui.Column(c).FillWidth().Gap(t.Space(2))
	files := ui.Local(host, "ant.extra.upload.files", func() []UploadFile { return nil })

	host.Children(func() {
		// 拖放区：ButtonBase 提供点击，同时是原生拖放目的地。
		zone := ui.ButtonBase(c).FillWidth().
			MinHeight(t.Space(12)).
			Radius(t.Radius).
			Justify(ui.Center).AlignItems(ui.Center)
		zone.Children(func() {
			body := ui.Column(c).Gap(t.Space(1)).AlignItems(ui.Center)
			body.Children(func() {
				ui.Text(c, "⬆").FontSize(t.Rem(1.5)).TextColor(t.Accent)
				ui.Text(c, cfg.hint).TextColor(t.Text)
				if cfg.accept != "" {
					ui.Text(c, "支持 "+cfg.accept).TextColor(t.TextMuted).FontSize(t.Rem(0.86))
				}
			})
		})
		// 悬停高亮：文件拖入时边框与底色切换为品牌色。
		if zone.FileDragOver() {
			zone.Background(t.Accent.Mix(t.Background, 0.92)).Border(1, t.Accent)
		} else {
			zone.Background(t.Surface).Border(1, t.Border)
			zone.BorderStyle(ui.BorderDashed)
		}
		// 接受系统文件拖入。
		if drop, ok := ui.DropData(zone, transfer.DropOptions{
			Formats: []transfer.Format{transfer.FileList, transfer.URIList},
		}); ok {
			if paths, err := drop.Data.Files(); err == nil && len(paths) > 0 {
				addPaths(files, paths)
				fireUpload(cfg, files)
			}
		}
		// 点击打开系统文件对话框（阻塞至对话框关闭）。
		if zone.Clicked() {
			openDialog(cfg, files)
		}

		// 文件列表。
		for i, f := range *files {
			i := i
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
				Padding(t.Space(1), 0)
			row.Children(func() {
				ui.Text(c, "📄").TextColor(t.TextMuted)
				ui.Text(c, f.Name).TextColor(t.Text)
				x := ui.ButtonBase(c).Size(t.Space(4), t.Space(4)).Center()
				x.Children(func() {
					ui.Text(c, "✕").TextColor(t.TextMuted).FontSize(t.Rem(0.7))
				})
				if x.Clicked() {
					removeFile(files, i)
					fireUpload(cfg, files)
				}
			})
		}
	})
	return host
}

// openDialog 打开系统文件选择对话框并合并结果。
func openDialog(cfg *uploadConfig, files *[]UploadFile) {
	paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{
		Title:    "选择文件",
		Multiple: cfg.multiple,
	})
	if err != nil || len(paths) == 0 {
		return
	}
	addPaths(files, paths)
	fireUpload(cfg, files)
}

// addPaths 把路径加入列表（按路径去重）。
func addPaths(files *[]UploadFile, paths []string) {
	for _, p := range paths {
		f := UploadFile{Name: filepath.Base(p), Path: p}
		dup := false
		for _, existing := range *files {
			if existing.Path == f.Path {
				dup = true
				break
			}
		}
		if !dup {
			*files = append(*files, f)
		}
	}
}

// removeFile 移除第 i 项。
func removeFile(files *[]UploadFile, i int) {
	if i < 0 || i >= len(*files) {
		return
	}
	*files = append((*files)[:i], (*files)[i+1:]...)
}

// fireUpload 变化回调。
func fireUpload(cfg *uploadConfig, files *[]UploadFile) {
	if cfg.onChange != nil {
		cfg.onChange(*files)
	}
}
