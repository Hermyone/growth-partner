<template>
	<el-config-provider :size="getGlobalComponentSize" :locale="state.i18nLocale">
		<router-view v-show="isLock" />
		<LockScreen v-if="themeConfig.isLockScreen" />
		<Setings ref="setingsRef" v-show="isLock" />
		<CloseFull v-if="!themeConfig.isLockScreen" />
	</el-config-provider>
</template>

<script setup lang="ts" name="app">
	import { defineAsyncComponent, computed, ref, getCurrentInstance, onBeforeMount, onMounted, onUnmounted, nextTick, watch, reactive } from 'vue';
	import { useRoute } from 'vue-router';
	import { storeToRefs } from 'pinia';
	import { useTagsViewRoutes } from '/@/stores/tagsViewRoutes';
	import { useThemeConfig } from '/@/stores/themeConfig';
	import common from '/@/utils/common';
	import { Local, Session } from '/@/utils/storage';
	import setIntroduction from '/@/utils/setIconfont';
	import { loadPlugins } from "/@/plugins/index";

	// 引入组件
	const LockScreen = defineAsyncComponent(() => import('/@/layout/lockScreen/index.vue'));
	const Setings = defineAsyncComponent(() => import('/@/layout/navBars/breadcrumb/setings.vue'));
	const CloseFull = defineAsyncComponent(() => import('/@/layout/navBars/breadcrumb/closeFull.vue'));

	const { proxy } = <any>getCurrentInstance();
	const setingsRef = ref();
	const route = useRoute();
	const stores = useTagsViewRoutes();
	const storesThemeConfig = useThemeConfig();
	const { themeConfig } = storeToRefs(storesThemeConfig);
	const state = reactive({
		i18nLocale: null,
	});
	// 获取布局配置信息
	const getThemeConfig = computed(() => {
		return themeConfig.value;
	});
	// 是否锁屏
	const isLock = () => {
		return getThemeConfig.value.isLockScreen && getThemeConfig.value.lockScreenTime > 1;
	};
	// 获取全局组件大小
	const getGlobalComponentSize = computed(() => {
		return common.globalComponentSize();
	});
	// 布局配置弹窗打开
	const openSetingsDrawer = () => {
		setingsRef.value.openDrawer();
	};
	// 打开AI助手
	const openAssistant = () => {
		
	};
	// 设置初始化，防止刷新时恢复默认
	onBeforeMount(() => {
		// 设置批量第三方 icon 图标
		setIntroduction.cssCdn();
		// 设置批量第三方 js
		setIntroduction.jsCdn();
	});
	// 页面加载时
	onMounted(() => {
		nextTick(() => {
			// 监听布局配置弹窗点击打开
			proxy.mittBus.on('openSetingsDrawer', () => {
				openSetingsDrawer();
			});
			// 监听AI助手弹窗
			proxy.mittBus.on('openAssistant', () => {
				openAssistant();
			});
			// 设置 i18n，App.vue 中的 el-config-provider
			proxy.mittBus.on('getI18nConfig', (locale: string) => {
				(state.i18nLocale as string | null) = locale;
			});
			// 获取缓存中的布局配置
			if (Local.get('themeConfig')) {
				storesThemeConfig.setThemeConfig(Local.get('themeConfig'));
				document.documentElement.style.cssText = Local.get('themeConfigStyle');
			}
			// 获取缓存中的全屏配置
			if (Session.get('isTagsViewCurrenFull')) {
				stores.setCurrenFullscreen(Session.get('isTagsViewCurrenFull'));
			}

			// 加载插件
			loadPlugins();
		});
	});
	// 页面销毁时，关闭监听布局配置/i18n监听
	onUnmounted(() => {
		proxy.mittBus.off('openSetingsDrawer', () => {});
		proxy.mittBus.off('openAssistant', () => {});
		proxy.mittBus.off('getI18nConfig', () => {});
	});
	// 监听路由的变化，设置网站标题
	watch(
		() => route.path,
		() => {
			common.useTitle();
		},
		{
			deep: true,
		}
	);
</script>
