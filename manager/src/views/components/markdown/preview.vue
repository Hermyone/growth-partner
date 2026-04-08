<template>
	<q-dialog v-model="state.isOpen" width="1024px" :close-on-click-modal="false">
		<template #header>
			预览
		</template>
		<div class="markdown-body">
<!--			<v-md-preview :text="state.content"></v-md-preview>-->
      <p v-html="state.content"></p>
		</div>
	</q-dialog>
</template>

<script setup lang="ts" name="markdown-preview">
	import { reactive, getCurrentInstance, onMounted } from 'vue';

	const { proxy } = <any>getCurrentInstance();
	const state = reactive({
		content: '',
		isOpen: false,
	});

	const openDialog = (content: string) => {
		state.content = content;
		state.isOpen = true;
	}
	const closeDialog = () => {
		state.isOpen = false;
	};

	onMounted(() => {
	});

	defineExpose({
		openDialog,
	})
</script>

<style scoped lang="scss">
	//.markdown-body {
	//	padding: 0 15px;
	//	pre {
	//		background-color: var(--el-bg-color);
	//		color: var(--el-text-color-primary);
	//	}
	//}
</style>