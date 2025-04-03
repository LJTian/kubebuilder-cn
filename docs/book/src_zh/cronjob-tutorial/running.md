# 运行和部署控制器

### 可选步骤
如果选择对 API 定义进行任何更改，则在继续之前，可以使用以下命令生成清单，如自定义资源（CRs）或自定义资源定义（CRDs）：
```bash
make manifests
```

要测试控制器，请在本地针对集群运行它。
在继续之前，我们需要安装我们的 CRDs，如[快速入门](/quick-start.md)中所述。这将自动使用 controller-tools 更新 YAML 清单（如果需要）：

```bash
make install
```

现在我们已经安装了我们的 CRDs，我们可以针对集群运行控制器。这将使用我们连接到集群的任何凭据，因此我们暂时不需要担心 RBAC。

<aside class="note">

<h1>在本地运行 Webhook</h1>

如果要在本地运行 Webhook，您需要为提供 Webhook 生成证书，并将其放在正确的目录下（默认为`/tmp/k8s-webhook-server/serving-certs/tls.{crt,key}`）。

如果您没有运行本地 API 服务器，您还需要弄清楚如何将流量从远程集群代理到本地 Webhook 服务器。
因此，通常建议在进行本地代码运行测试循环时禁用 Webhook，如下所示。

</aside>

在另一个终端中运行

```bash
export ENABLE_WEBHOOKS=false
make run
```

您应该会看到有关控制器启动的日志，但它目前还不会执行任何操作。

此时，我们需要一个 CronJob 进行测试。让我们编写一个样本到 `config/samples/batch_v1_cronjob.yaml`，然后使用该样本：

```yaml
{{#include ./testdata/project/config/samples/batch_v1_cronjob.yaml}}
```

```bash
kubectl create -f config/samples/batch_v1_cronjob.yaml
```

此时，您应该会看到大量活动。如果观察更改，您应该会看到您的 CronJob 正在运行，并更新状态：

```bash
kubectl get cronjob.batch.tutorial.kubebuilder.io -o yaml
kubectl get job
```

现在我们知道它正在运行，我们可以在集群中运行它。停止 `make run` 命令，并运行

```bash
make docker-build docker-push IMG=<some-registry>/<project-name>:tag
make deploy IMG=<some-registry>/<project-name>:tag
```

<aside class="note">
<h1>注册表权限</h1>

此映像应发布在您指定的个人注册表中。并且需要具有从工作环境中拉取映像的访问权限。
如果上述命令无法正常工作，请确保您对注册表具有适当的权限。

</aside>

如果再次列出 CronJob，就像我们之前所做的那样，我们应该看到控制器再次正常运行！