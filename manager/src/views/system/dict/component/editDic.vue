<template>
	<el-dialog :title="(ruleForm.dictId !== 0 ? '修改' : '添加') + '字典'" v-model="isShowDialog" width="600px" :close-on-click-modal="false">
		<el-form :model="ruleForm" ref="formRef" :rules="rules" size="default" label-width="90px" status-icon>
			<el-form-item label="字典名称" prop="dictName">
				<el-input v-model="ruleForm.dictName" placeholder="请输入字典名称" clearable />
			</el-form-item>
			<el-form-item label="字典类型" prop="dictType">
				<el-input v-model="ruleForm.dictType" placeholder="请输入字典类型" clearable />
			</el-form-item>
			<el-form-item label="状态" prop="status">
				<el-radio-group v-model="ruleForm.status">
					<el-radio :value="1">启用</el-radio>
					<el-radio :value="0">禁用</el-radio>
				</el-radio-group>
			</el-form-item>
			<el-form-item label="备注" prop="remark">
				<el-input v-model="ruleForm.remark" type="textarea" placeholder="请输入内容" ></el-input>
			</el-form-item>
		</el-form>
		<template #footer>
			<span class="dialog-footer">
				<el-button @click="onCancel" class="dialog-footer-button">取消</el-button>
				<el-button type="primary" @click="onSubmit" class="dialog-footer-button">{{ ruleForm.dictId !== 0 ? '修改' : '添加' }}</el-button>
			</span>
		</template>
	</el-dialog>
</template>

<script lang="ts">
import { reactive, toRefs, defineComponent, ref, unref } from 'vue';
import { getType, addType, editType } from '/@/api/system/dict/type';
import { ElMessage } from 'element-plus';
import { DictStateEx, RuleFormDictStateEx } from '/@/types/model';

export default defineComponent({
	name: 'systemEditDic',
	setup(prop, { emit }) {
		const formRef = ref<HTMLElement | null>(null);
		const state = reactive<DictStateEx>({
			isShowDialog: false,
			ruleForm: {
				dictId: 0,
				dictName: '',
				dictType: '',
				status: 1,
				remark: '',
			},
			rules: {
				dictName: [{ required: true, message: '字典名称不能为空' }],
				dictType: [{ required: true, message: '字典类型不能为空' }],
			},
		});
		// 打开弹窗
		const openDialog = (row: RuleFormDictStateEx | null) => {
			resetForm();
			if (row) {
				getType(row.dictId).then((res: any) => {
					state.ruleForm = res.data.dictType;
				});
				state.ruleForm = row;
			}
			state.isShowDialog = true;
		};
		const resetForm = () => {
			state.ruleForm = {
				dictId: 0,
				dictName: '',
				dictType: '',
				status: 1,
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
					if (state.ruleForm.dictId !== 0) {
						//修改
						editType(state.ruleForm).then(() => {
							ElMessage.success('字典类型修改成功');
							// 关闭弹窗
							closeDialog(); 
							emit('typeList');
						});
					} else {
						//添加
						addType(state.ruleForm).then(() => {
							ElMessage.success('字典类型添加成功');
							// 关闭弹窗
							closeDialog(); 
							emit('typeList');
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
