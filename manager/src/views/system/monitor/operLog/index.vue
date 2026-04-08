<template>
	<q-container :cardMode="getCardMode()">
		<q-search :search="search" @search="onSearch" >
			<el-button icon="ele-Delete" type="danger" round :disabled="multiple" @click="handleDelete(null)" v-auth="'api/v1/system/sysOperLog/delete'">
				删除
			</el-button>
			<el-button icon="ele-Delete" type="danger" round class="ml10" @click="onRowClear()">
				清空日志
			</el-button>
		</q-search>

		<q-table
				ref="tableRef"
				row-key="operId"
				:data="state.tableData"
				:recordCount="state.total"
				@pageChange="onTablePageChange"
				@selection-change="handleSelectionChange"
		>
			<el-table-column label="日志编号" prop="operId" min-width="100px" fixed="left" />
			<el-table-column label="系统模块" prop="title" min-width="180px" show-overflow-tooltip/>
			<el-table-column label="请求方式" prop="requestMethod" min-width="100px" >
				<template #default="scope">
					<q-tag v-if="scope.row.requestMethod === 'GET'">{{ requestMethodFormat(scope.row) }} </q-tag>
					<q-tag type="success" v-if="scope.row.requestMethod === 'POST'">{{ requestMethodFormat(scope.row) }} </q-tag>
					<q-tag type="warning" v-if="scope.row.requestMethod === 'PUT'">{{ requestMethodFormat(scope.row) }} </q-tag>
					<q-tag type="danger" v-if="scope.row.requestMethod === 'DELETE'">{{ requestMethodFormat(scope.row) }} </q-tag>
				</template>	
			</el-table-column>
			<el-table-column label="操作人员" prop="operName" min-width="100px" show-overflow-tooltip/>
			<el-table-column label="请求URL" prop="operUrl" min-width="200px" show-overflow-tooltip />
			<el-table-column label="主机地址" prop="operIp" min-width="140px" show-overflow-tooltip/>
			<el-table-column label="操作地点" prop="operLocation" min-width="120px" show-overflow-tooltip/>
			<el-table-column label="操作时间" align="center" prop="operTime" min-width="180px">
				<template #default="scope">
					<span>{{ proxy.parseTime(scope.row.operTime, '{y}-{m}-{d} {h}:{i}:{s}') }}</span>
				</template>
			</el-table-column>
			<el-table-column label="操作" align="center" class-name="small-padding" width="120px" fixed="right">
				<template #default="scope">
					<el-button size="small" type="primary" text @click="handleView(scope.row)" v-auth="'api/v1/system/sysOperLog/view'">详情</el-button>
					<el-divider direction="vertical" />
					<el-button size="small" type="danger" text @click="handleDelete(scope.row)" v-auth="'api/v1/system/sysOperLog/delete'">删除</el-button>
				</template>
			</el-table-column>
		</q-table>

		<apiV1SystemSysOperLogDetail
			ref="detailRef"
			:requestMethodOptions="sys_oper_log_type"
		></apiV1SystemSysOperLogDetail>
	</q-container>
</template>

<script setup lang="ts" name="apiV1SystemSysOperLogList">
	import { ItemOptions } from '/@/api/items';
	import { reactive, onMounted, ref, defineAsyncComponent, computed, getCurrentInstance, toRaw } from 'vue';
	import { ElMessageBox, ElMessage, FormInstance } from 'element-plus';
	import { listSysOperLog, delSysOperLog, clearOperLog } from '/@/api/system/monitor/operLog';
	import { getCardMode } from '/@/utils/common';
	import { TableOperLogInfoData, TableOperLogsColumn, TableOperLogsModel } from '/@/types/model';

	// 引入组件
	const apiV1SystemSysOperLogDetail = defineAsyncComponent(() => import('/@/views/system/monitor/operLog/component/detail.vue'));

	const { proxy } = <any>getCurrentInstance();
	const loading = ref(false);
	const tableRef = ref();
	const detailRef = ref();

	// 非单个禁用
	const single = ref(true);
	// 非多个禁用
	const multiple = ref(true);

	// 字典选项数据
	const { sys_oper_log_type } = proxy.useDict('sys_oper_log_type');
	// deptNameOptions关联表数据
	const deptNameOptions = ref<Array<ItemOptions>>([]);
	const state = reactive<TableOperLogsModel>({
		ids: [],
		operIds: [],
		total: 0,
		loading: false,
		tableData: [],
		param: {
			pageNum: 1,
			pageSize: 10,
			title: undefined,
			requestMethod: undefined,
			operName: undefined,
			status: undefined,
			dateRange: [],
		},
	});

	const search = ref(
			[
				{ label: '系统模块', prop: 'title', placeholder: '请输入系统模块', required: false, type: 'input' },
				{
					label: '请求方式',
					prop: 'requestMethod',
					placeholder: '请选择',
					required: false,
					type: 'select',
					options: sys_oper_log_type,
				},
				{ label: '操作人员', prop: 'operName', placeholder: '请输入操作人员', required: false, type: 'input' },
				{ label: '操作时间', prop: 'dateRange', placeholder: '请输入操作时间', required: false, type: 'daterange', format:'YYYY-MM-DD'},
			]
	);

	// 搜索点击时表单回调
	const onSearch = (data: EmptyObjectType) => {
		state.param = Object.assign({}, state.param, { ...data });
		tableRef.value?.pageReset();
	};

	// 分页改变时回调
	const onTablePageChange = (page: TablePageType) => {
		state.param.pageNum = page.pageNum;
		state.param.pageSize = page.pageSize;
		getTableData();
	};

	// 页面加载时
	onMounted(() => {
		getTableData();
	});

	// 获取列表数据
	const getTableData = () => {
		listSysOperLog(state.param).then((res: any) => {
			let list = res.data.list ?? [];
			state.tableData = list ?? [];
			state.total = res.data.total;
		});
	};

	// 请求方式字典翻译
	const requestMethodFormat = (row: TableOperLogsColumn) => {
		return proxy.selectDictLabel(sys_oper_log_type.value, row.requestMethod);
	};
	// 多选框选中数据
	const handleSelectionChange = (selection: Array<TableOperLogInfoData>) => {
		state.operIds = selection.map((item) => item.operId);
		single.value = selection.length != 1;
		multiple.value = !selection.length;
	};

	const handleDelete = (row: TableOperLogsColumn) => {
		let msg = '确认要删除所选数据？';
		let ids: number[] = [];
		if (row) {
			msg = `确认要删除日志吗?`;
			ids = [row.operId];
		} else {
			ids = state.operIds;
		}
		if (ids.length === 0) {
			ElMessage.error('请选择要删除的数据。');
			return;
		}
		ElMessageBox.confirm(msg, '提示', {
			confirmButtonText: '确认',
			cancelButtonText: '取消',
			type: 'warning',
		})
			.then(() => {
				delSysOperLog(ids).then(() => {
					ElMessage.success('删除成功');
					getTableData();
				});
			})
			.catch(() => {});
	};
	// 清空日志
	const onRowClear = () => {
		ElMessageBox.confirm('确认要删除所选数据？', '提示', {
			confirmButtonText: '确认',
			cancelButtonText: '取消',
			type: 'warning',
		})
			.then(() => {
				clearOperLog().then(() => {
					ElMessage.success('清除成功');
					getTableData();
				});
			})
			.catch(() => {});
	};
	const handleView = (row: TableOperLogsColumn) => {
		detailRef.value.openDialog(toRaw(row));
	};
</script>

<style lang="scss" scoped>
.colBlock {
	display: block;
}
.colNone {
	display: none;
}
</style>
