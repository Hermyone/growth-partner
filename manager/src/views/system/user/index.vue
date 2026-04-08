<template>
	<q-container :cardMode="getCardMode()" class="system-user-container">
		<el-container>
			<el-aside width="250px">
				<el-container>
					<el-header>
						<el-input placeholder="请输入部门名称" v-model="filterText" clearable></el-input>
					</el-header>
					<el-main class="nopadding">
						<el-scrollbar>
							<el-tree
									ref="treeRef"
									node-key="deptId"
									class="folder-tree menu"
									:data="state.deptData"
									:props="state.deptProps"
									:expand-on-click-node="false"
									:highlight-current="true"
									:filter-node-method="deptFilterNode"
									:default-expanded-keys="state.expandedKeys"
									@node-click="handleNodeClick"
							/>
						</el-scrollbar>
					</el-main>
				</el-container>
			</el-aside>
			<el-container>
				<el-header class="system-user-container-header">
					<el-form :inline="true" :model="state.param" ref="queryRef">
						<el-form-item label="关键字" prop="keyWords">
							<el-input
									v-model="state.param.keyWords"
									placeholder="请输入用户账号或姓名"
									clearable
									size="default"
									style="width:200px;"
									@keyup.enter.native="userList"
									@clear="userList"
							/>
						</el-form-item>
						<el-form-item>
							<el-button icon="ele-Search" size="default" type="primary" round @click="userList">
								查询
							</el-button>
							<el-button icon="ele-Refresh" size="default" round @click="resetQuery(queryRef)">
								重置
							</el-button>
							<el-button icon="ele-Plus" size="default" round type="success" @click="onOpenAddUser">
								新增用户
							</el-button>
							<el-button icon="ele-Delete" size="default" round type="danger" @click="onRowDel(null)">
								删除用户
							</el-button>
						</el-form-item>
					</el-form>
				</el-header>

				<el-main style="padding: 0 0 15px 0">
					<q-table
							ref="tableRef"
							row-key="userName"
							style="height: 100%"
							:data="state.tableData"
							:recordCount="state.total"
							:isSerialNo="false"
							@pageChange="onTablePageChange"
					>
						<el-table-column label="头像" width="60">
							<template #default="scope">
								<el-avatar fit="fill" :src="getHeadImageUrl(scope.row.avatar)" size="small" v-if="scope.row.avatar !== ''"></el-avatar>
							</template>
						</el-table-column>
						<el-table-column prop="userName" label="账户名称" width="120" show-overflow-tooltip></el-table-column>
						<el-table-column prop="userNickname" label="用户昵称" width="150" show-overflow-tooltip></el-table-column>
						<el-table-column prop="sex" label="性别" width="80">
							<template #default="scope">
								<q-tag v-if="scope.row.sex === 1">{{ sexFormat(scope.row) }} </q-tag>
								<q-tag type="danger" v-if="scope.row.sex === 2">{{ sexFormat(scope.row) }} </q-tag>
								<q-tag type="info" v-if="scope.row.sex === 0">{{ sexFormat(scope.row) }} </q-tag>
							</template>							
						</el-table-column>
						<el-table-column prop="mobile" label="手机号" width="120" show-overflow-tooltip></el-table-column>
						<el-table-column prop="userStatus" label="用户状态" show-overflow-tooltip>
							<template #default="scope">
								<el-switch
										v-if="scope.row.userName !=='admin'"
										v-model="scope.row.userStatus"
										inline-prompt
										:active-value="1"
										:inactive-value="0"
										active-text="启"
										inactive-text="禁"
										@change="handleStatusChange(scope.row)"
								>
								</el-switch>
								<el-tag v-else>启用</el-tag>
							</template>
						</el-table-column>
						<el-table-column prop="dept.deptName" label="部门" width="120" show-overflow-tooltip></el-table-column>
						<el-table-column label="角色" align="left" prop="roleInfo" width="120" show-overflow-tooltip>
							<template #default="scope">
								<span v-for="(item, index) of scope.row.roleInfo" :key="'role-' + index"> {{ item.name + '   ' }} </span>
							</template>
						</el-table-column>

						<el-table-column prop="createdAt" label="创建时间" width="180" show-overflow-tooltip></el-table-column>
						<el-table-column label="操作" align="center" width="180" fixed="right">
							<template #default="scope">
								<div v-if="scope.row.userName !=='admin'">									
									<el-button size="small" text type="primary" @click="onOpenEditUser(scope.row)">修改</el-button>
									<el-divider direction="vertical" />
									<el-button size="small" text type="danger" @click="onRowDel(scope.row)">删除</el-button>
									<el-divider direction="vertical" />
									<el-button size="small" text type="primary" @click="handleResetPwd(scope.row)">重置</el-button>
								</div>
							</template>
						</el-table-column>
					</q-table>
				</el-main>
			</el-container>
		</el-container>
		<EditUser ref="editUserRef" :dept-data="state.deptData" :gender-data="sys_user_sex" @getUserList="userList" />
	</q-container>
</template>

