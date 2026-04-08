<template>
	<q-container :cardMode="getCardMode()">
		<q-search :search="search" @search="onSearch" >
			<el-button icon="ele-Delete" type="danger" round :disabled="multiple" @click="onRowDel(null)">
				删除
			</el-button>
			<el-button icon="ele-Delete" size="default" round type="danger" class="ml10" @click="onRowClear()">
				清空日志
			</el-button>
		</q-search>

		<q-table
				ref="tableRef"
				row-key="infoId"
				:data="state.tableData"
				:recordCount="state.total"
				:summaryMethod="getSummaries"
				:default-sort="{ prop: 'infoId', order: 'ascending' }"
				@pageChange="onTablePageChange"
				@selection-change="handleSelectionChange"
		>
			<el-table-column label="编号" fixed="left" prop="infoId" />
			<el-table-column label="登录名称" prop="loginName" min-width="120" show-overflow-tooltip/>
			<el-table-column label="登录地址" prop="ipaddr" width="130" show-overflow-tooltip />
			<el-table-column label="登录地点" prop="loginLocation" min-width="120" show-overflow-tooltip />
			<el-table-column label="浏览器" prop="browser" min-width="100" show-overflow-tooltip/>
			<el-table-column label="操作系统" min-width="120" prop="os" show-overflow-tooltip/>
			<el-table-column label="登录状态" align="center" prop="status" min-width="120" :formatter="statusFormat" >
				<template #default="scope">
					<el-tag type="success" v-if="scope.row.status === 1">成功</el-tag>
					<el-tag type="info" v-else>失败</el-tag>
				</template>
			</el-table-column>
			<el-table-column label="操作信息" prop="msg" min-width="150" />
			<el-table-column label="登录日期" prop="loginTime" width="180" />
			<el-table-column label="登录模块" prop="module" min-width="120" ></el-table-column>
		</q-table>
	</q-container>
</template>

<script setup lang="ts" name="apiV1SystemLoginLogList">
	import { reactive, onMounted, ref, getCurrentInstance, unref } from 'vue';
	import { ElMessageBox, ElMessage, FormInstance } from 'element-plus';
	import { logList, deleteLog, clearLog } from '/@/api/system/monitor/loginLog';
	import { getCardMode } from '/@/utils/common';
	import { TableLoginLogsColumn, TableLoginLogsModel } from '/@/types/model';

	const { proxy } = getCurrentInstance() as any;
	const { sys_login_status } = proxy.useDict('sys_login_status');
	const state = reactive<TableLoginLogsModel>({
		ids: [],
		total: 0,
		loading: false,
		tableData: [],
		param: {
			pageNum: 1,
			pageSize: 10,
			dateRange: [],
			status: '',
			ipaddr: '',
			loginLocation: '',
			userName: '',
		},
	});
	const tableRef = ref();

	// 非单个禁用
	const single = ref(true);
	// 非多个禁用
	const multiple = ref(true);

	const search = ref(
			[
				{ label: '登录IP', prop: 'ipaddr', placeholder: '请输入登录地址', required: false, type: 'input' },
				{ label: '登录地点', prop: 'loginLocation', placeholder: '请输入登录地点', required: false, type: 'input' },
				{ label: '用户名称', prop: 'userName', placeholder: '请输入用户名称', required: false, type: 'input' },
				{
					label: '状态',
					prop: 'status',
					placeholder: '请选择',
					required: false,
					type: 'select',
					options: sys_login_status,
				},
				{ label: '登录时间', prop: 'dateRange', placeholder: '请输入登录时间', required: false, type: 'daterange', format:'YYYY-MM-DD'},
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

	const getTableData = () => {
		logList(state.param).then((res: any) => {
			state.tableData =  res.data.list ?? [];
			state.total = res.data.total;
		});
	};
	// 删除日志
	const onRowDel = (row: TableLoginLogsColumn) => {
		let msg = '确认要删除所选数据？';
		let ids: number[] = [];
		if (row) {
			msg = `确认要删除日志【${row.loginName}】吗?`;
			ids = [row.infoId];
		} else {
			ids = state.ids;
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
				deleteLog(ids).then(() => {
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
				clearLog().then(() => {
					ElMessage.success('清除成功');
					getTableData();
				});
			})
			.catch(() => {});
	};
	const getSummaries = (param: any) => {
		const { columns, data } = param;
		const sums: any = [];
		columns.forEach((column: any, index: any) => {
			if (index === 1) {
				sums[index] = '合计';
				return;
			}
			if (index === 2) {
				sums[index] = <string>(<unknown>state.total);
				return;
			}
			// const values = data.map((item: any) => Number(item[column.property]));
			// if (!values.every((value: any) => isNaN(value))) {
			// 	sums[index] = values.reduce((prev: any, curr: any) => {
			// 		const value = Number(curr);
			// 		if (!isNaN(value)) {
			// 			return prev + curr;
			// 		} else {
			// 			return prev;
			// 		}
			// 	}, 0);
			// 	sums[index] += ' 元';
			// } else {
			// 	sums[index] = '';
			// }
		});
		return sums;
	};
	// 页面加载时
	onMounted(() => {
		getTableData();
	});

	// 多选框选中数据
	const handleSelectionChange = (selection: TableLoginLogsColumn[]) => {
		state.ids = selection.map((item) => item.infoId);
		single.value = selection.length != 1;
		multiple.value = !selection.length;
	};
	// 登录状态字典翻译
	const statusFormat = (row: TableLoginLogsColumn) => {
		return proxy.selectDictLabel(unref(sys_login_status), row.status);
	};
</script>
