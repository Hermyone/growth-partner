<template>
	<q-container :cardMode="getCardMode()" class="user-info-container">
		<el-container>
			<el-aside width="221px">
				<el-container>
					<el-header style="height: auto;display: block;">
						<div class="user-info-top">
							<div class="user-info-avatar">
								<el-button icon="ele-Edit" circle size="small" class="user-info-avatar-btn" @click="editHeadImage" ></el-button>
								<el-avatar :size="65" :src="imageUrl"></el-avatar>
							</div>
							<p>{{ currentTime }}，{{ personalForm.nickname }} </p>
							<p><el-tag effect="dark" round size="small" v-if="roles.length > 0" disable-transitions>{{ roles.join(',') }}</el-tag></p>
						</div>
					</el-header>
					<el-main router class="nopadding">
						<el-menu class="custom-menu" :default-active="page">
							<el-menu-item-group v-for="group in menu" :key="group.groupName" :title="group.groupName">
								<el-menu-item v-for="item in group.list" :key="item.component" :index="item.component" @click="openPage">
									<el-icon v-if="item.icon"><component :is="item.icon"/></el-icon>
									<template #title>
										<span>{{item.title}}</span>
									</template>
								</el-menu-item>
							</el-menu-item-group>
						</el-menu>
					</el-main>
				</el-container>
			</el-aside>

			<el-main>
				<Suspense>
					<template #default>
						<component :is="page"/>
					</template>
					<template #fallback>
						<el-skeleton :rows="3" />
					</template>
				</Suspense>
			</el-main>
		</el-container>

		<q-cropper ref="cropperRef" title="更换头像" @cropperChange="changeHeadImage" :uploadUrl="uploadUrl"></q-cropper>
	</q-container>
</template>

<script lang="ts">
	import {
		defineComponent,
		reactive,
		computed,
		getCurrentInstance,
		onMounted,
		toRefs,
		defineAsyncComponent,
		nextTick,
		onUnmounted,
	} from 'vue';
	import { formatAxis } from '/@/utils/formatTime';
	import { getPersonalInfo, setAvatar } from '/@/api/system/personal';
	import { getCardMode  } from '/@/utils/common';
	import sysConfig from '/@/config';
  import {Session} from "/@/utils/storage";
  import {useUserInfo} from "/@/stores/userInfo";
  import {ElMessage} from "element-plus";

	// 引入组件
	const account = defineAsyncComponent(() => import('./component/account.vue'));
	const logs = defineAsyncComponent(() => import('./component/logs.vue'));
	const space = defineAsyncComponent(() => import('./component/space.vue'));
	const password = defineAsyncComponent(() => import('./component/password.vue'));
	const notice = defineAsyncComponent(() => import('./component/notice.vue'));
	const seting = defineAsyncComponent(() => import('./component/seting.vue'));
	const about = defineAsyncComponent(() => import('./component/about.vue'));
	const updateLogs = defineAsyncComponent(() => import('/@/views/system/update/component/index.vue'));

	export default defineComponent({
		name: 'personal',
		components: { account, logs, space, password, notice, seting, about, updateLogs },
		setup() {
			const { proxy } = <any>getCurrentInstance();
			const state = reactive({
				imageUrl: '',
				uploadUrl: sysConfig.FILE_URL,
				personalForm: {
					nickname: '',
					describe: '',
				},
				roles: [],
				menu: [
					{
						groupName: "基本设置",
						list: [
							{
								icon: "ele-User",
								title: "账号信息",
								component: "account",
							},
							{
								icon: "ele-Lock",
								title: "账号安全",
								component: "password",
							},
							{
								icon: "ele-SetUp",
								title: "个人设置",
								component: "seting",
							},
							{
								icon: "ele-Bell",
								title: "通知设置",
								component: "notice",
							}
						]
					},
					{
						groupName: "数据管理",
						list: [
							{
								icon: "ele-Document",
								title: "日志信息",
								component: "logs",
							},
							{
								icon: "ele-Monitor",
								title: "存储空间信息",
								component: "space",
							}
						]
					},
					{
						groupName: "系统",
						list: [
              {
                icon: "ele-Clock",
                title: "升级日志",
                component: "updateLogs",
              },
							{
								icon: "ele-Warning",
								title: "关于",
								component: "about",
							},
						]
					},
				],
				page: "account",
			});

			const openPage = (item: any) => {
				state.page = item.index as any;
			};

			// 当前时间提示语
			const currentTime = computed(() => {
				return formatAxis(new Date());
			});

			// 初始化用户数据
			const initUserInfo = () => {
				getPersonalInfo().then((res: any) => {
					const user = res.data.user;
					// state.imageUrl = user.avatar;// '/images/logo.png';// proxy.getUpFileUrl(user.avatar);
          convertImageUrlToBase64(user.avatar).then(base64 => {
            state.imageUrl = base64
          })
					state.personalForm = {
						nickname: user.userNickname,
						describe: user.describe,
					};
					state.roles = res.data.roles;
				});
			};

      const convertImageUrlToBase64 = (url: string) => {
        return fetch(url)
            .then(response => response.blob())
            .then(blob => {
              return new Promise((resolve, reject) => {
                const reader = new FileReader();
                reader.onloadend = () => resolve(reader.result);
                reader.onerror = reject;
                reader.readAsDataURL(blob);
              });
            });
      }

			// 修改头像
			const editHeadImage = () => {
				proxy.$refs.cropperRef.openDialog(state.imageUrl);
			};

			// 保存头像
			const changeHeadImage = (imageData: any, res: any) => {
        if(res){
          setAvatar({url: res.data.url}).then(resp => {
            state.imageUrl = imageData;

            const userInfo = resp.data;
            // 存储用户信息到浏览器缓存
            Session.set('userInfo', userInfo);
            useUserInfo().setUserInfos();
            ElMessage.success('已更新');

            proxy.mittBus.emit('updateUserInfo');
          })
        }else{
          state.imageUrl = imageData;
        }
			};

			// 页面加载时
			onMounted(() => {
				nextTick(() => {
					initUserInfo();
					proxy.mittBus.on('updateUserInfo', () => {
						initUserInfo();
					});
				});
			});

			onUnmounted(() => {
				proxy.mittBus.off('updateUserInfo', ()=>{});
			});

			return {
				getCardMode,
				currentTime,
				openPage,
				editHeadImage,
				changeHeadImage,
				...toRefs(state),
			}
		},
	});
</script>

<style scoped lang="scss">
	.user-info-container{
		:deep(.el-card__body) {
			padding: 0 !important;
		}
		&-header{
			align-items: flex-start;
		}

		.user-info-top {
			text-align: center;
		}

		.user-info-top p {margin: 10px 0 10px 0;}
	}
	.user-info-avatar {
		display: flex;
		padding-left: 32px;
		&-btn {
			position: relative;
			top: 50px;
			left: 65px;
		}

		:deep(.el-button i.el-icon,.el-button i.iconfont,.el-button i.fa,.el-button--default i.iconfont,.el-button--default i.fa){
			font-size: 12px !important;
			margin-right: 0!important;
		}
	}
</style>