<!--
 * @Descripttion: 上传组件
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-11-15
 * @LastEditors:
 * @LastEditTime:
-->

<template>
	<div class="upload" :class="{'upload-round':round}" :style="style">
		<div v-if="file && file.status != 'success'" class="upload__uploading">
			<div class="upload__progress">
				<el-progress :percentage="file.percentage" :text-inside="true" :stroke-width="16"/>
			</div>
			<el-image class="image" :src="file.tempFile" fit="cover"></el-image>
		</div>
		<div v-if="file && file.status=='success'" class="upload__img">
			<el-image class="image" :src="file.url" :preview-src-list="[file.url]" fit="cover" hide-on-click-modal append-to-body :z-index="9999">
				<template #placeholder>
					<div class="upload__img-slot">
						Loading...
					</div>
				</template>
			</el-image>
			<div class="upload__img-actions" v-if="!disabled">
				<span class="del" @click="handleRemove()"><el-icon><el-icon-delete /></el-icon></span>
			</div>
		</div>
		<el-upload v-if="!file" class="uploader" ref="uploader"
			:auto-upload="cropper?false:autoUpload"
			:disabled="disabled"
			:show-file-list="showFileList"
			:action="action"
			:name="name"
			:data="data"
			:accept="accept"
			:limit="1"
			:http-request="request"
			:on-change="change"
			:before-upload="before"
			:on-success="success"
			:on-error="error"
			:on-exceed="handleExceed">
			<slot>
				<div class="el-upload--picture-card">
					<div class="file-empty">
						<el-icon><component :is="icon" /></el-icon>
						<h4 v-if="title">{{title}}</h4>
					</div>
				</div>
			</slot>
		</el-upload>
		<span style="display:none!important"><el-input v-model="value"></el-input></span>

		<q-cropper ref="cropperRef" title="裁切" :compress="compress" :aspectRatio="aspectRatio" @cropperChange="cropperSave"></q-cropper>
	</div>
</template>

