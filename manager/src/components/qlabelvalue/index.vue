<!--
 * @Descripttion: 标签+VALUE框组件
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2024-06-07
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <div ref="lvRef" class="label-value" style="height: 32px">
        <span :style="{ width: `${labelWidth}px`, textAlign: `${labelAlign}` }">{{ label }}{{ separator }}</span>
        <el-text :type="valueType" :tag="valueTag" truncated :style="{ textAlign: `${valueAlign}` }">
            <slot name="prefix"></slot>
            {{ value }}
            <slot name="suffix"></slot>
            <slot></slot>
        </el-text>
    </div>
</template>

<script lang="ts">
	import { defineComponent, onMounted, ref, watch } from 'vue';

	const props = {
        label: {
            type: String,
            default: '',
        },
        labelWidth: {
        	type: Number,
            default: 90,
        },
        labelAlign: {
            type: String,
            default: 'right',
        },
        value: {
        	type: String,
            default: '',
        },
        valueAlign: {
            type: String,
            default: 'left',
        },
        valueType: {
            type: String,
            default: '',
        },
        valueTag: {
            type: String,
            default: 'span',
        },
        separator: {
            type: String,
            default: '：',
        }
	};

	export default defineComponent({
		name: 'q-label-value',
		props: props,
		setup(props) {
            const lvRef = ref<Element>()
            onMounted(() => {
            });

            return {
                lvRef,
            }
        },
	});
</script>

<style scoped lang="scss">
    .label-value {
        display: flex;
        align-items: center;
        justify-content: flex-start;

        span {
            padding-right: 15px;
        }

        .el-text {
            width: calc(100% - 100px);
        }
    }
</style>

