## 入门指南

### 概述

通过遵循[Operator 模式][k8s-operator-pattern]，不仅可以提供所有预期的资源，还可以在执行时动态、以编程方式管理它们。为了说明这个想法，想象一下，如果有人意外更改了配置或者误删了某个资源；在这种情况下，操作员可以在没有任何人工干预的情况下进行修复。

### 示例项目

我们将创建一个示例项目，以便让您了解它是如何工作的。这个示例将会：

- 对账一个 Memcached CR - 代表着在集群上部署/管理的 Memcached 实例
- 创建一个使用 Memcached 镜像的 Deployment
- 不允许超过 CR 中定义的大小的实例
- 更新 Memcached CR 的状态

请按照以下步骤操作。

### 创建项目

首先，创建一个用于您的项目的目录，并进入该目录，然后使用 `kubebuilder` 进行初始化：

```shell
mkdir $GOPATH/memcached-operator
cd $GOPATH/memcached-operator
kubebuilder init --domain=example.com
```

### 创建 Memcached API (CRD)

接下来，我们将创建一个新的 API，负责部署和管理我们的 Memcached 解决方案。在这个示例中，我们将使用[Deploy Image 插件][deploy-image]来获取我们解决方案的全面代码实现。

```shell
kubebuilder create api --group cache \
  --version v1alpha1 \
  --kind Memcached \
  --image=memcached:1.4.36-alpine \
  --image-container-command="memcached,-m=64,-o,modern,-v" \
  --image-container-port="11211" \
  --run-as-user="1001" \
  --plugins="deploy-image/v1-alpha" \
  --make=false
```

### 理解 API

这个命令的主要目的是为 Memcached 类型生成自定义资源（CR）和自定义资源定义（CRD）。它使用 group `cache.example.com` 和 version `v1alpha1` 来唯一标识 Memcached 类型的新 CRD。通过利用 Kubebuilder 工具，我们可以为这些平台定义我们的 API 和对象。虽然在这个示例中我们只添加了一种资源类型，但您可以根据需要拥有尽可能多的 `Groups` 和 `Kinds`。简而言之，CRD 是我们自定义对象的定义，而 CR 是它们的实例。

### 定义您的 API

在这个示例中，可以看到 Memcached 类型（CRD）具有一些特定规格。这些是由 Deploy Image 插件构建的，用于管理目的的默认脚手架：

#### 状态和规格

`MemcachedSpec` 部分是我们封装所有可用规格和配置的地方，用于我们的自定义资源（CR）。此外，值得注意的是，我们使用了状态条件。这确保了对 Memcached CR 的有效管理。当发生任何更改时，这些条件为我们提供了必要的数据，以便在 Kubernetes 集群中了解此资源的当前状态。这类似于我们为 Deployment 资源获取的状态信息。

从：`api/v1alpha1/memcached_types.go`

```go
// MemcachedSpec 定义了 Memcached 的期望状态
type MemcachedSpec struct {
	// 插入其他规格字段 - 集群的期望状态
	// 重要：修改此文件后运行 "make" 以重新生成代码
	// Size 定义了 Memcached 实例的数量
	// 以下标记将使用 OpenAPI v3 schema 来验证该值
	// 了解更多信息：https://book.kubebuilder.io/reference/markers/crd-validation.html
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=3
	// +kubebuilder:validation:ExclusiveMaximum=false
	Size int32 `json:"size,omitempty"`

	// Port 定义了将用于使用镜像初始化容器的端口
	ContainerPort int32 `json:"containerPort,omitempty"`
}

// MemcachedStatus 定义了 Memcached 的观察状态
type MemcachedStatus struct {
	// 代表了 Memcached 当前状态的观察结果
	// Memcached.status.conditions.type 为："Available"、"Progressing" 和 "Degraded"
	// Memcached.status.conditions.status 为 True、False、Unknown 中的一个
	// Memcached.status.conditions.reason 的值应为驼峰字符串，特定条件类型的产生者可以为此字段定义预期值和含义，以及这些值是否被视为 API 的保证
	// Memcached.status.conditions.Message 是一个人类可读的消息，指示有关转换的详细信息
	// 了解更多信息：https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}
```

因此，当我们向此文件添加新规格并执行 `make generate` 命令时，我们使用 [controller-gen][controller-gen] 生成了 CRD 清单，该清单位于 `config/crd/bases` 目录下。

#### 标记和验证

