# 部署准入 Webhooks

## Kind 集群

建议在 [kind](../reference/kind.md) 集群中开发您的 Webhook，以便快速迭代。
为什么呢？

- 您可以在本地不到 1 分钟内启动一个多节点集群。
- 您可以在几秒钟内将其拆除。
- 您不需要将镜像推送到远程仓库。

## cert-manager

您需要按照 [这里](./cert-manager.md) 的说明安装 cert-manager 捆绑包。

## 构建您的镜像

运行以下命令在本地构建您的镜像。

```bash
make docker-build docker-push IMG=<some-registry>/<project-name>:tag
```

如果您使用的是 kind 集群，您不需要将镜像推送到远程容器注册表。您可以直接将本地镜像加载到指定的 kind 集群中：

```bash
kind load docker-image <your-image-name>:tag --name <your-kind-cluster-name>
```

## 部署 Webhooks

您需要通过 kustomize 启用 Webhook 和 cert manager 配置。
`config/default/kustomization.yaml` 现在应该如下所示：

```yaml
{{#include ./testdata/project/config/default/kustomization.yaml}}
```

而 `config/crd/kustomization.yaml` 现在应该如下所示：

```yaml
{{#include ./testdata/project/config/crd/kustomization.yaml}}
```

现在您可以通过以下命令将其部署到集群中：

```bash
make deploy IMG=<some-registry>/<project-name>:tag
```

等待一段时间，直到 Webhook Pod 启动并证书被提供。通常在 1 分钟内完成。

现在您可以创建一个有效的 CronJob 来测试您的 Webhooks。创建应该成功通过。

```bash
kubectl create -f config/samples/batch_v1_cronjob.yaml
```

您还可以尝试创建一个无效的 CronJob（例如，使用格式不正确的 schedule 字段）。您应该看到创建失败并带有验证错误。

<aside class="note warning">

<h1>引导问题</h1>

如果您为同一集群中的 Pod 部署 Webhook，请注意引导问题，因为 Webhook Pod 的创建请求将被发送到尚未启动的 Webhook Pod 本身。

为使其正常工作，您可以使用 [namespaceSelector]（如果您的 Kubernetes 版本为 1.9+）或使用 [objectSelector]（如果您的 Kubernetes 版本为 1.15+）来跳过自身。

</aside>

[namespaceSelector]: https://github.com/kubernetes/api/blob/kubernetes-1.14.5/admissionregistration/v1beta1/types.go#L189-L233
[objectSelector]: https://github.com/kubernetes/api/blob/kubernetes-1.15.2/admissionregistration/v1beta1/types.go#L262-L274