<!--
 * @Descripttion: 自定义下拉选择框组件, 支持拼音检索
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-10-05
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <el-select v-bind="$attrs" @visible-change="visibleChange" clearable>
        <div v-if="searchVisible" class="select-search">
            <el-input v-model="searchValue" size="small" placeholder="请输入过滤条件" clearable @input="filterMethod"></el-input>
        </div>
        <el-option v-for="field in optionsList" :key="field[`${valueField}`]" :label="field[`${labelField}`]" :value="field[`${valueField}`]" :disabled="isDisabled(field[`${valueField}`])">
            <div v-if="expendList">
                <span style="float: left">{{ field[labelField] }}</span>
                <span style="float: right; color: var(--el-text-color-disabled);font-size: 10px;">{{ field[valueField] }}</span>
            </div>
        </el-option>
        <template v-if="emptyVisible && (optionsList.length === 0 || optionsList_.length === 0)" #empty>
            <div class="select-placeholder">
                <span>没有符合条件的选项</span>
            </div>
        </template>
    </el-select>
</template>

<script lang="ts">
	import { ref, toRefs, defineComponent, reactive, onMounted, watch } from 'vue';
	import pinyin from './pinyin.js';

    const props = {
        options: { type: Array, default: () => [] },
        filter: { type: Array, default: [] as any },
        searchVisible: { type: Boolean, default: false },
        emptyVisible: { type: Boolean, default: true },
        labelField: { type: String, default: 'label' },
        valueField: { type: String, default: 'value' },
        expendList: { type: Boolean, default: false },
    };

    export default defineComponent({
        name: 'q-select',
        props: props,
        setup(props, { emit }) {
            const searchValue = ref('');
            const state = reactive({
                optionsList: [] as any,
                optionsList_: [] as any,
            });

            // 页面加载时
            onMounted(() => {
                state.optionsList = props.options;
                state.optionsList_ = [...props.options];
            });

            const filterMethod = (keyword: string) => {
                if(keyword){
                    state.optionsList = state.optionsList_;
                    state.optionsList = state.optionsList.filter((item: any) =>{
                            // pinyin.match(item[props.labelField], keyword)
                            let py = pinyin.getCamelChars(item[props.labelField]);
                            keyword = keyword.toUpperCase();
                            return item[props.valueField].includes(keyword) || item[props.labelField].includes(keyword) || py.includes(keyword)
                        }
                    );
                    if(state.optionsList.length == 0) state.optionsList = state.optionsList_;
                }else{
                    state.optionsList = state.optionsList_;
                }
            };

            const visibleChange = (isopen: boolean) => {
                if(isopen){
                	searchValue.value = '';
                    state.optionsList = state.optionsList_;
                }
            };

            const isDisabled = (key: any) => {
                if(props.filter?.find((item:any) => item.field[props.valueField] == key && !item.field.repeat)){
                    return true
                }else{
                    return false
                }
            };

            watch(
              () => props.options,
              (val) => {
              	    state.optionsList = val;
              	    state.optionsList_ = val;
              }, {
              	    deep: true
              }
            );

            return {
                searchValue,
                filterMethod,
                visibleChange,
                isDisabled,
                ...toRefs(state),
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
    .select-search {
        height: 32px;
        margin: 2px;
        border-bottom: 1px dotted var(--el-border-color);
        :deep(.el-input__wrapper) {
            box-shadow: 0 0 0 0 var(--el-input-border-color, var(--el-border-color)) inset;
            cursor: default;
            .el-input__inner {
                cursor: default !important;
            }
        }
    }
</style>