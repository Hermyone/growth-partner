<!--
 * @Descripttion: 选择用户对话框
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2024-08-16
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <DialogCommon ref="tableRef" v-model="showing" title="选择用户" width="1024" @onSelect="doSelect"
                  :API="proxy.$API.im.user_list"
                  :dict="tableDict.table_user_list"
                  :extendParams="params"
                  :selection="true"
                  respDataNode="userList"
    >
        <template #filter>
            <el-form :inline="true" :model="params" ref="queryRef" style="padding-left: 15px">
                <el-form-item label="部门" prop="deptId">
                    <el-cascader
                            :options="deptData"
                            :props="{ checkStrictly: true, emitPath: false, value: 'deptId', label: 'deptName' }"
                            placeholder="请选择部门"
                            clearable
                            class="w100"
                            v-model="params.deptId"
                    >
                        <template #default="{ node, data }">
                            <span>{{ data.deptName }}</span>
                            <span v-if="!node.isLeaf"> ({{ data.children.length }}) </span>
                        </template>
                    </el-cascader>
                </el-form-item>
                <el-form-item label="关键字" prop="keyWords">
                    <el-input
                            v-model="params.keyWords"
                            placeholder="请输入用户账号或姓名"
                            clearable
                            size="default"
                            style="width:200px;"
                            @clear="doUpdate"
                    />
                </el-form-item>
                <el-form-item>
                    <el-button icon="ele-Search" size="default" type="primary" @click="doUpdate">
                        查询
                    </el-button>
                    <el-button icon="ele-Refresh" size="default" @click="doReset(queryRef)">
                        重置
                    </el-button>
                </el-form-item>
            </el-form>
        </template>
        <template #columns>
            <el-table-column label="头像" width="60" fixed="left">
                <template #default="{row}">
                    <el-avatar fit="fill" :src="getHeadImageUrl(row.avatar)" size="small" v-if="row.avatar !== ''"></el-avatar>
                </template>
            </el-table-column>
            <el-table-column prop="sex" label="性别" width="80">
                <template #default="{row}">
                    <el-tag :color="selectDictColor(sys_user_sex, row.sex)">{{ selectDictLabel(sys_user_sex, row.sex) }} </el-tag>
                </template>
            </el-table-column>
            <el-table-column prop="dept.deptName" label="部门" width="120" show-overflow-tooltip></el-table-column>
            <el-table-column prop="roleInfo" label="角色" width="120" show-overflow-tooltip>
                <template #default="{row}">
                    <span v-for="(item, index) of row.roleInfo" :key="'role-' + index"> {{ item.name + '   ' }} </span>
                </template>
            </el-table-column>
        </template>
        <template #default>
			<span class="dialog-footer">
				<el-button @click="onCancel" class="dialog-footer-button">取消</el-button>
				<el-button type="primary" @click="onSubmit" class="dialog-footer-button">确认</el-button>
			</span>
        </template>
    </DialogCommon>
</template>

<script lang="ts">
	import { defineComponent, getCurrentInstance, onMounted, reactive, ref, toRefs } from 'vue';
    import { tableDict } from '/@/config/dicts/table';
	import DialogCommon from "./index.vue";
    import { selectDictColor, selectDictLabel } from '/@/utils/common';
	import { getDeptTree } from '/@/api/system/user/index';
    import { ElMessage, FormInstance } from 'element-plus';

	export default defineComponent({
		name: 'q-dialog-user',
        components: { DialogCommon },
		setup(props, { emit }) {
            const { proxy } = <any>getCurrentInstance();
            const { sys_user_sex } = proxy.useDict('sys_user_sex');
            const queryRef = ref();
            const tableRef = ref();
			const state = reactive({
                showing: false,
                currentRow: null,
                deptData: [],
                params: {
					deptId: '',
                    keyWords: '',
                }
			});

			const openDialog = () => {
                state.showing = true;
            };
            const closeDialog = () => {
                state.showing = false;
            };

            const doUpdate = () => {
                tableRef.value.getTableData();
            };

            const doReset = (formEl: FormInstance | undefined) => {
                if (!formEl) return;
                formEl.resetFields();
                doUpdate();
            };

            const doSelect = (val: any) => {
                state.currentRow =  val;
            };

            const getHeadImageUrl = (avatar : string) => {
                return proxy.getUpFileUrl(avatar);
            };

            const onCancel = () => {
                closeDialog();
            };

            const onSubmit = () => {
                let ids: number[] = [];
                ids = tableRef.value.getSelectlist('id');
            	if(ids.length > 0) {
                    emit("onSelectData", ids);
                    closeDialog();
                }else{
            		ElMessage.warning("没有选择用户");
                }
            };

            onMounted(() => {
                getDeptTree().then((res: any) => {
                    state.deptData = res.data.deps;
                })
            });

			return {
                proxy,
				queryRef,
                tableRef,
                sys_user_sex,
                tableDict,
                doSelect,
                onSubmit,
                onCancel,
                doUpdate,
                doReset,
                openDialog,
                closeDialog,
                getHeadImageUrl,
				...toRefs(state),
			}
		},
	});
</script>

<style scoped lang="scss">
</style>