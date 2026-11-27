## 变更内容

<!-- 简要说明这个 PR 做了什么、为什么需要 -->

## 关联 Issue

<!-- 例如：Closes #12 -->

## 变更类型

- [ ] 新控件 / 新能力（feat）
- [ ] 缺陷修复（fix）
- [ ] 主题 / 密度 / 视觉调整
- [ ] 示例 / Gallery
- [ ] 文档 / 构建 / CI
- [ ] 重构（refactor）

## 自测清单

- [ ] `gofmt -l .` 无输出
- [ ] `go vet ./...` 通过
- [ ] `go test -count=1 ./...` 通过
- [ ] `go run ./internal/check` 通过（widgets/extra 无裸颜色、裸间距）
- [ ] 新增/修改控件已补示例（`samples/<分组>/demo.go`）并在 `gallery/main.go` 注册
- [ ] 截图基线已更新：`go test ./gallery ./samples -count=1 -run 'TestGalleryBaselines|TestControlScreenshots'`

## 截图 / 录屏

<!-- 视觉相关变更请附前后对比；无界面变更可删除 -->

## 备注

<!-- 兼容性、破坏性变更、后续计划等 -->
