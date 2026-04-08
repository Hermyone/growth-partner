<!--
 * @Descripttion: 自定义表格组件
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-06-05
 * @LastEditors:
 * @LastEditTime:
-->

<template>
	<div class="table-container">
		<el-table
				ref="tableRef"
				:row-key="rowKey"
				:border="config.border"
				:size="config.size"
				:stripe="config.stripe"
				:highlight-current-row="highlightCurrentRow"
				:row-style="rowStyle"
				:cell-style="cellStyle"
				:show-header="showHeader"
				style="width: 100%"
				v-bind="$attrs"
				v-loading="config.loading"
				@selection-change="onSelectionChange"
				@row-click="rowClick"
				@select="select"
				@select-all="selectAll"
				@expand-change="expandChange"
		>
			<el-table-column type="selection" width="40" v-if="isSelection" fixed="left" />
			<el-table-column type="index" label="序号" width="60" v-if="isSerialNo" fixed="left" >
				<template #default="scope">
					<span>{{scope.$index+(page.pageNum - 1) * page.pageSize + 1}}</span>
				</template>
			</el-table-column>
			<slot></slot>

			<el-table-column min-width="1" v-if="autoColumnWidth"></el-table-column>
			<template #empty>
				<el-empty description="暂无数据" />
			</template>
		</el-table>
		<div class="table-footer mt15" v-show="footerVisible">
			<el-pagination
					v-model:current-page="page.pageNum"
					v-model:page-size="page.pageSize"
					v-show="paginationVisible"
					:pager-count="5"
					:page-sizes="config.pageSizes"
					:total="config.total"
					:layout="layout"
					:size="paginationSmall ? 'small' : 'default'"
					background
			>
			</el-pagination>
			<div class="table-footer-tool">
				<el-button v-show="footerButton && isSelection && (exportColumns.length > 0)" @click="onExportTable" icon="ele-Download" circle class="table-footer-tool-button"></el-button>
				<el-button v-show="footerButton" @click="onRefreshTable" icon="ele-Refresh" circle class="table-footer-tool-button"></el-button>
				<el-popover placement="top" title="设置" :width="350" trigger="click" :hide-after="0" >
					<template #reference>
						<el-button v-show="footerButton" icon="ele-Setting" circle class="table-footer-tool-button"></el-button>
					</template>
					<el-form label-width="80px" label-position="left" size="small" style="margin-left: 20px;">
						<el-form-item label="表格尺寸">
							<el-radio-group v-model="config.size" size="small" @change="dataChanged">
								<el-radio-button value="large">大</el-radio-button>
								<el-radio-button value="default">正常</el-radio-button>
								<el-radio-button value="small">小</el-radio-button>
							</el-radio-group>
						</el-form-item>
						<el-form-item label="表格样式" style="margin-top: -5px;">
							<el-checkbox v-model="config.border" label="纵向边框"></el-checkbox>
							<el-checkbox v-model="config.stripe" label="斑马纹"></el-checkbox>
						</el-form-item>
					</el-form>
				</el-popover>
			</div>
		</div>
	</div>
</template>

