<template>
	<el-menu
		router
		:default-active="state.defaultActive"
		background-color="transparent"
		:collapse="state.isCollapse"
		:unique-opened="getThemeConfig.isUniqueOpened"
		:collapse-transition="false"
		:class="state.menuTheme"
	>
		<template v-for="val in menuLists">
			<el-sub-menu :index="val.path" v-if="val.children && val.children.length > 0" :key="val.path">
				<template #title>
					<SvgIcon :name="val.meta.icon" :size="getIconSize" />
					<span>{{ $t(val.meta.title) }} </span>
				</template>
				<SubItem :chil="val.children" />
			</el-sub-menu>
			<template v-else>
				<el-menu-item :index="val.path" :key="val.path">
					<SvgIcon :name="val.meta.icon" :size="getIconSize"/>
					<template #title v-if="!val.meta.isLink || (val.meta.isLink && val.meta.isIframe)">
						<span>{{ $t(val.meta.title) }}</span>
					</template>
					<template #title v-else>
						<a :href="val.meta.isLink" target="_blank" rel="opener" class="w100">{{ $t(val.meta.title) }}</a>
					</template>
				</el-menu-item>
			</template>
		</template>
	</el-menu>
</template>

<script setup lang="ts" name="navMenuVertical">
	import { reactive, computed, onMounted, watch, defineAsyncComponent } from 'vue';
	import { useRoute, onBeforeRouteUpdate } from 'vue-router';
	import { storeToRefs } from 'pinia';
	import { useThemeConfig } from '/@/stores/themeConfig';

	// 引入组件
	const SubItem = defineAsyncComponent(() => import('/@/layout/navMenu/subItem.vue'));

	const props = defineProps({
		menuList: {
			type: Array,
			default: () => [],
		},
	});

	const storesThemeConfig = useThemeConfig();
	const { themeConfig } = storeToRefs(storesThemeConfig);
	const route = useRoute();
	const state = reactive({
		// 修复：https://gitee.com/lyt-top/vue-next-admin/issues/I3YX6G
		defaultActive: route.meta.isDynamic ? route.meta.isDynamicPath : route.path,
		isCollapse: false,
		menuTheme: themeConfig.value.layout === 'columns' ? 'columns-menu' : '',
	});
	// 获取父级菜单数据
	const menuLists = computed(() => {
		return <any>props.menuList;
	});
	// 获取布局配置信息
	const getThemeConfig = computed(() => {
		return themeConfig.value;
	});
	// 获取图标大小
	const getIconSize = computed(() => {
		if (getThemeConfig.value.layout === 'classic' || getThemeConfig.value.layout === 'defaults') {
			if (state.isCollapse) return 24;
			else return 16;
		} else {
			return 16;
		}
	});
	// 菜单高亮（详情时，父级高亮）
	const setParentHighlight = (currentRoute: any) => {
		const { path, meta } = currentRoute;
		const pathSplit = meta.isDynamic ? meta.isDynamicPath.split('/') : path.split('/');
		if (pathSplit.length >= 4 && meta.isHide) return pathSplit.splice(0, 3).join('/');
		else return path;
	};
	// 设置菜单的收起/展开
	watch(
		themeConfig.value,
		() => {
			document.body.clientWidth <= 1000 ? (state.isCollapse = false) : (state.isCollapse = themeConfig.value.isCollapse);
		},
		{
			immediate: true,
		}
	);
	// 监听菜单列表,如果列表为空,自动折叠菜单
	// watch(
	// 		menuLists,
	// 	() => {
	// 		if(props.menuList.length <= 0){
	// 			themeConfig.value.isCollapse = true;
	// 		}else{
	// 			themeConfig.value.isCollapse = false;
	// 		}
	// 	},
	// );
	// 页面加载时
	onMounted(() => {
		state.defaultActive = setParentHighlight(route);
		if(props.menuList.length <= 0){
			themeConfig.value.isCollapse = true;
		}
	});
	// 路由更新时
	onBeforeRouteUpdate((to) => {
		// 修复：https://gitee.com/lyt-top/vue-next-admin/issues/I3YX6G
		state.defaultActive = setParentHighlight(to);
		const clientWidth = document.body.clientWidth;
		if (clientWidth < 1000) themeConfig.value.isCollapse = false;
	});

	defineExpose({
		menuLists,
		getThemeConfig,
		getIconSize,
	});
</script>
