<!--
 * @Descripttion: 自定义表格组件，支持自定义导出、自定义样式等
 * @version: 2.0
 * @Author: OTQ
 * @Date: 2023-06-10
 * @LastEditors:
 * @LastEditTime:
-->

<template>
	<div class="table-container">
		<el-table
        ref="qTableRef"
        :key="toggleIndex"
				:row-key="rowKey"
				:size="config.size"
				:border="config.border"
				:stripe="config.stripe"
				:highlight-current-row="highlightCurrentRow"
				:row-style="rowStyle"
				:cell-style="cellStyle"
        :summary-method="remoteSummary?remoteSummaryMethod:(summaryMethod ? summaryMethod : getSummaries)"
				style="width: 100%;"
				v-bind="$attrs"
				v-loading="config.loading"
        @selection-change="onSelectionChange"
        @sort-change="sortChange"
        @filter-change="filterChange"
				@row-click="rowClick"
				@select="select"
				@select-all="selectAll"
		>
			<el-table-column width="16" fixed="left" v-if="quickCustomColumn" label-class-name="quick-column">
				<template #header>
					<el-popover
							trigger="click"
							placement="bottom-start"
					>
						<el-checkbox v-for="item in userColumn" v-model="item.show" :key="item.key" :label="item.title" size="small"/>
						<template #reference>
							<el-icon class="table-container-filter"><ele-Edit/></el-icon>
						</template>
					</el-popover>
				</template>
			</el-table-column>
			<el-table-column type="selection" width="40" v-if="isSelection" :selectable="selectableFunc" fixed="left"/>
			<el-table-column type="index" label="序号" width="60" v-if="isSerialNo" fixed="left" >
				<template #default="scope">
					<span>{{scope.$index+(page.pageNum - 1) * page.pageSize + 1}}</span>
				</template>
			</el-table-column>
			<slot></slot>
			<template v-for="(item, index) in userColumn" :key="index">
				<el-table-column
						v-if="item.show??true"
						:column-key="item.key"
						:label="item.title"
						:prop="item.key"
						:width="item.width"
						:min-width="item.minwidth"
						:align="item.align"
						:sortable="item.sortable"
						:fixed="item.fixed"
						:filters="item.filters"
						:filter-method="remoteFilter||!item.filters?null:filterHandler"
						:show-overflow-tooltip="item.showOverflowTooltip"
						filter-placement="bottom-end"
				>
					<template #default="scope">
						<slot :name="item.key" v-bind="scope">
							<div v-if="item.keylink">
								<el-button size="small" link type="primary" @click="onKeyLink(scope.row, item.key)">{{ formatCellValue(item, scope.row[item.key], scope.row) }}</el-button>
							</div>
							<div v-else-if="item.dicts">
								<el-tag :type="formatCellDictType(item, scope.row[item.key])">{{ formatCellValue(item, scope.row[item.key], scope.row) }}</el-tag>
							</div>
							<span v-else>{{ formatCellValue(item, scope.row[item.key], scope.row) }}</span>
						</slot>
					</template>
				</el-table-column>
			</template>
			<el-table-column min-width="1" v-if="autoColumnWidth"></el-table-column>
			<template #empty>
				<el-empty description="暂无数据" />
			</template>
		</el-table>
		<div class="table-footer mt15" v-show="footerVisible">
			<el-pagination
					v-model:current-page="page.pageNum"
					v-model:page-size="page.pageSize"
					:pager-count="5"
					:page-sizes="config.pageSizes"
					:total="config.total"
					:layout="layout"
					:size="paginationSmall ? 'small' : 'default'"
					background
			>
			</el-pagination>
			<div class="table-footer-tool">
				<el-button v-show="footerButton" @click="onExportTable" icon="ele-Download" circle class="table-footer-tool-button"></el-button>
				<el-button v-show="footerButton" @click="onRefreshTable" icon="ele-Refresh" circle class="table-footer-tool-button"></el-button>
				<el-popover
						ref="setRef"
						v-if="columns"
						placement="top"
						title="设置"
						trigger="click"
						:width="500"
						:hide-after="0"
				>
					<template #reference>
						<el-button v-show="footerButton" icon="ele-Setting" circle class="table-footer-tool-button"></el-button>
					</template>
					<columnSetting
							ref="columnSettingRef"
							@userChange="columnSettingChange"
							@save="SaveConfig"
							@back="ResetConfig"
							:column="userColumn"
							:optionButton="tableName !== ''"
					>
						<el-form label-width="80px" label-position="left" size="small" style="margin-left: 20px;">
							<el-form-item label="表格尺寸">
								<el-radio-group v-model="config.size" size="small">
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
					</columnSetting>
				</el-popover>
			</div>
		</div>

		<exportFile ref="exportRef" :tableName="tableName" :tableTitle="globalTitle" :tableData="tableData()" :rowKey="rowKey" :columns="columns" />
	</div>