<script lang="ts">
	import { defineComponent, reactive, nextTick, onMounted, toRefs, watch, getCurrentInstance } from 'vue';
	import { ElMessage } from 'element-plus';
	import table2excel from 'js-table2excel';
	import '/@/theme/tableTool.scss';
	import { useRenderTableTree } from '/@/components/qtable/treetable';
	import { useDebounceFn } from '@vueuse/core';

	// 定义父组件传过来的值
	const props = {
		// 行主键 KEY
		rowKey: { type: String, default: "key" },
		// 每页显示的记录数
		pageSize: { type: Number, default: 20 },
		// 可选择的每页显示数量
		pageSizes: { type: Array, default: () => [10, 20, 30, 40, 50] },
		// 分页布局
		layout: { type: String, default: 'total, sizes, prev, pager, next, jumper' },
		// 表格尺寸
		size: { type: String, default: "default" },
		// 是否显示纵向边框
		border: { type: Boolean, default: false },
		// 是否显示斑马纹
		stripe: { type: Boolean, default: false },
		// 是否显示footer
		footerVisible: { type: Boolean, default: true },
        // 是否显示footer中的操作按钮
        footerButton: { type: Boolean, default: true },
		// 表格导出时标题，默认的文件名
		globalTitle: { type: String, default: "table" },
		// 是否显示序号
		isSerialNo: { type: Boolean, default: false },
		// 是否可以勾选
		isSelection: { type: Boolean, default: true },
		// 记录总数
		recordCount: { type: Number, default: 0 },
		// 加载状态
		loading: { type: Boolean, default: false },
		// 需要导出的列
		exportColumns: { type: Array, default: () => [] },
		// 自动列宽
		autoColumnWidth: { type: Boolean, default: true },
		// 是否使用小型分页样式
		paginationSmall: { type: Boolean, default: true },
		// 是否显示分页
		paginationVisible: { type: Boolean, default: true },
		// 行单击事件回调函数
		rowClick: { type: Function, default: null },
		// 高亮当前行
		highlightCurrentRow: { type: Boolean, default: false },
		// 选中/全选回调事件
		select: { type: Function, default: null},
		selectAll: { type: Function, default: null},

		// 是否显示表头
		showHeader: { type: Boolean, default: true },

		// 自定义行列样式回调
		rowStyle: { type: Function, default: null },
		cellStyle: { type: Function, default: null },
	};

	export default defineComponent({
		name: 'q-table',
		props: props,
		setup(props, { emit }) {
			// 定义变量内容
			const { proxy } = <any>getCurrentInstance();
			const { renderTableTree } = useRenderTableTree();
			const state = reactive({
				page: {
					pageNum: 1,
					pageSize: props.pageSize,
				},
				selectlist: [] as any,
				config: {
					total: 0,
					loading: false,
					size: props.size,
					pageSizes: props.pageSizes,
					stripe: props.stripe,
					border: props.border,
				},
			});

			// 搜索时，分页还原成默认
			const pageReset = () => {
				state.page.pageNum = 1;
				state.page.pageSize = props.pageSize;
				updateData();
			};
			// 导出
			const onExportTable = () => {
				if (state.selectlist.length <= 0) return ElMessage.warning('请先选择要导出的数据');
				table2excel(props.exportColumns, state.selectlist, `${props.globalTitle}${new Date().toLocaleString()}`);
			};
			// 刷新
			const onRefreshTable = () => {
				updateData();
			};

			// 表格多选改变时
			const onSelectionChange = (selection: any) => {
				state.selectlist = selection;
				emit('selectData', state.selectlist);
			};

			// 反回选中的数据
			const getSelectlist = (keyField: string) => {
				let res = [] as any[];
				state.selectlist.forEach((item : any) => res.push(item[keyField]));
				return res;
			};

			// 更新数据
			const updateData = () => {
				state.selectlist = [];
				state.config.loading = true;
				emit('pageChange', state.page);

				setTimeout(() => {
					state.config.loading = false;
					emit('queryCompleted');
					dataChanged();
				}, 500);
			};

			const dataChanged = useDebounceFn(function(resp = null) {
				renderTableTree(proxy.$refs.tableRef.$el)
				emit('dataChanged', resp)
			}, 100)

			// 当前表格
			const table = () => {
				return proxy.$refs.tableRef;
			};

			const expandChange = () => {
				renderTableTree(proxy.$refs.tableRef.$el)
			}

			// 页面加载时
			onMounted(() => {
				nextTick(() => {
					updateData();
				});
			});

			// 监听双向绑定 recordCount 的变化
			watch(
					() => props.recordCount,
					() => {
						state.config.total = props.recordCount;
					}
			);
			watch(
					() => props.loading,
					() => {
						state.config.loading = props.loading;
					}
			);
			watch(
					() => state.page,
					() => {
						updateData()
					},{
						deep: true
					}
			);

			return {
				table,
				pageReset,
				updateData,
				onExportTable,
				onRefreshTable,
				onSelectionChange,
				getSelectlist,
				expandChange,
				dataChanged,
				...toRefs(state),
			}
		},
	});
</script>

<style scoped lang="scss">
	.table-container {
		display: flex;
		flex: 1;
		overflow: auto;
		flex-direction: column;
		.el-table {
			flex: 1;
		}
		.table-footer {
			display: flex;
			padding: 5px 15px;
			.table-footer-tool {
				flex: 1;
				display: flex;
				align-items: center;
				justify-content: flex-end;
				i {
					margin-right: 10px;
					cursor: pointer;
					color: var(--el-text-color-regular);
					&:last-of-type {
						margin-right: 0;
					}
				}
				&-button {
					width: 24px;
					height: 24px;
					padding-left: 12px;
				}
			}
		}
	}
</style>
