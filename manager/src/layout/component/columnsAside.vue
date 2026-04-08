<template>
	<div class="layout-columns-aside">
		<div class="layout-logo">
			<img :src="logoMini" class="layout-logo-medium-img" />
		</div>
		<el-scrollbar>
			<ul @mouseleave="onColumnsAsideMenuMouseleave()">
				<li
					v-for="(v, k) in columnsAsideList"
					:key="k"
					@click="onColumnsAsideMenuClick(v, k)"
					@mouseenter="onColumnsAsideMenuMouseenter(v, k)"
					:ref="
						(el) => {
							if (el) columnsAsideOffsetTopRefs[k] = el;
						}
					"
					:class="getItemClass(k)"
					:title="$t(v.meta.title)"
				>
					<div :class="themeConfig.columnsAsideLayout" v-if="!v.meta.isLink || (v.meta.isLink && v.meta.isIframe)">
						<SvgIcon :name="v.meta.icon" :size="24"/>
						<div class="columns-vertical-title font12">
							{{
								$t(v.meta.title) && $t(v.meta.title).length >= 4
									? $t(v.meta.title).substring(0, themeConfig.columnsAsideLayout === 'columns-vertical' ? 4 : 3)
									: $t(v.meta.title)
							}}
						</div>
					</div>
					<div :class="themeConfig.columnsAsideLayout" v-else>
						<a :href="v.meta.isLink" target="_blank">
							<SvgIcon :name="v.meta.icon" :size="24"/>
							<div class="columns-vertical-title font12">
								{{
									$t(v.meta.title) && $t(v.meta.title).length >= 4
										? $t(v.meta.title).substring(0, themeConfig.columnsAsideLayout === 'columns-vertical' ? 4 : 3)
										: $t(v.meta.title)
								}}
							</div>
						</a>
					</div>
				</li>
				<div ref="columnsAsideActiveRef" :class="themeConfig.columnsAsideStyle"></div>
			</ul>
		</el-scrollbar>
	</div>
</template>

<script lang="ts">
import { reactive, toRefs, ref, onMounted, nextTick, getCurrentInstance, watch, onUnmounted, defineComponent } from 'vue';
import { useRoute, useRouter, onBeforeRouteUpdate, RouteRecordRaw } from 'vue-router';
import { storeToRefs } from 'pinia';
import pinia from '/@/stores/index';
import { useRoutesList } from '/@/stores/routesList';
import { useThemeConfig } from '/@/stores/themeConfig';
import logoMini from '/@/assets/logo-mini.svg';

