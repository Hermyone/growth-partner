<!--
 * @Descripttion: 自定义树形下拉框组件, 支持拼音检索
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-10-06
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <el-tree-select v-bind="$attrs" :data="options" :filter-node-method="filterMethod" filterable clearable >
        <template v-if="emptyVisible" #empty>
            <div class="select-placeholder">
                <span>没有符合条件的选项</span>
            </div>
        </template>
    </el-tree-select>
</template>

<script lang="ts">
	import { ref, toRefs, defineComponent, reactive, onMounted, watch } from 'vue';
	import pinyin from './pinyin.js';

    const props = {
        options: { type: Array, default: () => [] },
        emptyVisible: { type: Boolean, default: true },
    };

    export default defineComponent({
        name: 'q-tree-select',
        props: props,
        setup(props) {
            const filterMethod = (value: string, data: any) => {
                let py = pinyin.getCamelChars(data.label);
                return data.label.includes(value) || py.includes(value.toUpperCase());
            };

            return {
                filterMethod,
            };
        },
    });
</script>

<style lang="scss" scoped>
    .select-placeholder {
        height: 28px;
        font-size: 12px;
        align-items: center;
        text-align: center;
        padding: 5px;
        color: var(--el-text-color-disabled);
        margin: 2px 0!important;
        border-bottom: 1px dotted var(--el-border-color);
    }
</style>