此外，值得注意的是，我们正在使用 `标记`，例如 `+kubebuilder:validation:Minimum=1`。这些标记有助于定义验证和标准，确保用户提供的数据 - 当他们为 Memcached 类型创建或编辑自定义资源时 - 得到适当的验证。有关可用标记的全面列表和详细信息，请参阅[标记文档][markers]。

观察 CRD 中的验证模式；此模式确保 Kubernetes API 正确验证应用的自定义资源（CR）：

从：[config/crd/bases/cache.example.com_memcacheds.yaml](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/getting-started/testdata/project/config/crd/bases/cache.example.com_memcacheds.yaml)

```yaml
description: MemcachedSpec 定义了 Memcached 的期望状态
properties:
  containerPort:
    description: Port 定义了将用于使用镜像初始化容器的端口
    format: int32
    type: integer
  size:
    description: 'Size 定义了 Memcached 实例的数量 以下标记将使用 OpenAPI v3 schema 来验证该值 了解更多信息：https://book.kubebuilder.io/reference/markers/crd-validation.html'
    format: int32
    maximum: 3 ## 从标记 +kubebuilder:validation:Maximum=3 生成
    minimum: 1 ## 从标记 +kubebuilder:validation:Minimum=1 生成
    type: integer
```

#### 自定义资源示例

位于 "config/samples" 目录下的清单作为可以应用于集群的自定义资源的示例。
在这个特定示例中，通过将给定资源应用到集群中，我们将生成一个大小为 1 的 Deployment 实例（参见 `size: 1`）。

从：[config/samples/cache_v1alpha1_memcached.yaml](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/getting-started/testdata/project/config/samples/cache_v1alpha1_memcached.yaml)

```yaml
{{#include ./getting-started/testdata/project/config/samples/cache_v1alpha1_memcached.yaml}}
```

### 对账过程

对账函数在确保资源和其规格之间基于其中嵌入的业务逻辑的同步方面起着关键作用。它的作用类似于循环，不断检查条件并执行操作，直到所有条件符合其实现。以下是伪代码来说明这一点：

```go
reconcile App {

  // 检查应用的 Deployment 是否存在，如果不存在则创建一个
  // 如果出现错误，则重新开始对账
  if err != nil {
    return reconcile.Result{}, err
  }

  // 检查应用的 Service 是否存在，如果不存在则创建一个
  // 如果出现错误，则重新开始对账
  if err != nil {
    return reconcile.Result{}, err
  }

  // 查找数据库 CR/CRD
  // 检查数据库 Deployment 的副本大小
  // 如果 deployment.replicas 的大小与 cr.size 不匹配，则更新它
  // 然后，从头开始对账。例如，通过返回 `reconcile.Result{Requeue: true}, nil`。
  if err != nil {
    return reconcile.Result{Requeue: true}, nil
  }
  ...

  // 如果循环结束时：
  // 所有操作都成功执行，对账就可以停止了
  return reconcile.Result{}, nil

}
```

#### 返回选项

以下是重新开始对账的一些可能返回选项：

- 带有错误：

```go
return ctrl.Result{}, err
```
- 没有错误：

```go
return ctrl.Result{Requeue: true}, nil
``` 

- 停止对账，使用（执行成功之后，或者不需要再进行对账）：

```go
return ctrl.Result{}, nil
```

- X 时间后重新开始对账：

```go
return ctrl.Result{RequeueAfter: nextRun.Sub(r.Now())}, nil
```

#### 在我们的示例中

当将自定义资源应用到集群时，有一个指定的控制器来管理 Memcached 类型。您可以检查其对账是如何实现的：

