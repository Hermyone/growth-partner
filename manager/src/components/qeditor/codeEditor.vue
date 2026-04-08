<!--
 * @Descripttion: 封装代码编辑器
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-07-10
 * @LastEditors:
 * @LastEditTime:
-->

<template>
	<div class="editor-container">
		<Codemirror class="editor-container-code"
				v-model="editor"
				:autofocus="true"
				:indent-with-tab="true"
				:extensions="extensions"
				:options="editOptions"
				@ready="handleReady"
				@copy="copyCode"
		/>
	</div>
</template>

<script lang="ts">
	import { computed, defineComponent, onMounted, reactive, ref, toRefs, nextTick, watch, shallowRef } from 'vue';
	import { Codemirror } from 'vue-codemirror';
	import { javascript } from '@codemirror/lang-javascript';
	import { json } from '@codemirror/lang-json';
	import { sql } from '@codemirror/lang-sql';
	import { oneDark } from '@codemirror/theme-one-dark';
	import ClipboardJS from "clipboard";
	import { ElMessage } from "element-plus";

	export default defineComponent({
		name: 'q-code',
		components: { Codemirror },
		props: {
			mode: {
				type: String,
				default: "javascript",
			},
			value: {
				type: String,
				default: "",
			},
			isDark: {
				type: Boolean,
				default: false,
			},
			readOnly: {
				type: Boolean,
				default: true,
			},
			options: {
				type: Object,
				default: () => {}
			},
		},

		setup(props, { emit }) {
			const editor = ref();
			const view = shallowRef();

			const state = reactive({
				extensions: [] as any,
				editOptions: {
					theme: "eclipse",			//主题样式
					styleActiveLine: true,		//高亮当前行
					lineNumbers: true,			//行号
					lineWrapping: false,		//自动换行
					tabSize: 4,					//Tab缩进
					indentUnit: 4,				//缩进单位
					indentWithTabs : true,		//自动缩进
					matchBrackets: true,		//括号匹配
					mode : props.mode,			//语言
					readOnly: props.readOnly,	//只读
					foldGutter: true,
					lint: true,
					gutters: [
						"CodeMirror-linenumbers",
						"CodeMirror-foldgutter",
						"CodeMirror-lint-markers",
					],
					...props.options,
				},
			});

			const handleReady = (payload: any) => {
				editor.value = props.value;
				view.value = payload.view;
			};

			const formatStr2Json = (v: string) => {
				return JSON.stringify(JSON.parse(v), null, '\t');
			};

			const init = () => {
				state.extensions = [];
				if(props.mode === "javascript"){
					state.extensions.push(javascript())
				}else if(props.mode === "json"){
					state.extensions.push(json())
				}else if(props.mode === "sql"){
					state.extensions.push(sql())
				};
				if(props.isDark)
					state.extensions.push(oneDark);
			};

			watch(
					()=>props.mode,
					()=>{
						init();
					}
			);

			watch(
					() => props.value,
					() => {
						if(props.mode === 'json'){
							editor.value = formatStr2Json(props.value);

						}else {
							editor.value = props.value;
						}
					}
			)

			const copyCode = () => {
				const clipboard = new ClipboardJS("#copyCode");
				clipboard.on("success", (e) => {
					ElMessage.success('代码已复制到剪切板，可粘贴');
					clipboard.destroy();
				});
				clipboard.on("error", (e) => {
					ElMessage.error("代码复制失败");
					clipboard.destroy();
				});
			};

			const saveCode = () => {
				emit("input", props.value);
			};

			// 页面加载时
			onMounted(() => {
				nextTick(() => {
					init();
				});
			});

			return {
				editor,
				copyCode,
				saveCode,
				handleReady,
				formatStr2Json,
				...toRefs(state),
			};
		},
	});
</script>

<style scoped lang="scss">
	.editor-container {
		width: calc(100% - 30px);
		height: 100%;
		margin: 0 15px;
		font-family: Helvetica Neue, Helvetica, PingFang SC, Hiragino Sans GB, Microsoft YaHei, sans-serif;
		font-weight: normal;
		font-size: 18px;
		line-height: 20px;
		letter-spacing: 0px;
		border: 1px solid var(--el-border-color);
		border-radius: var(--el-border-radius-base);
		&-code {
			height: 100%;
			width: 100%;
			:deep(.CodeMirror) {
			}
			:deep(.cm-gutter) {
				background-color: var(--el-bg-color) !important;
			}
		}
	}
</style>