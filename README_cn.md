> ⚠️ **重要通知：** `gcr.io/kubebuilder/` 下的镜像将很快不可用。
>
> **如果您的项目使用 `gcr.io/kubebuilder/kube-rbac-proxy`，** 将会受到影响。
>
> 如果无法拉取该镜像，您的项目可能无法正常工作。**请尽快迁移**，因为从2025年初开始，GCR将会消失。
>
> 项目 [kube-rbac-proxy](https://github.com/brancz/kube-rbac-proxy) 在Kubebuilder中已被停止使用，并通过Controller-Runtime的特性 [WithAuthenticationAndAuthorization](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/metrics/filters#WithAuthenticationAndAuthorization) 替代以提供类似的保护。
>
> 有关更多信息和指导，请查看讨论：[讨论链接](https://github.com/kubernetes-sigs/kubebuilder/discussions/3907)

[![Lint](https://github.com/kubernetes-sigs/kubebuilder/actions/workflows/lint.yml/badge.svg)](https://github.com/kubernetes-sigs/kubebuilder/actions/workflows/lint.yml)

[![单元测试](https://github.com/kubernetes-sigs/kubebuilder/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/kubernetes-sigs/kubebuilder/actions/workflows/unit-tests.yml)

[![Go报告卡](https://goreportcard.com/badge/sigs.k8s.io/kubebuilder)](https://goreportcard.com/report/sigs.k8s.io/kubebuilder)

[![覆盖状态](https://coveralls.io/repos/github/kubernetes-sigs/kubebuilder/badge.svg?branch=master)](https://coveralls.io/github/kubernetes-sigs/kubebuilder?branch=master)

[![最新发布](https://badgen.net/github/release/kubernetes-sigs/kubebuilder)](https://github.com/kubernetes-sigs/kubebuilder/releases)

## Kubebuilder

Kubebuilder 是一个用于构建 Kubernetes API 的框架，使用 [自定义资源定义（CRDs）](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions)。

类似于 *Ruby on Rails* 和 *SpringBoot* 等Web开发框架，Kubebuilder 提高了开发人员构建和发布 Kubernetes API 的速度，同时减少了复杂性。它建立在构建核心 Kubernetes API 的标准技术之上，提供简单的抽象，减少样板代码和繁琐工作。

Kubebuilder **不是**一个可以直接*复制粘贴*的示例，而是提供强大的库和工具，以简化从头开始构建和发布 Kubernetes API 的过程。它提供了一个插件架构，允许用户利用可选的助手和功能。有关更多信息，请参见 [插件部分][plugin-section]。

Kubebuilder 基于 [controller-runtime][controller-runtime] 和 [controller-tools][controller-tools] 库进行开发。

### Kubebuilder 也是一个库

Kubebuilder 是可扩展的，可以作为其他项目中的库使用。[Operator-SDK][operator-sdk] 是一个使用 Kubebuilder 作为库的良好示例。[Operator-SDK][operator-sdk] 使用插件功能来包含非 Go 操作，例如 operator-sdk 的 Ansible 和 Helm 基于语言的操作。

要了解更多信息，请参见 [如何创建自己的插件][your-own-plugins]。

### 安装

强烈建议您使用已发布的版本。发布的二进制文件可以在 [发布](https://github.com/kubernetes-sigs/kubebuilder/releases) 页面上找到。请按照 [说明](https://book.kubebuilder.io/quick-start.html#installation) 安装 Kubebuilder。

## 入门

请查看 [入门](https://book.kubebuilder.io/quick-start.html) 文档。

![快速入门](docs/gif/kb-demo.v3.11.1.svg)

另外，请确保查看 [部署镜像](./docs/book/src/plugins/available/deploy-image-plugin-v1-alpha.md) 插件。该插件允许用户搭建 API/控制器，以便在集群中部署和管理操作对象（镜像），遵循指南和最佳实践。它抽象了实现此目标的复杂性，同时允许用户自定义生成的代码。

## 文档

请查看 Kubebuilder [书籍](https://book.kubebuilder.io)。

## 资源

- Kubebuilder 书籍: [book.kubebuilder.io](https://book.kubebuilder.io)
- GitHub 仓库: [kubernetes-sigs/kubebuilder](https://github.com/kubernetes-sigs/kubebuilder)
- Slack 频道: [#kubebuilder](https://kubernetes.slack.com/messages/#kubebuilder)
- Google 群组: [kubebuilder@googlegroups.com](https://groups.google.com/forum/#!forum/kubebuilder)
- 设计文档: [designs](designs/)
- 插件: [plugins][plugin-section]

## 动机

构建 Kubernetes 工具和 API 涉及大量决策和样板代码的编写。

为了方便使用标准方法轻松构建 Kubernetes API 和工具，该框架提供了一系列 Kubernetes 开发工具，以最小化繁琐工作。

Kubebuilder 努力促进以下开发人员工作流程来构建 API：

1. 创建一个新的项目目录
2. 创建一个或多个资源 API 作为 CRD，并为资源添加字段
3. 在控制器中实现协调循环并监视其他资源
4. 通过在集群中运行进行测试（自动安装 CRD 并启动控制器）
5. 更新引导集成测试以测试新字段和业务逻辑
6. 从提供的 Dockerfile 构建并发布容器

## 范围

使用 CRD、控制器和 Admission Webhook 构建 API。

## 哲学

有关各种 Kubebuilder 项目的指导原则，请参见 [DESIGN.md](DESIGN.md)。

简而言之：

提供干净的库抽象，并附有清晰且有示例的 Go 文档。

- 优先使用 Go *接口* 和 *库*，而不是过度依赖 *代码生成*
- 优先使用 *代码生成*，而不是 *一次性初始化* 的存根
- 优先 *一次性初始化* 的存根，而不是分叉和修改的样板代码
- 永远不要分叉和修改样板代码

## 技术

- 在低级客户端库之上提供更高级的库
  - 保护开发人员免受低级库的破坏性更改
  - 从最小开始，逐步发现功能
  - 提供合理的默认值，并允许用户在存在时进行覆盖
- 提供代码生成器以维护无法通过接口解决的常见样板代码
  - 通过 `// +` 注释驱动
- 提供引导命令以初始化新包

## 版本管理和发布

请参见 [VERSIONING.md](VERSIONING.md)。

## 故障排除

- ### 错误和功能请求：
  如果您发现看似错误的内容，或者希望提出功能请求，请使用 [Github 问题跟踪系统](https://github.com/kubernetes-sigs/kubebuilder/issues)。在提交问题之前，请搜索现有问题，查看您的问题是否已经被覆盖。

- ### Slack
  对于实时讨论，您可以加入 [#kubebuilder](https://slack.k8s.io/#kubebuilder) Slack 频道。Slack 需要注册，但 Kubernetes 团队欢迎任何人注册。欢迎随时来询问任何问题。

## 贡献

非常欢迎贡献。维护者积极管理问题列表，并努力突出适合新手的问题。该项目遵循典型的 GitHub 拉取请求模型。有关更多详细信息，请参见 [CONTRIBUTING.md](CONTRIBUTING.md)。在开始任何工作之前，请在现有问题上评论或提交新问题。

## 支持的操作系统

目前，Kubebuilder 官方支持 macOS 和 Linux 平台。如果您使用 Windows 操作系统，可能会遇到问题。欢迎为支持 Windows 贡献代码。

## 版本兼容性和支持

由 Kubebuilder 创建的项目包含一个 `Makefile`，在项目创建时安装定义版本的工具。主要包含的工具有：

- [kustomize](https://github.com/kubernetes-sigs/kustomize)
- [controller-gen](https://github.com/kubernetes-sigs/controller-tools)
- [setup-envtest](https://github.com/kubernetes-sigs/controller-runtime/tree/main/tools/setup-envtest)

此外，这些项目还包括一个 `go.mod` 文件，指定依赖版本。Kubebuilder 依赖于 [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime) 及其 Go 和 Kubernetes 依赖项。因此，`Makefile` 和 `go.mod` 文件中定义的版本是经过测试、支持和推荐的版本。

每个小版本的 Kubebuilder 都与特定的小版本 client-go 进行测试。虽然 Kubebuilder 的小版本 *可能* 与其他 client-go 小版本或其他工具兼容，但这种兼容性并不保证、支持或经过测试。

Kubebuilder 所需的最低 Go 版本由其依赖项所需的最高最低 Go 版本决定。这通常与相应的 `k8s.io/*` 依赖项所需的最低 Go 版本保持一致。

兼容的 `k8s.io/*` 版本、client-go 版本和最低 Go 版本可以在为每个 [标签发布](https://github.com/kubernetes-sigs/kubebuilder/tags) 的每个项目中找到的 `go.mod` 文件中查看。

**示例：** 对于 `4.1.1` 发布，最低 Go 版本兼容性为 `1.22`。您可以参考标签发布 [v4.1.1](https://github.com/kubernetes-sigs/kubebuilder/tree/v4.1.1/testdata) 中的 testdata 目录中的示例，例如 `project-v4` 的 [go.mod](https://github.com/kubernetes-sigs/kubebuilder/blob/v4.1.1/testdata/project-v4/go.mod#L3) 文件。您还可以通过检查 [Makefile](https://github.com/kubernetes-sigs/kubebuilder/blob/v4.1.1/testdata/project-v4/Makefile#L160-L165) 来查看此发布所支持和测试的工具版本。

## 社区会议

以下会议每两周举行一次：

- Kubebuilder 会议

欢迎您参加。如需更多信息，请加入 [kubebuilder@googlegroups.com](https://groups.google.com/g/kubebuilder)。我们的团队每月在第一个星期四的太平洋时间 11:00 举行会议，讨论我们的进展和未来几周的计划。请注意，我们最近通过 Slack 更频繁地进行离线同步。不过，如果您添加议题，我们将按计划举行会议。此外，我们可以利用这个频道展示新功能。

[operator-sdk]: https://github.com/operator-framework/operator-sdk  
[plugin-section]: https://book.kubebuilder.io/plugins/plugins.html  
[controller-runtime]: https://github.com/kubernetes-sigs/controller-runtime  
[your-own-plugins]: https://book.kubebuilder.io/plugins/extending  
[controller-tools]: https://github.com/kubernetes-sigs/controller-tools  