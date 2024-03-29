# 运行 mdBook

kubebuilder 书籍使用 [mdBook](https://github.com/rust-lang-nursery/mdBook) 进行服务。如果您想在本地测试对书籍的更改，请按照以下步骤进行：

1. 按照[https://github.com/rust-lang-nursery/mdBook#installation](https://github.com/rust-lang-nursery/mdBook#installation)上的说明安装 mdBook。
2. 确保在 `$GOPATH` 中安装了 [controller-gen](https://pkg.go.dev/sigs.k8s.io/controller-tools/cmd/controller-gen)。
3. 进入 `docs/book` 目录。
4. 运行 `mdbook serve`。
5. 访问 [http://localhost:3000](http://localhost:3000)。

# 部署步骤

不需要手动操作来部署网站。

Kubebuilder 书籍网站部署在 Netlify 上。每个 PR 都有一个网站预览。一旦 PR 合并，网站将在 Netlify 上构建并部署。

如果您需要进一步的帮助，请告诉我。