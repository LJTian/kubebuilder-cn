/*
版权所有 2024 年 Kubernetes 作者。

根据 Apache 许可证 2.0 版（"许可证"）许可;
除非符合许可证的规定，否则您不得使用此文件。
您可以在以下网址获取许可证的副本：

    http://www.apache.org/licenses/LICENSE-2.0

除非适用法律要求或经书面同意，否则根据许可证分发的软件
按"原样"提供，不附带任何担保或条件，无论是明示的还是暗示的。
请查看许可证以了解特定语言下的权限和限制。
*/
// +kubebuilder:docs-gen:collapse=Apache License

/*
当我们在[上一章](/cronjob-tutorial/new-api.md)中使用 `kubebuilder create api` 创建 CronJob API 时，Kubebuilder 已经为您做了一些测试工作。
Kubebuilder 生成了一个 `internal/controller/suite_test.go` 文件，其中包含了设置测试环境的基本内容。

首先，它将包含必要的导入项。
*/

package controller

// 这些测试使用 Ginkgo（BDD 风格的 Go 测试框架）。请参考
// http://onsi.github.io/ginkgo/ 了解更多关于 Ginkgo 的信息。

// +kubebuilder:docs-gen:collapse=Imports

/*
现在，让我们来看一下生成的代码。
*/

var (
    cfg       *rest.Config
    k8sClient client.Client // 您将在测试中使用此客户端。
    testEnv   *envtest.Environment
    ctx       context.Context
    cancel    context.CancelFunc
)

func TestControllers(t *testing.T) {
    RegisterFailHandler(Fail)

    RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
    // 省略了一些设置代码
})

/*
Kubebuilder 还生成了用于清理 envtest 并在控制器目录中实际运行测试文件的样板函数。
您不需要修改这些函数。
*/

var _ = AfterSuite(func() {
    // 省略了一些清理代码
})

/*
现在，您的控制器在测试集群上运行，并且已准备好在您的 CronJob 上执行操作的客户端，我们可以开始编写集成测试了！
*/
