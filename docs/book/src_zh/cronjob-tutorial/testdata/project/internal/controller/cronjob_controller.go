/*
版权所有 2024 Kubernetes 作者。

根据 Apache 许可证 2.0 版（"许可证"）获得许可；
除非符合许可证的规定，否则您不得使用此文件。
您可以在以下网址获取许可证的副本

    http://www.apache.org/licenses/LICENSE-2.0

除非适用法律要求或书面同意，否则根据许可证分发的软件
将按"原样"分发，不附带任何明示或暗示的担保或条件。
请参阅许可证以了解特定语言下的权限和限制。
*/
// +kubebuilder:docs-gen:collapse=Apache License

/*
我们将从一些导入开始。您将看到我们需要比为我们自动生成的导入更多的导入。
我们将在使用每个导入时讨论它们。
*/

package controller

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/robfig/cron"
	kbatch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ref "k8s.io/client-go/tools/reference"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	batchv1 "tutorial.kubebuilder.io/project/api/v1"
)

/*
接下来，我们需要一个时钟，它将允许我们在测试中模拟时间。
*/

// CronJobReconciler 调和 CronJob 对象
type CronJobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Clock
}

/*
我们将模拟时钟以便在测试中更容易地跳转时间，"真实"时钟只是调用 `time.Now`。
*/
type realClock struct{}

func (_ realClock) Now() time.Time { return time.Now() }

// 时钟知道如何获取当前时间。
// 它可以用于测试中模拟时间。
type Clock interface {
	Now() time.Time
}

// +kubebuilder:docs-gen:collapse=Clock

/*
请注意，我们需要更多的 RBAC 权限 —— 因为我们现在正在创建和管理作业，所以我们需要为这些操作添加权限，
这意味着需要添加一些 [标记](/reference/markers/rbac.md)。
*/

//+kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=cronjobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=cronjobs/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get

/*
现在，我们进入控制器的核心——调和逻辑。
*/
var (
	scheduledTimeAnnotation = "batch.tutorial.kubebuilder.io/scheduled-at"
)

