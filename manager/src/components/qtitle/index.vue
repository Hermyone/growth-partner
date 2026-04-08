<!--
 * @Descripttion: 自定义标签组件
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-07-15
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <div class="title-style">
        <div class="title-header"></div>
        {{ title }}
    </div>
</template>

<script lang="ts">
	import { defineComponent } from 'vue';

	const props = {
        title: { type: String, default: "" },
	};

	export default defineComponent({
		name: 'q-title',
		props: props,

		setup(props) {
        },
	});
</script>

<style scoped lang="scss">
    .title-header {
        position: relative;
        padding-left: 10px;
        color: var(--el-text-color-regular);
        &::after {
            content: '';
            width: 2px;
            height: 15px;
            position: absolute;
            left: 0;
            top: 50%;
            transform: translateY(-50%);
            background: var(--el-color-primary);
        }
    }
    .title-style {
        border-bottom: 1px solid var(--el-border-color);
        font-size: var(--el-font-size-base);
        padding-left: 15px;
        padding-bottom: 15px;
        color: var(--el-color-primary);
        height: 35px;
        display: flex;
        align-items: center;
    }
</style>
