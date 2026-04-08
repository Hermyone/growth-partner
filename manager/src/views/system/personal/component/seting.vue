<template>
	<div class="q-small-title">语言主题</div>

	<div class="personal-seting-bar">
		<el-row>
			<el-col>
				<div class="personal-seting-bar-flex mt15" >
					<div class="personal-seting-bar-flex-label">{{ $t('message.user.langSwitching') }}</div>
					<div class="personal-seting-bar-flex-value">
						<el-select
								v-model="getThemeConfig.globalI18n"
								placeholder="请选择"
								size="default"
								style="width: 100px"
								@change="onLanguageChange"
						>
							<el-option label="简体中文" value="zh-cn" :disabled="state.disabledI18n === 'zh-cn'"></el-option>
							<el-option label="English" value="en" :disabled="state.disabledI18n === 'en'"></el-option>
							<el-option label="繁體中文" value="zh-tw" :disabled="state.disabledI18n === 'zh-tw'"></el-option>
						</el-select>
					</div>
				</div>
			</el-col>
			<el-divider />
			<el-col>
				<div class="personal-seting-bar-flex mt15">
					<div class="personal-seting-bar-flex-label">主题</div>
					<div class="personal-seting-bar-flex-value">
						<el-color-picker v-model="getThemeConfig.primary" size="default" :predefine="state.predefineColors" @change="onColorPickerChange"> </el-color-picker>
					</div>
				</div>
			</el-col>
			<el-col>
				<div class="personal-seting-bar-flex mt15">
					<div class="personal-seting-bar-flex-label">{{ $t('message.layout.fourIsDark') }}</div>
					<div class="personal-seting-bar-flex-value">
						<el-switch v-model="getThemeConfig.isIsDark" size="small" @change="onAddDarkChange"></el-switch>
					</div>
				</div>
			</el-col>
			<el-col>
				<div class="personal-seting-bar-flex mt15">
					<div class="personal-seting-bar-flex-label">{{ $t('message.layout.fourIsGrayscale') }}</div>
					<div class="personal-seting-bar-flex-value">
						<el-switch v-model="getThemeConfig.isGrayscale" size="small" @change="onAddFilterChange('grayscale')"></el-switch>
					</div>
				</div>
			</el-col>
			<el-col>
				<div class="personal-seting-bar-flex mt15">
					<div class="personal-seting-bar-flex-label">{{ $t('message.layout.fourIsInvert') }}</div>
					<div class="personal-seting-bar-flex-value">
						<el-switch v-model="getThemeConfig.isInvert" size="small" @change="onAddFilterChange('invert')"></el-switch>
					</div>
				</div>
			</el-col>

			<el-col>
				<el-button size="default" class="copy-config-btn-reset" type="info" @click="onResetConfigClick">
					<el-icon class="mr5">
						<ele-RefreshRight />
					</el-icon>
					{{ $t('message.layout.resetText') }}
				</el-button>
			</el-col>
		</el-row>
	</div>
	<div style="padding: 20px;">
		<el-button type="primary" link @click="onMoreConfigClick">
			{{ $t('message.layout.fiveTitle') }}
		</el-button>
	</div>
</template>

