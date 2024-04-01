# 部署 cert-manager

我们建议使用 [cert-manager](https://github.com/cert-manager/cert-manager) 为 Webhook 服务器提供证书。只要它们将证书放在所需的位置，其他解决方案也应该可以正常工作。

您可以按照 [cert-manager 文档](https://cert-manager.io/docs/installation/) 进行安装。

cert-manager 还有一个名为 [CA 注入器](https://cert-manager.io/docs/concepts/ca-injector/) 的组件，负责将 CA bundle 注入到 [`MutatingWebhookConfiguration`](https://pkg.go.dev/k8s.io/api/admissionregistration/v1#MutatingWebhookConfiguration) / [`ValidatingWebhookConfiguration`](https://pkg.go.dev/k8s.io/api/admissionregistration/v1#ValidatingWebhookConfiguration) 中。

为了实现这一点，您需要在 [`MutatingWebhookConfiguration`](https://pkg.go.dev/k8s.io/api/admissionregistration/v1#MutatingWebhookConfiguration) / [`ValidatingWebhookConfiguration`](https://pkg.go.dev/k8s.io/api/admissionregistration/v1#ValidatingWebhookConfiguration) 对象中使用一个带有键 `cert-manager.io/inject-ca-from` 的注释。注释的值应该指向一个现有的 [证书请求实例](https://cert-manager.io/docs/concepts/certificaterequest/)，格式为 `<证书命名空间>/<证书名称>`。

这是我们用于给 [`MutatingWebhookConfiguration`](https://pkg.go.dev/k8s.io/api/admissionregistration/v1#MutatingWebhookConfiguration) / [`ValidatingWebhookConfiguration`](https://pkg.go.dev/k8s.io/api/admissionregistration/v1#ValidatingWebhookConfiguration) 对象添加注释的 [kustomize](https://github.com/kubernetes-sigs/kustomize) 补丁：

```yaml
{{#include ./testdata/project/config/default/webhookcainjection_patch.yaml}}
```