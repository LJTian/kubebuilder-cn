# 教程：构建 CronJob

许多教程都以一些非常牵强的设置或一些用于传达基础知识的玩具应用程序开头，然后在更复杂的内容上停滞不前。相反，这个教程应该带您（几乎）完整地了解 Kubebuilder 的复杂性，从简单开始逐步构建到相当全面的内容。

我们假装（当然，这有点牵强）我们终于厌倦了在 Kubernetes 中使用非 Kubebuilder 实现的 CronJob 控制器的维护负担，我们想要使用 Kubebuilder 进行重写。

*CronJob* 控制器的任务（不是故意的双关语）是在 Kubernetes 集群上定期间隔运行一次性任务。它通过在 *Job* 控制器的基础上构建来完成这一点，*Job* 控制器的任务是运行一次性任务一次，并确保其完成。

我们不打算试图重写 Job 控制器，而是将其视为一个机会来了解如何与外部类型交互。

<aside class="note">

<h1>跟随 vs 跳过</h1>

请注意，本教程的大部分内容都是从位于书籍源目录中的可阅读的 Go 文件生成的：[docs/book/src/cronjob-tutorial/testdata][tutorial-source]。完整的可运行项目位于[project][tutorial-project-source]，而中间文件直接位于[testdata][tutorial-source]目录下。

[tutorial-source]: https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata

[tutorial-project-source]: https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project

</aside>

## 构建项目框架

如[快速入门](../quick-start.md)中所述，我们需要构建一个新项目的框架。确保您已经[安装了 Kubebuilder](../quick-start.md#installation)，然后构建一个新项目：

```bash
# 创建一个项目目录，然后运行初始化命令。
mkdir project
cd project
# 我们将使用 tutorial.kubebuilder.io 作为域，
# 因此所有 API 组将是 <group>.tutorial.kubebuilder.io。
kubebuilder init --domain tutorial.kubebuilder.io --repo tutorial.kubebuilder.io/project
```

<aside class="note">

您的项目名称默认为当前工作目录的名称。您可以传递 `--project-name=<dns1123-label-string>` 来设置不同的项目名称。

</aside>

现在我们已经有了一个项目框架，让我们来看看 Kubebuilder 到目前为止为我们生成了什么...

<aside class="note">

<h1>在 <code>$GOPATH</code> 中开发</h1>

如果您的项目是在 [`GOPATH`][GOPATH-golang-docs] 内初始化的，隐式调用的 `go mod init` 将为您插入模块路径。
否则，必须设置 `--repo=<module path>`。

如果对模块系统不熟悉，请阅读[Go 模块博文][go-modules-blogpost]。

</aside>

[GOPATH-golang-docs]: https://golang.org/doc/code.html#GOPATH
[go-modules-blogpost]: https://blog.golang.org/using-go-modules