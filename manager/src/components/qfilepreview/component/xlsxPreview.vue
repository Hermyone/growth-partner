<template>
	<div class="xlsx-preview" ref="backRef">
		<div id="xlsx-preview-body" ref="bodyRef" v-loading="loading"></div>
	</div>
</template>
   
<script lang="ts">
import { defineComponent, reactive, toRefs, onMounted, getCurrentInstance, nextTick } from 'vue';
import { ElMessage } from 'element-plus';
// import { loadPlugins } from "/@/plugins/index";
import LuckyExcel from "luckyexcel";

	export default defineComponent({
		name: 'q-xlsx-preview',
		setup() {
			const { proxy } = getCurrentInstance() as any;
			const state = reactive({
				loading: false,
				isLoadingPlugins: false,	// 插件是否加载成功
			});

			const localPreview = (file: any) => {
				if (!state.isLoadingPlugins) {
					ElMessage.warning('插件没有加载完成，请刷新后重试');
					return;
				}

				state.loading = true;
				LuckyExcel.transformExcelToLucky(file, (exportJson: any) => {
					if (exportJson.sheets === null || !exportJson.sheets.length) {
						ElMessage.warning("无法读取excel文件的内容，当前不支持xls文件！");
						return;
					}

					createSheet(exportJson.sheets, exportJson.info.name, exportJson.info.name.creator);
				});

				setTimeout(() => {
					state.loading = false;
				}, 1000);
			};

			const remotePreview = (url: string) => {
				state.loading = true;
				LuckyExcel.transformExcelToLuckyByUrl(
					url,
					"", (exportJson: any, luckysheetfile: any) => {
						// console.log(exportJson);
						// console.log(luckysheetfile);
						if (exportJson.sheets == null || exportJson.sheets.length == 0) {
							alert("文件读取失败!");
							return;
						}

						createSheet(exportJson.sheets, exportJson.info.name, exportJson.info.name.creator);
					}
				)

				setTimeout(() => {
					state.loading = false;
				}, 1000);
			};

			const createSheet = (data: any, title: string, creator: string) => {
				// 先销毁当前容器
				luckysheet.destroy();

				luckysheet.create({
					container: "xlsx-preview-body",	// 设定DOM容器的id
					showtoolbar: false, 			// 是否显示工具栏
					showinfobar: false,				// 是否显示顶部信息栏
					showstatisticBar: false, 		// 是否显示底部计数栏
					sheetBottomConfig: false, 		// sheet页下方的添加行按钮和回到顶部按钮配置
					allowEdit: false, 				// 是否允许前台编辑
					enableAddRow: false, 			// 是否允许增加行
					enableAddCol: false, 			// 是否允许增加列
					sheetFormulaBar: false, 		// 是否显示公式栏
					enableAddBackTop: false, 		// 返回头部按钮
					lang: "zh",						// 语言
					data: data,						// 表格内容
					title: title, 					// 表格标题
					userInfo: creator,
					showtoolbarConfig: {
						undoRedo: false,
					},
					showsheetbar: false,
				});
			};

			const open = (url: string) => {
				remotePreview(url);
			};

			// 页面加载时
			onMounted(() => {
				// loadPlugins()
				// .then(() => {
				// 	state.isLoadingPlugins = true;
				// 	// eslint-disable-next-line no-console
				// 	console.log('插件加载成功');
				// })
				// .catch(() => {
				// 	ElMessage.warning('插件加载失败，请刷新后重试');
				// });

				state.isLoadingPlugins = true;
			});

			return {
				localPreview,
				remotePreview,
				open,
				...toRefs(state),
			};
		},
	});
</script>

<style scoped lang="scss">
	.xlsx-preview {
		padding: 0!important;
		width: 100%;
		height: 100%;
	}
	#xlsx-preview-body {
		width: 100%;
		height: 100%;
		flex: 1;
		box-shadow: inset 0 0 1px var(--el-border-color-light) !important;
	}
</style>