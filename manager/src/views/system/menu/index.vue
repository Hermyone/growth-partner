<template>
	<q-container :cardMode="getCardMode()" class="back-container">
		<q-search :search="search" @search="onSearch" >
			<el-button icon="ele-Plus" size="default" type="success" round class="ml10" @click="onOpenAddMenu(null)" v-auth="'api/v1/system/menu/add'">
				新增菜单
			</el-button>
		</q-search>

		<q-table :recordCount="state.total" :data="state.menuTableData" row-key="path" :isSelection="false" :tree-props="{ children: 'children', hasChildren: 'hasChildren' }">
			<el-table-column prop="title" label="菜单名称" width="160" fixed show-overflow-tooltip />
			<el-table-column label="类型" show-overflow-tooltip width="80">
				<template #default="scope">
					<el-tag :type="scope.row.menuType === 0 ? 'danger' : scope.row.menuType === 1 ? 'success' : 'warning'" size="small">
						{{ scope.row.menuType === 0 ? '目录' : scope.row.menuType === 1 ? '菜单' : '按钮' }}
					</el-tag>
				</template>
			</el-table-column>
			<el-table-column prop="path" label="路由路径" show-overflow-tooltip></el-table-column>
			<el-table-column label="组件路径" show-overflow-tooltip>
				<template #default="scope">
					<span>{{ scope.row.component }}</span>
				</template>
			</el-table-column>
			<el-table-column label="api接口" show-overflow-tooltip>
				<template #default="scope">
					<span>{{ scope.row.name }}</span>
				</template>
			</el-table-column>
			<el-table-column label="排序" show-overflow-tooltip width="80">
				<template #default="scope">
					{{ scope.row.weigh }}
				</template>
			</el-table-column>
			<el-table-column label="图标" width="80" align="center">
				<template #default="scope">
					<SvgIcon :name="scope.row.icon" :size="20"/>
				</template>
			</el-table-column>
			<el-table-column prop="isHide" label="显示状态"  width="120" align="center">
				<template #default="scope">
					<el-tag type="success" v-if="scope.row.isHide === 0">显示</el-tag>
					<el-tag type="info" v-else>隐藏</el-tag>
				</template>
			</el-table-column>
			<el-table-column label="操作" width="180">
				<template #default="scope">
					<el-button
						v-if="scope.row.menuType !== 2"
						size="small"
						text
						type="primary"
						@click="onOpenAddMenu(scope.row)"
						v-auth="'api/v1/system/menu/add'"
						>
						新增
					</el-button>
					<el-divider direction="vertical" v-if="scope.row.menuType !== 2"/>
					<el-button size="small" text type="primary" @click="onOpenEditMenu(scope.row)" v-auth="'api/v1/system/menu/update'">修改</el-button>
					<el-divider direction="vertical" />
					<el-button size="small" text type="danger" @click="onTabelRowDel(scope.row)" v-auth="'api/v1/system/menu/delete'">删除</el-button>
				</template>
			</el-table-column>

			<template #empty>
				<el-empty description="暂无数据" />
			</template>
		</q-table>

		<EditMenu ref="editMenuRef" @menuList="menuList" :visibleOptions="sys_show_hide" :acType="acType" />
	</q-container>
</template>

<script setup lang="ts" name="apiV1SystemAuthMenuList">
	import { ref, reactive, onBeforeMount, defineAsyncComponent, getCurrentInstance, unref } from 'vue';
	import { ElMessageBox, ElMessage } from 'element-plus';
	import { delMenu, getMenuList } from '/@/api/system/menu';
	import { getCardMode } from '/@/utils/common';

	// 引入组件
	const EditMenu = defineAsyncComponent(() => import('/@/views/system/menu/component/editMenu.vue'));

	const editMenuRef = ref();
	const state = reactive({
		total: 0,
		queryParams: {
			title: '',
			component: '',
		},
		menuTableData: [],
	});
	const { proxy } = getCurrentInstance() as any;
	const { sys_show_hide } = proxy.useDict('sys_show_hide');
	const acType = ref('add');

	const search = ref(
			[
				{ label: '菜单名称', prop: 'title', placeholder: '请输入菜单名称', required: false, type: 'input' },
				{ label: '组件路径', prop: 'component', placeholder: '请输入组件路径', required: false, type: 'input' },
			]
	);

	// 搜索点击时表单回调
	const onSearch = (data: EmptyObjectType) => {
		state.queryParams = Object.assign({}, state.queryParams, { ...data });
		menuList();
	};

	// 打开新增菜单弹窗
	const onOpenAddMenu = (row: any) => {
		acType.value = 'add';
		editMenuRef.value.openDialog(row);
	};
	// 打开编辑菜单弹窗
	const onOpenEditMenu = (row: any) => {
		acType.value = 'edit';
		editMenuRef.value.openDialog(row);
	};
	// 删除当前行
	const onTabelRowDel = (row: any) => {
		ElMessageBox.confirm(`确认要删除菜单【${row.title}】吗?`, '提示', {
			confirmButtonText: '删除',
			cancelButtonText: '取消',
			type: 'warning',
		})
			.then(() => {
				delMenu(row.id).then(() => {
					ElMessage.success('删除成功');
					proxy.$refs['editMenuRef'].resetMenuSession();
					menuList();
				});
			})
			.catch(() => {});
	};
	const formatIsHide = (row: any) => {
		return proxy.selectDictLabel(unref(sys_show_hide), '' + row.isHide);
	};
	onBeforeMount(() => {
		menuList();
	});
	const handleQuery = () => {
		menuList();
	};
	const menuList = () => {
		getMenuList(state.queryParams).then((res) => {
			state.menuTableData = proxy.handleTree(res.data.rules ?? [], 'id', 'pid');
			state.total = res.data.rules.length;
		});
	};
</script>

<style scoped lang="scss">
	.back-container {
		:deep(.el-card__body) {
			padding-bottom: 0 !important;
		}
	}
</style>
