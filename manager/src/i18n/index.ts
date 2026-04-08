import { createI18n } from 'vue-i18n';
import zhcnLocale from 'element-plus/es/locale/lang/zh-cn';
import enLocale from 'element-plus/es/locale/lang/en';
import zhtwLocale from 'element-plus/es/locale/lang/zh-tw';

import nextZhcn from '/@/i18n/lang/zh-cn';
import nextEn from '/@/i18n/lang/en';
import nextZhtw from '/@/i18n/lang/zh-tw';

import pagesLoginZhcn from '/@/i18n/pages/login/zh-cn';
import pagesLoginEn from '/@/i18n/pages/login/en';
import pagesLoginZhtw from '/@/i18n/pages/login/zh-tw';

import pinia from '/@/stores/index';
import { storeToRefs } from 'pinia';
import { useThemeConfig } from '/@/stores/themeConfig';


// 定义语言国际化内容
/**
 * 说明：
 * /src/i18n/lang 下的 ts 为框架的国际化内容
 * /src/i18n/pages 下的 ts 为各界面的国际化内容
 */
 const messages = {
	[zhcnLocale.name]: {
		...zhcnLocale,
		message: {
			...nextZhcn,
			...pagesLoginZhcn,
		},
	},
	[enLocale.name]: {
		...enLocale,
		message: {
			...nextEn,
			...pagesLoginEn,
		},
	},
	[zhtwLocale.name]: {
		...zhtwLocale,
		message: {
			...nextZhtw,
			...pagesLoginZhtw,
		},
	},
};

// 读取 pinia 默认语言
const stores = useThemeConfig(pinia);
const { themeConfig } = storeToRefs(stores);

// 导出语言国际化
export const i18n = createI18n({
    silentTranslationWarn: true,
	missingWarn: false,
	silentFallbackWarn: true,
	fallbackWarn: false,
	locale: themeConfig.value.globalI18n,
	fallbackLocale: zhcnLocale.name,
	legacy: false,
	messages,
});

