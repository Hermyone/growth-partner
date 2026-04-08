<template>
	<q-container :cardMode="getCardMode()">
		<q-search ref="searchRef" :search="tableDict.table_system_message.search"  @search="onSearch" >
			<el-button icon="ele-Plus" size="default" round type="success" class="ml10" @click="onAdd">
				新增
			</el-button>
		</q-search>

		<q-table-v2
				ref="tableRef"
				:isSelection="false"
				:data="state.tableData"
				:recordCount="state.total"
				:row-key="tableDict.table_system_message.rowKey"
				:columns="tableDict.table_system_message.columns"
				:tableName="tableDict.table_system_message.tableName"
				:globalTitle="tableDict.table_system_message.title"
				@pageChange="onTablePageChange"
				@keyLink="onView"
		>
			<template #kind="{row}">
				<el-tag v-if="row.kind === 0" effect="dark" >系统广播</el-tag>
				<el-tag v-if="row.kind === 1" effect="dark" type="warning">站内信</el-tag>
			</template>
			<template #level="{row}">
				<el-tag v-if="row.level === 0" >一般消息</el-tag>
				<el-tag v-if="row.level === 1" type="danger">重要消息</el-tag>
			</template>
			<template #issend="{row}">
				<el-tag v-if="row.issend === 0" >未发送</el-tag>
				<el-tag v-if="row.issend === 1" type="success">已发送</el-tag>
			</template>
			<template #opt="{row}">
				<el-button v-if="row.issend === 0" size="small" text type="primary" @click="onEdit(row)">修改</el-button>
				<el-divider v-if="row.issend === 0" direction="vertical" />
				<el-button v-if="row.issend === 0" size="small" text type="danger" @click="onDel(row)">删除</el-button>
				<el-divider v-if="row.issend === 0" direction="vertical" />
				<el-button v-if="row.issend === 0" size="small" icon="ele-Position" type="primary" @click="onSend(row)">发送</el-button>
			</template>
		</q-table-v2>

		<SelectDialog ref="userRef" @onSelectData="onSendUser"/>
		<EditDialog ref="dlgRef" @onSuccess="getTableData"/>
		<ContentPreview ref="previewRef" />
	</q-container>
</template>

<script setup lang="ts">
	import { reactive, onMounted, ref, getCurrentInstance } from 'vue';
	import { getCardMode } from '/@/utils/common';
	import { tableDict } from '/@/config/dicts/table';
	import { ElMessage, ElMessageBox } from 'element-plus';
	import EditDialog from './component/editdialog.vue';
	import ContentPreview from '/@/views/components/markdown/preview.vue';
	import SelectDialog from '/@/views/components/dialog-table/dialog-user.vue';
	import { getMessageList, deleteMessage, sendMessage } from '/@/api/system/message';

	const userRef = ref();
	const tableRef = ref();
	const previewRef = ref();
	const { proxy } = <any>getCurrentInstance();
	const dlgRef = ref();
	const state = reactive({
		total: 0,
		loading: false,
		tableData: [],
		currMsgid: 0,
		params: {
			id: "",
			pageNum: 1,
			pageSize: 10,
		},
	});

	// 搜索点击时表单回调
	const onSearch = (data: any) => {
		state.params = Object.assign({}, state.params, { ...data });
		tableRef.value?.pageReset();
	};

	// 分页改变时回调
	const onTablePageChange = (page: TablePageType) => {
		state.params.pageNum = page.pageNum;
		state.params.pageSize = page.pageSize;
		getTableData();
	};

	const getTableData = () => {
		state.loading = false;

		getMessageList(state.params).then((resp: any) => {
			state.tableData = resp.data.list || [];
			state.total = resp.data.total || state.tableData.length;
		})

		setTimeout(() => {
			state.loading = false;
		}, 1000);
	};

	const onView = (row: any) => {
		previewRef.value.openDialog(row.content);
	};

	const onAdd = () => {
		dlgRef.value.openDialog();
	};

	const onEdit = (row: any) => {
		dlgRef.value.openDialog(row);
	};

	const onSend = (row: any) => {
		if(row.kind === 0){
			ElMessageBox.confirm('确认要发布消息吗?', '警告', {
				confirmButtonText: '确认',
				cancelButtonText: '取消',
				type: 'warning',
			})
					.then(() => {
						sendMessage({id: row.id}).then((resp: any) => {
							ElMessage.success('消息发布成功');
							getTableData();
						});
					})
					.catch(function () {
					});
		}else{
			state.currMsgid = row.id;
			userRef.value.openDialog();
		}
	};

	const onSendUser = (val: any) => {
		if(!val || val.length == 0){
			ElMessage.warning("没有选择需要发送的用户");
			return
		}

		sendMessage({id: state.currMsgid, ids: val}).then((resp: any) => {
			ElMessage.success('消息发布成功');
			getTableData();
		});
	};

	const onDel = (row: any) => {
		ElMessageBox.confirm('确认要删除消息吗?', '警告', {
			confirmButtonText: '确认',
			cancelButtonText: '取消',
			type: 'warning',
		})
				.then(() => {
					deleteMessage(row.id)
							.then(() => {
								ElMessage.success('数据删除成功');
								getTableData();
							})
				})
				.catch(function () {
				});
	};

	// 页面加载时
	onMounted(() => {
		// getTableData();
	});

</script>

<style scoped lang="scss">
</style>