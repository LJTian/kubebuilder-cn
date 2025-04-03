/*
版权所有 2024 年 Kubernetes 作者。

根据 Apache 许可证 2.0 版（以下简称“许可证”）获得许可；
除非符合许可证的规定，否则您不得使用此文件。
您可以在以下网址获取许可证的副本

    http://www.apache.org/licenses/LICENSE-2.0

除非适用法律要求或书面同意，根据许可证分发的软件是基于“按原样”分发的，
没有任何明示或暗示的担保或条件。请参阅许可证以获取有关特定语言管理权限和限制的详细信息。
*/
// +kubebuilder:docs-gen:collapse=Apache License

/*
 */

package v1

/*
 */

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// 注意：json 标记是必需的。您添加的任何新字段都必须具有字段的 json 标记以进行序列化。

// +kubebuilder:docs-gen:collapse=Imports

/*
 首先，让我们看一下我们的规范。正如我们之前讨论过的，规范保存*期望状态*，因此我们控制器的任何“输入”都在这里。

 从根本上讲，CronJob 需要以下几个部分：

 - 一个计划（CronJob 中的 *cron*）
 - 一个要运行的作业的模板（CronJob 中的 *job*）

 我们还希望有一些额外的内容，这些将使我们的用户生活更轻松：

 - 启动作业的可选截止时间（如果错过此截止时间，我们将等到下一个预定的时间）
 - 如果多个作业同时运行，应该怎么办（我们等待吗？停止旧的作业？两者都运行？）
 - 暂停运行 CronJob 的方法，以防出现问题
 - 对旧作业历史记录的限制

 请记住，由于我们从不读取自己的状态，我们需要有其他方法来跟踪作业是否已运行。我们可以使用至少一个旧作业来做到这一点。

 我们将使用几个标记（`// +comment`）来指定额外的元数据。这些将在生成我们的 CRD 清单时由 [controller-tools](https://github.com/kubernetes-sigs/controller-tools) 使用。
正如我们将在稍后看到的，controller-tools 还将使用 GoDoc 来形成字段的描述。
*/

// CronJobSpec 定义了 CronJob 的期望状态
type CronJobSpec struct {
	//+kubebuilder:validation:MinLength=0

	// Cron 格式的计划，请参阅 https://en.wikipedia.org/wiki/Cron。
	Schedule string `json:"schedule"`

	//+kubebuilder:validation:Minimum=0

	// 如果由于任何原因错过预定的时间，则作业启动的可选截止时间（以秒为单位）。错过的作业执行将被视为失败的作业。
	// +optional
	StartingDeadlineSeconds *int64 `json:"startingDeadlineSeconds,omitempty"`

	// 指定如何处理作业的并发执行。
	// 有效值包括：
	// - "Allow"（默认）：允许 CronJob 并发运行；
	// - "Forbid"：禁止并发运行，如果上一次运行尚未完成，则跳过下一次运行；
	// - "Replace"：取消当前正在运行的作业，并用新作业替换它
	// +optional
	ConcurrencyPolicy ConcurrencyPolicy `json:"concurrencyPolicy,omitempty"`

	// 此标志告诉控制器暂停后续执行，它不适用于已经启动的执行。默认为 false。
	// +optional
	Suspend *bool `json:"suspend,omitempty"`

	// 指定执行 CronJob 时将创建的作业。
	JobTemplate batchv1.JobTemplateSpec `json:"jobTemplate"`

	//+kubebuilder:validation:Minimum=0

	// 要保留的成功完成作业的数量。
	// 这是一个指针，用于区分明确的零和未指定的情况。
	// +optional
	SuccessfulJobsHistoryLimit *int32 `json:"successfulJobsHistoryLimit,omitempty"`

	//+kubebuilder:validation:Minimum=0

	// 要保留的失败完成作业的数量。
	// 这是一个指针，用于区分明确的零和未指定的情况。
	// +optional
	FailedJobsHistoryLimit *int32 `json:"failedJobsHistoryLimit,omitempty"`
}

/*
   我们定义了一个自定义类型来保存我们的并发策略。实际上，它在内部只是一个字符串，但该类型提供了额外的文档，并允许我们在类型而不是字段上附加验证，使验证更容易重用。
*/

// ConcurrencyPolicy 描述作业将如何处理。
// 只能指定以下并发答案中的一个。
// 如果没有指定以下策略之一，则默认答案是 AllowConcurrent。
// +kubebuilder:validation:Enum=Allow;Forbid;Replace
type ConcurrencyPolicy string

const (
	// AllowConcurrent 允许 CronJob 并发运行。
	AllowConcurrent ConcurrencyPolicy = "Allow"

	// ForbidConcurrent 禁止并发运行，如果上一个作业尚未完成，则跳过下一个运行。
	ForbidConcurrent ConcurrencyPolicy = "Forbid"

	// ReplaceConcurrent 取消当前正在运行的作业，并用新作业替换它。
	ReplaceConcurrent ConcurrencyPolicy = "Replace"
)

/*
   接下来，让我们设计我们的状态，其中包含观察到的状态。它包含我们希望用户或其他控制器能够轻松获取的任何信息。

   我们将保留一个正在运行的作业列表，以及我们成功运行作业的上次时间。请注意，我们使用 `metav1.Time` 而不是 `time.Time` 来获得稳定的序列化，如上文所述。
*/

// CronJobStatus 定义了 CronJob 的观察状态
type CronJobStatus struct {
	// 插入额外的状态字段 - 定义集群的观察状态
	// 重要提示：在修改此文件后，请运行“make”以重新生成代码

	// 指向当前正在运行的作业的指针列表。
	// +optional
	Active []corev1.ObjectReference `json:"active,omitempty"`

	// 作业最后成功调度的时间。
	// +optional
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`
}

/*
   最后，我们有我们已经讨论过的其余样板。如前所述，除了标记我们想要一个状态子资源，以便表现得像内置的 Kubernetes 类型一样，我们不需要更改这个。
*/

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CronJob 是 cronjobs API 的模式
type CronJob struct {
	/*
	 */
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

func init() {
	SchemeBuilder.Register(&CronJob{}, &CronJobList{})
}

//+kubebuilder:docs-gen:collapse=Root Object Definitions
