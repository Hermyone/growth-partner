<template>
	<q-dialog v-model="isOpen" width="1024px" :close-on-click-modal="false">
		<template #header>
			<div>
				{{ (isEdit ? '修改' : '添加') }}
			</div>
		</template>
		<el-form ref="formRef" :model="formData" :rules="rules" size="default" label-width="100px">
			<el-form-item label="消息标题" prop="title" class="w100">
				<el-input v-model="formData.title" placeholder="请输入消息标题" clearable/>
			</el-form-item>
			<el-form-item label="" prop="kind" class="w100">
				<el-radio-group v-model="formData.kind">
					<el-radio :value='0' size="large">系统广播</el-radio>
					<el-radio :value='1' size="large">站内信</el-radio>
				</el-radio-group>
				<el-switch style="padding-left: 30px"
						   v-model="formData.level"
						   inline-prompt
						   :active-value="1"
						   :inactive-value="0"
						   active-text="重要消息"
						   inactive-text="一般消息"
				/>
			</el-form-item>
			<el-form-item label="消息内容" prop="content" class="w100">
				<EditView ref="editRef" class="w100"></EditView>
			</el-form-item>
		</el-form>
		<template #footer>
			<span class="dialog-footer">
				<el-button @click="onCancel" class="dialog-footer-button">取消</el-button>
				<el-button type="primary" @click="onSubmit" class="dialog-footer-button" :loading="loading">保存</el-button>
			</span>
		</template>
	</q-dialog>
</template>

<script lang="ts">
	import { reactive, toRefs, defineComponent, ref, unref, nextTick, getCurrentInstance } from 'vue';
	import { ElMessage } from 'element-plus';
	import { addMessage, editMessage } from '/@/api/system/message';
	import EditView from '/@/views/components/markdown/editor.vue';

	export default defineComponent({
		components: { EditView },
		setup(props, { emit }) {
			const { proxy } = <any>getCurrentInstance();
			const formRef = ref<HTMLElement | null>(null);
			const editRef = ref();
			const state = reactive({
				loading: false,
				isOpen: false,
				isEdit: false,
				formData: {
					title: '',
					content: '',
					kind: 0,
					level: 0,
				} as any,
				rules: {
					title: [{ required: true, message: '消息标题不能为空' }],
				},
			});
			// 打开弹窗
			const openDialog = (row: any) => {
				state.formData = {kind: 0};
				state.isEdit = false;
				if (row) {
					state.isEdit = true;
					state.formData = row;
					setTimeout(() => {
						editRef.value?.setContent(row.content);
					}, 500)
				}
				state.isOpen = true;
			};
			const closeDialog = () => {
				state.isOpen = false;
			};

			// 取消
			const onCancel = () => {
				closeDialog();
			};
			// 保存
			const onSubmit = () => {
				let content = editRef.value.getContent();
				if(content == "") {
					ElMessage.error("消息内容不能为空");
					return;
				}
				const formWrap = unref(formRef) as any;
				if (!formWrap) return;
				formWrap.validate((valid: boolean) => {
					if (valid) {
						state.loading = true;
						state.formData.content = content;
						if (!state.isEdit) {
							//添加
							addMessage(state.formData)
									.then(() => {
										ElMessage.success('数据新增成功');
										closeDialog();
										emit('onSuccess');
									})
									.finally(() => {
										state.loading = false;
									});
						} else {
							//修改
							editMessage(state.formData)
									.then(() => {
										ElMessage.success('数据更新成功');
										closeDialog();
										emit('onSuccess');
									})
									.finally(() => {
										state.loading = false;
									});
						}
					}
				});
			};

			return {
				openDialog,
				onCancel,
				onSubmit,
				formRef,
				editRef,
				...toRefs(state),
			};
		},
});
</script>

<style scoped lang="scss">
</style>
