<template>
	<el-tabs v-model="activePage" class="personal-user-tabs">
		<el-tab-pane label="基本信息" name="first" >
			<el-form :model="state.personalForm" size="default" label-width="40px" class="mt20 mb35">
				<el-row :gutter="20">
          <el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4" class="mb20">
            <el-form-item label="账号">
              <el-tag>{{ state.personalForm.username }}</el-tag>
            </el-form-item>
          </el-col>
					<el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4" class="mb20">
						<el-form-item label="昵称">
							<el-input v-model="state.personalForm.nickname" placeholder="请输入昵称" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4" class="mb20">
						<el-form-item label="邮箱">
							<el-input v-model="state.personalForm.userEmail" placeholder="请输入邮箱" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4" class="mb20">
						<el-form-item label="签名">
							<el-input v-model="state.personalForm.describe" placeholder="请输入签名" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4" class="mb20">
						<el-form-item label="职业">
							<el-select v-model="state.personalForm.remark" placeholder="请选择职业" clearable class="w100">
								<el-option v-for="item in CareerInfo" :key="item.value" :label="item.label" :value="item.value"> </el-option>
							</el-select>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4" class="mb20">
						<el-form-item label="手机">
							<el-input v-model="state.personalForm.mobile" placeholder="请输入手机" clearable></el-input>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4" class="mb20">
						<el-form-item label="性别">
							<el-select v-model="state.personalForm.sex" placeholder="请选择性别" clearable class="w100">
                <el-option label="保密" value="0"></el-option>
								<el-option label="男" value="1"></el-option>
								<el-option label="女" value="2"></el-option>
							</el-select>
						</el-form-item>
					</el-col>
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
						<el-form-item>
							<el-button icon="ele-Select" type="primary" @click="handleUpload">
								更新个人信息
							</el-button>
						</el-form-item>
					</el-col>
				</el-row>
			</el-form>
		</el-tab-pane>
	</el-tabs>

	<el-card shadow="hover" class="mt15" header="个人信息">
		<div class="personal-user-card">
			<el-row>
				<el-col :span="24" class="personal-title mb18">
					{{ state.personalForm.nickname }} {{ state.personalForm.describe }}
					<el-tag class="ml10" type="info" closable v-for="tag in state.tags" :key="tag" @close="handleClose(tag)">
						{{ tag }}
					</el-tag>
					<el-input
							v-if="inputVisible"
							ref="inputRef"
							v-model="inputValue"
							class="ml10 w20"
							size="small"
							style="width: 80px;  height: 23px"
							@keyup.enter="handleInputConfirm"
							@blur="handleInputConfirm"
					/>
					<el-button v-else class="button-new-tag ml10" size="small" icon="ele-Plus" @click="showInput">
						添加
					</el-button>
				</el-col>
				<el-col :span="24">
					<el-row>
						<el-col :xs="24" :sm="8" class="personal-item mb6">
							<div class="personal-item-label">昵称：</div>
							<div class="personal-item-value">{{ state.personalForm.nickname }}</div>
						</el-col>
						<el-col :xs="24" :sm="16" class="personal-item mb6">
							<div class="personal-item-label">所属角色：</div>
							<div class="personal-item-value">{{ state.roles.join(',') }}</div>
						</el-col>
					</el-row>
				</el-col>
				<el-col :span="24">
					<el-row>
						<el-col :xs="24" :sm="8" class="personal-item mb6">
							<div class="personal-item-label">登录IP：</div>
							<div class="personal-item-value">{{ state.personalForm.lastLoginIp }}</div>
						</el-col>
						<el-col :xs="24" :sm="16" class="personal-item mb6">
							<div class="personal-item-label">登录时间：</div>
							<div class="personal-item-value">{{ state.personalForm.lastLoginTime }}</div>
						</el-col>
					</el-row>
				</el-col>
			</el-row>
		</div>
	</el-card>
</template>

