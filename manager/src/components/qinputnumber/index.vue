<!--
 * @Descripttion: 自定义数值型输入框组件，支持格式化显示
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-10-20
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <el-input class="numrule" v-bind="$attrs" type="number" v-thousands="{fmt: `${format}`, focus: `${focusAlign}`}">
        <template #suffix>
            {{ suffix }}
        </template>
    </el-input>
</template>

<script lang="ts">
	import { ref, toRefs, defineComponent, reactive, onMounted } from 'vue';

    const props = {
        suffix: { type: String, default: '' },
        format: { type: String, default: '0.##' },
        focusAlign: { type: String, default: 'center' },
    };

    export default defineComponent({
        name: 'q-input-number',
        props: props,
        setup(props, { emit }) {
            const state = reactive({

            });

            // 页面加载时
            onMounted(() => {

            });

            return {
                ...toRefs(state),
            };
        },
    });
</script>

<style lang="scss" scoped>
</style>