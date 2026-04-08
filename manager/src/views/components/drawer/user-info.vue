<template>
	<el-drawer v-model="state.isOpen" size="400px" append-to-body direction="rtl">
		<template #header>
			<div style="display: flex; align-items: center">
				<el-avatar fit="fill" :src="getHeadImageUrl(state.userData.avatar)"></el-avatar>
				<h4 style="padding-left: 5px">用户信息</h4>
			</div>
		</template>
		<el-main style="padding:10px; height: 100%">
      <div class="user-sign" v-if="state.userData.motto != ''">
        <div class="sign-arrow"></div>
        <span>{{state.userData.motto}}</span>
      </div>
			<el-descriptions :column="1" size="small" border>
				<el-descriptions-item label="用户编号">{{state.userData.userName}}</el-descriptions-item>
				<el-descriptions-item label="用户昵称">{{state.userData.userNickname}}</el-descriptions-item>
				<el-descriptions-item label="描述">{{state.userData.describe}}</el-descriptions-item>
				<el-descriptions-item label="性别">
					<q-tag v-if="state.userData.sex === 1">男</q-tag>
					<q-tag type="danger" v-if="state.userData.sex === 2">女</q-tag>
					<q-tag type="info" v-if="state.userData.sex === 0">保密</q-tag>
				</el-descriptions-item>
				<el-descriptions-item label="手机">{{state.userData.mobile}}</el-descriptions-item>
        <el-descriptions-item label="注册时间">{{state.userData.createdAt}}</el-descriptions-item>
			</el-descriptions>

      <el-button type="primary" plain round @click="onView()" :icon="Search" style="margin-top: 15px;">
        查看消息
      </el-button>
		</el-main>

    <MsgListDlg ref="msgRef" :is-group="0"/>
	</el-drawer>
</template>

<script setup lang="ts">
	import { getCurrentInstance, reactive, ref } from 'vue';
	import { getUserInfo } from '/@/api/system/user/index';
  import MsgListDlg from '/@/views/components/dialog/msg-dialog.vue';
  import {Search} from "@element-plus/icons-vue";

	const { proxy } = <any>getCurrentInstance();
  const msgRef = ref();
	const state = reactive({
		isOpen: false,
		params: {} as any,
		userData: {} as any,
		userid: '',
	});
	// 打开弹窗
	const openDialog = (row: any, userid: string) => {
		state.params = row;
		state.userid = userid;
		state.userData = {};
		laodData();
		state.isOpen = true;
	};
	// 关闭弹窗
	const closeDialog = () => {
		state.isOpen = false;
	};

  const onView = () => {
    msgRef.value?.openDialog(state.userid);
  }

	const laodData = () => {
		let userId;
		if(state.userid === undefined || state.userid === ''){
			userId = state.params.user_id || state.params.createdBy || state.params.updateBy || state.params.id || "";
		}else{
			userId = state.userid;
		}

		getUserInfo({id: userId}).then((resp: any) => {
			let data = resp.data.userList || [];
			if(data.length > 0) {
				state.userData = data[0]
			}
		});
	}

	const getHeadImageUrl = (avatar : string) => {
		return proxy.getUpFileUrl(avatar);
	};

	defineExpose({
		openDialog,
		closeDialog,
	})

</script>

<style scoped lang="scss">
	:deep(.el-descriptions__label.el-descriptions__cell.is-bordered-label){
		background: none;
    color: var(--el-text-color-disabled);
		font-weight: 400;
    border: none;
    border-bottom: 1px solid var(--el-border-color-lighter);
    width: 80px;
	}
  :deep(.el-descriptions__content.el-descriptions__cell.is-bordered-content){
    background: none;
    font-weight: 500;
    border: none;
    border-bottom: 1px solid var(--el-border-color-lighter);
  }
  .user-sign {
    min-height: 26px;
    border-radius: 5px;
    padding: 5px;
    line-height: 25px;
    background: var(--el-border-color-lighter);
    font-size: 12px;
    margin-bottom: 20px;
    position: relative;
    display: -webkit-box;
    -webkit-box-orient: vertical;

    span {
      color: var(--el-text-color-primary);
    }

    .sign-arrow {
      position: absolute;
      width: 0;
      height: 0;
      border-left: 5px solid transparent;
      border-right: 5px solid transparent;
      border-bottom: 10px solid var(--el-border-color-lighter);
      left: 15px;
      top: -6px;
    }
  }
</style>
