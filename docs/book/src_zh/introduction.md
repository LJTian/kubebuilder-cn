**注意：** 急于开始的读者可以直接前往[快速入门](quick-start.md)。

**正在使用 Kubebuilder v1 或 v2 吗？请查看 [v1](https://book-v1.book.kubebuilder.io) 或 [v2](https://book-v2.book.kubebuilder.io) 的旧版文档。**

## 适用对象

#### Kubernetes 用户

Kubernetes 用户将通过学习 API 被设计和实现的基本概念，深入了解 Kubernetes，并且将会开发出更深刻的认识。本书将教会读者如何开发自己的 Kubernetes API，以及核心 Kubernetes API 设计的原则。

包括：

- Kubernetes API 和资源的结构
- API 版本语义
- 自愈
- 垃圾回收和终结器
- 声明式 vs 命令式 API
- 基于级别 vs 基于边缘的 API
- 资源 vs 子资源

#### Kubernetes API 扩展开发者

API 扩展开发者将学习实现规范 Kubernetes API 背后的原则和概念，以及快速执行的简单工具和库。本书涵盖了扩展开发者常遇到的陷阱和误解。

包括：

- 如何将多个事件批量处理为单个协调调用
- 如何配置定期协调
- *即将推出*
    - 何时使用列表缓存 vs 实时查找
    - 垃圾回收 vs 终结器
    - 如何使用声明式 vs Webhook 验证
    - 如何实现 API 版本管理

## 为什么选择 Kubernetes API

Kubernetes API 为对象提供了一致和明确定义的端点，这些对象遵循一致和丰富的结构。

这种方法培育了一个丰富的工具和库生态系统，用于处理 Kubernetes API。

用户通过将对象声明为 *yaml* 或 *json* 配置，并使用常见工具来管理对象来使用这些 API。

将服务构建为 Kubernetes API 相比于普通的 REST，提供了许多优势，包括：

* 托管的 API 端点、存储和验证。
* 丰富的工具和 CLI，如 `kubectl` 和 `kustomize`。
* 对 AuthN 和细粒度 AuthZ 的支持。
* 通过 API 版本控制和转换支持 API 演进。
* 促进自适应/自愈 API 的发展，这些 API 可以持续响应系统状态的变化，而无需用户干预。
* Kubernetes 作为托管环境

开发人员可以构建并发布自己的 Kubernetes API，以安装到运行中的 Kubernetes 集群中。

## 贡献

如果您想要为本书或代码做出贡献，请先阅读我们的[贡献](https://github.com/kubernetes-sigs/kubebuilder/blob/master/CONTRIBUTING.md)指南。

## 资源

* 仓库: [sigs.k8s.io/kubebuilder](https://sigs.k8s.io/kubebuilder)
* Slack 频道: [#kubebuilder](http://slack.k8s.io/#kubebuilder)
* Google Group: [kubebuilder@googlegroups.com](https://groups.google.com/forum/#!forum/kubebuilder)