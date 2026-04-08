<!--
 * @Descripttion: 自定义表格表单
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-07-10
 * @LastEditors:
 * @LastEditTime:
-->

<template>
	<div class="table-container" ref="formTable">
		<el-table
				ref="table"
				:border="border"
				:size="size"
				:stripe="stripe"
				:data="data"
		>
			<el-table-column type="index" width="40" fixed="left" v-if="!hideAdd || !hideDelete">
				<template #header>
					<el-button v-if="!hideAdd" type="primary" icon="ele-Plus" size="small" circle style="margin-left: 6px" @click="addRow"></el-button>
				</template>
				<template #default="scope">
					<div :class="['table-container-handle', {'table-container-handle-delete':!hideDelete}]">
						<span>{{scope.$index + 1}}</span>
						<el-button v-if="!hideDelete" type="danger" icon="ele-Delete" size="small" circle @click="delRow(scope.row, scope.$index)"></el-button>
					</div>
				</template>
			</el-table-column>
			<el-table-column label="" width="30" v-if="dragSort">
				<template #default>
					<div class="move" style="cursor: move;">
						<svg class="move" style="cursor: move;" viewBox="64 64 896 896" focusable="false" data-icon="holder" width="1em" height="1em" fill="currentColor" aria-hidden="true">
							<path d="M300 276.5a56 56 0 1056-97 56 56 0 00-56 97zm0 284a56 56 0 1056-97 56 56 0 00-56 97zM640 228a56 56 0 10112 0 56 56 0 00-112 0zm0 284a56 56 0 10112 0 56 56 0 00-112 0zM300 844.5a56 56 0 1056-97 56 56 0 00-56 97zM640 796a56 56 0 10112 0 56 56 0 00-112 0z"></path>
						</svg>
					</div>
				</template>
			</el-table-column>
			<slot></slot>
			<template #empty>
				{{ placeholder }}
			</template>
		</el-table>
	</div>
</template>

<script lang="ts">
	import { defineComponent, reactive, toRefs, watch, onMounted, nextTick } from 'vue';
	import Sortable from 'sortablejs'

	// 定义父组件传过来的值
	const props = {
		// 表格尺寸
		size: { type: String, default: "default" },
		// 是否显示纵向边框
		border: { type: Boolean, default: false },
		// 是否显示斑马纹
		stripe: { type: Boolean, default: false },
		// 数据
		modelValue: { type: Array, default: () => [] },
		// 添加模板
		addTemplate: { type: Object, default: () => {} },
		placeholder: { type: String, default: "暂无数据" },
		// 是否也拖动
		dragSort: { type: Boolean, default: false },
		// 是否隐藏添加按钮
		hideAdd: { type: Boolean, default: false },
		// 是否隐藏删除按钮
		hideDelete: { type: Boolean, default: false }
	};

	export default defineComponent({
		name: 'q-form-table',
		props: props,
		setup(props, { emit }) {
			// 定义变量内容
			const state = reactive({
				data: [] as any,
			});

			// 页面加载时
			onMounted(() => {
				nextTick(() => {
					state.data = props.modelValue;
				});
			})

			watch(
					() => props.modelValue,
					() => {
						state.data = props.modelValue;
					}
			);
			watch(
					() => state.data,
					(header) => {
						emit('update:modelValue', state.data);
					},
					{
						deep: true
					}
			);

			// 增加一行
			const addRow = () => {
				const v = JSON.parse(JSON.stringify(props.addTemplate));
				state.data.push(v);
				emit('add', v);
			};

			// 删除一行
			const delRow = (row: any, index : number) => {
				const v = state.data.splice(index, 1);
				emit('delete', v);
			};

			// 插入一行
			const insertRow = (row: any) => {
				const v = row || JSON.parse(JSON.stringify(props.addTemplate));
				state.data.push(v);
				emit('insert', v);
			};

			return {
				addRow,
				delRow,
				insertRow,
				...toRefs(state),
			}
		},
		mounted(){
			if(props.dragSort)
				this.dropRow();
		},
		methods: {
			dropRow: function() {
				const _this = this;
				const tableRef : any = this.$refs.table;
				const formTableRef : any = _this.$refs.formTable;
				const tbody = tableRef.$el?.querySelector('.el-table__body-wrapper tbody');
				Sortable.create(tbody, {
					handle: '.move',
					animation: 300,
					ghostClass: 'ghost',
					onEnd({ newIndex, oldIndex }) {
						_this.data.splice(<number>newIndex, 0, _this.data.splice(<number>oldIndex, 1)[0]);
						const newArray = _this.data.slice(0);
						const tmpHeight = formTableRef.offsetHeight;
						formTableRef.style.setProperty('height', tmpHeight + 'px');
						_this.data = [];
						_this.$nextTick(() => {
							_this.data = newArray;
							_this.$nextTick(() => {
								formTableRef.style.removeProperty('height');
							});
						});
					},
				});
			},
		},
	});
</script>

<style scoped>
	.table-container {
		width: 100%;
		:deep(.el-button i.el-icon,.el-button i.iconfont,.el-button i.fa,.el-button--default i.iconfont,.el-button--default i.fa){
			font-size: 12px !important;
			margin-right: 0;
		}
		:deep(th){
			background: transparent!important;
		}
	}
	.table-container .table-container-handle {text-align: center;}
	.table-container .table-container-handle span {display: inline-block;}
	.table-container .table-container-handle button {display: none;}
	.table-container .hover-row .table-container-handle-delete span {display: none;}
	.table-container .hover-row .table-container-handle-delete button {display: inline-block;}
	.table-container .move {text-align: center;font-size: 14px;margin-top: 3px;}

	:deep(.el-table__inner-wrapper::before) {
		background: none;
	}
</style>