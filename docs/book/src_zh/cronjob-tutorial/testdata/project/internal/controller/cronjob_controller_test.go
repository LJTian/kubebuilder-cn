/*
根据 Apache 许可证 2.0 版（"许可证"）许可；
除非符合许可证的规定，否则您不得使用此文件。
您可以在以下网址获取许可证的副本：

    http://www.apache.org/licenses/LICENSE-2.0

除非适用法律要求或经书面同意，根据许可证分发的软件
按"原样"提供，不附带任何担保或条件，无论是明示的还是暗示的。
请查看许可证以了解特定语言下的权限和限制。
*/
// +kubebuilder:docs-gen:collapse=Apache License

/*
理想情况下，对于每个在 `suite_test.go` 中调用的控制器，我们应该有一个 `<kind>_controller_test.go`。
因此，让我们为 CronJob 控制器编写示例测试（`cronjob_controller_test.go`）。
*/

/*
和往常一样，我们从必要的导入项开始。我们还定义了一些实用变量。
*/
package controller

import (
	"context"
	"reflect"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	cronjobv1 "tutorial.kubebuilder.io/project/api/v1"
)

// +kubebuilder:docs-gen:collapse=Imports

/*
编写简单集成测试的第一步是实际创建一个 CronJob 实例，以便对其运行测试。
请注意，要创建 CronJob，您需要创建一个包含您的 CronJob 规范的存根 CronJob 结构。

请注意，当我们创建存根 CronJob 时，CronJob 还需要其所需的下游对象的存根。
如果没有下游的存根 Job 模板规范和下游的 Pod 模板规范，Kubernetes API 将无法创建 CronJob。
*/
var _ = Describe("CronJob controller", func() {

    // 为对象名称和测试超时/持续时间和间隔定义实用常量。
    const (
        CronjobName      = "test-cronjob"
        CronjobNamespace = "default"
        JobName          = "test-job"

        timeout  = time.Second * 10
        duration = time.Second * 10
        interval = time.Millisecond * 250
    )

    Context("当更新 CronJob 状态时", func() {
        It("当创建新的 Job 时，应增加 CronJob 的 Status.Active 计数", func() {
            By("创建一个新的 CronJob")
            ctx := context.Background()
            cronJob := &cronjobv1.CronJob{
                TypeMeta: metav1.TypeMeta{
                    APIVersion: "batch.tutorial.kubebuilder.io/v1",
                    Kind:       "CronJob",
                },
                ObjectMeta: metav1.ObjectMeta{
                    Name:      CronjobName,
                    Namespace: CronjobNamespace,
                },
                Spec: cronjobv1.CronJobSpec{
                    Schedule: "1 * * * *",
                    JobTemplate: batchv1.JobTemplateSpec{
                        Spec: batchv1.JobSpec{
                            // 为简单起见，我们只填写了必填字段。
                            Template: v1.PodTemplateSpec{
                                Spec: v1.PodSpec{
                                    // 为简单起见，我们只填写了必填字段。
                                    Containers: []v1.Container{
                                        {
                                            Name:  "test-container",
                                            Image: "test-image",
                                        },
                                    },
                                    RestartPolicy: v1.RestartPolicyOnFailure,
                                },
                            },
                        },
                    },
                },
            }
            Expect(k8sClient.Create(ctx, cronJob)).Should(Succeed())

            /*
            	创建完这个 CronJob 后，让我们检查 CronJob 的 Spec 字段是否与我们传入的值匹配。
            	请注意，由于 k8s apiserver 在我们之前的 `Create()` 调用后可能尚未完成创建 CronJob，我们将使用 Gomega 的 Eventually() 测试函数，而不是 Expect()，以便让 apiserver 有机会完成创建我们的 CronJob。

            	`Eventually()` 将重复运行作为参数提供的函数，直到
            	(a) 函数的输出与随后的 `Should()` 调用中的预期值匹配，或者
            	(b) 尝试次数 * 间隔时间超过提供的超时值。

            	在下面的示例中，timeout 和 interval 是我们选择的 Go Duration 值。
            */

            cronjobLookupKey := types.NamespacedName{Name: CronjobName, Namespace: CronjobNamespace}
            createdCronjob := &cronjobv1.CronJob{}

            // 我们需要重试获取这个新创建的 CronJob，因为创建可能不会立即发生。
            Eventually(func() bool {
                err := k8sClient.Get(ctx, cronjobLookupKey, createdCronjob)
                return err == nil
            }, timeout, interval).Should(BeTrue())
            // 让我们确保我们的 Schedule 字符串值被正确转换/处理。
            Expect(createdCronjob.Spec.Schedule).Should(Equal("1 * * * *"))
            /*
            	现在我们在测试集群中创建了一个 CronJob，下一步是编写一个测试，实际测试我们的 CronJob 控制器的行为。
            	让我们测试负责更新 CronJob.Status.Active 以包含正在运行的 Job 的 CronJob 控制器逻辑。
            	我们将验证当 CronJob 有一个活动的下游 Job 时，其 CronJob.Status.Active 字段包含对该 Job 的引用。

            	首先，我们应该获取之前创建的测试 CronJob，并验证它当前是否没有任何活动的 Job。
            	我们在这里使用 Gomega 的 `Consistently()` 检查，以确保在一段时间内活动的 Job 计数保持为 0。
            */
            By("检查 CronJob 是否没有活动的 Jobs")
            Consistently(func() (int, error) {
                err := k8sClient.Get(ctx, cronjobLookupKey, createdCronjob)
                if err != nil {
                    return -1, err
                }
                return len(createdCronjob.Status.Active), nil
            }, duration, interval).Should(Equal(0))
            /*
            		接下来，我们实际创建一个属于我们的 CronJob 的存根 Job，以及其下游模板规范。
            		我们将 Job 的状态的 "Active" 计数设置为 2，以模拟 Job 运行两个 Pod，这意味着 Job 正在活动运行。

            		然后，我们获取存根 Job，并将其所有者引用设置为指向我们的测试 CronJob。
            		这确保测试 Job 属于我们的测试 CronJob，并由其跟踪。
            	完成后，我们创建我们的新 Job 实例。
            */
            By("创建一个新的 Job")
            testJob := &batchv1.Job{
                ObjectMeta: metav1.ObjectMeta{
                    Name:      JobName,
                    Namespace: CronjobNamespace,
                },
                Spec: batchv1.JobSpec{
                    Template: v1.PodTemplateSpec{
                        Spec: v1.PodSpec{
                            // 为简单起见，我们只填写了必填字段。
                            Containers: []v1.Container{
                                {
                                    Name:  "test-container",
                                    Image: "test-image",
                                },
                            },
                            RestartPolicy: v1.RestartPolicyOnFailure,
                        },
                    },
                },
                Status: batchv1.JobStatus{
                    Active: 2,
                },
            }

            // 请注意，设置此所有者引用需要您的 CronJob 的 GroupVersionKind。
            kind := reflect.TypeOf(cronjobv1.CronJob{}).Name()
            gvk := cronjobv1.GroupVersion.WithKind(kind)

            controllerRef := metav1.NewControllerRef(createdCronjob, gvk)
            testJob.SetOwnerReferences([]metav1.OwnerReference{*controllerRef})
            Expect(k8sClient.Create(ctx, testJob)).Should(Succeed())
            /*
            		将此 Job 添加到我们的测试 CronJob 应该触发我们控制器的协调逻辑。
            	之后，我们可以编写一个测试，评估我们的控制器是否最终按预期更新我们的 CronJob 的 Status 字段！
            */
            By("检查 CronJob 是否有一个活动的 Job")
            Eventually(func() ([]string, error) {
                err := k8sClient.Get(ctx, cronjobLookupKey, createdCronjob)
                if err != nil {
                    return nil, err
                }

                names := []string{}
                for _, job := range createdCronjob.Status.Active {
                    names = append(names, job.Name)
                }
                return names, nil
            }, timeout, interval).Should(ConsistOf(JobName), "应在状态的活动作业列表中列出我们的活动作业 %s", JobName)
        })
    })

})

/*
编写完所有这些代码后，您可以再次在您的 `controllers/` 目录中运行 `go test ./...` 来运行您的新测试！
*/