export default defineComponent({
	name: 'layoutColumnsAside',
	setup() {
		const columnsAsideOffsetTopRefs: any = ref([]);
		const columnsAsideActiveRef = ref();
		const { proxy } = <any>getCurrentInstance();
		const stores = useRoutesList();
		const storesThemeConfig = useThemeConfig();
		const { routesList, isColumnsMenuHover, isColumnsNavHover } = storeToRefs(stores);
		const { themeConfig } = storeToRefs(storesThemeConfig);
		const route = useRoute();
		const router = useRouter();
		const state = reactive<ColumnsAsideState>({
			columnsAsideList: [],
			liIndex: 0,
			liOldIndex: null,
			liHoverIndex: null,
			liOldPath: null,
			difference: 0,
			routeSplit: [],
		});

		const getItemClass = (idx: number) => {
			let c1 = { 'layout-columns-active': state.liIndex === idx, 'layout-columns-hover': state.liHoverIndex === idx };
			let c2 = { 'layout-columns-arrow-active': state.liIndex === idx, 'layout-columns-hover': state.liHoverIndex === idx };
			return themeConfig.value.columnsAsideStyle == 'columns-arrow' ? c2 : c1;
		};

		// 设置菜单高亮位置移动
		const setColumnsAsideMove = (k: number) => {
			state.liIndex = k;
			columnsAsideActiveRef.value.style.top = `${columnsAsideOffsetTopRefs.value[k].offsetTop + state.difference}px`;
			let item = state.columnsAsideList[k];
			if(item) {
				proxy.mittBus.emit('setMenuTitle', item.meta.title);
			}
		};
		// 菜单高亮点击事件
		const onColumnsAsideMenuClick = (v: Object, k: number) => {
			//如果菜单关闭，显示出来
			if(themeConfig.value.isCollapse)
				themeConfig.value.isCollapse = false;
			// setColumnsAsideMove(k);
			let { path, redirect, children } = v as any;
			if (!children || (children.length == 0)) {
				if (redirect) router.push(redirect);
				else router.push(path);
			} else {
				state.liOldPath = path;
				state.liOldIndex = k;
				state.liHoverIndex = k;
				proxy.mittBus.emit('setSendColumnsChildren', setSendChildren(path));
			}			
		};
		// 鼠标移入时，显示当前的子级菜单
		const onColumnsAsideMenuMouseenter = (v: RouteRecordRaw, k: number) => {
			let { path } = v;
			state.liOldPath = path;
			state.liOldIndex = k;
			state.liHoverIndex = k;
			proxy.mittBus.emit('setSendColumnsChildren', setSendChildren(path));
			stores.setColumnsMenuHover(false);
			stores.setColumnsNavHover(true);
			proxy.mittBus.emit('setMenuTitle', v.meta?.title);
		};
		// 鼠标移走时，显示原来的子级菜单
		const onColumnsAsideMenuMouseleave = async () => {
			await stores.setColumnsNavHover(false);
			// 添加延时器，防止拿到的 store.state.routesList 值不是最新的
			setTimeout(() => {
				if (!isColumnsMenuHover && !isColumnsNavHover) {
					proxy.mittBus.emit('restoreDefault');
				}
			}, 100);
		};
		// 设置高亮动态位置
		const onColumnsAsideDown = (k: number) => {
			nextTick(() => {
				setColumnsAsideMove(k);
			});
		};
		// 设置/过滤路由（非静态路由/是否显示在菜单中）
		const setFilterRoutes = () => {
			state.columnsAsideList = filterRoutesFun(routesList.value);
			const resData: any = setSendChildren(route.path);
			if (Object.keys(resData).length <= 0) return false;
			onColumnsAsideDown(resData.item[0].k);
			proxy.mittBus.emit('setSendColumnsChildren', resData);
		};
		// 传送当前子级数据到菜单中
		const setSendChildren = (path: string) => {
			let currentData: any = {};
			state.columnsAsideList.map((v: any, k: number) => {
				if (v.path.startsWith(path) || path.startsWith(v.path)) {
					v['k'] = k;
					currentData['item'] = [{ ...v }];
					currentData['children'] = [{ ...v }];
					if (v.children) currentData['children'] = v.children;
				}
			});
			return currentData;
		};
		// 路由过滤递归函数
		const filterRoutesFun = (arr: Array<string>) => {
			return arr
				.filter((item: any) => !item.meta.isHide)
				.map((item: any) => {
					item = Object.assign({}, item);
					if (item.children) item.children = filterRoutesFun(item.children);
					return item;
				});
		};
		// tagsView 点击时，根据路由查找下标 columnsAsideList，实现左侧菜单高亮
		const setColumnsMenuHighlight = (path: string) => {
			state.routeSplit = path.split('/');
			state.routeSplit.shift();
			let routePath = `/${state.routeSplit[0]}`;
			const idx = path.lastIndexOf('/');
			if (idx > 0) {
				routePath = path.substring(0, idx);
			}
			var currentSplitRoute = state.columnsAsideList.find((v: RouteItem) => v.path === routePath);
			if (!currentSplitRoute) {
				// 用于处理静态路由判断
				const routeFirst = `/${state.routeSplit[0]}`;
				currentSplitRoute = state.columnsAsideList.find((v: RouteItem) => v.path === routeFirst);
				if (!currentSplitRoute) return false;
			}
			// 延迟拿值，防止取不到
			setTimeout(() => {
				onColumnsAsideDown((<any>currentSplitRoute).k);
			}, 0);
		};
		// 监听布局配置信息的变化，动态增加菜单高亮位置移动像素
		watch(
			pinia.state,
			(val) => {
				val.themeConfig.themeConfig.columnsAsideStyle === 'columnsRound' ? (state.difference = 3) : (state.difference = 0);
				if (!val.routesList.isColumnsMenuHover && !val.routesList.isColumnsNavHover) {
					state.liHoverIndex = null;
					proxy.mittBus.emit('setSendColumnsChildren', setSendChildren(route.path));
				} else {
					state.liHoverIndex = state.liOldIndex;
					if (!state.liOldPath) return false;
					proxy.mittBus.emit('setSendColumnsChildren', setSendChildren(state.liOldPath));
				}
			},
			{
				deep: true,
			}
		);
		// 页面加载时
		onMounted(() => {
			setFilterRoutes();
			// 销毁变量，防止鼠标再次移入时，保留了上次的记录
			proxy.mittBus.on('restoreDefault', () => {
				state.liOldIndex = null;
				state.liOldPath = null;
			});
		});
		// 页面卸载时
		onUnmounted(() => {
			proxy.mittBus.off('restoreDefault', () => {});
		});
		// 路由更新时
		onBeforeRouteUpdate((to) => {
			setColumnsMenuHighlight(to.path);
			proxy.mittBus.emit('setSendColumnsChildren', setSendChildren(to.path));
		});

		return {
			themeConfig,
			logoMini,
			getItemClass,
			columnsAsideOffsetTopRefs,
			columnsAsideActiveRef,
			onColumnsAsideDown,
			onColumnsAsideMenuClick,
			onColumnsAsideMenuMouseenter,
			onColumnsAsideMenuMouseleave,
			...toRefs(state),
		};
	},
});
</script>

