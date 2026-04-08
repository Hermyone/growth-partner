<template>
	<div class="markdown-body">
		<v-md-editor
				v-model="state.content"
				height="calc(100% - 50px)"
				left-toolbar="undo redo clear | h bold italic strikethrough quote | ul ol table hr | link "
				right-toolbar="preview fullscreen"
				:disabled-menus="[1]"
				:include-level="[1, 2, 3, 4, 5, 6]"
				@save="saveArticle"
				@upload-image="handleUploadImage"
		>
		</v-md-editor>
	</div>
</template>

<script setup lang="ts" name="markdown">
	import { reactive } from 'vue';

	const state = reactive({
		content: ""
	});

	const  getContent = () => {
		return state.content;
	}

	const  setContent = (content: string) => {
		state.content = content;
	}

	defineExpose({
		getContent,
		setContent,
	})

</script>

<style scoped lang="scss">
	.markdown-body {
		height: 100%;
		min-height: 300px;

		::-webkit-scrollbar {
			width: 5px;
			height: 5px;
		}
		:deep(.v-md-editor__toolbar) {
			display: flex;
			justify-content: space-between;
			align-items: center;
			height: 32px;
			margin-top: 10px;
			ul {
				padding: 5px;
				li:first-child {
				}
			}
		}
	}
</style>