<template>
	<el-dialog :title="(ruleForm.configId !== 0 ? '修改' : '添加') + '参数'" v-model="isShowDialog" width="600px" :close-on-click-modal="false">
		<el-form :model="ruleForm" ref="formRef" :rules="rules" size="default" label-width="90px">
			<el-form-item label="参数名称" prop="configName">
				<el-input v-model="ruleForm.configName" placeholder="请输入参数名称" />
			</el-form-item>
			<el-form-item label="参数键名" prop="configKey">
				<el-input v-model="ruleForm.configKey" placeholder="请输入参数键名" />
			</el-form-item>
			<el-form-item label="参数键值" prop="configValue">
				<el-input v-model="ruleForm.configValue" placeholder="请输入参数键值" />
			</el-form-item>
			<el-form-item label="系统内置" prop="configType">
				<el-radio-group v-model="ruleForm.configType">
					<el-radio v-for="dict in sysYesNoOptions" :key="dict.value" :value="dict.value">{{ dict.label }}</el-radio>
				</el-radio-group>
			</el-form-item>
			<el-form-item label="备注" prop="remark">
				<el-input v-model="ruleForm.remark" type="textarea" placeholder="请输入内容" />
			</el-form-item>
		</el-form>
		<template #footer>
			<span class="dialog-footer">
				<el-button @click="onCancel" class="dialog-footer-button">取消</el-button>
				<el-button type="primary" @click="onSubmit" class="dialog-footer-button">{{ ruleForm.configId !== 0 ? '修改' : '添加' }}</el-button>
			</span>
		</template>
	</el-dialog>
</template>

<script lang="ts">
import { reactive, toRefs, defineComponent, ref, unref } from 'vue';
import { ElMessage } from 'element-plus';
import { addConfig, editConfig, getConfig } from '/@/api/system/config';
interface RuleFormState {
	configId: number;
	configName: string;
	configKey: string;
	configValue: string;
	configType: string;
	remark: string;
}
interface DicState {
	isShowDialog: boolean;
	ruleForm: RuleFormState;
	rules: {};
}

export default defineComponent({
	name: 'systemEditDicData',
	props: {
		sysYesNoOptions: {
			type: Array,
			default: () => [],
		},
	},
	setup(prop, { emit }) {
		const formRef = ref<HTMLElement | null>(null);
		const state = reactive<DicState>({
			isShowDialog: false,
			ruleForm: {
				configId: 0,
				configName: '',
				configKey: '',
				configValue: '',
				configType: '0',
				remark: '',
			},
			rules: {
				configName: [{ required: true, message: '参数名称不能为空' }],
				configKey: [{ required: true, message: '参数键名不能为空' }],
				configValue: [{ required: true, message: '参数键值不能为空' }],
			},
		});
		// 打开弹窗
		const openDialog = (row: RuleFormState | null) => {
			resetForm();
			if (row) {
				getConfig(row.configId).then((res: any) => {
					const data: RuleFormState = res.data.data || {};
					data.configType = String(data.configType);
					state.ruleForm = data;
				});
				state.ruleForm = row;
			}
			state.isShowDialog = true;
		};
		const resetForm = () => {
			state.ruleForm = {
				configId: 0,
				configName: '',
				configKey: '',
				configValue: '',
				configType: '0',
				remark: '',
			};
		};
		// 关闭弹窗
		const closeDialog = () => {
			state.isShowDialog = false;
		};
		// 取消
		const onCancel = () => {
			closeDialog();
		};
		// 新增
		const onSubmit = () => {
			const formWrap = unref(formRef) as any;
			if (!formWrap) return;
			formWrap.validate((valid: boolean) => {
				if (valid) {
					if (state.ruleForm.configId !== 0) {
						//修改
						editConfig(state.ruleForm).then(() => {
							ElMessage.success('参数修改成功');
							closeDialog(); // 关闭弹窗
							emit('dataList');
						});
					} else {
						//添加
						addConfig(state.ruleForm).then(() => {
							ElMessage.success('参数添加成功');
							closeDialog(); // 关闭弹窗
							emit('dataList');
						});
					}
				}
			});
		};
		return {
			openDialog,
			closeDialog,
			onCancel,
			onSubmit,
			formRef,
			...toRefs(state),
		};
	},
});
</script>
