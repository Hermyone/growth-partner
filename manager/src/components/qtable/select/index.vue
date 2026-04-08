<!--
 * @Descripttion: 自定义下拉选择表格组件
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-07-13
 * @LastEditors:
 * @LastEditTime:
-->

<template>
	<el-select
			class="q-select-style"
			ref="selectRef"
			v-model="defaultValue"
			:size="size"
			:clearable="clearable"
			:multiple="multiple"
			:collapse-tags="collapseTags"
			:collapse-tags-tooltip="collapseTagsTooltip"
			:filterable="filterable"
			:placeholder="placeholder"
			:style="{width: width+'px'}"
			@remove-tag="removeTag"
			@visible-change="visibleChange"
			@clear="clear">
		<template #empty>
			<div class="q-select-style-header">
				<slot name="header"></slot>
			</div>
			<div class="q-select-style-table" :style="{width: tableWidth+'px', height: tableHeight+'px'}" v-loading="loading">
				<slot name="aside"></slot>
				<q-table ref="tableRef"
						:isSelection="multiple"
						:isSerialNo="serialno"
						:data="tableData"
						:recordCount="total"
						:footerButton="false"
						:paginationSmall="true"
						:rowClick="rowClick"
						:highlightCurrentRow="!multiple"
						:row-key="keyFieldName"
						@pageChange="onPageChange"
						@select="select"
						@selectAll="selectAll"
						@queryCompleted="onQueryCompleted"
						layout="total, sizes, prev, pager, next"
						size="small"
						style="height: 100%"
				>
					<slot></slot>
				</q-table>
			</div>
		</template>
	</el-select>
</template>

