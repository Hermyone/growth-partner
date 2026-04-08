<template>
	<div class="q-small-title">事务待办</div>
	<div class="personal-notice-bar">
		<el-row>
			<el-col><el-checkbox>有待办事务时,自动打开待办窗口</el-checkbox></el-col>
			<el-col><el-checkbox>有待办事务时,短信推送</el-checkbox></el-col>
			<el-col><el-checkbox>有待办事务时,微信推送</el-checkbox></el-col>
		</el-row>
	</div>
</template>

<script setup lang="ts" name="personal-notice">
	import { reactive, onMounted } from 'vue';

	const state = reactive({

	});

	// 页面加载时
	onMounted(() => {

	});
</script>

<style scoped lang="scss">
	.personal-notice-bar {
		width: 300px;
		padding-left: 20px;
		&-flex {
			display: flex;
			align-items: center;
			margin-bottom: 2px;
			&-label {
				flex: 1;
				color: var(--el-text-color-primary);
			}
		}
	}
</style>

