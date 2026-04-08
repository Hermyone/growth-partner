<template>
	<el-dialog v-model="isOpen" width="1024px" :close-on-click-modal="false" append-to-body>
		<template #header>
			<span>消息列表</span>
		</template>

    <el-container>
      <el-aside width="220px" v-if="isGroup == '0'">
        <div>
          <el-radio-group v-model="userWay" @change="userWayChanged">
            <el-radio-button value="0">TA的会话</el-radio-button>
            <el-radio-button value="1">TA的联系人</el-radio-button>
          </el-radio-group>

          <el-scrollbar style="margin-top: 15px;">
            <el-tree
                ref="treeRef"
                node-key="fromUser"
                :data="userList"
                :props="userProps"
                :expand-on-click-node="false"
                :highlight-current="true"
                @node-click="handleNodeClick"
            >
              <template #default="{node, data}">
                  <div class="user-item">
                    <el-avatar fit="fill" :src="data.avatar"></el-avatar>
                    <span>{{ data.displayName }}</span>
                  </div>
              </template>
            </el-tree>
          </el-scrollbar>
        </div>
      </el-aside>
      <el-main>
        <div style="padding-left: 15px;">
          <el-form :inline="true" :model="params" ref="formRef">
            <el-form-item label="关键字" prop="keyWords">
              <el-input
                  v-model="params.keyWords"
                  placeholder="输入查询关键字"
                  clearable
                  size="default"
                  style="width:200px;"
                  @keyup.enter.native="getTableData"
                  @clear="getTableData"
              />
            </el-form-item>
            <el-form-item label="是否已读" prop="is_read" style="width: 200px;">
              <el-select v-model="params.is_read">
                <el-option label="未读" value='0' />
                <el-option label="已读" value='1' />
              </el-select>
            </el-form-item>

            <el-form-item>
              <el-button icon="ele-Search" size="default" type="primary" round @click="getTableData">
                查询
              </el-button>
              <el-button icon="ele-Refresh" size="default" round @click="reset(formRef)">
                重置
              </el-button>
            </el-form-item>
          </el-form>
        </div>

        <el-tabs v-model="params.type" @tab-change="getTableData">
          <el-tab-pane label="所有消息" name="all"></el-tab-pane>
          <el-tab-pane label="文本" name="text"></el-tab-pane>
          <el-tab-pane label="图片" name="image"></el-tab-pane>
          <el-tab-pane label="音频" name="voice"></el-tab-pane>
          <el-tab-pane label="视频" name="video"></el-tab-pane>
          <el-tab-pane label="文件" name="file"></el-tab-pane>
        </el-tabs>

        <q-table-v2
            ref="tableRef"
            :isSelection="false"
            :data="tableData"
            :recordCount="total"
            :pageSize="10"
            :footerButton="false"
            :row-key="tableDict.table_im_message_ex.rowKey"
            :columns="tableDict.table_im_message_ex.columns"
            :tableName="tableDict.table_im_message_ex.tableName"
            :globalTitle="tableDict.table_im_message_ex.title"
            @pageChange="onTablePageChange"
        >
          <template #type="{row}">
            <el-tag effect="dark" :color="selectDictColor(msg_type, row.type)">{{ selectDictLabel(msg_type, row.type) }} </el-tag>
          </template>
          <template #is_read="{row}">
            <el-tag v-if="row.is_read == 0" type="danger">未读</el-tag>
            <el-tag v-else>已读</el-tag>
          </template>
          <template #content="{row}">
            <img v-if="row.type === 'image'" style="width: 64px;height: 64px;" :src="row.content">
            <span v-if="row.type === 'text'" v-html="row.content"></span>
          </template>
        </q-table-v2>
      </el-main>
    </el-container>
	</el-dialog>
</template>

<script lang="ts">
	import { reactive, toRefs, defineComponent, ref, unref, getCurrentInstance, onMounted } from 'vue';
  import {ElMessage, FormInstance} from 'element-plus';
  import {tableDict} from "/@/config/dicts/table";
  import {selectDictColor, selectDictLabel} from "/@/utils/common";

  const props = {
    isGroup: {
      type: Number,
      default: 0,
    },
    toUser:{
      type: String,
      default: '',
    }
  };

