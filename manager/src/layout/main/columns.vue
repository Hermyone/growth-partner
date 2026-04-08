<template>
	<el-container class="layout-container">
		<ColumnsAside />
		<el-container class="layout-columns-warp layout-container-view h100">
			<LayoutAside />
			<el-scrollbar ref="layoutScrollbarRef" :class="{ 'layout-backtop': isFixedHeader }">
				<LayoutHeader v-if="isFixedHeader" />
				<LayoutMain />
			</el-scrollbar>
		</el-container>
	</el-container>
</template>

<script setup lang="ts" name="layoutColumns">
	import { computed, ref } from 'vue';
	import { storeToRefs } from 'pinia';
	import { useThemeConfig } from '/@/stores/themeConfig';

	import LayoutAside from '/@/layout/component/aside.vue';
	import LayoutHeader from '/@/layout/component/header.vue';
	import LayoutMain from '/@/layout/component/main.vue';
	import ColumnsAside from '/@/layout/component/columnsAside.vue';

	const storesThemeConfig = useThemeConfig();
	const { themeConfig } = storeToRefs(storesThemeConfig);
	const isFixedHeader = computed(() => {
		return themeConfig.value.isFixedHeader;
	});
</script>