从：[internal/controller/memcached_controller.go](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/getting-started/testdata/project/internal/controller/memcached_controller.go)
```go
func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// 获取 Memcached 实例
	// 目的是检查是否在集群上应用了 Memcached 类型的自定义资源
	// 如果没有，我们将返回 nil 以停止对账过程
	memcached := &examplecomv1alpha1.Memcached{}
	err := r.Get(ctx, req.NamespacedName, memcached)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// 如果找不到自定义资源，通常意味着它已被删除或尚未创建
			// 这样，我们将停止对账过程
			log.Info("未找到 memcached 资源。忽略，因为对象可能已被删除")
			return ctrl.Result{}, nil
		}
		// 读取对象时出错 - 重新排队请求
		log.Error(err, "获取 memcached 失败")
		return ctrl.Result{}, err
	}

	// 当没有状态可用时，让我们将状态设置为 Unknown
	if memcached.Status.Conditions == nil || len(memcached.Status.Conditions) == 0 {
		meta.SetStatusCondition(&memcached.Status.Conditions,
			metav1.Condition{
				Type: typeAvailableMemcached,
				Status: metav1.ConditionUnknown,
				Reason: "对账中",
				Message: "开始对账"
			})
		if err = r.Status().Update(ctx, memcached); err != nil {
			log.Error(err, "更新 Memcached 状态失败")
			return ctrl.Result{}, err
		}

		// 更新状态后，让我们重新获取 memcached 自定义资源
		// 以便我们在集群上拥有资源的最新状态，并且避免
		// 引发错误 "对象已被修改，请将您的更改应用到最新版本，然后重试"
		// 如果我们尝试在后续操作中再次更新它，这将重新触发对账过程
		if err := r.Get(ctx, req.NamespacedName, memcached); err != nil {
			log.Error(err, "重新获取 memcached 失败")
			return ctrl.Result{}, err
		}
	}

	// 添加 finalizer。然后，我们可以定义在删除自定义资源之前应执行的一些操作。
	// 更多信息：https://kubernetes.io/docs/concepts/overview/working-with-objects/finalizers
	if !controllerutil.ContainsFinalizer(memcached, memcachedFinalizer) {
		log.Info("为 Memcached 添加 Finalizer")
		if ok := controllerutil.AddFinalizer(memcached, memcachedFinalizer); !ok {
			log.Error(err, "无法将 finalizer 添加到自定义资源")
			return ctrl.Result{Requeue: true}, nil
		}

		if err = r.Update(ctx, memcached); err != nil {
			log.Error(err, "更新自定义资源以添加 finalizer 失败")
			return ctrl.Result{}, err
		}
	}

	// 检查是否标记要删除 Memcached 实例，这通过设置删除时间戳来表示。
	isMemcachedMarkedToBeDeleted := memcached.GetDeletionTimestamp() != nil
	if isMemcachedMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(memcached, memcachedFinalizer) {
			log.Info("在删除 CR 之前执行 Finalizer 操作")

			// 在这里添加一个状态 "Downgrade"，以反映该资源开始其终止过程。
			meta.SetStatusCondition(&memcached.Status.Conditions, 
				metav1.Condition{
					Type: typeDegradedMemcached,
					Status: metav1.ConditionUnknown,
				},
				Reason: "Finalizing",
				Message: fmt.Sprintf("执行自定义资源的 finalizer 操作：%s ", memcached.Name)})

			if err := r.Status().Update(ctx, memcached); err != nil {
				log.Error(err, "更新 Memcached 状态失败")
				return ctrl.Result{}, err
			}

			// 执行在删除 finalizer 之前需要的所有操作，并允许
			// Kubernetes API 删除自定义资源。
			r.doFinalizerOperationsForMemcached(memcached)

			// TODO（用户）：如果您在 doFinalizerOperationsForMemcached 方法中添加操作
			// 那么您需要确保一切顺利，然后再删除和更新 Downgrade 状态
			// 否则，您应该在此重新排队。

			// 在更新状态前重新获取 memcached 自定义资源
			// 以便我们在集群上拥有资源的最新状态，并且避免
			// 引发错误 "对象已被修改，请将您的更改应用到最新版本，然后重试"
			// 如果我们尝试在后续操作中再次更新它，这将重新触发对账过程
			if err := r.Get(ctx, req.NamespacedName, memcached); err != nil {
				log.Error(err, "重新获取 memcached 失败")
				return ctrl.Result{}, err
			}

			meta.SetStatusCondition(&memcached.Status.Conditions,
				metav1.Condition{
					Type: typeDegradedMemcached,
					Status: metav1.ConditionTrue,
					Reason: "Finalizing",
					Message: fmt.Sprintf("自定义资源 %s 的 finalizer 操作已成功完成", memcached.Name)
				})

			if err := r.Status().Update(ctx, memcached); err != nil {
				log.Error(err, "更新 Memcached 状态失败")
				return ctrl.Result{}, err
			}

			log.Info("成功执行操作后移除 Memcached 的 Finalizer")
			if ok := controllerutil.RemoveFinalizer(memcached, memcachedFinalizer); !ok {
				log.Error(err, "移除 Memcached 的 finalizer 失败")
				return ctrl.Result{Requeue: true}, nil
			}

			if err := r.Update(ctx, memcached); err != nil {
				log.Error(err, "移除 Memcached 的 finalizer 失败")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// 检查部署是否已经存在，如果不存在则创建新的
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: memcached.Name, Namespace: memcached.Namespace}, found)
	if err != nil && apierrors.IsNotFound(err) {
		// 定义一个新的部署
		dep, err := r.deploymentForMemcached(memcached)
		if err != nil {
			log.Error(err, "为 Memcached 定义新的 Deployment 资源失败")

			// 以下实现将更新状态
			meta.SetStatusCondition(&memcached.Status.Conditions, metav1.Condition{
				Type: typeAvailableMemcached,
				Status: metav1.ConditionFalse,
				Reason: "对账中",
				Message: fmt.Sprintf("为自定义资源创建 Deployment 失败 (%s): (%s)", memcached.Name, err)})

			if err := r.Status().Update(ctx, memcached); err != nil {
				log.Error(err, "更新 Memcached 状态失败")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, err
		}

		log.Info("创建新的 Deployment",
			"Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		if err = r.Create(ctx, dep); err != nil {
			log.Error(err, "创建新的 Deployment 失败",
				"Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}

		// 部署成功创建
		// 我们将重新排队对账，以便确保状态
		// 并继续进行下一步操作
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	} else if err != nil {
		log.Error(err, "获取 Deployment 失败")
		// 让我们返回错误以重新触发对账
		return ctrl.Result{}, err
	}

	// CRD API 定义了 Memcached 类型具有 MemcachedSpec.Size 字段
	// 以设置集群上所需的 Deployment 实例数量。
	// 因此，以下代码将确保 Deployment 大小与我们对账的自定义资源的 Size spec 相同。
	size := memcached.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		if err = r.Update(ctx, found); err != nil {
			log.Error(err, "更新 Deployment 失败",
				"Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)

			// 在更新状态前重新获取 memcached 自定义资源
			// 以便我们在集群上拥有资源的最新状态，并且避免
			// 引发错误 "对象已被修改，请将您的更改应用到最新版本，然后重试"
			// 如果我们尝试在后续操作中再次更新它，这将重新触发对账过程
			if err := r.Get(ctx, req.NamespacedName, memcached); err != nil {
				log.Error(err, "重新获取 memcached 失败")
				return ctrl.Result{}, err
			}

			// 以下实现将更新状态
			meta.SetStatusCondition(&memcached.Status.Conditions, 
				metav1.Condition{
					Type: typeAvailableMemcached,
					Status: metav1.ConditionFalse,
					Reason: "调整大小",
					Message: fmt.Sprintf("更新自定义资源的大小失败 (%s): (%s)", memcached.Name, err)
				})

			if err := r.Status().Update(ctx, memcached); err != nil {
				log.Error(err, "更新 Memcached 状态失败")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, err
		}

		// 现在，我们更新大小后，希望重新排队对账
		// 以便确保我们拥有资源的最新状态
		// 并帮助确保集群上的期望状态
		return ctrl.Result{Requeue: true}, nil
	}

	// 以下实现将更新状态
	meta.SetStatusCondition(&memcached.Status.Conditions,
		metav1.Condition{
			Type: typeAvailableMemcached,
			Status: metav1.ConditionTrue,
			Reason: "对账中",
			Message: fmt.Sprintf("为自定义资源创建 %d 个副本的 Deployment 成功", memcached.Name, size)
		})

	if err := r.Status().Update(ctx, memcached); err != nil {
		log.Error(err, "更新 Memcached 状态失败")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
```