<style scoped lang="scss">
.layout-columns-aside {
	width: 70px;
	height: 100%;
	border: none;
	background: var(--next-bg-columnsMenuBar);
	.layout-logo {
		width: 100%;
		height: 60px;
		display: flex;
		align-items: center;
		justify-content: center;
		&-medium-img {
			width: 40px;
		}
	}
	ul {
		position: relative;
		li {
			color: var(--next-bg-columnsMenuBarColor);
			width: 100%;
			height: 60px;
			text-align: center;
			display: flex;
			cursor: pointer;
			position: relative;
			z-index: 1;
			.columns-vertical {
				margin: auto;
				.columns-vertical-title {
					padding-top: 1px;
				}
			}
			.columns-horizontal {
				display: flex;
				height: 60px;
				width: 100%;
				align-items: center;
				padding: 0 5px;
				i {
					margin-right: 3px;
				}
				a {
					display: flex;
					.columns-horizontal-title {
						padding-top: 1px;
					}
				}
			}
			a {
				text-decoration: none;
				color: var(--next-bg-columnsMenuBarColor);
			}
		}
		.layout-columns-active {
			color: var(--el-color-primary-light-8) !important;
			transition: 0.3s ease-in-out;
		}
		.layout-columns-arrow-active {
			color: var(--el-color-primary) !important;
		}
		.layout-columns-hover {
			color: var(--el-color-primary);
			a {
				color: var(--el-color-primary);
			}
		}
		.columns-round {
			background: var(--el-color-primary);
			color: var(--el-color-white);
			position: absolute;
			left: 50%;
			top: 2px;
			height: 56px;
			width: 65px;
			transform: translateX(-50%);
			z-index: 0;
			transition: 0.3s ease-in-out;
			border-radius: 5px;
		}
		.columns-card {
			@extend .columns-round;
			top: 0;
			height: 60px;
			width: 100%;
			border-radius: 0;
		}
		.columns-arrow {
			position: absolute;
			left: 85%;
			margin-top: 20px;
			transition: 0.3s ease-in-out;
			z-index: 0;
			border-right: 15px solid var(--el-bg-color);
			border-bottom: 15px solid transparent;
			border-top: 15px solid transparent;
		}
	}
}
</style>
