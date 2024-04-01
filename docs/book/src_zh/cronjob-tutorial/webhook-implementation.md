# 实现默认值/验证 webhook

如果你想为你的 CRD 实现[准入 webhook](../reference/admission-webhook.md)，你需要做的唯一事情就是实现 `Defaulter` 和（或）`Validator` 接口。

Kubebuilder 会为你处理其余工作，比如

1. 创建 webhook 服务器。
2. 确保服务器已添加到 manager 中。
3. 为你的 webhook 创建处理程序。
4. 在服务器中为每个处理程序注册一个路径。

首先，让我们为我们的 CRD（CronJob）生成 webhook 框架。我们需要运行以下命令，带有 `--defaulting` 和 `--programmatic-validation` 标志（因为我们的测试项目将使用默认值和验证 webhook）：

```bash
kubebuilder create webhook --group batch --version v1 --kind CronJob --defaulting --programmatic-validation
```

这将为你生成 webhook 函数，并在你的 `main.go` 中为你的 webhook 将其注册到 manager 中。

<aside class="note">

<h1>支持旧的集群版本</h1>

与你的 Go webhook 实现一起创建的默认 WebhookConfiguration 清单使用 API 版本 `v1`。如果你的项目意图支持早于 v1.16 的 Kubernetes 集群版本，请设置 `--webhook-version v1beta1`。查看[webhook 参考文档][webhook-reference]获取更多信息。

[webhook-reference]: /reference/webhook-overview.md#supporting-older-cluster-versions

</aside>

{{#literatego ./testdata/project/api/v1/cronjob_webhook.go}}