<script setup lang="ts" name="personal-account">
	import { reactive, ref, onMounted, getCurrentInstance, nextTick } from 'vue';
	import { getPersonalInfo, editPersonal } from '/@/api/system/personal';
	import { Session } from '/@/utils/storage';
	import { ElMessage, ElInput } from 'element-plus';
	import { useUserInfo } from '/@/stores/userInfo';
	import { CareerInfo, recommendList } from '../mock';
	import Cookies from 'js-cookie';

	const { proxy } = <any>getCurrentInstance();
	const inputValue = ref('');
	const inputVisible = ref(false);
	const inputRef = ref<InstanceType<typeof ElInput>>()
	const activePage = ref('first')
	const state = reactive({
		deptName: '',
		recommendList,
		roles: [],
		personalForm: {
      username: '',
			nickname: '',
			userEmail: '',
			describe: '',
			mobile: '',
			sex: '',
			remark: '',
			// avatar: '',
			lastLoginIp: '',
			lastLoginTime: '',
			tags: '',
		},
		tags: []
	});

	// 初始化用户数据
	const initUserInfo = () => {
		getPersonalInfo().then((res: any) => {
			const user = res.data.user;
			state.personalForm = {
        username: user.userName,
				nickname: user.userNickname,
				userEmail: user.userEmail,
				describe: user.describe,
				mobile: user.mobile,
				sex: String(user.sex),
				remark: user.remark,
				// avatar: user.avatar,
				lastLoginIp: user.lastLoginIp,
				lastLoginTime: user.lastLoginTime,
				tags: user.tags,
			};
			if(user.tags && user.tags !== "") state.tags = user.tags?.split(',');
			state.deptName = res.data.deptName;
			state.roles = res.data.roles;
		});
	};

	// 提交信息变更
	const handleUpload = () => {
		state.personalForm.tags = state.tags.join(',');
		editPersonal(state.personalForm).then((res: any) => {
			const userInfo = res.data.userInfo;
			// userInfo.avatar = proxy.getUpFileUrl(userInfo.avatar);
			// 存储 token 到浏览器缓存
			Session.set('token', res.data.token);
			// 存储用户信息到浏览器缓存
			Session.set('userInfo', userInfo);
			Cookies.set('username', userInfo.userName);
			Cookies.set('usernickname', userInfo.userNickname);
			useUserInfo().setUserInfos();
			ElMessage.success('已更新');

			proxy.mittBus.emit('updateUserInfo');
		});
	};

	const handleClose = (tag : string) => {
		state.tags.splice(state.tags.indexOf(tag), 1);
	};

	const showInput = () => {
		inputVisible.value = true
		nextTick(() => {
			inputRef.value!.input!.focus()
		})
	};

	const handleInputConfirm = () => {
		if (inputValue.value) {
			state.tags.push(inputValue.value);
		}
		inputVisible.value = false;
		inputValue.value = '';
	};

	// 页面加载时
	onMounted(() => {
		initUserInfo();
	});
</script>

<style scoped lang="scss">
	@use '../../../../theme/mixins/index.scss' as mixins;
	.personal-user-card {
		flex: 1;
		padding: 15px;
		.personal-title {
			font-size: 18px;
			@include mixins.text-ellipsis(1);
		}
		.personal-item {
			display: flex;
			align-items: center;
			font-size: 13px;
			line-height: 26px;
			.personal-item-label {
				color: var(--el-text-color-secondary);
				@include mixins.text-ellipsis(1);
			}
			.personal-item-value {
				@include mixins.text-ellipsis(1);
			}
		}
	}
	.personal-recommend-row {
		.personal-recommend-col {
			.personal-recommend {
				position: relative;
				height: 100px;
				border-radius: 3px;
				overflow: hidden;
				cursor: pointer;
				&:hover {
					i {
						right: 0px !important;
						bottom: 0px !important;
						transition: all ease 0.3s;
					}
				}
				i {
					position: absolute;
					right: -10px;
					bottom: -10px;
					font-size: 70px;
					transform: rotate(-30deg);
					transition: all ease 0.3s;
				}
				.personal-recommend-auto {
					padding: 15px;
					position: absolute;
					left: 0;
					top: 5%;
					color: var(--next-color-white);
					.personal-recommend-msg {
						font-size: 12px;
						margin-top: 10px;
					}
				}
			}
		}
	}
	.personal-user-tabs {
		padding: 0 !important;
	}
</style>