<!--
 * @Descripttion: 多个文件上传组件
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-11-15
 * @LastEditors:
 * @LastEditTime:
-->

<template>
	<div class="upload-multiple">
		<el-upload ref="uploader" list-type="picture-card"
			:auto-upload="autoUpload"
			:disabled="disabled"
			:action="action"
			:name="name"
			:data="data"
			:http-request="request"
			v-model:file-list="defaultFileList"
			:show-file-list="showFileList"
			:accept="accept"
			:multiple="multiple"
			:limit="limit"
			:before-upload="before"
			:on-success="success"
			:on-error="error"
			:on-preview="handlePreview"
			:on-exceed="handleExceed">
			<slot>
				<el-icon><el-icon-plus/></el-icon>
			</slot>
			<template #tip>
				<div v-if="tip" class="el-upload__tip">{{tip}}</div>
			</template>
			<template #file="{ file }">
				<div class="upload-list-item">
					<el-image class="el-upload-list__item-thumbnail" :src="file.url" fit="cover" :preview-src-list="preview" :initial-index="preview.findIndex(n=>n==file.url)" hide-on-click-modal append-to-body :z-index="9999">
						<template #placeholder>
							<div class="upload-multiple-image-slot">
								Loading...
							</div>
						</template>
					</el-image>
					<div v-if="!disabled && file.status=='success'" class="upload__item-actions">
						<span class="del" @click="handleRemove(file)"><el-icon><el-icon-delete /></el-icon></span>
					</div>
					<div v-if="file.status=='ready' || file.status=='uploading'" class="upload__item-progress">
						<el-progress :percentage="file.percentage" :text-inside="true" :stroke-width="16"/>
					</div>
				</div>
			</template>
		</el-upload>
		<span style="display:none!important"><el-input v-model="value"></el-input></span>
	</div>
</template>

<script>
	import { ElMessage, ElNotification, ElMessageBox } from 'element-plus';
	import { uploadConfig, formatArr, toArr, toStr } from './config';
	import Sortable from 'sortablejs';

	export default {
		name: 'q-upload-multiple',
		props: {
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
				default: () => {
				}
			},
			name: {
				type: String,
				default: uploadConfig.filename
			},
			data: {
				type: Object,
				default: () => {
				}
			},
			accept: {
				type: String,
				default: "image/gif, image/jpeg, image/png"
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
			multiple: {
				type: Boolean,
				default: true
			},
			disabled: {
				type: Boolean,
				default: false
			},
			draggable: {
				type: Boolean,
				default: false
			},
			onSuccess: {
				type: Function,
				default: () => {
					return true
				}
			}
		},
		data(){
			return {
				value: "",
				defaultFileList: []
			}
		},
		watch:{
			modelValue(val){
				if(Array.isArray(val)){
					if (JSON.stringify(val) != JSON.stringify(formatArr(this.defaultFileList))) {
						this.defaultFileList = val
						this.value = val
					}
				}else{
					if (val != toStr(this.defaultFileList)) {
						this.defaultFileList = toArr(val)
						this.value = val
					}
				}
			},
			defaultFileList: {
				handler(val){
					this.$emit('update:modelValue', Array.isArray(this.modelValue) ? formatArr(val) : toStr(val))
					this.value = toStr(val)
				},
				deep: true
			}
		},
		computed: {
			preview(){
				return this.defaultFileList.map(v => v.url)
			}
		},
		mounted() {
			this.defaultFileList = Array.isArray(this.modelValue) ? this.modelValue : toArr(this.modelValue)
			this.value = this.modelValue
			if(!this.disabled && this.draggable){
				this.rowDrop()
			}
		},
		methods: {
			// 拖拽
			rowDrop(){
				const _this = this
				const itemBox = this.$refs.uploader.$el.querySelector('.el-upload-list')
				Sortable.create(itemBox, {
					handle: ".el-upload-list__item",
					animation: 200,
					ghostClass: "ghost",
					onEnd({ newIndex, oldIndex }) {
						const tableData = _this.defaultFileList
						const currRow = tableData.splice(oldIndex, 1)[0]
						tableData.splice(newIndex, 0, currRow)
					}
				})
			},
			before(file){
				if(!['image/jpeg','image/png','image/gif'].includes(file.type)){
					ElMessage.warning(`选择的文件类型 ${file.type} 非图像类文件`);
					return false;
				}
				const maxSize = file.size / 1024 / 1024 < this.maxSize;
				if (!maxSize) {
					ElMessage.warning(`上传文件大小不能超过 ${this.maxSize}MB!`);
					return false;
				}
			},
			success(res, file){
				var os = this.onSuccess(res, file)
				if(os!=undefined && os==false){
					return false
				}
				var response = uploadConfig.parseData(res)
				file.name = response.fileName
				file.url = response.src
			},
			error(err){
				ElNotification({
					title: '上传失败',
					message: err,
					type: 'error',
				});
			},
			beforeRemove(uploadFile){
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
			},
			handleRemove(file){
				this.$refs.uploader.handleRemove(file)
				//this.defaultFileList.splice(this.defaultFileList.findIndex(item => item.uid===file.uid), 1)
			},
			handleExceed(){
				ElMessage.warning(`当前设置最多上传 ${this.limit} 个文件，请移除后再上传!`)
			},
			handlePreview(uploadFile){
				window.open(uploadFile.url)
			},
			request(param){
				if(!this.apiObj) return;

				this.apiObj.post(param, {
					onUploadProgress: e => {
						const complete = parseInt(((e.loaded / e.total) * 100) | 0, 10)
						param.onProgress({percent: complete})
					}
				}).then(res => {
					var response = uploadConfig.parseData(res);
					if(response.code == uploadConfig.successCode){
						param.onSuccess(res)
					}else{
						param.onError(response.msg || "未知错误")
					}
				}).catch(err => {
					param.onError(err)
				})
			}
		}
	}
</script>

<style scoped>
	.el-form-item.is-error .upload-multiple:deep(.el-upload--picture-card) {border-color: var(--el-color-danger);}
	:deep(.el-upload-list__item) {transition:none;border-radius: 0;}
	.upload-multiple:deep(.el-upload-list__item.el-list-leave-active) {position: static!important;}
	.upload-multiple:deep(.el-upload--picture-card) {border-radius: 0;}
	.upload-list-item {width: 100%;height: 100%;position: relative;}
	.upload-multiple .el-image {display: block;}
	.upload-multiple .el-image:deep(img) {-webkit-user-drag: none;}
	.upload-multiple-image-slot {display: flex;justify-content: center;align-items: center;width: 100%;height: 100%;font-size: 12px;}
	.upload-multiple .el-upload-list__item:hover .upload__item-actions{display: block;}
	.upload__item-actions {position: absolute;top:0;right: 0;display: none;}
	.upload__item-actions span {display: flex;justify-content: center;align-items: center;width: 25px;height:25px;cursor: pointer;color: #fff;}
	.upload__item-actions span i {font-size: 12px;}
	.upload__item-actions .del {background: var(--el-color-primary);}
	.upload__item-progress {position: absolute;width: 100%;height: 100%;top: 0;left: 0;background-color: var(--el-overlay-color-lighter);}
</style>
