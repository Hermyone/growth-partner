<!--
 * @Descripttion: 文件预览组件,支持word,excel,pdf,image
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-10-09
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <el-dialog class="dialog-style" ref="dialogRef" v-model="dialogVisible" :fullscreen="true" v-bind="$attrs" :show-close="true">
        <template #header>
            <slot name="header">
                <span class="dialog-style-title">{{ title }}</span>
            </slot>
        </template>

        <component :is="previews[fileType]" ref="previewRef" />

        <template #footer>
            <div style="float: left">{{ url }}</div>
        </template>
    </el-dialog>
</template>

<script lang="ts">
	import { defineAsyncComponent, defineComponent, getCurrentInstance, reactive, toRefs } from 'vue';
	import { ElMessage } from 'element-plus';

    const props = {
        title: { type: String, default: "预览" },
    };

	export default defineComponent({
		name: 'q-file-preview',
        props: props,
        setup(props, { emit }) {
            const previews: any = {
                docx: defineAsyncComponent(() => import('./component/docxPreview.vue')),
                xlsx: defineAsyncComponent(() => import('./component/xlsxPreview.vue')),
                pdf:  defineAsyncComponent(() => import('./component/pdfPreview.vue')),
            };

            const { proxy } = <any>getCurrentInstance();
            const state = reactive({
                fileType: '',
                url: '',
                dialogVisible: false,
            });

            // 打开文件 type:‘docx’，‘xlsx’，‘pdf’
            const open = (url: string, type: string) => {
            	state.fileType = type;
            	state.url = url;
                state.dialogVisible = true;

                setTimeout(() => {
                	if(proxy.$refs.previewRef){
						proxy.$refs.previewRef.open(url);
                    }else{
                		ElMessage.error('未找到该类型文件的预览插件')
                    }

                }, 500);
            };

            return {
                previews,
                open,
                ...toRefs(state),
            }
        },
	});
</script>

<style scoped lang="scss">
    .dialog-style {
        &-title {
            color: var(--el-text-color-primary);
            font-size: var(--el-dialog-title-font-size);
            font-weight: bold;
        }

        :deep(.el-dialog).is-fullscreen {
            display: flex;
            flex-direction: column;
            top: 0 !important;
            left: 0 !important;
        }
        :deep(.el-dialog).is-fullscreen .el-dialog__header {}
        :deep(.el-dialog).is-fullscreen .el-dialog__body {
            flex: 1;
            overflow: auto;
            max-height: calc(100vh - 15px) !important;
        }
    }
</style>