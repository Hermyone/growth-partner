<template>
	<div v-show="isShowLockScreen" class="layout-lock-screen-root">
		<div class="layout-lock-screen-mask"></div>
		<div class="layout-lock-screen-img" :class="{ 'layout-lock-screen-filter': isShowLoockLogin }"></div>
		<div class="layout-lock-screen-tint" aria-hidden="true"></div>
		<div class="layout-lock-screen">
			<div
				class="layout-lock-screen-date"
				ref="layoutLockScreenDateRef"
				@mousedown="onDown"
				@mousemove="onMove"
				@mouseup="onEnd"
				@touchstart.stop="onDown"
				@touchmove.stop="onMove"
				@touchend.stop="onEnd"
			>
				<div class="layout-lock-screen-date-box">
					<div class="layout-lock-screen-date-box-time">
						{{ time.hm }}<span class="layout-lock-screen-date-box-minutes">{{ time.s }}</span>
					</div>
					<div class="layout-lock-screen-date-box-info">{{ time.mdq }}</div>
				</div>
				<div class="layout-lock-screen-date-top">
					<SvgIcon name="ele-Top" />
					<div class="layout-lock-screen-date-top-text">上滑解锁</div>
				</div>
				<p class="layout-lock-screen-hint" aria-hidden="true">温暖守护 · 成长相伴</p>
			</div>
			<transition name="el-zoom-in-center">
				<div v-show="isShowLoockLogin" class="layout-lock-screen-login">
					<div class="layout-lock-screen-login-card">
						<div class="layout-lock-screen-login-box">
							<div class="layout-lock-screen-login-box-img">
								<img :src="userInfos.avatar" alt="" />
							</div>
							<div class="layout-lock-screen-login-box-name">
								{{ userInfos.userName === '' ? 'common' : userInfos.userName }}
							</div>
							<div class="layout-lock-screen-login-box-value">
								<el-input
									type="password"
									size="large"
									placeholder="请输入密码"
									ref="layoutLockScreenInputRef"
									v-model="lockScreenPassword"
									@keyup.enter.stop="onLockScreenSubmit()"
								>
									<template #append>
										<div @click="onLockScreenSubmit">
											<el-icon class="el-input__icon" style="font-size: 24px; font-weight: 700;">
												<ele-Right />
											</el-icon>
										</div>
									</template>
								</el-input>
							</div>
						</div>
					</div>
					<!-- <div class="layout-lock-screen-login-icon">
						<SvgIcon name="ele-Lock" :size="32" @click="onLockScreen" />
					</div> -->
				</div>
			</transition>
		</div>
	</div>
</template>

<script lang="ts">
import { nextTick, onMounted, reactive, toRefs, ref, onUnmounted, getCurrentInstance, defineComponent } from 'vue';
import { formatDate } from '/@/utils/formatTime';
import { Local } from '/@/utils/storage';
import { storeToRefs } from 'pinia';
import { useThemeConfig } from '/@/stores/themeConfig';
import { useUserInfo } from '/@/stores/userInfo';
import { ElMessage } from 'element-plus';

