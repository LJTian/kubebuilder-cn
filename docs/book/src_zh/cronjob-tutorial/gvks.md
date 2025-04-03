# 组、版本和类型

实际上，在开始创建我们的 API 之前，我们应该稍微谈一下术语。

当我们在 Kubernetes 中讨论 API 时，我们经常使用 4 个术语：*groups*（组）、*versions*（版本）、*kinds*（类型）和*resources*（资源）。

## 组和版本

在 Kubernetes 中，*API Group*（API 组）简单地是相关功能的集合。每个组都有一个或多个*versions*（版本），正如其名称所示，允许我们随着时间的推移改变 API 的工作方式。

## 类型和资源

每个 API 组-版本包含一个或多个 API 类型，我们称之为*kinds*（类型）。虽然一个类型在不同版本之间可能会改变形式，但每种形式都必须能够以某种方式存储其他形式的所有数据（我们可以将数据存储在字段中，或者在注释中）。这意味着使用较旧的 API 版本不会导致较新的数据丢失或损坏。有时，同一类型可能由多个资源返回。例如，*pods*（Pod）资源对应于*Pod*类型。然而，有时相同的类型可能由多个资源返回。例如，*Scale*类型由所有规模子资源返回，比如*deployments/scale*或*replicasets/scale*。这就是允许 Kubernetes HorizontalPodAutoscaler 与不同资源交互的原因。然而，对于自定义资源定义（CRD），每种类型将对应于单个资源。

请注意，资源始终以小写形式存在，并且按照惯例是类型的小写形式。

## 这如何对应到 Go 语言？

当我们提到特定组-版本中的一种类型时，我们将其称为*GroupVersionKind*（GVK）。资源也是如此。正如我们将很快看到的那样，每个 GVK 对应于包中的给定根 Go 类型。

现在我们术语明晰了，我们可以*实际地*创建我们的 API！

## 那么，我们如何创建我们的 API？

在接下来的[添加新 API](../cronjob-tutorial/new-api.html)部分中，我们将检查工具如何帮助我们使用命令`kubebuilder create api`创建我们自己的 API。

这个命令的目标是为我们的类型创建自定义资源（CR）和自定义资源定义（CRD）。要进一步了解，请参阅[使用自定义资源定义扩展 Kubernetes API][kubernetes-extend-api]。

## 但是，为什么要创建 API？

新的 API 是我们向 Kubernetes 介绍自定义对象的方式。Go 结构用于生成包括我们数据模式以及跟踪新类型名称等数据的 CRD。然后，我们可以创建我们自定义对象的实例，这些实例将由我们的[controllers][controllers]管理。

我们的 API 和资源代表着我们在集群中的解决方案。基本上，CRD 是我们定制对象的定义，而 CR 是它的一个实例。

## 啊，你有例子吗？

让我们想象一个经典的场景，目标是在 Kubernetes 平台上运行应用程序及其数据库。然后，一个 CRD 可以代表应用程序，另一个可以代表数据库。通过创建一个 CRD 描述应用程序，另一个描述数据库，我们不会伤害封装、单一责任原则和内聚等概念。损害这些概念可能会导致意想不到的副作用，比如扩展、重用或维护方面的困难，仅举几例。

这样，我们可以创建应用程序 CRD，它将拥有自己的控制器，并负责创建包含应用程序的部署以及创建访问它的服务等工作。类似地，我们可以创建一个代表数据库的 CRD，并部署一个负责管理数据库实例的控制器。

## 呃，那个 Scheme 是什么？

我们之前看到的`Scheme`只是一种跟踪给定 GVK 对应的 Go 类型的方式（不要被其[godocs](https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime?tab=doc#Scheme)所压倒）。

例如，假设我们标记`"tutorial.kubebuilder.io/api/v1".CronJob{}`类型属于`batch.tutorial.kubebuilder.io/v1` API 组（隐含地表示它具有类型`CronJob`）。

然后，我们可以根据来自 API 服务器的一些 JSON 构造一个新的`&CronJob{}`，其中说

```json
{
    "kind": "CronJob",
    "apiVersion": "batch.tutorial.kubebuilder.io/v1",
    ...
}
```

或者在我们提交一个`&CronJob{}`进行更新时，正确查找组-版本。

[kubernetes-extend-api]: https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/
[controllers]: ../cronjob-tutorial/controller-overview.md