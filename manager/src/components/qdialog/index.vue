<!--
 * @Descripttion: 自定义对话框
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-07-24
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <el-dialog ref="dialogRef" class="dialog-style" v-model="dialogVisible" :fullscreen="isFullscreen" :close-on-click-modal="!modalFlag" append-to-body v-bind="$attrs" :show-close="false">
        <template #header>
            <slot name="header">
                <span class="dialog-style-title">{{ title }}</span>
            </slot>

            <div class="dialog-style-header">
                <button v-if="showFullscreen" aria-label="fullscreen" type="button" @click="setFullscreen">
                    <el-icon v-if="isFullscreen"><ele-BottomLeft /></el-icon>
                    <el-icon v-else><ele-FullScreen /></el-icon>
                </button>
                <button v-if="showClose" aria-label="close" type="button" @click="closeDialog" >
                    <el-icon><ele-CloseBold /></el-icon>
                </button>
            </div>
        </template>
        <div v-loading="loading">
            <slot></slot>
        </div>
        <template #footer>
            <slot name="footer"></slot>
        </template>
    </el-dialog>
</template>

<script lang="ts">
	import { defineComponent, onMounted, reactive, toRefs, watch } from 'vue';

	const props = {
        modelValue: { type: Boolean, default: false },
        modalFlag: { type: Boolean, default: true },
        title: { type: String, default: "" },
        showClose: { type: Boolean, default: true },
        showFullscreen: { type: Boolean, default: true },
        loading: { type: Boolean, default: false },
    };

	export default defineComponent({
        name: 'q-dialog',
        props: props,
		setup(props, { emit }) {

        	const state = reactive({
                dialogVisible: false,
                isFullscreen: false,
            });

        	const closeDialog = () => {
        		emit('close');
        	    state.dialogVisible = false;
            };

        	const setFullscreen = () => {
        	    state.isFullscreen = !state.isFullscreen;
            };

		    onMounted(() => {
                state.dialogVisible = props.modelValue;
            });

		    watch(
              () => props.modelValue,
              () => {
                    state.dialogVisible = props.modelValue;
                    if(state.dialogVisible)
                    	state.isFullscreen = false;
              }
            );

		    return {
                setFullscreen,
                closeDialog,
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

    &-header {
        position: absolute;
        top: var(--el-dialog-padding-primary);
        right: var(--el-dialog-padding-primary);

        button {
            padding: 5px;
            background: transparent;
            border: none;
            outline: none;
            cursor: pointer;
            font-size: var(--el-message-close-size);
            color: var(--el-color-info);
        }
        button:hover {
            color: var(--el-color-primary);
        }
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
        max-height: calc(90vh - 20px) !important
    }
    :deep(.el-dialog).is-fullscreen .el-dialog__footer {
        padding-bottom: 10px;
        border-top: 1px solid var(--el-border-color-light);
    }
}
</style>