export default defineComponent({
	name: 'layoutLockScreen',
	setup() {
		const { proxy } = <any>getCurrentInstance();
		const layoutLockScreenInputRef = ref();
		const storesThemeConfig = useThemeConfig();
		const { themeConfig } = storeToRefs(storesThemeConfig);
		const storesUser = useUserInfo();
		const { userInfos } = storeToRefs(storesUser);
		const state = reactive<LockScreenState>({
			transparency: 1,
			downClientY: 0,
			moveDifference: 0,
			isShowLoockLogin: false,
			isFlags: false,
			querySelectorEl: '',
			time: {
				hm: '',
				s: '',
				mdq: '',
			},
			setIntervalTime: 0,
			isShowLockScreen: false,
			isShowLockScreenIntervalTime: 0,
			lockScreenPassword: '',
			headImage: '',
		});
		// 鼠标按下
		const onDown = (down: any) => {
			state.isFlags = true;
			state.downClientY = down.touches ? down.touches[0].clientY : down.clientY;
		};
		// 鼠标移动
		const onMove = (move: any) => {
			if (state.isFlags) {
				const el = <HTMLElement>state.querySelectorEl;
				const opacitys = (state.transparency -= 1 / 200);
				if (move.touches) {
					state.moveDifference = move.touches[0].clientY - state.downClientY;
				} else {
					state.moveDifference = move.clientY - state.downClientY;
				}
				if (state.moveDifference >= 0) return false;
				el.setAttribute('style', `top:${state.moveDifference}px;cursor:pointer;opacity:${opacitys};`);
				if (state.moveDifference < -400) {
					el.setAttribute('style', `top:${-el.clientHeight}px;cursor:pointer;transition:all 0.3s ease;`);
					state.moveDifference = -el.clientHeight;
					setTimeout(() => {
						el && el.parentNode?.removeChild(el);
					}, 300);
				}
				if (state.moveDifference === -el.clientHeight) {
					state.isShowLoockLogin = true;
					layoutLockScreenInputRef.value.focus();
				}
			}
		};
		// 鼠标松开
		const onEnd = () => {
			state.isFlags = false;
			state.transparency = 1;
			if (state.moveDifference >= -400) {
				(<HTMLElement>state.querySelectorEl).setAttribute('style', `top:0px;opacity:1;transition:all 0.3s ease;`);
			}
		};
		// 获取要拖拽的初始元素
		const initGetElement = () => {
			nextTick(() => {
				state.querySelectorEl = proxy.$refs.layoutLockScreenDateRef;
			});
		};
		// 时间初始化
		const initTime = () => {
			state.time.hm = formatDate(new Date(), 'HH:MM');
			state.time.s = formatDate(new Date(), 'SS');
			state.time.mdq = formatDate(new Date(), 'mm月dd日，WWW');
		};
		// 时间初始化定时器
		const initSetTime = () => {
			initTime();
			state.setIntervalTime = window.setInterval(() => {
				initTime();
			}, 1000);
		};
		// 锁屏时间定时器
		const initLockScreen = () => {
			if (themeConfig.value.isLockScreen) {
				state.isShowLockScreenIntervalTime = window.setInterval(() => {
					if (themeConfig.value.lockScreenTime <= 1) {
						state.isShowLockScreen = true;
						setLocalThemeConfig();
						return false;
					}
					themeConfig.value.lockScreenTime--;
				}, 1000);
			} else {
				clearInterval(state.isShowLockScreenIntervalTime);
			}
		};
		// 存储布局配置
		const setLocalThemeConfig = () => {
			themeConfig.value.isDrawer = false;
			Local.set('themeConfig', themeConfig.value);
		};
		// 锁屏
		const onLockScreen = () => {
			window.location.reload();
		};

		// 密码输入点击事件
		const onLockScreenSubmit = () => {
			if(state.lockScreenPassword === themeConfig.value.lockPassword) {
				themeConfig.value.isLockScreen = false;
				themeConfig.value.lockScreenTime = 30;
				setLocalThemeConfig();
			} else {
				ElMessage.error('密码不正确');
			}
		};
		// 页面加载时
		onMounted(() => {
			initGetElement();
			initSetTime();
			initLockScreen();
		});
		// 页面卸载时
		onUnmounted(() => {
			window.clearInterval(state.setIntervalTime);
			window.clearInterval(state.isShowLockScreenIntervalTime);
		});
		return {
			layoutLockScreenInputRef,
			userInfos,
			onDown,
			onMove,
			onEnd,
			onLockScreenSubmit,
			onLockScreen,
			...toRefs(state),
		};
	},
});
</script>

<style scoped lang="scss">
/* 与登录页一致：阳光橙 #FF9F43 · 天空蓝 #48DBFB · 薄荷绿 #1DD1A1 · 柔和粉 #FF9FF3 */
.layout-lock-screen-root {
	--growth-orange: #ff9f43;
	--growth-sky: #48dbfb;
	--growth-mint: #1dd1a1;
	--growth-pink: #ff9ff3;
}

