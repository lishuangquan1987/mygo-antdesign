# visual-audit REFERENCE：工具详解、antd v5 检查表、坑位

## 1. scripts/statecapture 详解与扩展方法

### 运行

```bash
go run ./scripts/statecapture [-dark] [-out images]
```

每个 spec 渲染一块 480×240 画布：`theme.Apply` 后用 `ui.Row.Padding(48).AlignItems(ui.Center)` 居中放置控件。

- 默认态：`Frame()` 后直接保存 `<file>.png`。
- hover：`tt.Find(label)` 定位探针文本 → 取包围盒中心 → `tt.Move(x, y)` → `Frame()` → `<file>-hover.png`。
- pressed：`tt.Press(x, y)`（按下不释放）→ `Frame()` → `<file>-pressed.png` → `tt.CancelPress()`。
- focus：`tt.ClickAt(x, y)` 点击后聚焦 → `Frame()` → `<file>-focus.png`。
- `label` 为空串的 spec（禁用/加载中）只出默认态。

### 新增控件 spec（模板）

编辑 `scripts/statecapture/main.go`，追加一个 `xxxSpecs()` 并在 `specs()` 中 append：

```go
func checkboxSpecs() []spec {
	label := "同意协议"
	return []spec{
		{file: "checkbox-default", label: label,
			view: func(c *ui.Context) { widgets.Checkbox(c, label) }},
		{file: "checkbox-checked", label: label,
			view: func(c *ui.Context) { v := true; widgets.Checkbox(c, label, widgets.Checked(&v)) }},
		{file: "checkbox-disabled", label: "",
			view: func(c *ui.Context) { widgets.Checkbox(c, label, widgets.CheckboxDisabled(true)) }},
	}
}
```

要点：
- `view` 创建即挂载（与 samples 写法一致），带状态的控件用局部变量/`ui.Local`。
- 探针 `label` 必须是控件内**静态渲染的文本**；placeholder 等非文本节点 `Find` 不到（见第 5 节坑位）。
- 变体清单应覆盖 antd 该控件的全部 type + danger + 主要状态（禁用/加载中/选中态）。

### 像素级对比（确认色阶方向）

```bash
go run /tmp/cmp.go a.png b.png   # 或临时写脚本：输出差异包围盒与中心像素颜色
```

示例结论（已验证）：`button-primary` 默认中心 #1677FF，hover #1369E0（加深）、pressed 更深 —— 是否违规见第 3 节 B-1。

## 2. gallery/audit_test.go 详解

`TestGalleryAuditScreenshots` 由 `GALLERY_CAPTURE=1` 门控（不设则 Skip，不影响常规 `go test ./...`）：

```bash
GALLERY_CAPTURE=1 go test ./gallery -run TestGalleryAuditScreenshots
GALLERY_CAPTURE=1 GALLERY_CAPTURE_DIR=/tmp/shots go test ./gallery -run TestGalleryAuditScreenshots
```

- `registerAllDemos()` 的注册清单**必须与 gallery/main.go 保持同步**，新增控件示例分组时两边都要加。
- 输出：首页 `gallery-home-{light,dark}.png` + 每控件页 `gallery-<Name>-{light,dark}.png`，1280×800。
- 页面内容超出视口时只截到视口；需要看底部时临时调大 `height` 再跑。

## 3. antd v5 规范检查表（逐状态核对用）

### A. 全局 token

- [ ] 主色 `colorPrimary` = `#1677FF`（品牌色切换除外）
- [ ] 成功 `#52C41A` / 警告 `#FAAD14` / 错误 `#FF4D4F` / 信息 `#1677FF`
- [ ] 文字主色 rgba(0,0,0,0.88)、次要 rgba(0,0,0,0.65)、禁用 rgba(0,0,0,0.25)；暗色对应白系
- [ ] 圆角 `borderRadius` = 6px（按钮、输入框等）
- [ ] 字号 14px；控件高度 large 40 / middle 32 / small 24（紧凑密度按 docs/05 缩放）

### B. 交互状态规则（重点，逐张对比）

1. **实底按钮（primary/danger 实底）**：antd v5 亮色下 **hover 提亮**（#4096ff）、**pressed 加深**（#0958d9）；暗色主题 hover 同样更亮。若截图显示亮色 hover 加深，属违规 → 修 `internal/variant.go` 的 `ButtonColors`/`AccentHover` 映射或 theme token。
2. **描边按钮（default/dashed）**：hover 时**文字与描边变主色、底色不变**；pressed 同理加深。danger 描边 hover 出浅红底（colorErrorBg）。
3. **text 按钮**：默认透明底主文字（danger 红字）；hover 出浅填充（SurfaceHover）；pressed 更深填充。
4. **link 按钮**：透明底主色文字；hover 加深/提亮主色阶，且应有下划线语义（按现实现为准，记录差异）。
5. **禁用**：antd 灰字 rgba(0,0,0,0.25) + 浅灰底 #f5f5f5 + 描边 #d9d9d9；禁用不响应 hover/pressed（截图不得出现交互色）。
6. **加载中**：Spinner + 主色，文字保留；loading 期间不可点击。
7. **focus 环**：antd 用 2px 主色外扩光环（box-shadow）；核对 focus 截图是否有明显聚焦指示。

### C. 布局度量

- [ ] 按钮内边距、图标-文字间距符合 docs/05 数值表
- [ ] hover/pressed 的色阶过渡自然，无突变脏色（裸 #hex 或残留旧色）

## 4. Gallery 截图审查清单

逐张检查（亮/暗各一遍）：

- [ ] 侧栏：控件列表完整、无截断、选中态高亮清晰、搜索框正常
- [ ] 示例卡片：标题/描述/预览区/代码块排版整齐，无文字溢出或半行截断
- [ ] 控制条：亮暗/紧凑/尺寸/品牌色切换控件无错位
- [ ] 暗色页无亮色残留（背景、卡片、边框全部跟随主题）
- [ ] 紧凑密度下无元素重叠（可改 `g.density` 后补截图核对）
- [ ] 各控件页示例预览中控件本身符合第 3 节规范
- [ ] 无大面积空白、异常留白或元素重叠

发现问题 → 修代码（布局改控件文件、排版改 gallery/）→ 重跑截图 → 直至合格。

## 5. 已知坑位

1. **placeholder 不可作为探针**：`ui.Tester.Find` 只匹配文本节点，Input 的 placeholder `Find` 不到 → 用 `InputLabel` 的 label 文本定位，点击坐标取 label 右侧输入区或直接接受跳过（工具会打印 `[跳过交互态]`）。
2. **focus 是点击后聚焦态**：`ClickAt` 会同时触发 hover + active + focus，截图是复合态；键盘 focus-visible（Tab 聚焦）目前无法模拟，记录差异即可。
3. **动画帧不稳定**：Spinner/加载动画逐帧像素不同，不做逐字节对比，也勿纳入 golden 基线（参照 samples/golden_test.go 的 `animated` 排除表）。
4. **字体跨平台差异**：截图逐字节比对只在同平台有效；goldens 带 manifest.json 平台标记（见 docs/07 第 8 节）。人工审读不受影响。
5. **不要提交 images/**：已 gitignore；需要留档的结论写进提交说明或 docs。
