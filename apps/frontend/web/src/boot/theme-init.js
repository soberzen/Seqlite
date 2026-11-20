(() => {
  let isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

  try {
    const theme = window.localStorage.getItem('theme');

    if (theme === 'light' || theme === 'dark') {
      isDark = theme === 'dark';
    }
  } catch {
    // 本地存储不可用时，继续跟随系统。
  }

  document.documentElement.classList.toggle('dark', isDark);
})();
