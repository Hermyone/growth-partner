<!--
 * @Descripttion: 具有el-input样式的标签框，内容支持单击事件
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2024-07-25
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <div class="custom-label">
        <div class="label-content" tabindex="0" @click.stop="handleClickContent">
            <span>{{ text }}</span>
            <el-icon v-if="icon"><Position /></el-icon>
        </div>
    </div>
</template>

<script lang="ts">
	import { defineComponent, reactive, toRefs } from 'vue';
	import { Position } from '@element-plus/icons-vue';

    const props = {
        text: {
            type: String,
            default: ''
        },
        icon: {
        	type: Boolean,
            default: false
        }
	};

	export default defineComponent({
		name: 'q-link-label',
        components: { Position },
		props: props,
		setup(props, {emit}) {
            const state = reactive({
            });

            const handleClickContent = () => {
                emit("click");
            };

            return {
                handleClickContent,
                ...toRefs(state),
            }
        },
	});
</script>

<style scoped lang="scss">
    .custom-label {
        border: 1px solid var(--el-border-color);
        border-radius: 4px;
        padding: 5px 10px;
        position: relative;
        height: 32px;
    }

    .label-content {
        user-select: none;
        cursor: pointer;
        transition: color 0.3s ease;
        display: flex;
        justify-content: space-between;
        height: 100%;
        align-items: center;

        span {
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }
        .el-icon {
            padding-top: 8px;
        }
    }

    .label-content:hover {
        color: var(--el-color-primary);
    }
</style>
