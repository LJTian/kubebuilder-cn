# 快速入门

本快速入门指南将涵盖以下内容：

- [创建项目](#create-a-project)
- [创建 API](#create-an-api)
- [本地运行](#test-it-out)
- [集群运行](#run-it-on-the-cluster)

## 先决条件

- [go](https://golang.org/dl/) 版本 v1.20.0+
- [docker](https://docs.docker.com/install/) 版本 17.03+。
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) 版本 v1.11.3+。
- 访问 Kubernetes v1.11.3+ 集群。

<aside class="note">
<h1>版本和支持性</h1>

Kubebuilder 创建的项目包含一个 Makefile，在创建时将安装工具的版本。这些工具包括：
- [kustomize](https://github.com/kubernetes-sigs/kustomize)
- [controller-gen](https://github.com/kubernetes-sigs/controller-tools)

在 `Makefile` 和 `go.mod` 文件中定义的版本是经过测试的版本，因此建议使用指定的版本。

</aside>

## 安装

安装 [kubebuilder](https://sigs.k8s.io/kubebuilder)：

```bash
# 下载 kubebuilder 并在本地安装。
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)"
chmod +x kubebuilder && mv kubebuilder /usr/local/bin/
```

<aside class="note">
<h1>使用主分支</h1>

您可以通过从 `https://go.kubebuilder.io/dl/master/$(go env GOOS)/$(go env GOARCH)` 安装来使用主分支的快照。

</aside>

<aside class="note">
<h1>启用 shell 自动补全</h1>

Kubebuilder 通过命令 `kubebuilder completion <bash|fish|powershell|zsh>` 提供自动补全支持，可以节省大量输入。有关详细信息，请参阅[自动补全](./reference/completion.md)文档。

</aside>

## 创建项目

创建一个目录，然后在其中运行 init 命令以初始化一个新项目。以下是一个示例。

```bash
mkdir -p ~/projects/guestbook
cd ~/projects/guestbook
kubebuilder init --domain my.domain --repo my.domain/guestbook
```

<aside class="note">
<h1>在 $GOPATH 中开发</h1>

如果您的项目在 [`GOPATH`][GOPATH-golang-docs] 中初始化，隐式调用的 `go mod init` 将为您插入模块路径。
否则，必须设置 `--repo=<module path>`。

如果对模块系统不熟悉，请阅读[Go 模块博文][go-modules-blogpost]。

</aside>


## 创建 API

运行以下命令以创建一个名为 `webapp/v1` 的新 API（组/版本），并在其中创建一个名为 `Guestbook` 的新 Kind（CRD）：

```bash
kubebuilder create api --group webapp --version v1 --kind Guestbook
```

<aside class="note">
<h1>按键选项</h1>

如果按 `y` 键创建资源 [y/n] 和创建控制器 [y/n]，则会创建文件 `api/v1/guestbook_types.go`，其中定义了 API，
以及 `internal/controllers/guestbook_controller.go`，其中实现了此 Kind（CRD）的调和业务逻辑。

</aside>


**可选步骤：** 编辑 API 定义和调和业务逻辑。有关更多信息，请参阅[设计 API](/cronjob-tutorial/api-design.md)和[控制器概述](cronjob-tutorial/controller-overview.md)。

如果您正在编辑 API 定义，可以使用以下命令生成诸如自定义资源（CRs）或自定义资源定义（CRDs）之类的清单：
```bash
make manifests
```

<details><summary>点击此处查看示例。<tt>(api/v1/guestbook_types.go)</tt></summary>
<p>

```go
// GuestbookSpec 定义了 Guestbook 的期望状态
type GuestbookSpec struct {
	// 插入其他规范字段 - 集群的期望状态
	// 重要提示：在修改此文件后运行 "make" 以重新生成代码

	// 实例数量
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	Size int32 `json:"size"`

	// GuestbookSpec 配置的 ConfigMap 名称
	// +kubebuilder:validation:MaxLength=15
	// +kubebuilder:validation:MinLength=1
	ConfigMapName string `json:"configMapName"`

	// +kubebuilder:validation:Enum=Phone;Address;Name
	Type string `json:"alias,omitempty"`
}

// GuestbookStatus 定义了 Guestbook 的观察状态
type GuestbookStatus struct {
	// 插入其他状态字段 - 定义集群的观察状态
	// 重要提示：在修改此文件后运行 "make" 以重新生成代码

	// 活动的 Guestbook 节点的 PodName
	Active string `json:"active"`

	// 待机的 Guestbook 节点的 PodNames
	Standby []string `json:"standby"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// Guestbook 是 guestbooks API 的架构
type Guestbook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GuestbookSpec   `json:"spec,omitempty"`
	Status GuestbookStatus `json:"status,omitempty"`
}
```

</p>
</details>


## 测试

您需要一个 Kubernetes 集群来运行。您可以使用 [KIND](https://sigs.k8s.io/kind) 获取本地集群进行测试，或者运行在远程集群上。

<aside class="note">
<h1>使用的上下文</h1>

您的控制器将自动使用 kubeconfig 文件中的当前上下文（即 `kubectl cluster-info` 显示的集群）。

</aside>

将 CRD 安装到集群中：
```bash
make install
```

为了快速反馈和代码级调试，运行您的控制器（这将在前台运行，如果要保持运行状态，请切换到新终端）：
```bash
make run
```

## 安装自定义资源实例

如果按 `y` 键创建资源 [y/n]，则会在您的样本中为您的 CRD 创建一个 CR（如果已更改 API 定义，请确保先编辑它们）：

```bash
kubectl apply -k config/samples/
```

## 在集群上运行

当您的控制器准备好打包并在其他集群中进行测试时。

构建并将您的镜像推送到 `IMG` 指定的位置：

```bash
make docker-build docker-push IMG=<some-registry>/<project-name>:tag
```

使用由 `IMG` 指定的镜像将控制器部署到集群中：

```bash
make deploy IMG=<some-registry>/<project-name>:tag
```

<aside class="note">
<h1>注册表权限</h1>

此镜像应发布在您指定的个人注册表中。如果上述命令不起作用，需要有权限从工作环境中拉取镜像。

<h1>RBAC 错误</h1>

如果遇到 RBAC 错误，您可能需要授予自己 cluster-admin 权限，或者以管理员身份登录。请参阅 [使用 Kubernetes RBAC 的先决条件 GKE 集群 v1.11.x 和更旧版本][pre-rbc-gke]，这可能是您的情况。

</aside>

## 卸载 CRD

从集群中删除您的 CRD：

```bash
make uninstall
```

## 卸载控制器

从集群中卸载控制器：

```bash
make undeploy
```

## 下一步

现在，查看[架构概念图][architecture-concept-diagram]以获得更好的概述，并跟随[CronJob 教程][cronjob-tutorial]，以便通过开发演示示例项目更好地了解其工作原理。

<aside class="note">
<h1> 使用 Deploy Image 插件生成 API 和控制器代码 </h1>

确保您查看[Deploy Image](https://book.kubebuilder.io/plugins/deploy-image-plugin-v1-alpha.html)插件。此插件允许用户生成 API/控制器以部署和管理集群中的操作数（镜像），遵循指南和最佳实践。它在抽象实现此目标的复杂性的同时，允许用户定制生成的代码。

</aside>

[pre-rbc-gke]: https://cloud.google.com/kubernetes-engine/docs/how-to/role-based-access-control#iam-rolebinding-bootstrap
[cronjob-tutorial]: https://book.kubebuilder.io/cronjob-tutorial/cronjob-tutorial.html
[GOPATH-golang-docs]: https://golang.org/doc/code.html#GOPATH
[go-modules-blogpost]: https://blog.golang.org/using-go-modules
[envtest]: https://book.kubebuilder.io/reference/testing/envtest.html
[architecture-concept-diagram]: architecture.md
[kustomize]: https://github.com/kubernetes-sigs/kustomize
