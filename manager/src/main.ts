import { createApp } from 'vue';
import App from './app.vue'
import pinia from '/@/stores/index';
import { i18n } from '/@/i18n/index';
import router from './router';
import stores from './stores';
import ElementPlus from 'element-plus';
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import 'element-plus/dist/index.css';
import '/@/theme/index.scss';
import '/@/assets/iconfont/iconfont.css';

import init from '/@/init';

const app = createApp(App)

app.use(pinia)
app.use(router)
app.use(stores)
app.use(ElementPlus, { locale: zhCn })
app.use(i18n)
app.use(init)

// 挂载app
app.mount('#app')

 