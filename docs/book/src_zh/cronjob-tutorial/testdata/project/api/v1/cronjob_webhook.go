/*
版权所有 2024 年 Kubernetes 作者。

根据 Apache 许可证 2.0 版进行许可;
除非符合许可证的规定，否则您不得使用此文件。
您可以在以下网址获取许可证副本：

    http://www.apache.org/licenses/LICENSE-2.0

除非适用法律要求或书面同意，否则根据许可证分发的软件
按"原样"分发，没有任何担保或条件，无论是明示的还是暗示的。
请查看许可证以获取特定语言的权限和限制。
*/
// +kubebuilder:docs-gen:collapse=Apache 许可证

package v1

import (
	"github.com/robfig/cron"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	validationutils "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:docs-gen:collapse=Go 导入

/*
接下来，我们为 Webhook 设置一个日志记录器。
*/

var cronjoblog = logf.Log.WithName("cronjob-resource")

/*
然后，我们使用管理器设置 Webhook。
*/

// SetupWebhookWithManager 将设置管理器以管理 Webhook
func (r *CronJob) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

/*
请注意，我们使用 kubebuilder 标记生成 Webhook 清单。
此标记负责生成一个变更 Webhook 清单。

每个标记的含义可以在[这里](/reference/markers/webhook.md)找到。
*/

//+kubebuilder:webhook:path=/mutate-batch-tutorial-kubebuilder-io-v1-cronjob,mutating=true,failurePolicy=fail,groups=batch.tutorial.kubebuilder.io,resources=cronjobs,verbs=create;update,versions=v1,name=mcronjob.kb.io,sideEffects=None,admissionReviewVersions=v1

/*
我们使用 `webhook.Defaulter` 接口为我们的 CRD 设置默认值。
将自动提供一个调用此默认值的 Webhook。

`Default` 方法应该改变接收器，设置默认值。
*/

var _ webhook.Defaulter = &CronJob{}

// Default 实现了 webhook.Defaulter，因此将为该类型注册 Webhook
func (r *CronJob) Default() {
	cronjoblog.Info("默认值", "名称", r.Name)

	if r.Spec.ConcurrencyPolicy == "" {
		r.Spec.ConcurrencyPolicy = AllowConcurrent
	}
	if r.Spec.Suspend == nil {
		r.Spec.Suspend = new(bool)
	}
	if r.Spec.SuccessfulJobsHistoryLimit == nil {
		r.Spec.SuccessfulJobsHistoryLimit = new(int32)
		*r.Spec.SuccessfulJobsHistoryLimit = 3
	}
	if r.Spec.FailedJobsHistoryLimit == nil {
		r.Spec.FailedJobsHistoryLimit = new(int32)
		*r.Spec.FailedJobsHistoryLimit = 1
	}
}

/*
此标记负责生成一个验证 Webhook 清单。
*/

//+kubebuilder:webhook:verbs=create;update;delete,path=/validate-batch-tutorial-kubebuilder-io-v1-cronjob,mutating=false,failurePolicy=fail,groups=batch.tutorial.kubebuilder.io,resources=cronjobs,versions=v1,name=vcronjob.kb.io,sideEffects=None,admissionReviewVersions=v1

/*
我们可以对我们的 CRD 进行超出声明性验证的验证。
通常，声明性验证应该足够了，但有时更复杂的用例需要复杂的验证。

例如，我们将在下面看到，我们使用此功能来验证格式良好的 cron 调度，而不是编写一个长正则表达式。

如果实现了 `webhook.Validator` 接口，将自动提供一个调用验证的 Webhook。

`ValidateCreate`、`ValidateUpdate` 和 `ValidateDelete` 方法预期在创建、更新和删除时验证其接收器。
我们将 `ValidateCreate` 与 `ValidateUpdate` 分开，以允许像使某些字段不可变这样的行为，这样它们只能在创建时设置。
我们还将 `ValidateDelete` 与 `ValidateUpdate` 分开，以允许在删除时进行不同的验证行为。
在这里，我们只为 `ValidateCreate` 和 `ValidateUpdate` 使用相同的共享验证。在 `ValidateDelete` 中不执行任何操作，因为我们不需要在删除时验证任何内容。
*/

var _ webhook.Validator = &CronJob{}

// ValidateCreate 实现了 webhook.Validator，因此将为该类型注册 Webhook
func (r *CronJob) ValidateCreate() (admission.Warnings, error) {
	cronjoblog.Info("验证创建", "名称", r.Name)

	return nil, r.validateCronJob()
}

// ValidateUpdate 实现了 webhook.Validator，因此将为该类型注册 Webhook
func (r *CronJob) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	cronjoblog.Info("验证更新", "名称", r.Name)

	return nil, r.validateCronJob()
}

// ValidateDelete 实现了 webhook.Validator，因此将为该类型注册 Webhook
func (r *CronJob) ValidateDelete() (admission.Warnings, error) {
	cronjoblog.Info("验证删除", "名称", r.Name)

	// TODO（用户）：在对象删除时填充您的验证逻辑。
	return nil, nil
}

/*
我们验证 CronJob 的名称和规范。
*/

func (r *CronJob) validateCronJob() error {
	var allErrs field.ErrorList
	if err := r.validateCronJobName(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := r.validateCronJobSpec(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "batch.tutorial.kubebuilder.io", Kind: "CronJob"},
		r.Name, allErrs)
}

/*
一些字段通过 OpenAPI 模式进行声明性验证。
您可以在[API 设计](api-design.md)部分找到 kubebuilder 验证标记（以`// +kubebuilder:validation`为前缀）。
您可以通过运行`controller-gen crd -w`来找到所有 kubebuilder 支持的用于声明验证的标记，
或者在[这里](/reference/markers/crd-validation.md)找到它们。
*/

func (r *CronJob) validateCronJobSpec() *field.Error {
	// 来自 Kubernetes API 机制的字段助手帮助我们返回结构化良好的验证错误。
	return validateScheduleFormat(
		r.Spec.Schedule,
		field.NewPath("spec").Child("schedule"))
}

/*
   我们需要验证 [cron](https://en.wikipedia.org/wiki/Cron) 调度是否格式良好。
*/

func validateScheduleFormat(schedule string, fldPath *field.Path) *field.Error {
	if _, err := cron.ParseStandard(schedule); err != nil {
		return field.Invalid(fldPath, schedule, err.Error())
	}
	return nil
}

/*
   验证字符串字段的长度可以通过验证模式进行声明性验证。
   但是，`ObjectMeta.Name` 字段是在 apimachinery 仓库的一个共享包中定义的，因此我们无法使用验证模式进行声明性验证。
*/

func (r *CronJob) validateCronJobName() *field.Error {
	if len(r.ObjectMeta.Name) > validationutils.DNS1035LabelMaxLength-11 {
		// 作业名称长度为 63 个字符，与所有 Kubernetes 对象一样（必须适合 DNS 子域）。当创建作业时，cronjob 控制器会在 cronjob 后附加一个 11 个字符的后缀（`-$TIMESTAMP`）。作业名称长度限制为 63 个字符。因此，cronjob 名称长度必须小于等于 63-11=52。如果我们不在这里验证这一点，那么作业创建将在稍后失败。
		return field.Invalid(field.NewPath("metadata").Child("name"), r.Name, "必须不超过 52 个字符")
	}
	return nil
}

// +kubebuilder:docs-gen:collapse=验证对象名称