#### 观察集群上的变化

该控制器持续地观察与该类型相关的任何事件。因此，相关的变化会立即触发控制器的对账过程。值得注意的是，我们已经实现了 `watches` 特性。[(更多信息)][watches]。这使我们能够监视与创建、更新或删除 Memcached 类型的自定义资源相关的事件，以及由其相应控制器编排和拥有的 Deployment。请注意以下代码：

```go
// SetupWithManager 使用 Manager 设置控制器。
// 请注意，也将监视 Deployment 以确保其在集群中处于期望的状态
func (r *MemcachedReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
    For(&examplecomv1alpha1.Memcached{}). // 为 Memcached 类型创建监视
    Owns(&appsv1.Deployment{}). // 为其控制器拥有的 Deployment 创建监视
    Complete(r)
}
```

<aside class="note">
<h1>为 Deployment 设置 ownerRef</h1>

请注意，当我们创建用于运行 Memcached 镜像的 Deployment 时，我们正在设置引用：

```go
// 为 Deployment 设置 ownerRef
// 更多信息：https://kubernetes.io/docs/concepts/overview/working-with-objects/owners-dependents/
if err := ctrl.SetControllerReference(memcached, dep, r.Scheme); err != nil {
    return nil, err
}

```

</aside>

### 设置 RBAC 权限

