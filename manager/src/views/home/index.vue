<template>
	<div class="home-container">
    <h1>{{ config.APP_NAME }}</h1>
    <span> ver: {{ config.APP_VER }}</span>
	</div>
</template>

<script lang="ts">
	import {
		toRefs,
		reactive,
		defineComponent,
		onMounted,
		watch,
		nextTick,
		getCurrentInstance,
	} from 'vue';
	import { storeToRefs } from 'pinia';
	import { useThemeConfig } from '/@/stores/themeConfig';
  import config from '/@/config';

export default defineComponent({
	name: 'home',
	setup() {
		const storesThemeConfig = useThemeConfig();
		const { themeConfig } = storeToRefs(storesThemeConfig);
		const state = reactive({
		});

		// 页面加载时
		onMounted(async() => {
		});

		// 监听 vuex 中是否开启深色主题
		watch(
			() => themeConfig.value.isIsDark,
			(isIsDark) => {
				nextTick(() => {
				});
			},
			{
				deep: true,
				immediate: true,
			}
		);
		watch(
				()=>themeConfig.value.isCollapse,
				(isCollapse)=>{
				},
				{
					deep: true,
				}
		);

		return {
      config,
			...toRefs(state),
		};
	},
});
</script>

<style scoped lang="scss">
  .home-container {
    background-image: url("/@/assets/frame.png");
    background-size: contain;
    background-repeat: no-repeat;
    background-position: center;
    background-color: var(--el-bg-color);
    width: 100%;
    height: calc(100vh - 120px);
    min-height: 600px;
    h1 {
      font-size: 36px;
      line-height: 60px;
      padding-top: 100px;
      padding-left: 100px;
      color: var(--el-text-color-disabled);
      opacity: 0.5;
    }
    span {
      padding-left: 100px;
      color: var(--el-text-color-disabled);
      opacity: 0.5;
    }
  }
</style>
