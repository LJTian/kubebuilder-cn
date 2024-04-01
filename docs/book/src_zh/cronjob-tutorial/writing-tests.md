# 编写控制器测试

测试 Kubernetes 控制器是一个庞大的主题，而 kubebuilder 为您生成的样板测试文件相对较少。

为了引导您了解 Kubebuilder 生成的控制器的集成测试模式，我们将回顾我们在第一个教程中构建的 CronJob，并为其编写一个简单的测试。

基本方法是，在生成的 `suite_test.go` 文件中，您将使用 envtest 创建一个本地 Kubernetes API 服务器，实例化和运行您的控制器，然后编写额外的 `*_test.go` 文件使用 [Ginkgo](http://onsi.github.io/ginkgo) 进行测试。

如果您想调整您的 envtest 集群的配置，请参阅 [为集成测试配置 envtest](../reference/envtest.md) 部分以及 [`envtest 文档`](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/envtest?tab=doc)。

## 测试环境设置

{{#literatego ../cronjob-tutorial/testdata/project/internal/controller/suite_test.go}}

## 测试控制器行为

{{#literatego ../cronjob-tutorial/testdata/project/internal/controller/cronjob_controller_test.go}}

上面的状态更新示例演示了一个用于自定义 Kind 与下游对象的一般测试策略。到目前为止，您希望已经学会了以下测试控制器行为的方法：

* 设置您的控制器在 envtest 集群上运行
* 编写用于创建测试对象的存根
* 隔离对象的更改以测试特定的控制器行为

## 高级示例

有更复杂的示例使用 envtest 严格测试控制器行为。示例包括：

* Azure Databricks Operator：查看他们完全完善的 [`suite_test.go`](https://github.com/microsoft/azure-databricks-operator/blob/0f722a710fea06b86ecdccd9455336ca712bf775/controllers/suite_test.go) 以及该目录中的任何 `*_test.go` 文件 [比如这个](https://github.com/microsoft/azure-databricks-operator/blob/0f722a710fea06b86ecdccd9455336ca712bf775/controllers/secretscope_controller_test.go)。