export default defineComponent({
  props: props,
	setup(props, { emit }) {
		const { proxy } = <any>getCurrentInstance();
    const formRef = ref();
    const treeRef = ref();
    const { msg_type } = proxy.useDict('msg_type');
		const state = reactive({
			loading: false,
			isOpen: false,
      total: 0,
      tableData: [],
      params: {
        id: "",
        type: "all",
        is_manager: "Y",
        keyWords: "",
        is_read: "",
        is_group: props.isGroup,
        toContactId: props.toUser,
        pageNum: 1,
        pageSize: 10,
      },
      userWay: 0,
      oldUser: '',
      userList: [] as any,
      userProps: {
        id: 'fromUser',
        children: 'children',
        label: 'displayName',
      }
		});

		// 打开弹窗
		const openDialog = (to_user: string) => {
      state.oldUser = to_user;
      state.params.toContactId = to_user;
      state.params.type = "all";
			state.isOpen = true;
      state.userWay = 0;

      if(props.isGroup == 0){
        //查好友
        proxy.$API.im.msg.user.post({user_id: to_user}).then((res: any) => {
          state.userList = res.data.list ?? [];
          if(state.userList.length > 0){
            setTimeout(() => {
              state.params.toContactId = state.userList[0].fromUser;
              treeRef.value?.setCurrentKey(state.userList[0].fromUser);
              getTableData();
            }, 500)
          }
        });
      }else{
        getTableData();
      }
		};

		const closeDialog = () => {
			state.isOpen = false;
		};

    // 分页改变时回调
    const onTablePageChange = (page: any) => {
      state.params.pageNum = page.pageNum;
      state.params.pageSize = page.pageSize;
      getTableData();
    }

    const reset = (formEl: FormInstance | undefined) => {
      if (!formEl) return;
      formEl.resetFields();
      getTableData();
    };

    const handleNodeClick = (data: any) => {
      state.params.toContactId = data.fromUser;
      getTableData();
    };

    const getTableData = () => {
      state.loading = true;

      proxy.$API.im.msg.query.post(state.params).then((resp: any) => {
        state.tableData = resp.data.list || [];
        state.total = resp.data.total || state.tableData.length;
      })

      setTimeout(() => {
        state.loading = false;
      }, 1000);
    }

    const userWayChanged = () => {
      state.userList = [];
      if(state.userWay == 0) {
        // 查会话用户
        proxy.$API.im.msg.user.post({user_id: state.oldUser}).then((res: any) => {
          state.userList = res.data.list ?? [];
        });
      }else{
        // 查联系人
        proxy.$API.im.friend.friends.post({is_invite: 1, user_id: state.oldUser, pageNum:1, pageSize: 9999}).then((res: any) => {
          let list = res.data.list ?? [];
          list.forEach((item:any) => {
            state.userList.push({
              fromUser: item.friend_user_id,
              displayName: item.displayName,
              avatar: item.avatar,
            })
          })
        });
      }

      if(state.userList.length > 0){
        setTimeout(() => {
          state.params.toContactId = state.userList[0].fromUser;
          treeRef.value?.setCurrentKey(state.userList[0].fromUser);
          getTableData();
        }, 500)
      }
    }

		onMounted(() => {
		});

		return {
			openDialog,
      closeDialog,
      onTablePageChange,
      getTableData,
      handleNodeClick,
      userWayChanged,
      reset,
      tableDict,
      msg_type,
      formRef,
      treeRef,
			...toRefs(state),
		};
	},
});
</script>

<style scoped lang="scss">
:deep(.el-tabs) {
  height: 40px;
}
.el-main {
  padding: 0;
}
.user-item {
  display: flex;
  align-items: center;
  gap: 15px;
}
:deep(.el-tree-node__content) {
  line-height: 50px;
  height: 50px;
}
:deep(.el-tabs__item) {
  width: 100px;
}
:deep(.el-tree-node__expand-icon.is-leaf) {
  width: 0 !important;
}
</style>
