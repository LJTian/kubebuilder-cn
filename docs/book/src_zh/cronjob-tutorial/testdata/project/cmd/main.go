/*
版权所有 2024 年 Kubernetes 作者。

根据 Apache 许可证 2.0 版（"许可证"）获得许可；
除非符合许可证的规定，否则您不得使用此文件。
您可以在以下网址获取许可证的副本：

    http://www.apache.org/licenses/LICENSE-2.0

除非适用法律要求或书面同意，根据许可证分发的软件是基于"按原样"的基础分发的，
没有任何明示或暗示的担保或条件。
请查看许可证以了解特定语言管理权限和限制。

*/
// +kubebuilder:docs-gen:collapse=Apache 许可证

package main

import (
	"crypto/tls"
	"flag"
	"os"

	// 导入所有 Kubernetes 客户端认证插件（例如 Azure、GCP、OIDC 等）
	// 以确保 exec-entrypoint 和 run 可以利用它们。
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	batchv1 "tutorial.kubebuilder.io/project/api/v1"
	"tutorial.kubebuilder.io/project/internal/controller"
	//+kubebuilder:scaffold:imports
)

// +kubebuilder:docs-gen:collapse=Imports

/*
要注意的第一个变化是，kubebuilder 已将新 API 组的包（`batchv1`）添加到我们的 scheme 中。
这意味着我们可以在我们的控制器中使用这些对象。

如果我们将使用任何其他 CRD，我们将不得不以相同的方式添加它们的 scheme。
诸如 Job 之类的内置类型通过 `clientgoscheme` 添加了它们的 scheme。
*/

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(batchv1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

/*
另一个发生变化的地方是，kubebuilder 已添加了一个块，调用我们的 CronJob 控制器的 `SetupWithManager` 方法。
*/

func main() {
	/*
	 */
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var secureMetrics bool
	var enableHTTP2 bool
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&secureMetrics, "metrics-secure", false,
		"If set the metrics endpoint is served securely")
	flag.BoolVar(&enableHTTP2, "enable-http2", false,
		"If set, HTTP/2 will be enabled for the metrics and webhook servers")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	// 如果 enable-http2 标志为 false（默认值），则应禁用 http/2
	// 由于其漏洞。更具体地说，禁用 http/2 将防止受到 HTTP/2 流取消和
	// 快速重置 CVE 的影响。更多信息请参见：
	// - https://github.com/advisories/GHSA-qppj-fm5r-hxr3
	// - https://github.com/advisories/GHSA-4374-p667-p6c8
	disableHTTP2 := func(c *tls.Config) {
		setupLog.Info("disabling http/2")
		c.NextProtos = []string{"http/1.1"}
	}

	tlsOpts := []func(*tls.Config){}
	if !enableHTTP2 {
		tlsOpts = append(tlsOpts, disableHTTP2)
	}

	webhookServer := webhook.NewServer(webhook.Options{
		TLSOpts: tlsOpts,
	})

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress:   metricsAddr,
			SecureServing: secureMetrics,
			TLSOpts:       tlsOpts,
		},
		WebhookServer:          webhookServer,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "80807133.tutorial.kubebuilder.io",
		// LeaderElectionReleaseOnCancel 定义了在 Manager 结束时领导者是否应主动下台
		//。这需要二进制文件在 Manager 停止后立即结束，否则，此设置是不安全的。设置这将显著
		// 加快自愿领导者过渡的速度，因为新领导者无需等待 LeaseDuration 时间。
		//
		// 在默认提供的脚手架中，程序在 Manager 停止后立即结束，因此可以启用此选项。
		// 但是，如果您正在执行或打算在 Manager 停止后执行任何操作，比如执行清理操作，
		// 那么它的使用可能是不安全的。
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// +kubebuilder:docs-gen:collapse=old stuff

	if err = (&controller.CronJobReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CronJob")
		os.Exit(1)
	}

	/*
		我们还将为我们的类型设置 webhooks，接下来我们将讨论它们。
		我们只需要将它们添加到 manager 中。由于我们可能希望单独运行 webhooks，
		或者在本地测试控制器时不运行它们，我们将它们放在一个环境变量后面。

		我们只需确保在本地运行时设置 `ENABLE_WEBHOOKS=false`。
	*/
	if os.Getenv("ENABLE_WEBHOOKS") != "false" {
		if err = (&batchv1.CronJob{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CronJob")
			os.Exit(1)
		}
	}
	//+kubebuilder:scaffold:builder

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
	// +kubebuilder:docs-gen:collapse=old stuff
}
