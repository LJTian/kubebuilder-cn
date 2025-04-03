# 创建一个新的 API

要创建一个新的 Kind（你有关注[上一章](./gvks.md#kinds-and-resources)的内容吗？）以及相应的控制器，我们可以使用 `kubebuilder create api` 命令：

```bash
kubebuilder create api --group batch --version v1 --kind CronJob
```

按下 `y` 键来选择 "Create Resource" 和 "Create Controller"。

第一次针对每个 group-version 调用此命令时，它将为新的 group-version 创建一个目录。

<aside class="note">

<h1>支持旧的集群版本</h1>

与你的 Go API 类型一起创建的默认 CustomResourceDefinition 清单使用 API 版本 `v1`。如果你的项目打算支持早于 v1.16 的 Kubernetes 集群版本，你必须设置 `--crd-version v1beta1`，并从 `CRD_OPTIONS` Makefile 变量中移除 `preserveUnknownFields=false`。详情请参阅[CustomResourceDefinition 生成参考][crd-reference]。

[crd-reference]: /reference/generating-crd.md#supporting-older-cluster-versions

</aside>

在这种情况下，将创建一个名为 [`api/v1/`](https://sigs.k8s.io/kubebuilder/docs/book/src/cronjob-tutorial/testdata/project/api/v1) 的目录，对应于 `batch.tutorial.kubebuilder.io/v1`（还记得我们从一开始设置的 `--domain` 吗？）。

它还将添加一个用于我们的 `CronJob` Kind 的文件，即 `api/v1/cronjob_types.go`。每次使用不同的 Kind 调用该命令时，它都会添加一个相应的新文件。

让我们看看我们得到了什么，然后我们可以开始填写它。

{{#literatego ./testdata/emptyapi.go}}

现在我们已经了解了基本结构，让我们继续填写它！