<script setup lang="ts" name="systemUser">
	import { reactive, onMounted, ref, defineAsyncComponent, watch, getCurrentInstance } from 'vue';
	import { ElMessageBox, ElMessage, ElTree, FormInstance } from 'element-plus';
	import { getUserList, getDeptTree, resetUserPwd, changeUserStatus, deleteUser } from '/@/api/system/user/index';
	import { TableUserModel } from '/@/types/model';
	import { getCardMode } from '/@/utils/common';

	// 引入组件
	const EditUser = defineAsyncComponent(() => import('/@/views/system/user/component/editUser.vue'));

	const { proxy } = <any>getCurrentInstance();
	const { sys_user_sex } = proxy.useDict('sys_user_sex');
	const tableRef = ref();
	const editUserRef = ref();
	const queryRef = ref();
	const filterText = ref('');
	const treeRef = ref<InstanceType<typeof ElTree>>();

	const state = reactive({
		ids: [],
		deptProps: {
			id: 'deptId',
			children: 'children',
			label: 'deptName',
		},
		deptData: [],
		expandedKeys: [] as any,
		total: 0,
		loading: false,
		tableData: [],
		param: {
			pageNum: 1,
			pageSize: 10,
			deptId: '',
			mobile: '',
			status: '',
			keyWords: '',
			dateRange: [],
		},
	});

	// 分页改变时回调
	const onTablePageChange = (page: TablePageType) => {
		state.param.pageNum = page.pageNum;
		state.param.pageSize = page.pageSize;
		userList();
	};

	// 初始化表格数据
	const initTableData = () => {
		getDeptTree().then((res: any) => {
			state.deptData = res.data.deps;

			if(state.deptData.length > 0) {
				setTimeout(() => {
					const firstNode: any = state.deptData[0];
					const firstNodeId = firstNode.deptId;
					const firstNodeParentId = firstNode.parent || null;

					// 设置默认展开的节点
					state.expandedKeys = [firstNodeId, firstNodeParentId];

					handleNodeClick(firstNode);
					treeRef.value?.setCurrentKey(firstNode.deptId);
				}, 500)
			}
		});
		userList();
	};
	const userList = () => {
		getUserList(state.param).then((res: any) => {
			state.tableData = res.data.userList ?? [];
			state.total = res.data.total;
		});
	};
	// 打开新增用户弹窗
	const onOpenAddUser = () => {
		editUserRef.value.openDialog();
	};
	// 打开修改用户弹窗
	const onOpenEditUser = (row: any) => {
		editUserRef.value.openDialog(row);
	};
	// 删除用户
	const onRowDel = (row: any) => {
		let msg = '确认要删除所选用户？';
		let ids: number[] = [];
		if (row) {
			msg = `确认要删除用户【${row.userName}】吗?`;
			ids = [row.id];
		} else {
			ids = tableRef.value.getSelectlist('id');
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
				deleteUser(ids).then(() => {
					ElMessage.success('删除成功');
					userList();
				});
			})
			.catch(() => {});
	};

	// 页面加载时
	onMounted(() => {
		initTableData();
	});

	watch(filterText, (val) => {
		treeRef.value!.filter(val);
	});
	const deptFilterNode = (value: string, data: any) => {
		if (!value) return true;
		return data.deptName.includes(value);
	};

	// 节点单击事件
	const handleNodeClick = (data: any) => {
		state.param.deptId = data.deptId;
		userList();
	};
	// 重置密码按钮操作
	const handleResetPwd = (row: any) => {
		ElMessageBox.prompt('请输入"' + row.userNickname + '"的新密码', '提示', {
			confirmButtonText: '确认',
			cancelButtonText: '取消',
		})
			.then(({ value }) => {
				if (!value || value == '') {
					ElMessage.success('密码不能为空');
					return;
				}
				resetUserPwd(row.id, value).then(() => {
					ElMessage.success('修改成功，新密码是：' + value);
				});
			})
			.catch(() => {});
	};
	// 用户状态修改
	const handleStatusChange = (row: any) => {
		let text = row.userStatus === 1 ? '启用' : '停用';
		ElMessageBox.confirm('确认要"' + text + '"："' + row.userName + '"用户吗?', '警告', {
			confirmButtonText: '确认',
			cancelButtonText: '取消',
			type: 'warning',
		})
			.then(function () {
				return changeUserStatus(row.id, row.userStatus);
			})
			.then(() => {
				ElMessage.success(text + '成功');
			})
			.catch(function () {
				row.userStatus = row.userStatus === 0 ? 1 : 0;
			});
	};
	// 重置按钮操作
	const resetQuery = (formEl: FormInstance | undefined) => {
		if (!formEl) return;
		formEl.resetFields();
		userList();
	};

	const sexFormat = (row : any) => {
		return proxy.selectDictLabel(sys_user_sex.value, row.sex);
	};

	const getHeadImageUrl = (avatar : string) => {
		return proxy.getUpFileUrl(avatar);
	};

</script>

<style scoped lang="scss">
	.system-user-container{
		:deep(.el-card__body) {
			padding: 0 !important;
		}
		&-header{
			align-items: flex-start;
		}
	}
</style>