<script lang="ts">
	import { defineComponent, reactive, toRefs, watch, nextTick, onMounted, getCurrentInstance } from 'vue';
	import QTable from '../v1/index.vue';

	// 定义父组件传过来的值
	const props = {
		modelValue: null,
		// 默认显示宽度
		width: {type: Number, default: 250},
		// q-table表格的分页回调函数
		onPageChange: { type: Function, default: null},
		// q-table数据属性
		tableData: { type: Object, default: []},
		// q-table数据行总数，用于分页
		total: { type: Number, default: 0},
		// 组件大小
		size: { type: String, default: "default" },
		// 站位文字
		placeholder: { type: String, default: "请选择" },
		// 显示清空选项
		clearable: { type: Boolean, default: false },
		// 是否可多选
		multiple: { type: Boolean, default: false },
		// 是否显示序号
		serialno: { type: Boolean, default: false },
		// 是否可筛选
		filterable: { type: Boolean, default: false },
		// 多选时是否将选中值按文字的形式展示, 当鼠标悬停于折叠标签的文本时，是否显示所有选中的标签。 要使用此属性，collapse-tags属性必须设定为 true
		collapseTags: { type: Boolean, default: false },
		collapseTagsTooltip: { type: Boolean, default: false },
		// 表格高度和宽度
		tableHeight: {type: Number, default: 300},
		tableWidth: {type: Number, default: 400},
		props: { type: Object, default: () => {} },
		// 加载中标记
		loading: { type: Boolean, default: false },
		// 需要显示列字段名
		labelFieldName: { type: String, default: 'ID' },
		// 主键列字段名
		keyFieldName: { type: String, default: 'key' },
		// 格式化显示
		formatText: { type: Function, default: null },
	};

	export default defineComponent({
		name: 'q-select-table',
		components: { QTable },
		props: props,
		setup(props, { emit }) {
			// 定义变量内容
			const { proxy } = <any>getCurrentInstance();
			const state = reactive({
				defaultValue: [] as any,
			});

			const tableGrid = () => {
				return proxy.$refs.tableRef.value;
			};

			const visibleChange = (val : boolean) => {
				if(val){
					proxy.$refs.tableRef.pageReset();
				}else{
					autoCurrentLabel();
				}
			};

			const clear = () => {
				emit('update:modelValue', state.defaultValue);
			};

			const removeTag = (tag: any) => {
				var row = findRowByKey(tag[props.keyFieldName]);
				proxy.$refs.tableRef.table().toggleRowSelection(row, false);
				emit('update:modelValue', state.defaultValue);
			};

			const findRowByKey = (value: any) => {
				return props.tableData.find((item:any) => item[props.keyFieldName] === value);
			};

			const autoCurrentLabel = () => {
				nextTick(() => {
					if(state.defaultValue){
						if(props.multiple){
							proxy.$refs.selectRef.selected.forEach((item:any) => {
								item.currentLabel = item.value[props.labelFieldName]
							});
						}else{
							if(props.formatText){
								proxy.$refs.selectRef.selectedLabel = props.formatText(state.defaultValue);
							}else{
								proxy.$refs.selectRef.selectedLabel = state.defaultValue[props.labelFieldName];
							}
						}
					}
				});
			};

			const select = (rows: any, row: any) => {
				var isSelect = rows.length && rows.indexOf(row) !== -1
				if(isSelect){
					state.defaultValue.push(row)
				}else{
					state.defaultValue.splice(state.defaultValue.findIndex((item:any) => item[props.keyFieldName] == row[props.keyFieldName]), 1)
				}
				autoCurrentLabel()
				emit('update:modelValue', state.defaultValue);
				emit('change', state.defaultValue);
			};

			const selectAll = (rows: any) => {
				var isAllSelect = rows.length > 0
				if(isAllSelect){
					rows.forEach((row:any) => {
						var isHas = state.defaultValue.find((item:any) => item[props.keyFieldName] == row[props.keyFieldName])
						if(!isHas){
							state.defaultValue.push(row)
						}
					});
				}else{
					props.tableData.forEach((row:any) => {
						var isHas = state.defaultValue.find((item:any) => item[props.keyFieldName] == row[props.keyFieldName])
						if(isHas){
							state.defaultValue.splice(state.defaultValue.findIndex((item:any) => item[props.keyFieldName] == row[props.keyFieldName]), 1)
						}
					});
				}
				autoCurrentLabel();
				emit('update:modelValue', state.defaultValue);
				emit('change', state.defaultValue);
			};

			// 查询完成，显示默认值
			const onQueryCompleted = () => {
				nextTick(() => {
					if(state.defaultValue){
						if(props.multiple){
							state.defaultValue.forEach((row:any) => {
								var setrow = props.tableData.filter((item:any) => item[props.keyFieldName]===row[props.keyFieldName] )
								if(setrow.length > 0){
									proxy.$refs.tableRef.table().toggleRowSelection(setrow[0], true);
								}
							})
						}else{
							var setrow = props.tableData.filter((item: any) => item[props.keyFieldName] === state.defaultValue[props.keyFieldName])
							proxy.$refs.tableRef.table().setCurrentRow(setrow[0]);
						}

						proxy.$refs.tableRef.table().setScrollTop(0);
					}
				});
			};

			const rowClick = (row: any) => {
				if(!props.multiple){
					state.defaultValue = row;
					proxy.$refs.selectRef.blur();
					autoCurrentLabel();
					emit('update:modelValue', state.defaultValue);
					emit('change', state.defaultValue);
				}
			};

			watch(
				() => props.modelValue,
				(handler) => {
					state.defaultValue = props.modelValue;
					autoCurrentLabel();
				},
				{
					deep: true,
				}
			);

			// 页面加载时
			onMounted(() => {
				state.defaultValue = props.modelValue || [];
				autoCurrentLabel();
			})

			return {
				tableGrid,
				removeTag,
				clear,
				visibleChange,
				select,
				selectAll,
				onQueryCompleted,
				rowClick,
				...toRefs(state),
			}
		},
	});
</script>

<style scoped lang="scss">
	.q-select-style {
		&-header {
			padding: 15px 15px 0 15px;
		}
		&-table {
			display: flex;
			align-items: center;
			justify-content: space-between;
		}
	}
</style>