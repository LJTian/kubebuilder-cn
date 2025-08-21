document.addEventListener('DOMContentLoaded', () => {
  try {
    const main = document.querySelector('main');
    if (!main) return;

    // 若已插入则不重复插入
    if (document.getElementById('global-banner')) return;

    const banner = document.createElement('div');
    banner.id = 'global-banner';
    banner.setAttribute('role', 'note');
    banner.innerHTML = '欢迎阅读 Kubebuilder 中文文档（试运行），内容仍在完善中。';

    // 插入到正文区域顶部
    main.prepend(banner);
  } catch (err) {
    // 静默失败，避免影响页面
    // eslint-disable-next-line no-console
    console && console.warn && console.warn('global banner inject failed:', err);
  }
});

