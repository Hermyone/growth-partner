<template>
	<el-main class="layout-main" :style="themeConfig.isFixedHeader ? `height: calc(100% - ${setMainHeight})` : `minHeight: calc(100% - ${setMainHeight})`">
		<el-scrollbar
			ref="layoutMainScrollbarRef"
			:class="{
				'layout-scrollbar':
					(!isClassicOrTransverse && !state.currentRouteMeta.isLink && !state.currentRouteMeta.isIframe) ||
					(!isClassicOrTransverse && state.currentRouteMeta.isLink && !state.currentRouteMeta.isIframe),
			}"
		>
			<LayoutParentView
				:style="{
					padding: !isClassicOrTransverse || (state.currentRouteMeta.isLink && state.currentRouteMeta.isIframe) ? '0' : '15px',
					transition: 'padding 0.3s ease-in-out',
				}"
			/>
			<FooterView v-if="themeConfig.isFooter" />
		</el-scrollbar>
	</el-main>
</template>

<script setup lang="ts" name="layoutMain">
	import { reactive, getCurrentInstance, watch, onMounted, computed, ref } from 'vue';
	import { useRoute } from 'vue-router';
	import { storeToRefs } from 'pinia';
	import { ElScrollbar } from 'element-plus';
	import { useThemeConfig } from '/@/stores/themeConfig';
	import { NextLoading } from '/@/utils/loading';
	import { useTagsViewRoutes } from '/@/stores/tagsViewRoutes';

	// 引入组件
	import LayoutParentView from '/@/layout/routerView/parent.vue';
	import FooterView from '/@/layout/footer/index.vue';

	// 定义接口来定义对象的类型
	interface MainState {
		headerHeight: string | number;
		currentRouteMeta: any;
	}

	const { proxy } = <any>getCurrentInstance();
	const storesThemeConfig = useThemeConfig();
	const { themeConfig } = storeToRefs(storesThemeConfig);
	const route = useRoute();
	const layoutMainScrollbarRef = ref<InstanceType<typeof ElScrollbar>>();
	const state = reactive<MainState>({
		headerHeight: '',
		currentRouteMeta: {},
	});
	const storesTagsViewRoutes = useTagsViewRoutes();
	const { isTagsViewCurrenFull } = storeToRefs(storesTagsViewRoutes);

	// 判断布局
	const isClassicOrTransverse = computed(() => {
		const { layout } = themeConfig.value;
		return layout === 'classic' || layout === 'transverse';
	});
	// 设置主内容区的高度
	const setMainHeight = computed(() => {
		if (isTagsViewCurrenFull.value) return '0px';
		const { isTagsview, layout } = themeConfig.value;
		if (isTagsview && layout !== 'classic') return '85px';
		else return '51px';
	});
	// 设置 main 的高度
	const initHeaderHeight = () => {
		const bool = state.currentRouteMeta.isLink && state.currentRouteMeta.isIframe;
		let { isTagsview } = themeConfig.value;
		if (isTagsview) return (state.headerHeight = bool ? `86px` : `115px`);
		else return (state.headerHeight = `80px`);
	};
	// 初始化获取当前路由 meta，用于设置 iframes padding
	const initGetMeta = () => {
		state.currentRouteMeta = route.meta;
	};
	// 页面加载前
	onMounted(async () => {
		await initGetMeta();
		initHeaderHeight();
		NextLoading.done(600);
	});
	// 监听路由变化
	watch(
		() => route.path,
		() => {
			state.currentRouteMeta = route.meta;
			const bool = state.currentRouteMeta.isLink && state.currentRouteMeta.isIframe;
			state.headerHeight = bool ? `86px` : `115px`;
			proxy.$refs.layoutMainScrollbarRef.update();
		}
	);
	// 监听 themeConfig 配置文件的变化，更新菜单 el-scrollbar 的高度
	watch(
		themeConfig,
		(val) => {
			state.currentRouteMeta = route.meta;
			const bool = state.currentRouteMeta.isLink && state.currentRouteMeta.isIframe;
			state.headerHeight = val.isTagsview ? (bool ? `86px` : `115px`) : '51px';
			proxy.$refs?.layoutMainScrollbarRef?.update();
		},
		{
			deep: true,
		}
	);

	// 暴露变量
	defineExpose({
		layoutMainScrollbarRef,
	});
</script>
