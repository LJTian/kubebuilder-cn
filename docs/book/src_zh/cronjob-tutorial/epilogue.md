到目前为止，我们已经实现了相当全面的CronJob控制器，充分利用了Kubebuilder的大多数功能，并使用envtest为控制器编写了测试。

如果您想了解更多内容，请前往[多版本教程](/multiversion-tutorial/tutorial.md)，了解如何向项目添加新的API版本。

此外，您可以尝试以下步骤：我们将很快在教程中介绍这些内容：

- 为 `kubectl get` 命令添加[额外的打印列](/reference/generating-crd.md#additional-printer-columns)，以改善自定义资源在 `kubectl get` 命令输出中的显示。

希望这次翻译更符合您的需求。