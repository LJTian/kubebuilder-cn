[![Lint](https://github.com/kubernetes-sigs/kubebuilder/actions/workflows/lint.yml/badge.svg)](https://github.com/kubernetes-sigs/kubebuilder/actions/workflows/lint.yml)
[![Unit tests](https://github.com/kubernetes-sigs/kubebuilder/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/kubernetes-sigs/kubebuilder/actions/workflows/unit-tests.yml)
[![Go Report Card](https://goreportcard.com/badge/sigs.k8s.io/kubebuilder)](https://goreportcard.com/report/sigs.k8s.io/kubebuilder)
[![Coverage Status](https://coveralls.io/repos/github/kubernetes-sigs/kubebuilder/badge.svg?branch=master)](https://coveralls.io/github/kubernetes-sigs/kubebuilder?branch=master)
[![Latest release](https://badgen.net/github/release/kubernetes-sigs/kubebuilder)](https://github.com/kubernetes-sigs/kubebuilder/lreleases)

## Kubebuilder

Kubebuilder 是一个使用[自定义资源定义 (CRDs)](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions)构建 Kubernetes API 的框架。

类似于 *Ruby on Rails* 和 *SpringBoot* 等 Web 开发框架，Kubebuilder 增加了开发人员的速度，减少了构建和发布 Kubernetes API 时的复杂性。它建立在构建核心 Kubernetes API 所使用的规范技术之上，提供了简单的抽象，减少了样板代码和重复劳动。

Kubebuilder 并不是一个可以*复制粘贴*的示例，而是提供了强大的库和工具，用于从头开始简化构建和发布 Kubernetes API。它提供了一个插件架构，允许用户利用可选的辅助工具和功能。要了解更多信息，请参阅[插件部分][plugin-section]。

Kubebuilder 是在 [controller-runtime][controller-runtime] 和 [controller-tools][controller-tools] 库之上开发的。

### Kubebuilder 也是一个库

Kubebuilder 是可扩展的，可以作为其他项目中的库使用。[Operator-SDK][operator-sdk] 是使用 Kubebuilder 作为库的一个很好的例子。[Operator-SDK][operator-sdk] 使用插件功能来包含非 Go 运算符，例如 operator-sdk 的 Ansible 和基于 Helm 的语言运算符。

要了解更多信息，请参阅[如何创建您自己的插件][your-own-plugins]。

### 安装

强烈建议您使用发布版本。发布的二进制文件可在[发布页面](https://github.com/kubernetes-sigs/kubebuilder/releases)上获得。请按照[说明](https://book.kubebuilder.io/quick-start.html#installation)安装 Kubebuilder。

## 入门指南

请查看[入门指南](https://book.kubebuilder.io/quick-start.html)文档。

![快速入门](docs/gif/kb-demo.v3.11.1.svg)

此外，请确保查看[部署镜像](https://book.kubebuilder.io/plugins/deploy-image-plugin-v1-alpha.html)插件。该插件允许用户生成 API/Controller，以便在集群上部署和管理操作数（镜像），遵循指南和最佳实践。它在抽象实现此目标的复杂性的同时，允许用户定制生成的代码。

## 文档

查看 Kubebuilder [书籍](https://book.kubebuilder.io)。

## 资源

- Kubebuilder 书籍（汉化版）：[kubebuilder.cn](https://kubebuilder.cn)
- Kubebuilder 书籍：[book.kubebuilder.io](https://book.kubebuilder.io)
- GitHub 仓库：[kubernetes-sigs/kubebuilder](https://github.com/kubernetes-sigs/kubebuilder)
- GitHub 汉化版仓库：[LJTian/kubebuilder-cn](https://github.com/LJTian/kubebuilder-cn)
- Slack 频道：[#kubebuilder](https://slack.k8s.io/#kubebuilder)
- Google Group：[kubebuilder@googlegroups.com](https://groups.google.com/forum/#!forum/kubebuilder)
- 设计文档：[designs](designs/)
- 插件：[plugins][plugin-section]

## 动机

___
kubebuilder.cn 汉化版的初衷：
1. 学习 kubebuilder 相关知识。
2. 没找到相对翻译较好且内容比较新的网站。
3. 稍微做点贡献。 

*此部分内容非官网内容同步*
___
构建 Kubernetes 工具和 API 涉及做出许多决策并编写大量样板代码。

为了便于使用规范方法轻松构建 Kubernetes API 和工具，该框架提供了一系列 Kubernetes 开发工具，以最小化重复劳动。

Kubebuilder 试图为构建 API 的开发者工作流程提供帮助：

1. 创建一个新的项目目录
2. 创建一个或多个资源 API 作为 CRDs，然后向资源添加字段
3. 在控制器中实现调和循环并监视其他资源
4. 通过针对集群运行测试（自动安装 CRDs 并自动启动控制器）
5. 更新引导集成测试以测试新字段和业务逻辑
6. 从提供的 Dockerfile 构建和发布容器

## 范围

使用 CRDs、控制器和准入 Webhook 构建 API。

## 哲学

请参阅 [DESIGN.md](DESIGN.md)，了解各种 Kubebuilder 项目的指导原则。

简而言之：

提供清晰和良好示例的库抽象。

- 更倾向于使用 go *接口* 和 *库*，而不是依赖于*代码生成*
- 更倾向于使用*代码生成*，而不是*一次性初始化*存根
- 更倾向于*一次性初始化*存根，而不是分叉和修改样板代码
- 永远不要分叉和修改样板代码

## 技术

- 在低级别客户端库之上提供更高级别的库
    - 保护开发人员免受低级别库的重大更改
    - 从最小开始，并提供功能的渐进式发现
    - 提供合理的默认值，并允许用户在存在时进行覆盖
- 提供代码生成器来维护无法通过接口解决的常见样板代码
    - 由 `//+` 注释驱动
- 提供引导命令以初始化新包

## 版本控制和发布

请参阅 [VERSIONING.md](VERSIONING.md)。

## 故障排除

- ### 错误和功能请求：
  如果您遇到看起来像错误的问题，或者您想提出功能请求，请使用[Github 问题跟踪系统](https://github.com/kubernetes-sigs/kubebuilder/issues)。
  在提交问题之前，请搜索现有问题，以查看您的问题是否已经有解决方案。

- ### Slack
  对于实时讨论，您可以加入 [#kubebuilder](https://slack.k8s.io/#kubebuilder) slack 频道。Slack 需要注册，但 Kubernetes 团队对任何人的注册都是开放邀请。欢迎随时前来提问。

## 贡献

非常感谢您的贡献。维护人员积极管理问题列表，并尝试突出适合新手的问题。
该项目遵循典型的 GitHub 拉取请求模型。有关更多详细信息，请参阅[CONTRIBUTING.md](CONTRIBUTING.md)。
在开始任何工作之前，请在现有问题上发表评论，或者提出新问题。

## 支持性说明

目前，Kubebuilder 正式支持 OSX 和 Linux 平台。
因此，如果您使用的是 Windows 操作系统，可能会遇到问题。欢迎为支持 Windows 而做出贡献。

### 苹果芯片

苹果芯片（`darwin/arm64`）支持始于 `go/v4` 插件。

## 社区会议

以下会议每两周举行一次：

- Kubebuilder 会议

欢迎参加。有关更多信息，请加入 [kubebuilder@googlegroups.com](https://groups.google.com/g/kubebuilder)。
每月，我们的团队在每个月的第一个星期四的 11:00 PT（太平洋时间）举行会议，讨论我们的进展并制定未来几周的计划。

[operator-sdk]: https://github.com/operator-framework/operator-sdk
[plugin-section]: https://book.kubebuilder.io/plugins/plugins.html
[controller-runtime]: https://github.com/kubernetes-sigs/controller-runtime
[your-own-plugins]: https://book.kubebuilder.io/plugins/creating-plugins.html
[controller-tools]: https://github.com/kubernetes-sigs/controller-tools
