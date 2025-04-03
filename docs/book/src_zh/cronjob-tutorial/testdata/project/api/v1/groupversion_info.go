/*
版权所有 2024 年 Kubernetes 作者。

根据 Apache 许可，版本 2.0 进行许可（"许可"）；
除非遵守许可，否则您不得使用此文件。
您可以在以下网址获取许可的副本：

    http://www.apache.org/licenses/LICENSE-2.0

除非适用法律要求或书面同意，根据许可分发的软件是基于"按原样"的基础分发的，
不附带任何明示或暗示的担保或条件。
请参阅许可以获取特定语言下的权限和限制。
*/
// +kubebuilder:docs-gen:collapse=Apache License

/*
首先，我们有一些 *包级别* 的标记，表示此包中有 Kubernetes 对象，并且此包表示组 `batch.tutorial.kubebuilder.io`。
`object` 生成器利用前者，而 CRD 生成器则利用后者从此包中生成正确的 CRD 元数据。
*/

// Package v1 包含了 batch v1 API 组的 API Schema 定义
// +kubebuilder:object:generate=true
// +groupName=batch.tutorial.kubebuilder.io
package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

/*
然后，我们有一些通常有用的变量，帮助我们设置 Scheme。
由于我们需要在我们的控制器中使用此包中的所有类型，有一个方便的方法将所有类型添加到某个 `Scheme` 中是很有帮助的（也是惯例）。SchemeBuilder 为我们简化了这一过程。
*/

var (
	// GroupVersion 是用于注册这些对象的组版本
	GroupVersion = schema.GroupVersion{Group: "batch.tutorial.kubebuilder.io", Version: "v1"}

	// SchemeBuilder 用于将 go 类型添加到 GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme 将此组版本中的类型添加到给定的 scheme 中。
	AddToScheme = SchemeBuilder.AddToScheme
)
