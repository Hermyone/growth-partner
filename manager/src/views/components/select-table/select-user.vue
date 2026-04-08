<!--
 * @Descripttion: 选择用户
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2024-11-14
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <q-select-table
            class="select-table"
            v-bind="$attrs"
            clearable
            :tableData="tableData"
            :total="total"
            :onPageChange="onTablePageChange"
            :tableWidth="650"
            labelFieldName="userNickname"
            keyFieldName="id"
            @change="change"
    >
        <template #header>
            <el-form :inline="true" :model="params">
                <el-form-item prop="keyWords">
                    <el-input
                            v-model="params.keyWords"
                            placeholder="请输入查询关键字"
                            clearable
                            size="small"
                            style="width:200px;"
                            @keyup.enter.native="getTableData"
                            @clear="getTableData"
                    />
                </el-form-item>
                <el-form-item>
                    <el-button icon="ele-Search" size="small" type="primary" @click="getTableData">
                        查询
                    </el-button>
                    <el-button icon="ele-Refresh" size="small" @click="resetQuery">
                        重置
                    </el-button>
                </el-form-item>
            </el-form>
        </template>
        <el-table-column label="头像" width="50">
            <template #default="{row}">
                <el-avatar fit="fill" :src="getHeadImageUrl(row.avatar)" size="small" v-if="row.avatar !== ''"></el-avatar>
            </template>
        </el-table-column>
        <el-table-column prop="userName" label="账户名称" width="100" show-overflow-tooltip></el-table-column>
        <el-table-column prop="userNickname" label="用户昵称" width="150" show-overflow-tooltip></el-table-column>
        <el-table-column prop="sex" label="性别" width="60">
            <template #default="{row}">
                <el-tag :color="selectDictColor(sys_user_sex, row.sex)">{{ selectDictLabel(sys_user_sex, row.sex) }} </el-tag>
            </template>
        </el-table-column>
        <el-table-column prop="mobile" label="手机号" width="100" show-overflow-tooltip></el-table-column>
    </q-select-table>
</template>

<script lang="ts">
	import { defineComponent, getCurrentInstance, onMounted, reactive, ref, toRefs } from 'vue';
  import { selectDictLabel, selectDictColor } from '/@/utils/common';

	const props = {
	};

	export default defineComponent({
		name: 'q-select-table-firm',
		props: props,
		setup(props, { emit }) {
            const { proxy } = <any>getCurrentInstance();
            const { sys_user_sex } = proxy.useDict('sys_user_sex');
			const state = reactive({
                loading: false,
                total: 0,
                deptProps: {
                    id: 'deptId',
                    children: 'children',
                    label: 'deptName',
                },
                expandedKeys: [] as any,
                deptData: [],
                tableData: [],
                params: {
                    deptId: '',
                    pageNum: 1,
                    pageSize: 10,
                    keyWords: '',
                },
			});

      const onTablePageChange = (page: TablePageType) => {
          state.params.pageNum = page.pageNum;
          state.params.pageSize = page.pageSize;
          getTableData();
      };

      const getTableData = () => {
          state.loading = false;

          proxy.$API.im.user_list.get(state.params).then((resp: any) => {
              state.tableData = resp.data.userList || [];
              state.total = resp.data.total || state.tableData.length;
          })

          setTimeout(() => {
              state.loading = false;
          }, 1000);
      };

      const resetQuery = () => {
          state.params.keyWords = '';
          getTableData();
      };

      const handleNodeClick = (data: any) => {
          state.params.deptId = data.deptId;
          getTableData();
      };

      const getHeadImageUrl = (avatar : string) => {
          return proxy.getUpFileUrl(avatar);
      };

      const change = (val: any) => {
        console.log("q-select-table:change", val);
      };

      onMounted(() => {
      })

			return {
        getTableData,
        onTablePageChange,
        resetQuery,
        getHeadImageUrl,
        sys_user_sex,
        handleNodeClick,
        change,
				...toRefs(state),
			}
		},
	});

</script>

<style scoped lang="scss">
    .select-table {
        width: 100%;
    }
    .select-dept {
        width: 200px;
        border: 1px solid var(--el-border-color);
    }
</style>