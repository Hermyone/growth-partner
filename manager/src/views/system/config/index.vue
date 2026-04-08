<template>
	<q-container :cardMode="getCardMode()">
		<q-search :search="search" @search="onSearch" >
			<el-button icon="ele-Plus" size="default" round type="success" class="ml10" @click="onOpenAddDic">
				新增参数
			</el-button>
			<el-button icon="ele-Delete" size="default" round type="danger" class="ml10" @click="onRowDel(null)">
				删除参数
			</el-button>
		</q-search>

		<q-table ref="tableRef"
				 row-key="configId"
				 :data="state.tableData"
				 :recordCount="state.total"
				 @pageChange="onTablePageChange"
		>
			<el-table-column label="参数主键" prop="configId" />
			<el-table-column label="参数名称" width="200" prop="configName" show-overflow-tooltip />
			<el-table-column label="参数键名" width="200" prop="configKey" show-overflow-tooltip />
			<el-table-column label="参数键值" width="200" prop="configValue" show-overflow-tooltip />
			<el-table-column label="系统内置" align="center" width="100" prop="configType" >
				<template #default="scope">
					<el-tag type="success" v-if="scope.row.configType === 1">是</el-tag>
					<el-tag type="info" v-else>否</el-tag>
				</template>
			</el-table-column>
			<el-table-column label="创建时间" align="center" prop="createdAt" width="180" />
			<el-table-column label="备注" width="300" prop="remark" show-overflow-tooltip />
			<el-table-column label="操作" align="center" width="150" fixed="right">
				<template #default="scope">
					<el-button size="small" text type="primary" @click="onOpenEditDic(scope.row)">修改</el-button>
					<el-divider direction="vertical" />
					<el-button size="small" text type="danger" @click="onRowDel(scope.row)">删除</el-button>
				</template>
			</el-table-column>
		</q-table>

		<EditConfig ref="editDicRef" @dataList="getTableData" :sysYesNoOptions="sys_yes_no" />
	</q-container>
</template>

<script setup lang="ts" name="apiV1SystemDictDataList">
	import { reactive, onMounted, ref, defineAsyncComponent, unref, getCurrentInstance, computed } from 'vue';
	import { ElMessageBox, ElMessage, FormInstance } from 'element-plus';
	import { deleteConfig, getConfigList } from '/@/api/system/config';
	import { getCardMode } from '/@/utils/common';
	import { TableType } from '/@/types/global';
	import { TableConfigColumn, TableConfigModel } from '/@/types/model';

	// 引用组件
	const EditConfig = defineAsyncComponent(() => import('/@/views/system/config/component/editConfig.vue'));

	const { proxy } = getCurrentInstance() as any;
	const tableRef = ref();
	const editDicRef = ref();
	const { sys_yes_no } = proxy.useDict('sys_yes_no');
	const state = reactive<TableConfigModel>({
		ids: [],
		total: 0,
		loading: false,
		tableData: [],
		param: {
			dateRange: [],
			pageNum: 1,
			pageSize: 10,
			configName: '',
			configKey: '',
			configType: '',
		},
	});

	const search = ref(
			[
				{ label: '参数名称', prop: 'configName', placeholder: '请输入参数名称', required: false, type: 'input' },
				{ label: '参数键名', prop: 'configKey', placeholder: '请输入参数键名', required: false, type: 'input' },
				{
					label: '系统内置',
					prop: 'configType',
					placeholder: '请选择',
					required: false,
					type: 'select',
					options: sys_yes_no,
				},
				{ label: '创建时间', prop: 'dateRange', placeholder: '请输入创建时间', required: false, type: 'daterange', format:'YYYY-MM-DD'},
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
		getConfigList(state.param).then((res: any) => {
			state.tableData = res.data.list ?? [];
			state.total = res.data.total;
		});
	};
	// 打开新增字典弹窗
	const onOpenAddDic = () => {
		editDicRef.value.openDialog();
	};
	// 打开修改字典弹窗
	const onOpenEditDic = (row: TableConfigColumn) => {
		editDicRef.value.openDialog(row);
	};
	// 删除字典
	const onRowDel = (row: TableConfigColumn) => {
		let msg = '确认要删除所选数据？';
		let ids: number[] = [];
		if (row) {
			msg = `确认要删除参数【${row.configName}】吗?`;
			ids = [row.configId];
		} else {
			ids = tableRef.value.getSelectlist('configId');
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
				deleteConfig(ids).then(() => {
					ElMessage.success('删除成功');
					getTableData();
				});
			})
			.catch(() => {});
	};
	// 页面加载时
	onMounted(() => {
		getTableData();
	});

</script>
