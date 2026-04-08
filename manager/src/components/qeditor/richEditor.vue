<!--
 * @Descripttion: 封装富文本编辑器
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-07-10
 * @LastEditors:
 * @LastEditTime:
-->

<template>
	<div class="editor-container">
		<Toolbar :editor="editorRef" :mode="mode" :defaultConfig="toolbarConfig" />
		<Editor
			:mode="mode"
			:defaultConfig="editorConfig"
			v-model="editorVal"
			@onCreated="handleCreated"
			@onChange="handleChange"
		/>
	</div>
</template>

<script lang="ts">
import '@wangeditor/editor/dist/css/style.css';
import { reactive, toRefs, shallowRef, watch, defineComponent, onBeforeUnmount } from 'vue';
import { IDomEditor } from '@wangeditor/editor';
import { Toolbar, Editor } from '@wangeditor/editor-for-vue';

export default defineComponent({
	name: 'q-rich',
	components: {Toolbar, Editor},
	props: {
		// 是否禁用
		disable: {
			type: Boolean,
			default: () => false,
		},
		// 内容框默认 placeholder
		placeholder: {
			type: String,
			default: () => '请输入内容...',
		},
		// 模式，可选 <default|simple>，默认 default
		mode: {
			type: String,
			default: () => 'default',
		},
		// 双向绑定，用于获取 editor.getHtml()
		getHtml: String,
		// 双向绑定，用于获取 editor.getText()
		getText: String,
	},
	setup(props, { emit }) {
		// 定义变量内容
		const editorRef = shallowRef();
		const state = reactive({
			editorConfig: {
				placeholder: props.placeholder,
				showLinkImg: false,
				showLinkImgAlt: false,
				showLinkImgHref: false,
				readOnly: false,
			},
      toolbarConfig: {
        toolbarKeys: [
          'headerSelect',
          'bold',
          'italic',
          'underline',
          'through',
          '|',
          'color',
          'bgColor',
          '|',
          'fontSize',
          'fontFamily',
          'lineHeight',
          '|',
          'bulletedList',
          'numberedList',
          'todo',
          '|',
          'justifyLeft',
          'justifyRight',
          'justifyCenter',
          'justifyJustify',
          '|',
          'emotion',
          'insertLink',
          '|',
          'undo',
          'redo',
          '|',
          'fullScreen'
        ],
        excludeKeys: [
          'bgColor',
          'blockquote',
          'codeBlock',
          'emotion',
          'fontFamily',
          'headerSelect',
          'uploadImage',
        ]
      },
			editorVal: props.getHtml,
		});

		// 编辑器回调函数
		const handleCreated = (editor: IDomEditor) => {
			editorRef.value = editor;
		};

		// 编辑器内容改变时
		const handleChange = (editor: IDomEditor) => {
			// const allMenuKeys = editor.getAllMenuKeys();
			// console.log("888", allMenuKeys);

			emit('update:getHtml', editor.getHtml());
			emit('update:getText', editor.getText());
		};

		// 页面销毁时
		onBeforeUnmount(() => {
			const editor = editorRef.value;
			if (editor == null) return;
			editor.destroy();
		});

		// 监听是否禁用改变
		watch(
			() => props.disable,
			(bool) => {
				const editor = editorRef.value;
				if (editor == null) return;
				bool ? editor.disable() : editor.enable();
			},
			{
				deep: true,
			}
		);

		// 监听双向绑定值改变，用于回显
		watch(
			() => props.getHtml,
			(val) => {
				state.editorVal = val;
			},
			{
				deep: true,
			}
		);
	
		return {
			editorRef,
			handleCreated,
			handleChange,
			onBeforeUnmount,
			...toRefs(state),
		};
	},
});
</script>

<style scoped lang="scss">
	.editor-container {
		height: 100%;
		width: 100%;
		
		// 隐藏视频和表格相关的工具栏按钮
		:deep(.w-e-toolbar) {
			.w-e-menu[data-w-e-type="insertVideo"],
			.w-e-menu[data-w-e-type="insertTable"],
			.w-e-menu[data-w-e-type="table"],
			.w-e-menu[data-w-e-type="video"] {
				display: none !important;
			}
		}
	}
</style>