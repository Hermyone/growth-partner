<!--
 * @Descripttion: 自定义图形下拉框组件
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-12-26
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <el-select v-bind="$attrs" :data="images" class="custom">
        <el-option v-for="item in images"
                   :key="item.path"
                   :label="item.name"
                   :value="item.path"
        >
            <div class="custom-item" :class="{'large': imagesize==='large'}">
                <img :src="item.path"/>
                <span v-if="imagesize === 'small'">{{ item.name }}</span>
            </div>
        </el-option>
    </el-select>
</template>

<script lang="ts">
	import { ref, defineComponent, watch } from 'vue';

    const props = {
        images: { type: Array, default: () => [] },
        // 大小 small|large
        imagesize: { type: String, default: 'small' }
    };

    export default defineComponent({
        name: 'q-image-select',
        props: props,
        setup(props) {
            return {
            };
        },
    });
</script>

<style lang="scss" scoped>
    .el-select-dropdown__item {
        height: 60px;
        padding: 0 5px;
    }
    .custom {
        &-item {
            height: 100%;
            display: flex;
            align-items: center;
            img {
                width: 45px;
                height: 45px;
                border: none;
                border-radius: 5px;
                opacity: 0.8;
                &:hover {
                    border: 1px solid var(--el-color-primary);
                    opacity: 1;
                }
            }
            span {
               height: 100%;
               padding: 10px;
               color: var(--el-text-color-primary);
               &:hover{
                   color: var(--el-color-primary);
               }
            }

            &.large {
                img {
                    width: 100%;
                    height: 100%;
                    border: none;
                    border-radius: 5px;
                    padding: 2px;
                    opacity: 0.8;
                    &:hover {
                        border: 1px solid var(--el-color-primary);
                        opacity: 1;
                    }
                }
            }
        }
    }
</style>