现在通过 [RBAC markers][rbac-markers] 配置了 [RBAC 权限][k8s-rbac]，用于生成和更新 `config/rbac/` 中的清单文件。这些标记可以在每个控制器的 `Reconcile()` 方法中找到（并应该被定义），请看我们示例中的实现方式：

```go
//+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
```

重要的是，如果您希望添加或修改 RBAC 规则，可以通过更新或添加控制器中的相应标记来实现。在进行必要的更改后，运行 `make generate` 命令。这将促使 [controller-gen][controller-gen] 刷新位于 `config/rbac` 下的文件。

<aside class="note">
<h1>在 config/rbac 下生成 RBAC 规则</h1>

对于每个类型，Kubebuilder 将生成具有查看和编辑权限的脚手架规则（例如 `memcached_editor_role.yaml` 和 `memcached_viewer_role.yaml`）。
当您使用 `make deploy IMG=myregistery/example:1.0.0` 部署解决方案时，这些规则不会应用于集群。
这些规则旨在帮助系统管理员知道在授予用户组权限时应允许什么。

</aside>

### Manager（main.go）

[Manager][manager] 在监督控制器方面扮演着至关重要的角色，这些控制器进而使集群端的操作成为可能。如果您检查 `cmd/main.go` 文件，您会看到以下内容：

```go
...
    mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
        Scheme:                 scheme,
        Metrics:                metricsserver.Options{BindAddress: metricsAddr},
        HealthProbeBindAddress: probeAddr,
        LeaderElection:         enableLeaderElection,
        LeaderElectionID:       "1836d577.testproject.org",
        // LeaderElectionReleaseOnCancel 定义了领导者在 Manager 结束时是否应主动放弃领导权。
        // 这要求二进制在 Manager 停止时立即结束，否则此设置是不安全的。设置此选项显著加快主动领导者转换的速度，
        // 因为新领导者无需等待 LeaseDuration 时间。
        //
        // 在提供的默认脚手架中，程序在 Manager 停止后立即结束，因此启用此选项是可以的。但是，
        // 如果您正在进行任何操作，例如在 Manager 停止后执行清理操作，那么使用它可能是不安全的。
        // LeaderElectionReleaseOnCancel: true,
    })
    if err != nil {
        setupLog.Error(err, "无法启动 Manager")
        os.Exit(1)
    }
```

上面的代码片段概述了 Manager 的配置[选项][options-manager]。虽然我们在当前示例中不会更改这些选项，但了解其位置以及初始化您的基于 Operator 的镜像的过程非常重要。Manager 负责监督为您的 Operator API 生成的控制器。

### 检查在集群中运行的项目

此时，您可以执行 [快速入门][quick-start] 中突出显示的命令。通过执行 `make build IMG=myregistry/example:1.0.0`，您将为项目构建镜像。出于测试目的，建议将此镜像发布到公共注册表。这样可以确保轻松访问，无需额外的配置。完成后，您可以使用 `make deploy IMG=myregistry/example:1.0.0` 命令将镜像部署到集群中。

## 下一步

- 要深入了解开发解决方案，请考虑阅读提供的教程。
- 要了解优化您的方法的见解，请参阅[最佳实践][best-practices]文档。

[k8s-operator-pattern]: https://kubernetes.io/docs/concepts/extend-kubernetes/operator/
[controller-runtime]: https://github.com/kubernetes-sigs/controller-runtime
[group-kind-oh-my]: ./cronjob-tutorial/gvks.md
[controller-gen]: ./reference/controller-gen.md
[markers]: ./reference/markers.md
[watches]: ./reference/watching-resources.md
[rbac-markers]: ./reference/markers/rbac.md
[k8s-rbac]: https://kubernetes.io/docs/reference/access-authn-authz/rbac/
[manager]: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/manager
[options-manager]: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/manager#Options
[quick-start]: ./quick-start.md
[best-practices]: ./reference/good-practices.md
