package extra

import (
	"path/filepath"
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/transfer"
	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/icon"
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
	exts     []string // 允许的扩展名（含点，如 ".png"）；空表示不限
	maxCount int      // 最多文件数；0 表示不限
}

// UploadMultiple 允许一次选择多个文件。
func UploadMultiple(b bool) UploadOption { return func(c *uploadConfig) { c.multiple = b } }

// UploadAccept 提示的文件类型说明（如 "图片或 PDF"）。
func UploadAccept(desc string) UploadOption { return func(c *uploadConfig) { c.accept = desc } }

// UploadExtensions 限定可选文件类型（antd accept）：扩展名可带点（".png"）
// 也可不带（"png"），不区分大小写；不在列表中的文件会被忽略。
func UploadExtensions(exts ...string) UploadOption {
	return func(c *uploadConfig) { c.exts = append(c.exts, exts...) }
}

// UploadMaxCount 最多文件数（antd maxCount）：达到上限后新增文件被忽略。
func UploadMaxCount(n int) UploadOption { return func(c *uploadConfig) { c.maxCount = n } }

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
			MinHeight(t.Space(24)).
			Padding(t.Space(6)).
			Radius(t.Radius).
			Justify(ui.Center).AlignItems(ui.Center)
		zone.Children(func() {
			body := ui.Column(c).Gap(t.Space(2)).AlignItems(ui.Center)
			body.Children(func() {
				icon.Icon(c, icon.Upload, t.Space(10), icon.WithColor(t.Accent))
				ui.Text(c, cfg.hint).TextColor(t.Text).FontSize(t.Rem(1.14))
				if cfg.accept != "" {
					ui.Text(c, "支持 "+cfg.accept).TextColor(t.TextMuted).FontSize(t.FontSize)
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
				addPaths(cfg, files, paths)
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
				Padding(t.Space(1), t.Space(2)).
				Radius(t.Radius * 0.67)
			if row.Hovered() {
				row.Background(t.SurfaceHover)
			}
			row.Children(func() {
				icon.Icon(c, icon.File, t.Space(4), icon.WithColor(t.TextMuted))
				ui.Text(c, f.Name).TextColor(t.Text).FontSize(t.FontSize)
				x := ui.ButtonBase(c).Size(t.Space(5), t.Space(5)).Radius(t.Radius * 0.67).Center().Label("移除文件")
				if x.Hovered() {
					x.Background(t.Danger.Mix(t.Surface, 0.86))
				}
				x.Children(func() {
					icon.Icon(c, icon.Close, t.Space(2.5), icon.WithColor(t.TextMuted))
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
	addPaths(cfg, files, paths)
	fireUpload(cfg, files)
}

// addPaths 把路径加入列表（按路径去重，按扩展名与数量上限过滤）。
func addPaths(cfg *uploadConfig, files *[]UploadFile, paths []string) {
	for _, p := range paths {
		if cfg.maxCount > 0 && len(*files) >= cfg.maxCount {
			break
		}
		if !allowedExt(cfg, p) {
			continue
		}
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

// allowedExt 判断扩展名是否在允许列表内（不区分大小写；"png" 与 ".png" 等价）。
func allowedExt(cfg *uploadConfig, path string) bool {
	if len(cfg.exts) == 0 {
		return true
	}
	ext := strings.ToLower(filepath.Ext(path))
	for _, want := range cfg.exts {
		w := strings.ToLower(want)
		if !strings.HasPrefix(w, ".") {
			w = "." + w
		}
		if w == ext {
			return true
		}
	}
	return false
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
