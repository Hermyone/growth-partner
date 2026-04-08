<template>
	<div class="docx-preview">
		<div id="docx-preview-body" v-loading="loading"></div>
	</div>
</template>
   
<script lang="ts">
import { toRefs, reactive, defineComponent } from 'vue';
import { renderAsync } from "docx-preview";
import axios from 'axios';

	export default defineComponent({
		name: 'q-docx-preview',
		setup() {
			const state = reactive({
				loading: false,
				docxOptions: {
					className: "kaimo-docx", 			// string：默认和文档样式类的类名/前缀
					inWrapper:  true, 					// boolean：启用围绕文档内容的包装器渲染
					ignoreWidth: false, 				// boolean：禁用页面的渲染宽度
					ignoreHeight: false, 				// boolean：禁止渲染页面高度
					ignoreFonts: false, 				// boolean：禁用字体渲染
					breakPages: true, 					// boolean：在分页符上启用分页
					ignoreLastRenderedPageBreak: true, 	// boolean：在 lastRenderedPageBreak 元素上禁用分页
					experimental: false, 				// boolean：启用实验功能（制表符停止计算）
					trimXmlDeclaration: true, 			// boolean：如果为true，解析前会从​​ xml 文档中移除 xml 声明
					useBase64URL: false, 				// boolean：如果为true，图片、字体等会转为base 64 URL，否则使用URL.createObjectURL
					useMathMLPolyfill: false, 			// boolean：包括用于 chrome、edge 等的 MathML polyfill。
					showChanges: false, 				// boolean：启用文档更改的实验性渲染（插入/删除）
					debug: false, 						// boolean：启用额外的日志记录
				}
			});

			const localPreview = (buffer: any) => {
				state.loading = true;
				docxRender(buffer);
				setTimeout(() => {
					state.loading = false;
				}, 1000);
			};

			const remotePreview = (url: string) => {
				state.loading = true;
				axios({
					method: "get",
					responseType: "blob",
					url: url,
				}).then((response) => {
					docxRender(response.data);
				});
				setTimeout(() => {
					state.loading = false;
				}, 2000);
			};

			// 渲染docx
			const docxRender = (buffer: any) => {
				let bodyContainer : any = document.getElementById("docx-preview-body");
				renderAsync(
					buffer, 					// Blob | ArrayBuffer | Uint8Array, 可以是 JSZip.loadAsync 支持的任何类型
					bodyContainer, 				// HTMLElement 渲染文档内容的元素,
					undefined,		// HTMLElement, 用于呈现文档样式、数字、字体的元素。如果为 null，则将使用 bodyContainer。
					state.docxOptions 			// 配置
				).then(res => {
					// eslint-disable-next-line no-console
					console.log('docx 文件加载成功');
					state.loading = false;
				})
			};

			const open = (url: string) => {
				remotePreview(url);
			};

			return {
				localPreview,
				remotePreview,
				docxRender,
				open,
				...toRefs(state),
			};
		},
	});
</script>

<style scoped lang="scss">
	.docx-preview {
		:deep(.docx-wrapper) {
			  background-color: #fff;
			  padding: 0!important;
		  }
		:deep(.docx-wrapper > section.docx) {
			width: 100% !important;
			padding: 0!important;
			min-height: auto !important;
			box-shadow: none;
		}
	}
</style>