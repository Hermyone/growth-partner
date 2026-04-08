<!--
 * @Descripttion: 图像裁切
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-07-26
 * @LastEditors:
 * @LastEditTime:
-->

<template>
	<el-dialog v-model="isShowDialog" width="769px">
		<template #header="{ close, titleId, titleClass }">
			<div class="cropper-title">
				<h4>{{ title }}</h4>
				<el-tooltip effect="dark" content="重置" placement="bottom">
					<el-button size="small" icon="ele-RefreshLeft" circle @click="onReset" />
				</el-tooltip>

				<el-upload
						ref="uploadRef"
						v-show="uploadUrl !==''"
						:action="uploadUrl"
						:show-file-list="false"
						:on-success="handleAvatarSuccess"
						:on-error="handleAvatarError"
						:before-upload="beforeAvatarUpload"
				>
					<el-tooltip effect="dark" content="更换" placement="bottom">
						<el-button size="small" icon="ele-Plus" circle />
					</el-tooltip>
				</el-upload>
			</div>
		</template>
		<div class="cropper-warp" v-loading="loading">
			<div class="cropper-warp-left">
				<img :src="cropperImg" class="cropper-warp-left-img" />
			</div>
			<div class="cropper-warp-right">
				<div class="cropper-warp-right-title">预览</div>
				<div class="cropper-warp-right-item">
					<div class="cropper-warp-right-value">
						<img :src="cropperImgBase64" class="cropper-warp-right-value-img" />
					</div>
					<div class="cropper-warp-right-label">96 x 96</div>
				</div>
				<div class="cropper-warp-right-item">
					<div class="cropper-warp-right-value">
						<img :src="cropperImgBase64" class="cropper-warp-right-value-img cropper-size" />
					</div>
					<div class="cropper-warp-right-label">72 x 72</div>
				</div>
				<div class="cropper-warp-right-item">
					<div class="cropper-warp-right-value">
						<img :src="cropperImgBase64" class="cropper-warp-right-value-img cropper-size48" />
					</div>
					<div class="cropper-warp-right-label">48 x 48</div>
				</div>
			</div>
		</div>
		<template #footer>
			<span class="dialog-footer">
				<el-button @click="onCancel" size="default">取消</el-button>
				<el-button type="primary" @click="onSubmit" size="default">确认</el-button>
			</span>
		</template>
	</el-dialog>
</template>

