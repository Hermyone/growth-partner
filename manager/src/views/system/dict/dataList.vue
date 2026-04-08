<template>
	<q-container :cardMode="getCardMode()">
		<q-search ref="searchRef" :search="tableDict.tableDictList.search"  @search="onSearch" >
			<el-button icon="ele-Plus" size="default" type="success" class="ml10" @click="onOpenAddDic">
				新增字典
			</el-button>
			<el-button icon="ele-Delete" size="default" type="danger" class="ml10" @click="onRowDel(null)">
				删除字典
			</el-button>
		</q-search>

		<q-table-v2
				 ref="tableRef"
				 :data="state.tableData"
				 :recordCount="state.total"
				 :row-key="tableDict.tableDictList.rowKey"
				 :columns="tableDict.tableDictList.columns"
				 :tableName="tableDict.tableDictList.tableName"
				 :globalTitle="tableDict.tableDictList.title"
				 @pageChange="onTablePageChange"
		>
			<template #status="scope">
				<el-tag type="success" v-if="scope.row.status">启用</el-tag>
				<el-tag type="info" v-else>禁用</el-tag>
			</template>
			<template #opt="{row}">
				<el-button size="small" text type="primary" @click="onOpenEditDic(row)">修改</el-button>
				<el-divider direction="vertical" />
				<el-button size="small" text type="danger" @click="onRowDel(row)">删除</el-button>
			</template>
		</q-table-v2>

		<EditDic ref="editDicRef" @dataList="getTableData" :dict-type="state.param.dictType" />
	</q-container>
</template>

<script setup lang="ts" name="apiV1SystemDictDataList">
	import { reactive, onMounted, ref, defineAsyncComponent, nextTick } from 'vue';
	import { ElMessageBox, ElMessage, FormInstance } from 'element-plus';
	import { getDataList, deleteData } from '/@/api/system/dict/data';
	import { useRoute } from 'vue-router';
	import { getCardMode } from '/@/utils/common';
	import { TableDictTypeColumn, TableDictTypeModel } from '/@/types/model';
	import { tableDict } from '/@/config/dicts/table';

	// 引用组件
	const EditDic = defineAsyncComponent(() => import('/@/views/system/dict/component/editDicData.vue'));

	const route = useRoute();
	const tableRef = ref();
	const editDicRef = ref();
	const searchRef = ref();
	const state = reactive<TableDictTypeModel>({
		ids: [],
		total: 0,
		loading: false,
		tableData: [],
		param: {
			pageNum: 1,
			pageSize: 10,
			dictLabel: '',
			dictType: '',
			status: '',
		},
	});

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
		getDataList(state.param).then((res: any) => {
			state.tableData = res.data.list ?? [];
			state.total = res.data.total;
		});
	};
	// 打开新增字典弹窗
	const onOpenAddDic = () => {
		editDicRef.value.openDialog();
	};
	// 打开修改字典弹窗
	const onOpenEditDic = (row: TableDictTypeColumn) => {
		editDicRef.value.openDialog(row);
	};
	// 删除字典
	const onRowDel = (row: TableDictTypeColumn) => {
		let msg = '确认要删除所选数据？';
		let ids: number[] = [];
		if (row) {
			msg = `确认要删除字典【${row.dictLabel}】吗?`;
			ids = [row.dictCode];
		} else {
			ids = tableRef.value.getSelectlist('dictCode');
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
				deleteData(ids).then(() => {
					ElMessage.success('删除成功');
					getTableData();
				});
			})
			.catch(() => {});
	};
	// 页面加载时
	onMounted(() => {
		const dictType = route.params && route.params.dictType;
		state.param.dictType = <string>dictType;
		getTableData();

		nextTick(() => {
			searchRef.value.setDefaultValue('dictType', state.param.dictType);

			// 根据元素ID给元素赋值
			// const dict_type : any = document.getElementById("dictType");
			// if (dict_type) {
			// 	dict_type.value = state.param.dictType;
			// 	searchRef.value.$forceUpdate();
			// };
		});
	});

</script>
