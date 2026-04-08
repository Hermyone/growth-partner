<!--
 * @Descripttion: 文件上传组件
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-11-15
 * @LastEditors:
 * @LastEditTime:
-->

<template>
	<div class="upload-file">
		<el-upload
			:disabled="disabled"
			:auto-upload="autoUpload"
			:action="action"
			:name="name"
			:data="data"
			:http-request="request"
			v-model:file-list="defaultFileList"
			:show-file-list="showFileList"
			:drag="drag"
			:accept="accept"
			:multiple="multiple"
			:limit="limit"
			:before-upload="before"
			:on-success="success"
			:on-error="error"
			:on-preview="handlePreview"
			:on-exceed="handleExceed">
			<slot>
				<el-button type="primary" :disabled="disabled">Click to upload</el-button>
			</slot>
			<template #tip>
				<div v-if="tip" class="el-upload__tip">{{tip}}</div>
			</template>
		</el-upload>
		<span style="display:none!important"><el-input v-model="value"></el-input></span>
	</div>
</template>

<script lang="ts">
	import { defineComponent, reactive, toRefs, watch, onMounted } from 'vue';
	import { ElMessage, ElNotification, ElMessageBox } from 'element-plus';
	import { uploadConfig, formatArr, toArr, toStr } from './config';

	const props = {
		modelValue: {
			type: [String, Array],
			default: ""
		},
		tip: {
			type: String,
			default: ""
		},
		action: {
			type: String,
			default: ""
		},
		apiObj: {
			type: Object,
			default: () => {}
		},
		name: {
			type: String,
			default: uploadConfig.filename
		},
		data: {
			type: Object,
			default: () => {}
			},
		accept: {
			type: String,
			default: ""
		},
		maxSize: {
			type: Number,
			default: uploadConfig.maxSize
		},
		limit: {
			type: Number,
			default: 0
		},
		autoUpload: {
			type: Boolean,
			default: true
		},
		showFileList: {
			type: Boolean,
			default: true
		},
		drag: {
			type: Boolean,
			default: false
		},
		multiple: {
			type: Boolean,
			default: true
		},
		disabled: {
			type: Boolean,
			default: false
		},
		onSuccess: {
			type: Function,
			default: () => { return true }
		}
	};

	export default defineComponent({
		name: 'q-upload-file',
		props: props,

		setup(props, { emit }) {
			const state = reactive({
				value: "" as any,
				defaultFileList: [] as any
			});

			watch(
				() => props.modelValue,
				(val) => {
					if(Array.isArray(val)){
						if (JSON.stringify(val) != JSON.stringify(formatArr(state.defaultFileList))) {
							state.defaultFileList = val
							state.value = val
						}
					}else{
						if (val != toStr(state.defaultFileList)) {
							state.defaultFileList = toArr(val)
							state.value = val
						}
					}
				}
			);

			watch(
					() => state.defaultFileList,
					(val) => {
						emit('update:modelValue', Array.isArray(props.modelValue) ? formatArr(val) : toStr(val))
						state.value = toStr(val)
					},{
						deep: true
					}
			);

			onMounted(() => {
				state.defaultFileList = Array.isArray(props.modelValue) ? props.modelValue : toArr(props.modelValue);
				state.value = props.modelValue;
			});

			const before = (file: any) => {
				const maxSize = file.size / 1024 / 1024 < props.maxSize;
				if (!maxSize) {
					ElMessage.warning(`上传文件大小不能超过 ${props.maxSize}MB!`);
					return false;
				}
			};

			const success = (res: any, file: any) => {
				let os = props.onSuccess(res, file);
				if(os!=undefined && os==false){
					return false
				}
				let response = uploadConfig.parseData(res);
				file.name = response.fileName;
				file.url = response.src;
			};

			const error = (err: any) => {
				ElNotification({
					title: '上传失败',
					message: err,
					type: 'error',
				});
			};

			const beforeRemove = (uploadFile: any) => {
				return  ElMessageBox.confirm(`是否移除 ${uploadFile.name} ?`, '提示', {
					confirmButtonText: '确认',
					cancelButtonText: '取消',
					type: 'warning',
				})
						.then(() => {
							return true
						})
						.catch(() => {
							return false
						});
			};

			const handleExceed = () => {
				ElMessage.warning(`当前设置最多上传 ${props.limit} 个文件，请移除后再上传!`)
			};

			const handlePreview = (uploadFile: any) => {
				window.open(uploadFile.url);
			};

			const request = (param: any) => {
				if(!props.apiObj) return;

				props.apiObj.post(param, {
					onUploadProgress:  (Progress: { loaded: number; total: number }) => {
						let complete = Math.round((Progress.loaded * 100) / Progress.total)
						param.onProgress({percent: complete})
					}
				}).then((res:any) => {
					var response = uploadConfig.parseData(res);
					if(response.code == uploadConfig.successCode){
						param.onSuccess(res)
					}else{
						param.onError(response.msg || "未知错误")
					}
				}).catch((err:any) => {
					param.onError(err)
				})
			};

			return {
				before,
				success,
				error,
				beforeRemove,
				handleExceed,
				handlePreview,
				request,
				...toRefs(state),
			}
		},
	});
</script>

<style scoped>
	.el-form-item.is-error .upload-file:deep(.el-upload-dragger) { border-color: var(--el-color-danger );}
	.upload-file { width: 100%; }
	.upload-file:deep(.el-upload-list__item) {transition: none !important;}
</style>
