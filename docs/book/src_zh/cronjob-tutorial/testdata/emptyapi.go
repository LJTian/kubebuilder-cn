/*
版权所有 2022。

根据 Apache 许可证 2.0 版（"许可证"）获得许可；
除非符合许可证的规定，否则您不得使用此文件。
您可以在下面的网址获取许可证的副本

    http://www.apache.org/licenses/LICENSE-2.0

除非适用法律要求或书面同意，否则根据许可证分发的软件
以"原样"为基础分发，没有任何明示或暗示的保证或条件。
请参阅特定语言管理权限和限制的许可证。
*/
// +kubebuilder:docs-gen:collapse=Apache License

/*
我们从简单的开始：我们导入 `meta/v1` API 组，它通常不是单独公开的，而是包含所有 Kubernetes Kind 的公共元数据。
*/
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
接下来，我们为我们的 Kind 的 Spec 和 Status 定义类型。Kubernetes 通过协调期望的状态（`Spec`）与实际的集群状态（其他对象的`Status`）和外部状态，然后记录它观察到的内容（`Status`）来运行。因此，每个 *functional* 对象都包括 spec 和 status。一些类型，比如 `ConfigMap` 不遵循这种模式，因为它们不编码期望的状态，但大多数类型都是这样。
*/
// 编辑此文件！这是你拥有的脚手架！
// 注意：json 标记是必需的。您添加的任何新字段必须具有 json 标记，以便对字段进行序列化。

// CronJobSpec 定义了 CronJob 的期望状态
type CronJobSpec struct {
	// 插入其他的 Spec 字段 - 集群的期望状态
	// 重要提示：在修改此文件后运行 "make" 以重新生成代码
}

// CronJobStatus 定义了 CronJob 的观察状态
type CronJobStatus struct {
	// 插入其他的状态字段 - 定义集群的观察状态
	// 重要提示：在修改此文件后运行 "make" 以重新生成代码
}

/*
接下来，我们定义与实际 Kinds 对应的类型，`CronJob` 和 `CronJobList`。
`CronJob` 是我们的根类型，描述了 `CronJob` 类型。与所有 Kubernetes 对象一样，它包含 `TypeMeta`（描述 API 版本和 Kind），
还包含 `ObjectMeta`，其中包含名称、命名空间和标签等信息。

`CronJobList` 简单地是多个 `CronJob` 的容器。它是用于批量操作（如 LIST）的 Kind。

一般情况下，我们从不修改它们中的任何一个 -- 所有的修改都在 Spec 或 Status 中进行。

这个小小的 `+kubebuilder:object:root` 注释称为标记。我们稍后会看到更多这样的标记，但要知道它们作为额外的元数据，告诉 [controller-tools](https://github.com/kubernetes-sigs/controller-tools)（我们的代码和 YAML 生成器）额外的信息。
这个特定的标记告诉 `object` 生成器，这个类型表示一个 Kind。然后，`object` 生成器为我们生成了 [runtime.Object](https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime?tab=doc#Object) 接口的实现，这是所有表示 Kinds 的类型必须实现的标准接口。
*/

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CronJob 是 cronjobs API 的架构
type CronJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CronJobSpec   `json:"spec,omitempty"`
	Status CronJobStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CronJobList 包含 CronJob 的列表
type CronJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CronJob `json:"items"`
}

/*
最后，我们将 Go 类型添加到 API 组中。这使我们可以将此 API 组中的类型添加到任何 [Scheme](https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime?tab=doc#Scheme) 中。
*/
func init() {
	SchemeBuilder.Register(&CronJob{}, &CronJobList{})
}