</template>

<script lang="ts">
	import { defineComponent, reactive, watch, nextTick, toRefs, getCurrentInstance } from 'vue';
	import Cookies from 'js-cookie';
	import { TableColumnCtx, ElMessage } from 'element-plus';
	import Sortable from 'sortablejs';
	import columnSetting from './columnSetting.vue';
	import exportFile from './exportFile.vue';
	import '/@/theme/tableTool.scss';

	// 定义父组件传过来的值
	const props = {
		// 表格名称
		tableName: { type: String, default: "" },
		// 行主键 KEY
		rowKey: { type: String, default: "key" },
		// 定义列
		columns: { type: Array, default: () => [] },
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
    // 远程排序
    remoteSort: { type: Boolean, default: false },
		// 远程过滤
		remoteFilter: { type: Boolean, default: false },
    // 远程合计
    remoteSummary: { type: Boolean, default: false },
		summary: { type: Object, default: {} },
    // 本地合计
    summaryMethod: { type: Function, default: null },
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
		// 自动列宽
		autoColumnWidth: { type: Boolean, default: true },
		// 是否使用小型分页样式
		paginationSmall: { type: Boolean, default: true },
		// 行单击事件回调函数
		rowClick: { type: Function, default: null },
		// 高亮当前行
		highlightCurrentRow: { type: Boolean, default: false },
		// 选中/全选回调事件
		select: { type: Function, default: null},
		selectAll: { type: Function, default: null},

		// 自定义行列样式回调
		rowStyle: { type: Function, default: null },
		cellStyle: { type: Function, default: null },

		selectableFunc: { type: Function, default: null },

		// 是否显示快速自定义列
		quickCustomColumn: { type: Boolean, default: false }
	};

	export default defineComponent({
		name: 'q-table-v2',
		props: props,
		components: { columnSetting, exportFile },
		setup(props, { emit }) {
			// 定义变量内容
			const { proxy } = <any>getCurrentInstance();
			const state = reactive({
				page: {
					pageNum: 1,
					pageSize: props.pageSize,
				},
				selectlist: [] as any,
				toggleIndex: 0,
				userColumn: [] as any,
        sourceColumn: [] as any,
        prop: null,
        order: null,
				config: {
					total: 0,
					loading: props.loading,
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
				proxy.$refs.exportRef.open();
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
			// 清除
			const clear = () => {
				state.selectlist = [];
				proxy.$refs.qTableRef.clearSelection();
				proxy.$refs.qTableRef.clearSort();
				proxy.$refs.qTableRef.clearFilter();
			};

			// 当前表格
			const table = () => {
				return proxy.$refs.qTableRef;
			};

			// 表格数据
			const tableData = () => {
				return table()?.data || [];
			};

			// 更新数据
			const updateData = () => {
				clear();
				state.config.loading = true;
				emit('pageChange', state.page);

				setTimeout(() => {
					state.config.loading = false;
					emit('queryCompleted');
				}, 500);
			};

			// 格式化单元格数据
			const formatCellValue = (column: any, value: any, row: any) => {
				if(column.dicts){
					let v = '';
					column.dicts.forEach((item: any) => {
						if(item.value == value){
							v = item.text;
							return;
						}
					});
					return v;
				}
				if(column.formatter && typeof column.formatter === 'function'){
					return column.formatter(row, value, column);
				}
				return value;
			};

			const formatCellDictType = (column: any, value: any) => {
				if(column.dicts){
					let v = '';
					column.dicts.forEach((item: any) => {
						if(item.value == value){
							v = item.type;
							return;
						}
					});
					return v;
				}
			};

            //远程合计行处理
            const remoteSummaryMethod = (param: any) =>{
                const { columns } = param;
                const sums : any = [];
                columns.forEach((column: any, index: number) => {
                    if(props.quickCustomColumn && index == 1 || !props.quickCustomColumn && index == 0) {
                        sums[index] = '合计';
                        return;
                    }
                    const values =  props.summary[column.property];
                    if(values){
                        sums[index] = values;
                    }else{
                        sums[index] = '';
                    }
                })
                return sums;
            };

            // 排序事件
            const sortChange = (obj: any) => {
                if(!props.remoteSort){
                    return false;
                }
                if(obj.column && obj.prop){
                    state.prop = obj.prop;
					state.order = obj.order;
                }else{
                    state.prop = null;
                    state.order = null;
                }
                updateData();
            };

			// 过滤事件
			const filterChange = (filters: any) => {
				if(!props.remoteFilter){
					return false;
				}
				Object.keys(filters).forEach(key => {
					filters[key] = filters[key].join(',');
				})
				updateData();
			};

			// 本地过滤
			const filterHandler = (value : string, row : any, column : TableColumnCtx<any>) => {
				const property = column.property;
				return row[property] === value;
			};

			// 自定义列变化事件
			const columnSettingChange = (userColumn: any) => {
				state.userColumn = userColumn;
				// state.toggleIndex += 1;
			};

			// 主键超链接事件
			const onKeyLink = (row: any, field: string) => {
				emit('keyLink', row, field);
			};

			// 保存自定义配置
			const SaveConfig = (userColumn: any) => {
				proxy.$refs.columnSettingRef.isSave = true;
				try {
					new Promise(() => {
						setTimeout(()=>{
							Cookies.set(props.tableName, JSON.stringify(userColumn));
						},500);
					});
				}catch(error){
					ElMessage.error('保存失败');
				}

				ElMessage.success('保存成功');
				setTimeout(() => {
					proxy.$refs.columnSettingRef.isSave = false;
					proxy.$refs.setRef.hide();
				}, 500);
			};
			// 重置配置
			const ResetConfig = () => {
				proxy.$refs.columnSettingRef.isSave = true;
				try {
					new Promise(() => {
						setTimeout(()=>{
							Cookies.remove(props.tableName);
						},500);
					});
					state.userColumn = state.sourceColumn;
					proxy.$refs.columnSettingRef.usercolumn = JSON.parse(JSON.stringify(state.userColumn||[]));
				}catch(error){
					ElMessage.error('重置失败');
				}

				ElMessage.success('重置成功');
				setTimeout(() => {
					proxy.$refs.columnSettingRef.isSave = false;
				}, 500);
			};

      const getColumnSummary = (key: string) => {
        const column : any = props.columns.find((item:any) => item.key === key);
        return column ? column.summary : '';
      };

      const getSummaries = (param: any) => {
        const { columns, data } = param;
        const sums : any = [];
        columns.forEach((column: any, index: number) => {
          const kind = getColumnSummary(column.columnKey);
          if (kind === "static") {
            sums[index] = '合计';
          } else if (kind === "total") {
            sums[index] = <string>(<unknown>props.recordCount);
          } else if (kind === "sum") {
            const values: any = data.map((item: any) => Number(item[column.property]));
            if (!values.every((value: any) => isNaN(value))) {
              sums[index] = values.reduce((prev: any, curr: any) => {
                const value = Number(curr);
                if (!isNaN(value)) {
                  return prev + curr;
                } else {
                  return prev;
                }
              }, 0);
            } else {
              sums[index] = '';
            }
          }
        });

        return sums;
      };

			// 页面加载时
			onMounted: {
				state.userColumn = props.columns;
        state.sourceColumn = JSON.parse(JSON.stringify(props.columns));
				const customColumn : any = Cookies.get(props.tableName);
				if(customColumn){
					state.userColumn = JSON.parse(customColumn);
				}

				nextTick(() => {
					updateData();
				});
			}

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
				onExportTable,
				onSelectionChange,
				onRefreshTable,
				onKeyLink,
				sortChange,
				filterChange,
				filterHandler,
				formatCellValue,
				formatCellDictType,
				columnSettingChange,
				clear,
				tableData,
				pageReset,
				updateData,
        remoteSummaryMethod,
				SaveConfig,
				ResetConfig,
				getSelectlist,
        getColumnSummary,
        getSummaries,
				...toRefs(state),
			}
		},
	});
</script>

<style scoped lang="scss">
    :deep(.quick-column) {
		.cell {
			padding: 1px;
		}
    }
	.table-container {
		display: flex;
		flex: 1;
		overflow: auto;
		flex-direction: column;
		&-filter {
			width: 1em;
			height: 1em;
			margin-top: 2px;
			:hover {
				color: var(--el-color-primary);
			}
		}
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
