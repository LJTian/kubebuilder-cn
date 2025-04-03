# 简要说明：其他文件的内容是什么？

如果你浏览了 [`api/v1/`](https://sigs.k8s.io/kubebuilder/docs/book/src/cronjob-tutorial/testdata/project/api/v1) 目录中的其他文件，你可能会注意到除了 `cronjob_types.go` 外还有两个额外的文件：`groupversion_info.go` 和 `zz_generated.deepcopy.go`。

这两个文件都不需要进行编辑（前者保持不变，后者是自动生成的），但了解它们的内容是很有用的。

## `groupversion_info.go`

`groupversion_info.go` 包含有关组版本的常见元数据：

{{#literatego ./testdata/project/api/v1/groupversion_info.go}}

## `zz_generated.deepcopy.go`

`zz_generated.deepcopy.go` 包含了上述 `runtime.Object` 接口的自动生成实现，该接口标记了所有我们的根类型表示的 Kinds。

`runtime.Object` 接口的核心是一个深度复制方法 `DeepCopyObject`。

controller-tools 中的 `object` 生成器还为每个根类型及其所有子类型生成了另外两个方便的方法：`DeepCopy` 和 `DeepCopyInto`。