<script lang="ts">
	import { reactive, toRefs, nextTick, defineComponent, watch } from 'vue';
	import Cropper from 'cropperjs';
	import 'cropperjs/dist/cropper.css';
	import { ElMessage } from 'element-plus';

	const props = {
		title: { type: String, default: '' },
		// 图像压缩率
		compress: { type: Number, default: 1 },
		// 裁剪框的宽高比 0(null): 自由, 1/1 4/3 16/9
		aspectRatio:  { type: Number, default: NaN },
		// 上传URL地址
		uploadUrl: { type: String, default: ''},
	}

	export default defineComponent({
		name: 'q-cropper',
		props: props,
		setup(props, { emit }) {
			const state = reactive({
				isShowDialog: false,
				cropperImg: '',
				cropperImgBase64: '',
				cropper: null as any,
				loading: false,
        uploadResult: null
			});
			const cropperOptions : any = {
				viewMode: 1,
				dragMode: 'move',
				initialAspectRatio: 1,
				aspectRatio: props.aspectRatio,
				responsive: false,
				preview: '.before',
				background: true,
				autoCropArea: 0.6,
				zoomOnWheel: true,
				minCropBoxWidth: 24,
				minCropBoxHeight: 24,
				checkCrossOrigin: false,
				checkOrientation: false,
				crop: () => {
					state.cropperImgBase64 = (<any>state.cropper).getCroppedCanvas()?.toDataURL('image/jpeg', props.compress);
				},
			};
			// 打开弹窗
			const openDialog = (imgs: any) => {
				state.cropperImg = imgs;
				state.isShowDialog = true;
				nextTick(() => {
					initCropper();
				});
			};
			// 关闭弹窗
			const closeDialog = () => {
				state.isShowDialog = false;
			};
			// 取消
			const onCancel = () => {
				closeDialog();
			};
			// 重置
			const onReset = () => {
				state.cropper?.reset();
			};
			// 确认
			const onSubmit = () => {
				// state.cropperImgBase64 = state.cropper.getCroppedCanvas()?.toDataURL('image/jpeg', props.compress);
				emit('cropperChange', state.cropperImgBase64, state.uploadResult);
				closeDialog();
			};
			// 初始化cropperjs图片裁剪
			const initCropper = () => {
				if(state.cropper) return;

				const letImg: any = document.querySelector('.cropper-warp-left-img');
				(<any>state.cropper) = new Cropper(letImg, cropperOptions);
			};

			// 设置裁剪框的宽高比
			const setAspectRatio = (val : any) => {
				state.cropper?.setAspectRatio(val);
			};

			// 返回裁切数据
			const getCropData = (fun: any, type : string ='image/jpeg') => {
				fun(state.cropper?.getCroppedCanvas()?.toDataURL(type, props.compress));
			};
			const getCropBlob = (fun: any, type : string ='image/jpeg') => {
				state.cropper.getCroppedCanvas().toBlob((blob:any) => {
					fun(blob)
				}, type, props.compress);
			};
			const getCropFile = (fun: any, fileName : string = 'cropper.jpg', type : string ='image/jpeg') => {
				state.cropper.getCroppedCanvas().toBlob((blob:any) => {
					let file = new File([blob], fileName, {type: type})
					fun(file)
				}, type, props.compress);
			};

			// 上传成功
			const handleAvatarSuccess = (res:any, file:any) => {
				state.cropperImg = URL.createObjectURL(file.raw);
				state.cropper.replace(state.cropperImg);
				state.loading = false;
        state.uploadResult = res
			};
			// 上传失败
			const handleAvatarError = (err: any) => {
				state.loading = false;
				ElMessage.error("文件上传失败!");
				console.error('文件上传失败', err);
			};
			// 上传校验
			const beforeAvatarUpload = (file:any) => {
				const isJPG = (file.type === 'image/jpeg' || file.type === 'image/png');
				const isLt2M = file.size / 1024 / 1024 < 2;
				if (!isJPG) {
					ElMessage.error('上传头像图片只能是 JPG/PNG 格式!');
				}
				if (!isLt2M) {
					ElMessage.error('上传头像图片大小不能超过 2MB!');
				}
				return isJPG && isLt2M;
			};

			watch(
				() => props.aspectRatio,
				(val) => {
					state.cropper?.setAspectRatio(val);
				}
			);

			return {
				openDialog,
				closeDialog,
				onCancel,
				onReset,
				onSubmit,
				initCropper,
				setAspectRatio,
				getCropData,
				getCropBlob,
				getCropFile,
				handleAvatarSuccess,
				handleAvatarError,
				beforeAvatarUpload,
				...toRefs(state),
			};
		},
	});
</script>

<style scoped lang="scss">
	.cropper-title {
		display: flex;
		gap: 5px;
		h4 {
			font-size: var(--el-dialog-title-font-size);
			padding-right: 10px;
		}
		:deep(.el-button i.el-icon,.el-button i.iconfont,.el-button i.fa,.el-button--default i.iconfont,.el-button--default i.fa){
			font-size: 14px !important;
			margin-right: 0!important;
		}
	}
	.cropper-warp {
		display: flex;
		.cropper-warp-left {
			position: relative;
			display: inline-block;
			height: 350px;
			flex: 1;
			border: 1px solid var(--el-border-color);
			background: var(--el-color-white);
			overflow: hidden;
			background-repeat: no-repeat;
			cursor: move;
			border-radius: var(--el-border-radius-base);
			.cropper-warp-left-img {
				width: 100%;
				height: 100%;
			}
		}
		.cropper-warp-right {
			width: 150px;
			height: 350px;
			.cropper-warp-right-title {
				text-align: center;
				height: 20px;
				line-height: 20px;
				font-weight: bold;
			}
			.cropper-warp-right-item {
				margin: 10px 0;
				.cropper-warp-right-value {
					display: flex;
					.cropper-warp-right-value-img {
						width: 96px;
						height: 96px;
						border-radius: var(--el-border-radius-circle);
						margin: auto;
					}
					.cropper-size {
						width: 72px;
						height: 72px;
					}
					.cropper-size48 {
						width: 48px;
						height: 48px;
					}
				}
				.cropper-warp-right-label {
					text-align: center;
					font-size: 12px;
					color: var(--el-text-color-primary);
					height: 30px;
					line-height: 30px;
				}
			}
		}
	}
</style>