.layout-lock-screen-fixed {
	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
}

.layout-lock-screen-filter {
	filter: blur(6px);
}

.layout-lock-screen-mask {
	background: linear-gradient(165deg, #fff9f2 0%, #e8f9ff 50%, #f5fffb 100%);
	@extend .layout-lock-screen-fixed;
	z-index: 999990;
}

.layout-lock-screen-img {
	@extend .layout-lock-screen-fixed;
	background-image: url('/@/assets/login-growth-bg.svg');
	background-size: cover;
	background-position: center;
	z-index: 999991;
}

.layout-lock-screen-tint {
	@extend .layout-lock-screen-fixed;
	z-index: 999991;
	pointer-events: none;
	background: linear-gradient(
		125deg,
		rgba(255, 159, 67, 0.1) 0%,
		rgba(72, 219, 251, 0.12) 42%,
		rgba(29, 209, 161, 0.07) 100%
	);
}

.layout-lock-screen {
	@extend .layout-lock-screen-fixed;
	z-index: 999992;

	&-date {
		position: absolute;
		left: 0;
		top: 0;
		width: 100%;
		height: 100%;
		z-index: 999993;
		user-select: none;
		color: #fff;

		&-box {
			position: absolute;
			left: 28px;
			bottom: 52px;
			text-shadow:
				0 2px 20px rgba(0, 80, 100, 0.18),
				0 1px 3px rgba(0, 0, 0, 0.12);

			&-time {
				font-size: clamp(56px, 14vw, 100px);
				font-weight: 600;
				letter-spacing: -0.02em;
				line-height: 1.05;
				background: linear-gradient(115deg, #ffffff 0%, #fff8f0 35%, #e8fbff 100%);
				-webkit-background-clip: text;
				background-clip: text;
				color: transparent;
				filter: drop-shadow(0 2px 12px rgba(72, 219, 251, 0.25));
			}

			&-info {
				font-size: clamp(22px, 4.5vw, 40px);
				margin-top: 8px;
				font-weight: 500;
				color: rgba(255, 255, 255, 0.96);
				text-shadow: 0 2px 16px rgba(0, 60, 80, 0.2);
			}

			&-minutes {
				font-size: clamp(14px, 2.5vw, 18px);
				font-weight: 500;
				opacity: 0.95;
				vertical-align: super;
				margin-left: 4px;
			}
		}

		&-top {
			width: 44px;
			height: 44px;
			line-height: 44px;
			border-radius: 100%;
			border: 2px solid rgba(255, 255, 255, 0.65);
			background: linear-gradient(145deg, rgba(255, 159, 67, 0.35), rgba(72, 219, 251, 0.3));
			backdrop-filter: blur(8px);
			-webkit-backdrop-filter: blur(8px);
			box-shadow: 0 4px 20px rgba(255, 159, 67, 0.2);
			color: #fff;
			opacity: 0.92;
			position: absolute;
			right: 28px;
			bottom: 52px;
			text-align: center;
			overflow: hidden;
			transition: all 0.3s ease;

			i {
				transition: all 0.3s ease;
			}

			&-text {
				opacity: 0;
				position: absolute;
				top: 150%;
				font-size: 12px;
				color: #fff;
				left: 50%;
				line-height: 1.2;
				transform: translate(-50%, -50%);
				transition: all 0.3s ease;
				width: 40px;
				text-shadow: 0 1px 4px rgba(0, 0, 0, 0.15);
			}

			&:hover {
				border-color: rgba(255, 255, 255, 0.9);
				background: linear-gradient(145deg, rgba(255, 159, 67, 0.5), rgba(72, 219, 251, 0.45));
				box-shadow:
					0 0 0 2px rgba(255, 159, 243, 0.35),
					0 8px 28px rgba(72, 219, 251, 0.25);
				opacity: 1;

				i {
					transform: translateY(-44px);
				}

				.layout-lock-screen-date-top-text {
					opacity: 1;
					top: 50%;
				}
			}
		}
	}

	.layout-lock-screen-hint {
		position: absolute;
		left: 50%;
		bottom: 24px;
		transform: translateX(-50%);
		margin: 0;
		font-size: 13px;
		letter-spacing: 0.12em;
		color: rgba(255, 255, 255, 0.88);
		text-shadow: 0 1px 8px rgba(0, 80, 90, 0.2);
		white-space: nowrap;
	}
}

.layout-lock-screen-login {
	position: relative;
	z-index: 999994;
	width: 100%;
	height: 100%;
	left: 0;
	top: 0;
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	padding: 24px;

	&-card {
		width: 100%;
		max-width: 400px;
		padding: 36px 32px 40px;
		background: rgba(255, 255, 255, 0.92);
		backdrop-filter: blur(14px);
		-webkit-backdrop-filter: blur(14px);
		border-radius: 24px;
		border: 1px solid rgba(255, 159, 67, 0.28);
		box-shadow:
			0 4px 28px rgba(72, 219, 251, 0.14),
			0 12px 48px rgba(255, 159, 67, 0.1),
			inset 0 1px 0 rgba(255, 255, 255, 0.95);
		position: relative;
		overflow: hidden;

		&::before {
			content: '';
			position: absolute;
			top: 0;
			left: 0;
			right: 0;
			height: 4px;
			background: linear-gradient(
				90deg,
				var(--growth-orange) 0%,
				var(--growth-sky) 34%,
				var(--growth-mint) 68%,
				var(--growth-pink) 100%
			);
		}
	}

	&-box {
		text-align: center;
		margin: auto;

		&-img {
			width: 180px;
			height: 180px;
			margin: auto;
			padding: 4px;
			border-radius: 100%;
			background: linear-gradient(135deg, var(--growth-orange), var(--growth-sky), var(--growth-mint));
			box-shadow: 0 8px 32px rgba(255, 159, 67, 0.25);

			img {
				width: 100%;
				height: 100%;
				border-radius: 100%;
				object-fit: cover;
				display: block;
				border: 3px solid rgba(255, 255, 255, 0.95);
			}
		}

		&-name {
			font-size: 26px;
			font-weight: 600;
			margin: 18px 0 24px;
			background: linear-gradient(100deg, var(--growth-orange) 0%, #e8892e 38%, var(--growth-sky) 100%);
			-webkit-background-clip: text;
			background-clip: text;
			color: transparent;
		}

		&-value {
			width: 100%;

			:deep(.el-input-group) {
				width: 100%;
				max-width: 320px;
				margin: 0 auto;
			}

			:deep(.el-input__wrapper) {
				border-radius: 12px 0 0 12px !important;
				box-shadow: 0 0 0 1px rgba(72, 219, 251, 0.35) inset !important;
				font-size: 18px;
				min-height: 48px;
			}

			:deep(.el-input__wrapper):hover {
				box-shadow: 0 0 0 1px rgba(255, 159, 67, 0.45) inset !important;
			}

			:deep(.el-input__wrapper.is-focus) {
				box-shadow:
					0 0 0 1px var(--growth-orange) inset,
					0 0 0 3px rgba(72, 219, 251, 0.2) !important;
			}

			:deep(.el-input-group__append) {
				padding: 0;
				border: none !important;
				background: linear-gradient(135deg, var(--growth-orange) 0%, #ffb35a 50%, var(--growth-sky) 100%) !important;
				box-shadow: none !important;
				border-radius: 0 12px 12px 0 !important;
				cursor: pointer;
				transition: filter 0.2s ease;

				&:hover {
					filter: brightness(1.06);
				}

				.el-input__icon {
					color: #fff !important;
					margin: 0 16px;
				}
			}
		}
	}

	&-icon {
		position: absolute;
		right: 30px;
		bottom: 30px;

		i {
			font-size: 20px;
			margin-left: 15px;
			cursor: pointer;
			opacity: 0.8;

			&:hover {
				opacity: 1;
			}
		}
	}
}
</style>