<script setup lang="ts" name="personal-seting">
	import { reactive, onMounted, getCurrentInstance, computed } from 'vue';
	import { storeToRefs } from 'pinia';
	import { useThemeConfig } from '/@/stores/themeConfig';
	import { ElMessage, ElMessageBox } from 'element-plus';
	import { getDarkColor, getLightColor } from '/@/utils/theme';
	import { Local } from '/@/utils/storage';
	import common from '/@/utils/common';
	import { resetPwdPersonal } from '/@/api/system/personal';

	const { proxy } = <any>getCurrentInstance();
	const storesThemeConfig = useThemeConfig();
	const { themeConfig } = storeToRefs(storesThemeConfig);
	const state = reactive({
		disabledI18n: 'zh-cn',
		predefineColors:[
			'#409eff',
			'#0A69F6',
			'#07cd5a',
			'#14B47D',
			'#ff8c00',
			'#ff4500',
		],
	});

	// 获取布局配置信息
	const getThemeConfig = computed(() => {
		return themeConfig.value;
	});

	// 初始化言语国际化
	const initI18n = () => {
		switch (Local.get('themeConfig').globalI18n) {
			case 'zh-cn':
				state.disabledI18n = 'zh-cn';
				setI18nConfig('zh-cn');
				break;
			case 'en':
				state.disabledI18n = 'en';
				setI18nConfig('en');
				break;
			case 'zh-tw':
				state.disabledI18n = 'zh-tw';
				setI18nConfig('zh-tw');
				break;
		}
	};

	// 设置 element plus 组件的国际化
	const setI18nConfig = (locale: string) => {
		proxy.mittBus.emit('getI18nConfig', proxy.i18n.global.messages.value[locale]);
	};

	// 语言切换
	const onLanguageChange = (lang: string) => {
		Local.remove('themeConfig');
		themeConfig.value.globalI18n = lang;
		Local.set('themeConfig', themeConfig.value);
		proxy.$i18n.locale = lang;
		initI18n();
		common.useTitle();
	};

	// 全局主题
	const onColorPickerChange = () => {
		if (!getThemeConfig.value.primary) return ElMessage.warning('全局主题 primary 颜色值不能为空');
		// 颜色加深
		document.documentElement.style.setProperty('--el-color-primary-dark-2', `${getDarkColor(getThemeConfig.value.primary, 0.1)}`);
		document.documentElement.style.setProperty('--el-color-primary', getThemeConfig.value.primary);
		// 颜色变浅
		for (let i = 1; i <= 9; i++) {
			document.documentElement.style.setProperty(`--el-color-primary-light-${i}`, `${getLightColor(getThemeConfig.value.primary, i / 10)}`);
		}
		setDispatchThemeConfig();
	};

	// 触发 store 布局配置更新
	const setDispatchThemeConfig = () => {
		setLocalThemeConfig();
		setLocalThemeConfigStyle();
	};

	// 存储布局配置
	const setLocalThemeConfig = () => {
		Local.remove('themeConfig');
		Local.set('themeConfig', getThemeConfig.value);
	};
	// 存储布局配置全局主题样式（html根标签）
	const setLocalThemeConfigStyle = () => {
		Local.set('themeConfigStyle', document.documentElement.style.cssText);
	};
	// 深色模式
	const onAddDarkChange = () => {
		const body = document.documentElement as HTMLElement;
		if (getThemeConfig.value.isIsDark) body.setAttribute('data-theme', 'dark');
		else body.setAttribute('data-theme', '');
	};
	// 灰色模式/色弱模式
	const onAddFilterChange = (attr: string) => {
		if (attr === 'grayscale') {
			if (getThemeConfig.value.isGrayscale) getThemeConfig.value.isInvert = false;
		} else {
			if (getThemeConfig.value.isInvert) getThemeConfig.value.isGrayscale = false;
		}
		const cssAttr =
				attr === 'grayscale' ? `grayscale(${getThemeConfig.value.isGrayscale ? 1 : 0})` : `invert(${getThemeConfig.value.isInvert ? '80%' : '0%'})`;
		const appEle: any = document.body;
		appEle.setAttribute('style', `filter: ${cssAttr}`);
		setLocalThemeConfig();
	};

	// 一键恢复默认
	const onResetConfigClick = () => {
		ElMessageBox.confirm('确认要还原默认设置吗？', '提示', {
			confirmButtonText: '确认',
			cancelButtonText: '取消',
			type: 'warning',
		})
				.then(() => {
					Local.clear();
					window.location.reload();
				})
				.catch(() => {
				});
	};

	// 更多设置
	const onMoreConfigClick = () => {
		proxy.mittBus.emit('openSetingsDrawer');
	};

</script>

<style scoped lang="scss">
	.personal-seting-bar {
		padding-left: 20px;
		&-flex {
			display: flex;
			align-items: center;
			margin-bottom: 2px;
			width: 300px;
			&-label {
				flex: 1;
				color: var(--el-text-color-primary);
			}
		}
	}
	.copy-config-btn-reset {
		width: 300px;
		margin: 20px 0 0;
	}
</style>

