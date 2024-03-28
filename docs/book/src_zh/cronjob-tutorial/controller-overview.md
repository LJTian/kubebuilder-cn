# 控制器的内容是什么？

控制器是 Kubernetes 和任何操作者的核心。

控制器的工作是确保对于任何给定的对象，世界的实际状态（集群状态，以及可能是 Kubelet 的运行容器或云提供商的负载均衡器等外部状态）与对象中的期望状态相匹配。每个控制器专注于一个*根* Kind，但可能会与其他 Kinds 交互。

我们称这个过程为*调和*。

在 controller-runtime 中，实现特定 Kind 的调和逻辑称为[*Reconciler*](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/reconcile?tab=doc)。调和器接受一个对象的名称，并返回我们是否需要再次尝试（例如在出现错误或周期性控制器（如 HorizontalPodAutoscaler）的情况下）。

{{#literatego ./testdata/emptycontroller.go}}

现在我们已经看到了调和器的基本结构，让我们填写 `CronJob` 的逻辑。