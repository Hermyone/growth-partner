<template>
	<div class="layout-padding layout-link-container">
		<div class="layout-padding-auto layout-padding-view">
			<div class="layout-link-warp">
				<i class="layout-link-icon iconfont icon-xingqiu"></i>
				<div class="layout-link-msg">页面 "{{ $t(currentRouteMeta.title) }}" 已在新窗口中打开</div>
			
				<!-- <el-button class="mt30" round size="default" @click="onClosePage">
					<i class="iconfont icon-lianjie"></i>
					<span>关闭</span>
				</el-button> -->
			</div>
		</div>
	</div>
</template>

<script lang="ts">
import { defineComponent, toRefs, reactive, computed, watch, getCurrentInstance } from 'vue';
import { useRoute, RouteMeta } from 'vue-router';
import { storeToRefs } from 'pinia';
import { useThemeConfig } from '/@/stores/themeConfig';

interface LinkViewRouteMeta extends RouteMeta {
	isLink: string;
	title: string;
}

export default defineComponent({
	name: 'layoutLinkView',
	setup() {
		const { proxy } = <any>getCurrentInstance();
		const storesThemeConfig = useThemeConfig();
		const { themeConfig } = storeToRefs(storesThemeConfig);
		const route = useRoute();
		const state = reactive<LinkViewState>({
			currentRouteMeta: {
				isLink: '',
				title: '',
			},
		});
		// 设置 link 的高度
		const setLinkHeight = computed(() => {
			let { isTagsview } = themeConfig.value;
			if (isTagsview) return `115px`;
			else return `80px`;
		});
		// 关闭
		const onClosePage = () => {
			proxy.mittBus.emit('onCurrentContextmenuClick', Object.assign({}, { contextMenuClickId: 1, ...route }));
		};
		// 监听路由的变化，设置内容
		watch(
			() => route.path,
			() => {
				state.currentRouteMeta = <LinkViewRouteMeta>route.meta;
			},
			{
				immediate: true,
			}
		);
		return {
			setLinkHeight,
			onClosePage,
			...toRefs(state),
		};
	},
});
</script>

<style scoped lang="scss">
.layout-link-container {
	.layout-link-warp {
		margin: auto;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		i.layout-link-icon {
			position: relative;
			font-size: 100px;
			color: var(--el-color-primary);
			&::after {
				content: '';
				position: absolute;
				left: 50px;
				top: 0;
				width: 15px;
				height: 100px;
				background: linear-gradient(
					rgba(255, 255, 255, 0.01),
					rgba(255, 255, 255, 0.01),
					rgba(255, 255, 255, 0.01),
					rgba(255, 255, 255, 0.05),
					rgba(255, 255, 255, 0.05),
					rgba(255, 255, 255, 0.05),
					rgba(235, 255, 255, 0.5),
					rgba(255, 255, 255, 0.05),
					rgba(255, 255, 255, 0.05),
					rgba(255, 255, 255, 0.05),
					rgba(255, 255, 255, 0.01),
					rgba(255, 255, 255, 0.01),
					rgba(255, 255, 255, 0.01)
				);
				transform: rotate(-15deg);
				animation: toRight 5s linear infinite;
			}
		}
		.layout-link-msg {
			font-size: 12px;
			color: var(--next-bg-topBarColor);
			opacity: 0.7;
			margin-top: 15px;
		}
	}
}
</style>