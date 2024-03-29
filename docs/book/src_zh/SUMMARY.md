# 概要

[介绍](./introduction.md)

[架构](./architecture.md)

[快速入门](./quick-start.md)

[入门指南](./getting-started.md)

---

- [教程：构建 CronJob](cronjob-tutorial/cronjob-tutorial.md)

    - [基本项目包含什么？](./cronjob-tutorial/basic-project.md)
    - [每个旅程都需要一个起点，每个程序都需要一个主函数](./cronjob-tutorial/empty-main.md)
    - [组、版本和类型，哦！](./cronjob-tutorial/gvks.md)
    - [添加新的 API](./cronjob-tutorial/new-api.md)
    - [设计 API](./cronjob-tutorial/api-design.md)

        - [简短的插曲：其他文件中包含什么？](./cronjob-tutorial/other-api-files.md)

    - [控制器包含什么？](./cronjob-tutorial/controller-overview.md)
    - [实现一个控制器](./cronjob-tutorial/controller-implementation.md)

        - [你提到了主函数？](./cronjob-tutorial/main-revisited.md)

    - [实现默认值设置/验证 Webhook](./cronjob-tutorial/webhook-implementation.md)
    - [运行和部署控制器](./cronjob-tutorial/running.md)

        - [部署 cert-manager](./cronjob-tutorial/cert-manager.md)
        - [部署 Webhooks](./cronjob-tutorial/running-webhook.md)

    - [编写测试](./cronjob-tutorial/writing-tests.md)

    - [结语](./cronjob-tutorial/epilogue.md)

- [教程：多版本 API](./multiversion-tutorial/tutorial.md)

    - [改变事物](./multiversion-tutorial/api-changes.md)
    - [中枢、辐射和其他轮子的隐喻](./multiversion-tutorial/conversion-concepts.md)
    - [实现转换](./multiversion-tutorial/conversion.md)

        - [并设置 Webhooks](./multiversion-tutorial/webhooks.md)

    - [部署和测试](./multiversion-tutorial/deployment.md)

- [教程：组件配置](./component-config-tutorial/tutorial.md)

    - [改变事物](./component-config-tutorial/api-changes.md)
    - [定义您的配置](./component-config-tutorial/define-config.md)

    - [使用自定义类型](./component-config-tutorial/custom-type.md)

        - [添加新的配置类型](./component-config-tutorial/config-type.md)
        - [更新主函数](./component-config-tutorial/updating-main.md)
        - [定义您的自定义配置](./component-config-tutorial/define-custom-config.md)

---

- [迁移](./migrations.md)

    - [旧版（v3.0.0 之前）](./migration/legacy.md)
        - [Kubebuilder v1 与 v2](migration/legacy/v1vsv2.md)

            - [迁移指南](./migration/legacy/migration_guide_v1tov2.md)

        - [Kubebuilder v2 与 v3](migration/legacy/v2vsv3.md)

            - [迁移指南](migration/legacy/migration_guide_v2tov3.md)
            - [通过更新文件进行迁移](migration/legacy/manually_migration_guide_v2_v3.md)
    - [从 v3.0.0 开始使用插件](./migration/v3-plugins.md)
        - [go/v3 与 go/v4](migration/v3vsv4.md)

            - [迁移指南](migration/migration_guide_gov3_to_gov4.md)
            - [通过更新文件进行迁移](migration/manually_migration_guide_gov3_to_gov4.md)
    - [单组到多组](./migration/multi-group.md)

- [项目升级助手](./reference/rescaffold.md)

---

- [参考资料](./reference/reference.md)

    - [生成 CRD](./reference/generating-crd.md)
    - [使用终结器](./reference/using-finalizers.md)
    - [良好实践](./reference/good-practices.md)
    - [触发事件](./reference/raising-events.md)
    - [监视资源](./reference/watching-resources.md)
        - [由操作员管理的资源](./reference/watching-resources/operator-managed.md)
        - [外部管理的资源](./reference/watching-resources/externally-managed.md)
    - [Kind 集群](reference/kind.md)
    - [什么是 Webhook？](reference/webhook-overview.md)
        - [准入 Webhook](reference/admission-webhook.md)
        - [核心类型的 Webhook](reference/webhook-for-core-types.md)
    - [配置/代码生成的标记](./reference/markers.md)

        - [CRD 生成](./reference/markers/crd.md)
        - [CRD 验证](./reference/markers/crd-validation.md)
        - [CRD 处理](./reference/markers/crd-processing.md)
        - [Webhook](./reference/markers/webhook.md)
        - [对象/DeepCopy](./reference/markers/object.md)
        - [RBAC](./reference/markers/rbac.md)

    - [controller-gen CLI](./reference/controller-gen.md)
    - [自动补全](./reference/completion.md)
    - [构件](./reference/artifacts.md)
    - [平台支持](./reference/platform.md)

    - [子模块布局](./reference/submodule-layouts.md)
    - [使用外部类型/API](./reference/using_an_external_type.md)

    - [配置 EnvTest](./reference/envtest.md)

    - [指标](./reference/metrics.md)

        - [参考资料](./reference/metrics-reference.md)

    - [Makefile 助手](./reference/makefile-helpers.md)
    - [项目配置](./reference/project-config.md)

---

- [插件][plugins]

    - [可用插件](./plugins/available-plugins.md)
        - [用于脚手架项目](./plugins/to-scaffold-project.md)
            - [go/v2（已弃用）](./plugins/go-v2-plugin.md)
            - [go/v3（已弃用）](./plugins/go-v3-plugin.md)
            - [go/v4（默认初始化脚手架）](./plugins/go-v4-plugin.md)
        - [添加可选功能](./plugins/to-add-optional-features.md)
            - [declarative/v1（已弃用）](./plugins/declarative-v1.md)
            - [grafana/v1-alpha](./plugins/grafana-v1-alpha.md)
            - [deploy-image/v1-alpha](./plugins/deploy-image-plugin-v1-alpha.md)
        - [用于其他工具的扩展](./plugins/to-be-extended.md)
            - [kustomize/v1（已弃用）](./plugins/kustomize-v1.md)
            - [kustomize/v2](./plugins/kustomize-v2.md)
    - [扩展 CLI](./plugins/extending-cli.md)
    - [创建自己的插件](./plugins/creating-plugins.md)
    - [测试您自己的插件](./plugins/testing-plugins.md)
    - [插件版本控制](./plugins/plugins-versioning.md)
    - [创建外部插件](./plugins/external-plugins.md)

---

[常见问题解答](./faq.md)

[附录：TODO 落地页](./TODO.md)

[plugins]: ./plugins/plugins.md