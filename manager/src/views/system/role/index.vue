<template>
	<q-container :cardMode="getCardMode()">
		<q-search :search="search" @search="onSearch" >
			<el-button icon="ele-Plus" size="default" round type="success" class="ml10" @click="onOpenAddRole">
				新增角色
			</el-button>
		</q-search>

		<q-table ref="tableRef"
				 :data="state.tableData"
				 :recordCount="state.total"
				 :isSerialNo="true"
				 :isSelection="false"
				 @pageChange="onTablePageChange"
		>
			<el-table-column prop="name" label="角色名称" show-overflow-tooltip></el-table-column>
			<el-table-column prop="listOrder" label="排序" show-overflow-tooltip></el-table-column>
			<el-table-column prop="status" label="角色状态" show-overflow-tooltip>
				<template #default="scope">
					<el-tag type="success" v-if="scope.row.status === 1">启用</el-tag>
					<el-tag type="info" v-else>禁用</el-tag>
				</template>
			</el-table-column>
			<el-table-column prop="remark" label="角色描述" show-overflow-tooltip></el-table-column>
			<el-table-column prop="createdAt" label="创建时间" show-overflow-tooltip></el-table-column>
			<el-table-column label="操作" align="center" width="150" fixed="right">
				<template #default="scope">
					<el-button size="small" text type="primary" @click="onOpenEditRole(scope.row)">修改</el-button>
					<el-divider direction="vertical" />
					<el-button size="small" text type="danger" @click="onRowDel(scope.row)">删除</el-button>
				</template>
			</el-table-column>
		</q-table>

		<EditRole ref="editRoleRef" @getRoleList="getTableData" />
	</q-container>
</template>

<script setup lang="ts" name="apiV1SystemRoleList">
	import { reactive, onMounted, ref, defineAsyncComponent, toRaw, getCurrentInstance } from 'vue';
	import { ElMessageBox, ElMessage } from 'element-plus';
	import { deleteRole, getRoleList } from '/@/api/system/role';
	import { getCardMode } from '/@/utils/common';
	import { TableRoleColumn, TableRoleModel } from '/@/types/model';

	// 引入组件
	const EditRole = defineAsyncComponent(() => import('/@/views/system/role/component/editRole.vue'));

	const { proxy } = getCurrentInstance() as any;
	const tableRef = ref();
	const editRoleRef = ref();
	const state = reactive<TableRoleModel>({
		ids: [],
		total: 0,
		loading: false,
		tableData: [],
		param: {
			roleName: '',
			roleStatus: '',
			pageNum: 1,
			pageSize: 10,
		},
	});

	const search = ref(
			[
				{ label: '角色名称', prop: 'roleName', placeholder: '请输入角色名称', required: false, type: 'input' },
				{
					label: '状态',
					prop: 'roleStatus',
					placeholder: '请选择',
					required: false,
					type: 'select',
					options: [
						{ label:'启用', value:'1' },
						{ label:'禁用', value:'0' },
					],
				},
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
		const data: Array<TableRoleColumn> = [];
		getRoleList(state.param).then((res) => {
			const list = res.data.list ?? [];
			list.map((item: TableRoleColumn) => {
				data.push({
					id: item.id,
					status: item.status,
					listOrder: item.listOrder,
					name: item.name,
					remark: item.remark,
					dataScope: item.dataScope,
					createdAt: item.createdAt,
				});
			});
			state.tableData = data;
			state.total = res.data.total;
		});
	};
	// 打开新增角色弹窗
	const onOpenAddRole = () => {
		editRoleRef.value.openDialog();
	};
	// 打开修改角色弹窗
	const onOpenEditRole = (row: Object) => {
		editRoleRef.value.openDialog(toRaw(row));
	};

	// 删除角色
	const onRowDel = (row: any) => {
		ElMessageBox.confirm(`确认要删除角色【${row.name}】吗?`, '提示', {
			confirmButtonText: '确认',
			cancelButtonText: '取消',
			type: 'warning',
		})
			.then(() => {
				deleteRole(row.id).then(() => {
					ElMessage.success('删除成功');
					proxy.$refs['editRoleRef'].resetMenuSession();
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
