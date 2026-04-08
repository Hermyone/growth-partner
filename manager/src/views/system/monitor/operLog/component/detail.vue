<template>
	<!-- 操作日志详情抽屉 -->
	<div class="system-sysOperLog-detail">
		<el-drawer v-model="isShowDialog" size="600px" direction="rtl">
			<template #header>
				<h4>日志详情</h4>
			</template>
			<el-main style="padding:15px;">
				<el-descriptions :column="1" border size="small">
					<el-descriptions-item label="日志编号">{{ formData.operId }}</el-descriptions-item>
					<el-descriptions-item label="操作人员">{{ formData.operName }}</el-descriptions-item>
					<el-descriptions-item label="操作方法">{{ formData.method }}</el-descriptions-item>
					<el-descriptions-item label="系统模块">{{ formData.title }}</el-descriptions-item>
					<el-descriptions-item label="日志时间">{{ proxy.parseTime(formData.operTime, '{y}-{m}-{d} {h}:{i}:{s}') }}</el-descriptions-item>
					<el-descriptions-item label="操作地点">{{ formData.operLocation }} ({{ formData.operIp }}) </el-descriptions-item>
				</el-descriptions>

				<el-collapse style="margin-top: 20px;" v-model="activeCollapse">
					<el-collapse-item title="详细" name="more">
						<p class="code">
							请求方式：
							<q-tag v-if="formData.requestMethod === 'GET'">{{ proxy.getOptionValue(formData.requestMethod, requestMethodOptions, 'value', 'label') }} </q-tag>
							<q-tag type="success" v-if="formData.requestMethod === 'POST'">{{ proxy.getOptionValue(formData.requestMethod, requestMethodOptions, 'value', 'label') }} </q-tag>
							<q-tag type="warning" v-if="formData.requestMethod === 'PUT'">{{ proxy.getOptionValue(formData.requestMethod, requestMethodOptions, 'value', 'label') }} </q-tag>
							<q-tag type="danger" v-if="formData.requestMethod === 'DELETE'">{{ proxy.getOptionValue(formData.requestMethod, requestMethodOptions, 'value', 'label') }} </q-tag>
							<br/>请求URL：<br/>
							{{ formData.operUrl }}<br/>
							<br/>请求参数：<br/>
							{{ formData.operParam }}<br/>
							<br>错误消息：<br/>
							{{ formData.errorMsg }}<br/>
						</p>
					</el-collapse-item>
				</el-collapse>
			</el-main>
		</el-drawer>
	</div>
</template>
<script lang="ts">
import { reactive, toRefs, defineComponent, ref, getCurrentInstance } from 'vue';
import { getSysOperLog } from '/@/api/system/monitor/operLog';
import { TableOperLogEditState, TableOperLogInfoData } from '/@/types/model';

export default defineComponent({
	name: 'apiV1SystemSysOperLogDetail',
	components: {},
	props: {
		requestMethodOptions: {
			type: Array,
			default: () => [],
		},
	},
	setup(props, { emit }) {
		const { proxy } = <any>getCurrentInstance();
		const formRef = ref<HTMLElement | null>(null);
		const activeCollapse = ref("more");
		const state = reactive<TableOperLogEditState>({
			loading: false,
			isShowDialog: false,
			formData: {
				operId: undefined,
				title: undefined,
				businessType: undefined,
				method: undefined,
				requestMethod: undefined,
				operatorType: undefined,
				operName: undefined,
				deptName: undefined,
				operUrl: undefined,
				operIp: undefined,
				operLocation: undefined,
				operParam: undefined,
				jsonResult: undefined,
				status: false,
				errorMsg: undefined,
				operTime: undefined,
				linkedSysOperLogSysDept: {
					deptId: undefined, 		// 部门id
					deptName: undefined, 	// 部门名称
				},
			},
			rules: {},
		});
		// 打开弹窗
		const openDialog = (row?: TableOperLogInfoData) => {
			resetForm();
			if (row) {
				getSysOperLog(row.operId!).then((res: any) => {
					const data = res.data;
					state.formData = data;
				});
			}
			state.isShowDialog = true;
		};
		// 关闭弹窗
		const closeDialog = () => {
			state.isShowDialog = false;
		};

		const resetForm = () => {
			state.formData = {
				operId: undefined,
				title: undefined,
				businessType: undefined,
				method: undefined,
				requestMethod: undefined,
				operatorType: undefined,
				operName: undefined,
				deptName: undefined,
				operUrl: undefined,
				operIp: undefined,
				operLocation: undefined,
				operParam: undefined,
				jsonResult: undefined,
				status: false,
				errorMsg: undefined,
				operTime: undefined,
				linkedSysOperLogSysDept: {
					deptId: undefined, 		// 部门id
					deptName: undefined, 	// 部门名称
				},
			};
		};
		//关联sys_dept表选项
		const getSysDeptItemsDeptName = () => {
			emit('getSysDeptItemsDeptName');
		};
		return {
			proxy,
			openDialog,
			closeDialog,
			formRef,
			activeCollapse,
			getSysDeptItemsDeptName,
			...toRefs(state),
		};
	},
});

</script>
<style scoped>
.system-sysOperLog-detail :deep(.el-form-item--large .el-form-item__label) {
	font-weight: bolder;
}
.pic-block {
	margin-right: 8px;
}
.file-block {
	width: 100%;
	border: 1px solid var(--el-border-color);
	border-radius: 6px;
	cursor: pointer;
	position: relative;
	overflow: hidden;
	transition: var(--el-transition-duration-fast);
	margin-bottom: 5px;
	padding: 3px 6px;
}
.ml-2 {
	margin-right: 5px;
}
.code {
	background: #848484;
	padding:15px;
	color: #fff;
	font-size: 12px;
	border-radius: 4px;
}
</style>
