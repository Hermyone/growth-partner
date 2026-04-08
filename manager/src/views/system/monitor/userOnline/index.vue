<template>
	<q-container :cardMode="getCardMode()">
		<q-search :search="search" @search="onSearch" >
			<el-button size="default" round icon="ele-Delete" type="danger" class="ml10" @click="onRowDel(null)">强制退出</el-button>
		</q-search>

		<q-table
				ref="tableRef"
				row-key="uuid"
				:data="state.tableData"
				:recordCount="state.total"
				:isSerialNo="true"
				@pageChange="onTablePageChange"
		>
			<el-table-column prop="uuid" label="会话编号" width="300" show-overflow-tooltip></el-table-column>
			<el-table-column prop="userName" label="登录名称" show-overflow-tooltip></el-table-column>
			<el-table-column prop="ip" label="主机" show-overflow-tooltip></el-table-column>
			<el-table-column prop="explorer" label="浏览器" show-overflow-tooltip></el-table-column>
			<el-table-column label="操作系统" align="center" prop="os" show-overflow-tooltip/>
			<el-table-column prop="createTime" align="center" label="创建时间" width="180"></el-table-column>
			<el-table-column label="操作" align="center" width="100" fixed="right">
				<template #default="scope">
					<el-button size="small" type="danger" text icon="ele-RemoveFilled" v-if="buttonVisible(scope.row.createTime)" @click="onRowDel(scope.row)">强退</el-button>
				</template>
			</el-table-column>
		</q-table>
	</q-container>
</template>

<script setup lang="ts" name="apiV1SystemOnlineList">
	import { reactive, onMounted, ref } from 'vue';
	import { ElMessageBox, ElMessage, FormInstance } from 'element-plus';
	import { forceLogout, listSysUserOnline } from '/@/api/system/monitor/userOnline';
	import { getCardMode, getCurrentDate, parseTime } from '/@/utils/common';
	import { TableOnlineColumn, TableOnlinedModel } from '/@/types/model';

	const tableRef = ref();
	const state = reactive<TableOnlinedModel>({
		ids: [],
		total: 0,
		loading: false,
		tableData: [],
		param: {
			ipaddr: '',
			userName: '',
			pageNum: 1,
			pageSize: 10,
		},
	});

	const search = ref(
			[
				{ label: '登录IP', prop: 'ipaddr', placeholder: '请输入登录地址', required: false, type: 'input' },
				{ label: '用户名称', prop: 'userName', placeholder: '请输入用户名称', required: false, type: 'input' },
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
		listSysUserOnline(state.param).then((res) => {
			state.tableData = res.data.list ?? [];
			state.total = res.data.total;
		});
	};

	// 删除
	const onRowDel = (row: TableOnlineColumn) => {
		let msg = '确认要强制退出用户登录？';
		let ids: number[] = [];
		if (row) {
			msg = `将强制用户下线，是否继续?`;
			ids = [row.id];
		} else {
			ids = tableRef.value.getSelectlist('id');
		}
		if (ids.length === 0) {
			ElMessage.error('请选择要强制退出登录的用户。');
			return;
		}
		ElMessageBox.confirm(msg, '提示', {
			confirmButtonText: '确认',
			cancelButtonText: '取消',
			type: 'warning',
		})
			.then(() => {
				forceLogout(ids).then(() => {
					ElMessage.success('退出成功');
					getTableData();
				});
			})
			.catch(() => {});
	};

	const buttonVisible = (time: string) => {
		return getCurrentDate() === parseTime(time, '{y}-{m}-{d}');
	};

	// 页面加载时
	onMounted(() => {
		getTableData();
	});
</script>
