<template>
	<el-card shadow="hover" header="修改密码" class="user-password">
		<div style="padding: 15px;">
			<el-alert title="密码更新成功后，建议您使用新密码重新登录。" type="warning" show-icon style="margin-bottom: 15px;"/>
			<el-form ref="formRef" :model="state.password" :rules="state.rules" label-width="100px">
				<el-form-item label="当前密码" prop="userPassword" style="height: 60px;">
					<el-input v-model="state.password.userPassword" type="password" show-password placeholder="请输入当前密码" style="height: 32px;"></el-input>
					<div class="el-form-item-msg">必须提供当前登录用户密码才能进行更改</div>
				</el-form-item>
				<el-form-item label="新密码" prop="newPassword" style="height: 80px;">
					<el-input v-model="state.password.newPassword" type="password" show-password placeholder="请输入新密码" style="height: 32px;"></el-input>
					<div class="el-form-item-msg">请输入包含英文、数字的8位以上密码</div>
					<q-passwordstrength v-model="state.password.newPassword"></q-passwordstrength>
				</el-form-item>
				<el-form-item label="确认新密码" prop="confirmNewPassword">
					<el-input v-model="state.password.confirmNewPassword" type="password" show-password placeholder="请再次输入新密码" style="height: 32px;"></el-input>
				</el-form-item>
				<el-form-item>
					<el-button type="primary" @click="savePassword">修改密码</el-button>
				</el-form-item>
			</el-form>
		</div>
	</el-card>

<!--	<div style="padding: 30px 0 !important;">-->
<!--		<div class="q-small-title">账号安全</div>-->

<!--		<div class="personal-edit-safe-box">-->
<!--			<div class="personal-edit-safe-item">-->
<!--				<div class="personal-edit-safe-item-left">-->
<!--					<div class="personal-edit-safe-item-left-label">密保手机</div>-->
<!--					<div class="personal-edit-safe-item-left-value">已绑定手机：132****4108</div>-->
<!--				</div>-->
<!--				<div class="personal-edit-safe-item-right">-->
<!--					<el-button text type="primary">立即修改</el-button>-->
<!--				</div>-->
<!--			</div>-->
<!--		</div>-->
<!--		<div class="personal-edit-safe-box">-->
<!--			<div class="personal-edit-safe-item">-->
<!--				<div class="personal-edit-safe-item-left">-->
<!--					<div class="personal-edit-safe-item-left-label">密保问题</div>-->
<!--					<div class="personal-edit-safe-item-left-value">已设置密保问题，账号安全大幅度提升</div>-->
<!--				</div>-->
<!--				<div class="personal-edit-safe-item-right">-->
<!--					<el-button text type="primary">立即设置</el-button>-->
<!--				</div>-->
<!--			</div>-->
<!--		</div>-->
<!--		<div class="personal-edit-safe-box">-->
<!--			<div class="personal-edit-safe-item">-->
<!--				<div class="personal-edit-safe-item-left">-->
<!--					<div class="personal-edit-safe-item-left-label">绑定QQ</div>-->
<!--					<div class="personal-edit-safe-item-left-value">已绑定QQ：110****566</div>-->
<!--				</div>-->
<!--				<div class="personal-edit-safe-item-right">-->
<!--					<el-button text type="primary">立即设置</el-button>-->
<!--				</div>-->
<!--			</div>-->
<!--		</div>-->
<!--	</div>-->

</template>

<script setup lang="ts" name="personal-password">
	import { reactive, onMounted, ref, unref } from 'vue';
	import { ElMessage, ElMessageBox } from 'element-plus';
	import { resetPwdPersonal } from '/@/api/system/personal';
	// import Cookies from 'js-cookie';
	// import crypto from '/@/utils/crypto';

	const formRef = ref<HTMLElement | null>(null);
	const state = reactive({
		rules: {
			userPassword: [ { required: true, message: '请输入当前密码'} ],
			newPassword: [ { required: true, message: '请输入新密码'} ],
			confirmNewPassword: [
				{ required: true, message: '请再次输入新密码'},
				{validator: (rule: any, value: any, callback: any) => {
						if (value !== state.password.newPassword) {
							callback(new Error('两次输入密码不一致'));
						}else{
							callback();
						}
					}}
			],
		},
		password: {
			userPassword: '',
			newPassword: '',
			confirmNewPassword: '',
		},
	});

	// 保存密码
	const savePassword = () => {
		const formWrap = unref(formRef) as any;
		if (!formWrap) return;
		formWrap.validate((valid: boolean) => {
			if(valid) {
				// let curPass = Cookies.get('password');
				// if(curPass !== crypto.MD5(state.password.userPassword)){
				// 	ElMessage.error('当前密码不正确');
				// 	state.password.userPassword = '';
				// 	return;
				// }

        if(state.password.userPassword.length < 6 || state.password.newPassword.length < 6){
          ElMessage.error('密码长度不能小于6位');
          return;
        }

				ElMessageBox.confirm('确认要修改密码吗？', '提示', {
					confirmButtonText: '确认',
					cancelButtonText: '取消',
					type: 'warning',
				})
						.then(() => {
							resetPwdPersonal({ oldpassword: state.password.userPassword, password: state.password.newPassword }).then(() => {
								ElMessage.success('密码修改成功');
								state.password.userPassword = '';
								state.password.newPassword = '';
								state.password.confirmNewPassword = '';
							});
						})
						.catch(() => {
						});
			}
		});
	};

	// 页面加载时
	onMounted(() => {

	});
</script>

<style scoped lang="scss">
	@use '../../../../theme/mixins/index.scss' as mixins;
	.user-password {
		padding: 0 15px;
	}
	.personal-edit-safe-box {
		border-bottom: 1px solid var(--el-border-color-light, #ebeef5);
		padding: 15px 0;
		.personal-edit-safe-item {
			width: 100%;
			display: flex;
			align-items: center;
			justify-content: space-between;
			.personal-edit-safe-item-left {
				flex: 1;
				overflow: hidden;
				.personal-edit-safe-item-left-label {
					color: var(--el-text-color-regular);
					margin-bottom: 5px;
				}
				.personal-edit-safe-item-left-value {
					color: var(--el-text-color-secondary);
					@include mixins.text-ellipsis(1);
					margin-right: 15px;
				}
			}
		}
		&:last-of-type {
			padding-bottom: 0;
			border-bottom: none;
		}
	}
</style>

