---
name: visual-audit
description: 本仓库（mygo-antdesign）的控件视觉截图与 antd 规范审查工具链。用 scripts/statecapture 生成每个控件每种交互状态（默认/悬停/按下/聚焦/禁用/加载中）的 PNG 到 images/，逐状态比对 antd v5 规范并修复违规代码；用 gallery 的 TestGalleryAuditScreenshots 生成 Gallery 首页与每个控件页截图并审查布局瑕疵。Use when 用户要求 控件状态截图、视觉审查、截图比对 antd 规范、检查控件样式违规、gallery 页面截图审查、生成控件图片。
---

# visual-audit：控件状态截图 + antd 规范审查

三条工作流，可独立执行，也可按 A→B→C 完整走一遍。**全程遵守 AGENTS.md 第 2 节样式铁律**：颜色只来自主题、间距只用 `t.Space()/t.Rem()`、状态映射集中在 `internal/variant.go`。

## 工作流 A：生成控件状态截图

```bash
go run ./scripts/statecapture                # 亮色，输出 images/
go run ./scripts/statecapture -dark          # 暗色
go run ./scripts/statecapture -out .audit    # 自定义输出目录
```

- 每个变体输出：`<控件>-<变体>.png`（默认态）、`-hover`、`-pressed`、`-focus`（模拟交互）；禁用/加载中只有默认态一张。
- 命名示例：`images/button-primary.png`、`images/button-primary-hover.png`。
- **新增控件后必须在此登记**：编辑 `scripts/statecapture/main.go` 的 `specs()`，仿照 `buttonSpecs()` 追加一组 spec（变体构造 + 交互探针文本），见 [REFERENCE.md](REFERENCE.md) 第 1 节。

## 工作流 B：逐状态比对 antd 规范并修复

1. 用图片查看工具逐张审读 `images/` 下的状态截图（默认 vs hover vs pressed 逐张对比，可用像素对比确认色阶方向）。
2. 对照 [REFERENCE.md](REFERENCE.md) 第 3 节的 antd v5 规范检查表逐项核对：颜色、色阶方向（hover 提亮 / pressed 加深）、描边、圆角、尺寸、状态行为。
3. **发现违规立即修代码**，修复路径：
   - 状态颜色映射不对 → 改 `internal/variant.go`（不要在控件里写分支）；
   - 基础 token 值不对 → 改 `theme/` 的 AntToken 映射；
   - 布局/度量不对 → 改对应控件文件。
4. 修复后重新跑工作流 A 复核，直至全部符合规范；再跑基线回归并更新基线（见第 5 节）。

## 工作流 C：Gallery 页面截图与布局审查

```bash
GALLERY_CAPTURE=1 go test ./gallery -run TestGalleryAuditScreenshots
```

- 输出 `images/gallery/gallery-home-{light,dark}.png` 与每个控件页 `gallery-<控件名>-{light,dark}.png`（1280×800）。
- 逐张审读，检查项见 [REFERENCE.md](REFERENCE.md) 第 4 节：对齐、截断、留白、主题残留、密度错乱等。发现问题 → 修代码 → 重新截图，直至无瑕疵。

## 5. 收尾（强制，同 AGENTS.md 第 5 节）

```bash
gofmt -l . && go vet ./... && go test ./... && go run ./internal/check
go test ./gallery -run TestGalleryBaselines   # 改了控件/主题/示例时必跑
```

- 截图基线不一致且属预期改动 → 删除 `samples/golden/` 对应文件重新生成并提交。
- 一个任务一个提交：`feat(widgets): ...` / `fix(theme): ...` / `chore(scripts): ...`，提交后立即 `git push origin main`。
- `images/` 已加入 .gitignore，是审查工作产物，**不要提交**。

## 注意事项

- `scripts/statecapture` 与 `gallery/audit_test.go` 是本 skill 的执行器，属仓库正式代码，随首次提交入库并持续维护。
- 交互态模拟基于 `ui.Tester` 的 `Find/Move/Press/ClickAt`；focus 是「点击后」的聚焦态（非键盘 focus-visible）。
- 含动画的变体（如 loading 中的 Spinner）逐帧像素不稳定，只截默认态即可。
- 更多细节、antd 检查表全文、坑位说明见 [REFERENCE.md](REFERENCE.md)。
