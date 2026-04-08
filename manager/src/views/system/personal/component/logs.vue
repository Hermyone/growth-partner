<template>
	<el-alert title="最近7天的登录和操作日志" type="info" show-icon style="margin-bottom: 15px;"/>
<!--	<h4>最近7天的登录和操作日志</h4>-->
	<el-tabs v-model="activePage" class="personal-logs-tabs" @tab-change="handleClick">
		<el-tab-pane label="登录日志" name="login" >
			<q-table
					:data="state.login.tableData"
					:recordCount="state.login.total"
					:isSelection="false"
					:pageSize="10"
					size="small"
					@pageChange="onLoginTablePageChange"
			>
				<el-table-column label="编号" sortable prop="infoId" />
<!--				<el-table-column label="登录地址" prop="ipaddr" width="130" show-overflow-tooltip />-->
				<el-table-column label="登录地点" prop="loginLocation" min-width="120" show-overflow-tooltip sortable/>
				<el-table-column label="浏览器" prop="browser" min-width="100" sortable/>
				<el-table-column label="操作系统" min-width="120" prop="os" sortable/>
				<el-table-column label="登录状态" align="center" prop="status" min-width="120" :formatter="statusFormat" sortable>
					<template #default="scope">
						<el-tag type="success" v-if="scope.row.status === 1">成功</el-tag>
						<el-tag type="danger" v-else>失败</el-tag>
					</template>
				</el-table-column>
				<el-table-column label="登录日期" align="center" prop="loginTime" width="180" sortable/>
			</q-table>
		</el-tab-pane>

		<el-tab-pane label="操作日志" name="option" >
			<q-table
					ref="tableRef"
					:data="state.oper.tableData"
					:recordCount="state.oper.total"
					:isSelection="false"
					:pageSize="10"
					size="small"
					@pageChange="onOperTablePageChange"
			>
				<el-table-column label="编号" sortable prop="operId" min-width="100px" />
				<el-table-column label="系统模块" prop="title" min-width="100px" sortable/>
				<el-table-column label="请求方式" prop="requestMethod" min-width="100px" sortable>
					<template #default="scope">
						<q-tag v-if="scope.row.requestMethod === 'GET'">{{ requestMethodFormat(scope.row) }} </q-tag>
						<q-tag type="success" v-if="scope.row.requestMethod === 'POST'">{{ requestMethodFormat(scope.row) }} </q-tag>
						<q-tag type="warning" v-if="scope.row.requestMethod === 'PUT'">{{ requestMethodFormat(scope.row) }} </q-tag>
						<q-tag type="danger" v-if="scope.row.requestMethod === 'DELETE'">{{ requestMethodFormat(scope.row) }} </q-tag>
					</template>	
				</el-table-column>
<!--				<el-table-column label="操作人员" prop="operName" min-width="100px" />-->
<!--				<el-table-column label="部门名称" prop="deptName" min-width="150px" />-->
				<el-table-column label="请求URL" prop="operUrl" min-width="200px" show-overflow-tooltip  sortable/>
<!--				<el-table-column label="主机地址" prop="operIp" min-width="100px" />-->
<!--				<el-table-column label="操作地点" prop="operLocation" min-width="100px" />-->
				<el-table-column label="操作时间" align="center" prop="operTime" min-width="180px" sortable>
					<template #default="scope">
						<span>{{ proxy.parseTime(scope.row.operTime, '{y}-{m}-{d} {h}:{i}:{s}') }}</span>
					</template>
				</el-table-column>
				<el-table-column label="操作" align="center" class-name="small-padding" width="80px" fixed="right">
					<template #default="scope">
						<el-button size="small" type="primary" text @click="handleView(scope.row)" v-auth="'api/v1/system/sysOperLog/view'">详情</el-button>
					</template>
				</el-table-column>
			</q-table>

			<apiV1SystemSysOperLogDetail
					ref="detailRef"
					:requestMethodOptions="sys_oper_log_type"
			></apiV1SystemSysOperLogDetail>
		</el-tab-pane>
	</el-tabs>
</template>

<script setup lang="ts" name="personal-logs">
	import { reactive, onMounted, ref, unref, getCurrentInstance, toRaw, defineAsyncComponent } from 'vue';
	import type { TabPaneName } from 'element-plus'
	import { TableLoginLogsColumn, TableOperLogsColumn } from '/@/types/model';
	import { logList } from '/@/api/system/monitor/loginLog';
	import { listSysOperLog } from '/@/api/system/monitor/operLog';
  import { formatDate } from "/@/utils/formatTime";

  const apiV1SystemSysOperLogDetail = defineAsyncComponent(() => import('/@/views/system/monitor/operLog/component/detail.vue'));

	const detailRef = ref();
	const { proxy } = <any>getCurrentInstance();
	const activePage = ref('login');
	const state = reactive({
		login: {
			total: 0,
			tableData: [],
			param: {
				pageNum: 1,
				pageSize: 10,
				dateRange: [] as string[],
			}
		},
		oper: {
			total: 0,
			tableData: [],
			param: {
				pageNum: 1,
				pageSize: 10,
				dateRange: [] as string[],
			}
		},
	});
	// 字典选项数据
	const { sys_oper_log_type } = proxy.useDict('sys_oper_log_type');
	const { admin_login_status } = proxy.useDict('admin_login_status');

	const onOperTablePageChange = (page: TablePageType) => {
		state.oper.param.pageNum = page.pageNum;
		state.oper.param.pageSize = page.pageSize;
		getOperTableData();
	};

	// 请求方式字典翻译
	const requestMethodFormat = (row: TableOperLogsColumn) => {
		return proxy.selectDictLabel(sys_oper_log_type.value, row.requestMethod);
	};

	const getOperTableData = () => {
		listSysOperLog(state.oper.param).then((res: any) => {
			state.oper.tableData = res.data.list ?? [];
			state.oper.total = res.data.total;
		});
	};

	const onLoginTablePageChange = (page: TablePageType) => {
		state.login.param.pageNum = page.pageNum;
		state.login.param.pageSize = page.pageSize;
		getLoginTableData();
	};

	// 登录状态字典翻译
	const statusFormat = (row: TableLoginLogsColumn) => {
		return proxy.selectDictLabel(unref(admin_login_status), row.status);
	};

	const getLoginTableData = () => {
		logList(state.login.param).then((res: any) => {
			state.login.tableData =  res.data.list ?? [];
			state.login.total = res.data.total;
		});
	};

	const handleClick = (tab: TabPaneName) => {
		if(tab === 'login') getLoginTableData()
		else getOperTableData();
	};

	const handleView = (row: TableOperLogsColumn) => {
		detailRef.value.openDialog(toRaw(row));
	};

	// 页面加载时
	onMounted(() => {
		// 查最近7天的数据
		var date = new Date();
		var edate = formatDate(date, 'YYYY-mm-dd');
		date.setDate(date.getDate()-7);
		var sdate = formatDate(date, 'YYYY-mm-dd');

		state.login.param.dateRange.push(sdate);
		state.login.param.dateRange.push(edate);
		state.oper.param.dateRange.push(sdate);
		state.oper.param.dateRange.push(edate);

		getLoginTableData();
	});
</script>

<style scoped lang="scss">
	.personal-logs-tabs {
		:deep(.el-tabs__item) {
			width: 100px;
		}
	}
</style>

