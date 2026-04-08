/* eslint-disable no-console */
// 封装加载 src 文件方法，加载 cdn 文件
export function asynLoad(src : string, isCss = false) {
    return new Promise(res => {
      let el: any;
      if (isCss) {
        el = document.createElement('link');
        el.rel = 'stylesheet';
        el.href = src;
      } else {
        el = document.createElement('script');
        el.src = src;
      }
      document.documentElement.appendChild(el);
      el.onload = el.onreadystatechange = function() {
        if (
          !this.readyState ||
          this.readyState == 'loaded' ||
          this.readyState == 'complete'
        ) {
          res(true);
        }
        this.onload = this.onreadystatechange = null;
      };
    });
  }
  
  export const loadPlugins = () => {
  }