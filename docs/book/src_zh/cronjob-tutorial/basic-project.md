# 基本项目结构包含什么？

在构建新项目的框架时，Kubebuilder 为我们提供了一些基本的样板文件。

## 构建基础设施

首先是构建项目的基本基础设施：

<details><summary><code>go.mod</code>：与我们的项目匹配的新 Go 模块，具有基本依赖项</summary>

```go
{{#include ./testdata/project/go.mod}}
```
</details>

<details><summary><code>Makefile</code>：用于构建和部署控制器的 Make 目标</summary>

```makefile
{{#include ./testdata/project/Makefile}}
```
</details>

<details><summary><code>PROJECT</code>：用于构建新组件的 Kubebuilder 元数据</summary>

```yaml
{{#include ./testdata/project/PROJECT}}
```
</details>

## 启动配置

我们还在[`config/`](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project/config)目录下获得了启动配置。目前，它只包含了[Kustomize](https://sigs.k8s.io/kustomize) YAML 定义，用于在集群上启动我们的控制器，但一旦我们开始编写控制器，它还将包含我们的自定义资源定义、RBAC 配置和 Webhook 配置。

[`config/default`](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project/config/default) 包含了一个[Kustomize 基础配置](https://github.com/kubernetes-sigs/kubebuilder/blob/master/docs/book/src/cronjob-tutorial/testdata/project/config/default/kustomization.yaml)，用于在标准配置中启动控制器。

每个其他目录都包含一个不同的配置部分，重构为自己的基础配置：

- [`config/manager`](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project/config/manager)：在集群中将您的控制器作为 Pod 启动

- [`config/rbac`](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project/config/rbac)：在其自己的服务帐户下运行您的控制器所需的权限

## 入口点

最后，但肯定不是最不重要的，Kubebuilder 为我们的项目生成了基本的入口点：`main.go`。让我们接着看看...