<script>
	import { genFileId, ElMessage, ElNotification } from 'element-plus';
	import { uploadConfig } from './config';

	export default {
		name: 'q-upload',
		props: {
			modelValue: { type: String, default: "" },
			height: {type: Number, default: 148},
			width: {type: Number, default: 148},
			title: { type: String, default: "" },
			icon: { type: String, default: "el-icon-plus" },
			action: { type: String, default: "" },
			apiObj: { type: Object, default: () => {} },
			name: { type: String, default: uploadConfig.filename },
			data: { type: Object, default: () => {} },
			accept: { type: String, default: "image/gif, image/jpeg, image/png" },
			maxSize: { type: Number, default: uploadConfig.maxSize },
			limit: { type: Number, default: 1 },
			autoUpload: { type: Boolean, default: true },
			showFileList: { type: Boolean, default: false },
			disabled: { type: Boolean, default: false },
			round: { type: Boolean, default: false },
			onSuccess: { type: Function, default: () => { return true } },

			cropper: { type: Boolean, default: false },
			compress: {type: Number, default: 1},
			aspectRatio:  {type: Number, default: NaN}
		},
		data() {
			return {
				value: "",
				file: null,
				style: {
					width: this.width + "px",
					height: this.height + "px"
				},
				cropperFile: null
			}
		},
		watch:{
			modelValue(val){
				this.value = val
				this.newFile(val)
			},
			value(val){
				this.$emit('update:modelValue', val)
			}
		},
		mounted() {
			this.value = this.modelValue
			this.newFile(this.modelValue)
		},
		methods: {
			newFile(url){
				if(url){
					this.file = {
						status: "success",
						url: url
					}
				}else{
					this.file = null
				}
			},
			cropperSave(){
				this.$refs.cropperRef.getCropFile(file => {

					file.uid = this.cropperFile.uid
					this.cropperFile.raw = file

					this.file = this.cropperFile
					this.file.tempFile = URL.createObjectURL(this.file.raw)
					this.$refs.uploader.submit()

				}, this.cropperFile.name, this.cropperFile.type)
			},
			cropperClosed(){
				URL.revokeObjectURL(this.cropperFile.tempCropperFile)
				delete this.cropperFile.tempCropperFile
			},
			handleRemove(){
				this.clearFiles()
			},
			clearFiles(){
				URL.revokeObjectURL(this.file.tempFile)
				this.value = ""
				this.file = null
				this.$nextTick(()=>{
					this.$refs.uploader.clearFiles()
				})
			},
			change(file,files){
				if(files.length > 1){
					files.splice(0, 1)
				}
				if(this.cropper && file.status=='ready'){
					const acceptIncludes = ["image/gif", "image/jpeg", "image/png"].includes(file.raw.type)
					if(!acceptIncludes){
						ElNotification({
							title: '上传文件警告',
							message: '选择的文件非图像类文件'
						})
						return false
					}

					this.cropperFile = file
					this.cropperFile.tempCropperFile = URL.createObjectURL(file.raw)
					this.$refs.cropperRef.openDialog(this.cropperFile.tempCropperFile)
					return false
				}
				this.file = file
				if(file.status=='ready'){
					file.tempFile = URL.createObjectURL(file.raw)
				}
			},
			before(file){
				const acceptIncludes = this.accept.replace(/\s/g,"").split(",").includes(file.type)
				if(!acceptIncludes){
					ElNotification({
						title: '上传文件警告',
						message: '选择的文件非图像类文件'
					})
					this.clearFiles()
					return false
				}
				const maxSize = file.size / 1024 / 1024 < this.maxSize;
				if (!maxSize) {
					ElMessage.warning(`上传文件大小不能超过 ${this.maxSize}MB!`);
					this.clearFiles()
					return false
				}
			},
			handleExceed(files){
				const file = files[0]
				file.uid = genFileId()
				this.$refs.uploader.handleStart(file)
			},
			success(res, file){
				//释放内存删除blob
				URL.revokeObjectURL(file.tempFile)
				delete file.tempFile
				var os = this.onSuccess(res, file)
				if(os!=undefined && os==false){
					this.$nextTick(() => {
						this.file = null
						this.value = ""
					})
					return false
				}
				var response = uploadConfig.parseData(res)
				file.url = response.src
				this.value = file.url
			},
			error(err){
				this.$nextTick(()=>{
					this.clearFiles()
				})
				ElNotification({
					title: '上传失败',
					message: err,
					type: 'error',
				});
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
	.el-form-item.is-error .upload .el-upload--picture-card {border-color: var(--el-color-danger);}
	.upload .el-upload--picture-card {border-radius: 0;}

	.upload .uploader,.upload:deep(.el-upload) {width: 100%;height: 100%;}

	.upload__img {width: 100%;height: 100%;position: relative;}
	.upload__img .image {width: 100%;height: 100%;}
	.upload__img-actions {position: absolute;top:0;right: 0;display: none;}
	.upload__img-actions span {display: flex;justify-content: center;align-items: center;width: 25px;height:25px;cursor: pointer;color: #fff;}
	.upload__img-actions span i {font-size: 12px;}
	.upload__img-actions .del {background: var(--el-color-primary);}
	.upload__img:hover .upload__img-actions {display: block;}
	.upload__img-slot {display: flex;justify-content: center;align-items: center;width: 100%;height: 100%;font-size: 12px;background-color: var(--el-fill-color-lighter);}

	.upload__uploading {width: 100%;height: 100%;position: relative;}
	.upload__progress {position: absolute;width: 100%;height: 100%;display: flex;justify-content: center;align-items: center;background-color: var(--el-overlay-color-lighter);z-index: 1;padding:10px;}
	.upload__progress .el-progress {width: 100%;}
	.upload__uploading .image {width: 100%;height: 100%;}

	.upload .file-empty {width: 100%;height: 100%;display: flex;justify-content: center;align-items: center;flex-direction: column;}
	.upload .file-empty i {font-size: 28px;}
	.upload .file-empty h4 {font-size: 12px;font-weight: normal;color: #8c939d;margin-top: 8px;}

	.upload.upload-round {border-radius: 50%;overflow: hidden;}
	.upload.upload-round .el-upload--picture-card {border-radius: 50%;}
	.upload.upload-round .upload__img-actions {top: auto;left: 0;right: 0;bottom: 0;}
	.upload.upload-round .upload__img-actions span {width: 100%;}
</style>
