/*
版权所有 2022 年 Kubernetes 作者。

根据 Apache 许可，版本 2.0（"许可"）获得许可；
除非符合许可的规定，否则您不得使用此文件。
您可以在以下网址获取许可的副本：

    http://www.apache.org/licenses/LICENSE-2.0

除非适用法律要求或书面同意，否则根据许可分发的软件将按"原样"分发，
不附带任何明示或暗示的担保或条件。
请参阅许可，了解特定语言下的权限和限制。
*/
// +kubebuilder:docs-gen:collapse=Apache License

/*

我们的包从一些基本的导入开始。特别是：

- 核心的 [controller-runtime](https://pkg.go.dev/sigs.k8s.io/controller-runtime?tab=doc) 库
- 默认的 controller-runtime 日志记录，[Zap](https://pkg.go.dev/go.uber.org/zap)（稍后会详细介绍）

*/

package main

import (
	"flag"
	"os"

	// 导入所有 Kubernetes 客户端认证插件（例如 Azure、GCP、OIDC 等）
	// 以确保 exec-entrypoint 和 run 可以利用它们。
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	// +kubebuilder:scaffold:imports
)

/*
每组控制器都需要一个
[*Scheme*](https://book.kubebuilder.io/cronjob-tutorial/gvks.html#err-but-whats-that-scheme-thing)，
它提供了 Kinds 与它们对应的 Go 类型之间的映射。稍后在编写 API 定义时，我们将更详细地讨论 Kinds，所以稍后再谈。
*/
var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	//+kubebuilder:scaffold:scheme
}

/*
此时，我们的主函数相当简单：

- 我们为指标设置了一些基本的标志。

- 我们实例化了一个
[*manager*](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/manager?tab=doc#Manager)，
它负责运行我们所有的控制器，并设置了共享缓存和客户端到 API 服务器的连接（请注意我们告诉 manager 关于我们的 Scheme）。

- 我们运行我们的 manager，它反过来运行我们所有的控制器和 Webhooks。
manager 被设置为在接收到优雅关闭信号之前一直运行。
这样，当我们在 Kubernetes 上运行时，我们会在 Pod 优雅终止时表现良好。

虽然目前我们没有任何东西要运行，但记住`+kubebuilder:scaffold:builder`注释的位置——很快那里会变得有趣起来。

*/

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Metrics: server.Options{
			BindAddress: metricsAddr,
		},
		WebhookServer:          webhook.NewServer(webhook.Options{Port: 9443}),
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "80807133.tutorial.kubebuilder.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	/*
		注意，`Manager` 可以通过以下方式限制所有控制器将监视资源的命名空间：
	*/

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Cache: cache.Options{
			DefaultNamespaces: map[string]cache.Config{
				namespace: {},
			},
		},
		Metrics: server.Options{
			BindAddress: metricsAddr,
		},
		WebhookServer:          webhook.NewServer(webhook.Options{Port: 9443}),
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "80807133.tutorial.kubebuilder.io",
	})

	/*
		上面的示例将把项目的范围更改为单个`Namespace`。在这种情况下，建议将提供的授权限制为此命名空间，
		方法是将默认的`ClusterRole`和`ClusterRoleBinding`替换为`Role`和`RoleBinding`。
		有关更多信息，请参阅 Kubernetes 关于使用 [RBAC 授权](https://kubernetes.io/docs/reference/access-authn-authz/rbac/) 的文档。

		此外，还可以使用 [`DefaultNamespaces`](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/cache#Options)
		从 `cache.Options{}` 缓存特定一组命名空间中的对象：
	*/

	var namespaces []string // 名称空间列表
	defaultNamespaces := make(map[string]cache.Config)

	for _, ns := range namespaces {
		defaultNamespaces[ns] = cache.Config{}
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Cache: cache.Options{
			DefaultNamespaces: defaultNamespaces,
		},
		Metrics: server.Options{
			BindAddress: metricsAddr,
		},
		WebhookServer:          webhook.NewServer(webhook.Options{Port: 9443}),
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "80807133.tutorial.kubebuilder.io",
	})

	/*
		有关更多信息，请参阅 [`cache.Options{}`](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/cache#Options)
	*/

	// +kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
