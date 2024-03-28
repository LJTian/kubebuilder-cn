/*
版权所有 2022。

根据 Apache 许可，版本 2.0 进行许可（"许可"）；
除非遵守许可，否则您不得使用此文件。
您可以在以下网址获取许可的副本：

    http://www.apache.org/licenses/LICENSE-2.0

除非适用法律要求或书面同意，根据许可分发的软件是基于"按原样"的基础分发的，
不附带任何明示或暗示的担保或条件。
请参阅许可以获取特定语言下的权限和限制。
*/
// +kubebuilder:docs-gen:collapse=Apache License

/*
首先，我们从一些标准的导入开始。
与之前一样，我们需要核心的 controller-runtime 库，以及 client 包和我们的 API 类型包。
*/
package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	batchv1 "tutorial.kubebuilder.io/project/api/v1"
)

/*
接下来，kubebuilder 为我们生成了一个基本的 reconciler 结构。
几乎每个 reconciler 都需要记录日志，并且需要能够获取对象，因此这些都是开箱即用的。
*/

// CronJobReconciler reconciles a CronJob object
type CronJobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

/*
大多数控制器最终都会在集群上运行，因此它们需要 RBAC 权限，我们使用 controller-tools 的 [RBAC markers](/reference/markers/rbac.md) 来指定这些权限。这些是运行所需的最低权限。
随着我们添加更多功能，我们将需要重新审视这些权限。
*/

// +kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=cronjobs/status,verbs=get;update;patch

/*
`ClusterRole` manifest 位于 `config/rbac/role.yaml`，通过以下命令使用 controller-gen 从上述标记生成：
*/

// make manifests

/*
注意：如果收到错误，请运行错误中指定的命令，然后重新运行 `make manifests`。
*/

/*
`Reconcile` 实际上执行单个命名对象的对账。
我们的 [Request](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/reconcile?tab=doc#Request) 只有一个名称，但我们可以使用 client 从缓存中获取该对象。

我们返回一个空结果和没有错误，这表示 controller-runtime 我们已成功对账了此对象，并且在有变更之前不需要再次尝试。

大多数控制器需要一个记录句柄和一个上下文，因此我们在这里设置它们。

[context](https://golang.org/pkg/context/) 用于允许取消请求，以及可能的跟踪等功能。它是所有 client 方法的第一个参数。`Background` 上下文只是一个基本上没有任何额外数据或时间限制的上下文。

记录句柄让我们记录日志。controller-runtime 通过一个名为 [logr](https://github.com/go-logr/logr) 的库使用结构化日志。很快我们会看到，日志记录通过将键值对附加到静态消息上来实现。我们可以在我们的 reconciler 的顶部预先分配一些键值对，以便将它们附加到此 reconciler 中的所有日志行。
*/
func (r *CronJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// your logic here

	return ctrl.Result{}, nil
}

/*
最后，我们将此 reconciler 添加到 manager 中，以便在启动 manager 时启动它。

目前，我们只指出此 reconciler 作用于 `CronJob`。稍后，我们将使用这个来标记我们关心相关的对象。
*/

func (r *CronJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.CronJob{}).
		Complete(r)
}