// Reconcile 是主要的 Kubernetes 调和循环的一部分，旨在将集群的当前状态移动到期望的状态。
// TODO（用户）：修改 Reconcile 函数以比较 CronJob 对象指定的状态与实际集群状态，然后执行操作以使集群状态反映用户指定的状态。
//
// 有关更多详细信息，请查看此处的 Reconcile 和其结果：
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *CronJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	/*
		### 1: 通过名称加载 CronJob

		我们将使用我们的客户端获取 CronJob。所有客户端方法都以上下文（以允许取消）作为它们的第一个参数，
		并以对象本身作为它们的最后一个参数。Get 有点特殊，因为它以一个 [`NamespacedName`](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client?tab=doc#ObjectKey)
		作为中间参数（大多数没有中间参数，正如我们将在下面看到的）。

		许多客户端方法还在最后接受可变选项。
	*/
	var cronJob batchv1.CronJob
	if err := r.Get(ctx, req.NamespacedName, &cronJob); err != nil {
		log.Error(err, "无法获取 CronJob")
		// 我们将忽略未找到的错误，因为它们不能通过立即重新排队来修复（我们需要等待新的通知），并且我们可以在删除的请求中得到它们。
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	/*
		### 2: 列出所有活动作业，并更新状态

		为了完全更新我们的状态，我们需要列出此命名空间中属于此 CronJob 的所有子作业。
		类似于 Get，我们可以使用 List 方法列出子作业。请注意，我们使用可变选项设置命名空间和字段匹配（实际上是我们在下面设置的索引查找）。
	*/
	var childJobs kbatch.JobList
	if err := r.List(ctx, &childJobs, client.InNamespace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name}); err != nil {
		log.Error(err, "无法列出子作业")
		return ctrl.Result{}, err
	}

	/*
		<aside class="note">

		<h1>这个索引是什么意思？</h1>

		<p>调解程序获取由 cronjob 拥有的所有作业以获取状态。随着我们的 cronjob 数量的增加，
		查找这些作业可能会变得非常慢，因为我们必须对所有作业进行筛选。为了更高效地查找，
		这些作业将在控制器的名称上进行本地索引。在缓存作业对象上添加了 jobOwnerKey 字段。
		此键引用拥有的控制器，并充当索引。在本文档的后面，我们将配置管理器以实际上索引此字段。</p>

		</aside>

		一旦我们拥有所有我们拥有的作业，我们将它们分为活动、成功和失败的作业，并跟踪最近的运行时间，以便我们可以在状态中记录它。
		请记住，状态应该能够从世界的状态中重建，因此通常不建议从根对象的状态中读取。相反，您应该在每次运行时重新构建它。这就是我们将在这里做的事情。

		我们可以使用状态条件来检查作业是否"完成"，以及它是成功还是失败。我们将把这个逻辑放在一个辅助函数中，使我们的代码更清晰。
	*/

	// 找到活动作业列表
	var activeJobs []*kbatch.Job
	var successfulJobs []*kbatch.Job
	var failedJobs []*kbatch.Job
	var mostRecentTime *time.Time // 找到最近的运行时间，以便我们可以在状态中记录它

	/*
		我们认为作业"完成"，如果它具有标记为 true 的"Complete"或"Failed"条件。
		状态条件允许我们向对象添加可扩展的状态信息，其他人类和控制器可以检查这些信息以检查完成和健康等情况。
	*/
	isJobFinished := func(job *kbatch.Job) (bool, kbatch.JobConditionType) {
		for _, c := range job.Status.Conditions {
			if (c.Type == kbatch.JobComplete || c.Type == kbatch.JobFailed) && c.Status == corev1.ConditionTrue {
				return true, c.Type
			}
		}

		return false, ""
	}
	// +kubebuilder:docs-gen:collapse=isJobFinished

	/*
		我们将使用一个辅助函数从我们在作业创建时添加的注释中提取预定时间。
	*/
	getScheduledTimeForJob := func(job *kbatch.Job) (*time.Time, error) {
		timeRaw := job.Annotations[scheduledTimeAnnotation]
		if len(timeRaw) == 0 {
			return nil, nil
		}

		timeParsed, err := time.Parse(time.RFC3339, timeRaw)
		if err != nil {
			return nil, err
		}
		return &timeParsed, nil
	}
	// +kubebuilder:docs-gen:collapse=getScheduledTimeForJob

	for i, job := range childJobs.Items {
		_, finishedType := isJobFinished(&job)
		switch finishedType {
		case "": // 进行中
			activeJobs = append(activeJobs, &childJobs.Items[i])
		case kbatch.JobFailed:
			failedJobs = append(failedJobs, &childJobs.Items[i])
		case kbatch.JobComplete:
			successfulJobs = append(successfulJobs, &childJobs.Items[i])
		}

		// 我们将在注释中存储启动时间，因此我们将从活动作业中重新构建它。
		scheduledTimeForJob, err := getScheduledTimeForJob(&job)
		if err != nil {
			log.Error(err, "无法解析子作业的计划时间", "job", &job)
			continue
		}
		if scheduledTimeForJob != nil {
			if mostRecentTime == nil || mostRecentTime.Before(*scheduledTimeForJob) {
				mostRecentTime = scheduledTimeForJob
			}
		}
	}

	if mostRecentTime != nil {
		cronJob.Status.LastScheduleTime = &metav1.Time{Time: *mostRecentTime}
	} else {
		cronJob.Status.LastScheduleTime = nil
	}
	cronJob.Status.Active = nil
	for _, activeJob := range activeJobs {
		jobRef, err := ref.GetReference(r.Scheme, activeJob)
		if err != nil {
			log.Error(err, "无法引用活动作业", "job", activeJob)
			continue
		}
		cronJob.Status.Active = append(cronJob.Status.Active, *jobRef)
	}

	/*
		在这里，我们将记录我们观察到的作业数量，以便进行调试。请注意，我们不使用格式字符串，而是使用固定消息，并附加附加信息的键值对。这样可以更容易地过滤和查询日志行。
	*/
	log.V(1).Info("作业数量", "活动作业", len(activeJobs), "成功的作业", len(successfulJobs), "失败的作业", len(failedJobs))

	/*
			使用我们收集的数据，我们将更新我们的 CRD 的状态。
		就像之前一样，我们使用我们的客户端。为了专门更新状态子资源，我们将使用客户端的 `Status` 部分，以及 `Update` 方法。

		状态子资源会忽略对 spec 的更改，因此不太可能与任何其他更新冲突，并且可以具有单独的权限。
	*/
	if err := r.Status().Update(ctx, &cronJob); err != nil {
		log.Error(err, "无法更新 CronJob 状态")
		return ctrl.Result{}, err
	}

	/*
		一旦我们更新了我们的状态，我们可以继续确保世界的状态与我们在规范中想要的状态匹配。

		### 3: 根据历史限制清理旧作业

		首先，我们将尝试清理旧作业，以免留下太多作业。
	*/

	// 注意：删除这些是"尽力而为"的——如果我们在特定的作业上失败，我们不会重新排队只是为了完成删除。
	if cronJob.Spec.FailedJobsHistoryLimit != nil {
		sort.Slice(failedJobs, func(i, j int) bool {
			if failedJobs[i].Status.StartTime == nil {
				return failedJobs[j].Status.StartTime != nil
			}
			return failedJobs[i].Status.StartTime.Before(failedJobs[j].Status.StartTime)
		})
		for i, job := range failedJobs {
			if int32(i) >= int32(len(failedJobs))-*cronJob.Spec.FailedJobsHistoryLimit {
				break
			}
			if err := r.Delete(ctx, job, client.PropagationPolicy(metav1.DeletePropagationBackground)); client.IgnoreNotFound(err) != nil {
				log.Error(err, "无法删除旧的失败作业", "job", job)
			} else {
				log.V(0).Info("已删除旧的失败作业", "job", job)
			}
		}
	}

	if cronJob.Spec.SuccessfulJobsHistoryLimit != nil {
		sort.Slice(successfulJobs, func(i, j int) bool {
			if successfulJobs[i].Status.StartTime == nil {
				return successfulJobs[j].Status.StartTime != nil
			}
			return successfulJobs[i].Status.StartTime.Before(successfulJobs[j].Status.StartTime)
		})
		for i, job := range successfulJobs {
			if int32(i) >= int32(len(successfulJobs))-*cronJob.Spec.SuccessfulJobsHistoryLimit {
				break
			}
			if err := r.Delete(ctx, job, client.PropagationPolicy(metav1.DeletePropagationBackground)); err != nil {
				log.Error(err, "无法删除旧的成功作业", "job", job)
			} else {
				log.V(0).Info("已删除旧的成功作业", "job", job)
			}
		}
	}

	/*
		### 4: 检查我们是否被暂停

		如果此对象被暂停，我们不希望运行任何作业，所以我们将立即停止。
		如果我们正在运行的作业出现问题，我们希望暂停运行以进行调查或对集群进行操作，而不删除对象，这是很有用的。
	*/
	if cronJob.Spec.Suspend != nil && *cronJob.Spec.Suspend {
		log.V(1).Info("CronJob 已暂停，跳过")
		return ctrl.Result{}, nil
	}

	/*
		### 5: 获取下一个预定运行时间

		如果我们没有暂停，我们将需要计算下一个预定运行时间，以及我们是否有一个尚未处理的运行。
	*/

	/*
		我们将使用我们有用的 cron 库来计算下一个预定时间。
		我们将从我们的最后一次运行时间开始计算适当的时间，或者如果我们找不到最后一次运行，则从 CronJob 的创建开始计算。

		如果错过了太多的运行并且我们没有设置任何截止时间，那么我们将中止，以免在控制器重新启动或发生故障时引起问题。

		否则，我们将返回错过的运行（我们将只使用最新的），以及下一个运行，以便我们知道何时再次进行调和。
	*/
	getNextSchedule := func(cronJob *batchv1.CronJob, now time.Time) (lastMissed time.Time, next time.Time, err error) {
		sched, err := cron.ParseStandard(cronJob.Spec.Schedule)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("不可解析的调度 %q：%v", cronJob.Spec.Schedule, err)
		}

		// 为了优化起见，稍微作弊一下，从我们最后观察到的运行时间开始
		// 我们可以在这里重建这个，但是没有什么意义，因为我们刚刚更新了它。
		var earliestTime time.Time
		if cronJob.Status.LastScheduleTime != nil {
			earliestTime = cronJob.Status.LastScheduleTime.Time
		} else {
			earliestTime = cronJob.ObjectMeta.CreationTimestamp.Time
		}
		if cronJob.Spec.StartingDeadlineSeconds != nil {
			// 控制器将不会在此点以下调度任何内容
			schedulingDeadline := now.Add(-time.Second * time.Duration(*cronJob.Spec.StartingDeadlineSeconds))

			if schedulingDeadline.After(earliestTime) {
				earliestTime = schedulingDeadline
			}
		}
		if earliestTime.After(now) {
			return time.Time{}, sched.Next(now), nil
		}

		starts := 0

		// 我们将从最后一次运行时间开始，找到下一个运行时间
		for t := sched.Next(earliestTime); !t.After(now); t = sched.Next(t) {
			starts++
			if starts > 100 {
				return time.Time{}, time.Time{}, fmt.Errorf("错过了太多的运行")
			}
			lastMissed = t
		}

		return lastMissed, sched.Next(now), nil
	}

	lastMissed, nextRun, err := getNextSchedule(&cronJob, r.Now())
	if err != nil {
		log.Error(err, "无法计算下一个运行时间")
		return ctrl.Result{}, err
	}

	/*
		### 6: 创建下一个作业

		最后，我们将创建下一个作业，以便在下一个运行时间触发。
	*/

	// 我们将创建一个新的作业对象，并设置它的所有者引用以确保我们在删除时正确清理。
	newJob := &kbatch.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: cronJob.Name + "-",
			Namespace:    cronJob.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(&cronJob, batchv1.SchemeGroupVersion.WithKind("CronJob")),
			},
			Annotations: map[string]string{
				scheduledTimeAnnotation: nextRun.Format(time.RFC3339),
			},
		},
		Spec: cronJob.Spec.JobTemplate.Spec,
	}

	// 我们将等待我们的作业创建
	if err := r.Create(ctx, newJob); err != nil {
		log.Error(err, "无法创建作业")
		return ctrl.Result{}, err
	}

	log.V(0).Info("已创建新作业", "job", newJob)

	// 我们已经创建了一个新的作业，所以我们将在下一个运行时间重新排队。
	return ctrl.Result{RequeueAfter: nextRun.Sub(r.Now())}, nil
}

/*
现在我们已经实现了 CronJobReconciler 的 Reconcile 方法，我们需要在 manager 中注册它。

我们将在 manager 中注册一个新的控制器，用于管理 CronJob 对象。
*/
func (r *CronJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.CronJob{}).
		Owns(&kbatch.Job{}).
		Complete(r)
}
