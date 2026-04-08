<template>
    <el-dialog v-model="isOpen" width="600px" :close-on-click-modal="false">
        <template #header>
            上传图片
        </template>
        <q-upload style="width: 550px; height: 400px;"
                  v-model="url"
                  :apiObj="uploadApi"
                  :on-success="onUploadSuccess"
                  title="点击选择图像"
                  icon="el-icon-picture"
        >
        </q-upload>
        <template #footer>
			<span class="dialog-footer">
				<el-button @click="isOpen = false" class="dialog-footer-button">取消</el-button>
				<el-button type="primary" @click="onSubmit" class="dialog-footer-button" :loading="loading">{{caption}}</el-button>
			</span>
        </template>
    </el-dialog>
</template>

<script lang="ts">
	import { defineComponent, reactive, toRefs, ref } from 'vue';
    import { uplaodFile } from '/@/api/system/file';
    import { ElMessage } from 'element-plus';
		import Cookies from 'js-cookie';

    const props = {
        caption: {
            type: String,
            default: "保存",
        },
    };

	export default defineComponent({
		name: 'q-page-upload-image',
		props: props,
        setup(props, {emit}) {
            const userName = Cookies.get('username');
			const state = reactive({
                loading: false,
                isOpen: false,
                params: {} as any,
                uploadApi: {
                    name: "文件上传",
                    post: function(data:any, params = {}) {
                    	return uplaodFile(data.file, "upload/worker/", userName + "_" + getCurrentDate() + "_" + data.file.name)
                    },
                },
                url: "",
			});

			const getCurrentDate = () => {
                let date = new Date(Date.now());

                // 获取年、月、日
                const year = date.getFullYear();
                const month = String(date.getMonth() + 1).padStart(2, '0');
                const day = String(date.getDate()).padStart(2, '0');
                return `${year}${month}${day}`
            };

            const openDialog = (params: any) => {
                state.isOpen = true;
                state.url = "";
                state.params = params;
            };

            const onSubmit = () => {
            	if(state.url == ""){
            		ElMessage.error("请先上传图像");
            		return
                }
            	emit("onSubmit", state.params, state.url);
                state.isOpen = false;
            };

            const onUploadSuccess = (val: any) => {
            	state.url = val?.url || "";
            };

			return {
                openDialog,
                onSubmit,
                onUploadSuccess,
                getCurrentDate,
				...toRefs(state),
			}
		},
	});

</script>

<style scoped lang="scss">
</style>