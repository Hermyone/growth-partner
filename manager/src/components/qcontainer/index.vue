<!--
 * @Descripttion: 自定义容器
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-06-17
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <div class="system-custom-container" :class="{'layout-padding': !cardMode}">
        <el-card
                :shadow="cardShadow"
                :class="{'layout-padding-auto': !cardMode}"
                :style="{padding: cardMode ? `${contentPadding}px`:'0'}"
                :body-style="{ padding: `${boxPadding}!important` }"
                :header="header"
        >
            <template #header v-if="customHeader">
                <slot name="header"></slot>
            </template>
            <slot></slot>
        </el-card>
    </div>
</template>

<script lang="ts">
	import { defineComponent, computed } from 'vue';

	const props = {
		cardMode: {
            type: Boolean,
            default: true,
		},
        header: {
            type: String,
            default: "",
        },
        customHeader: {
		    type: Boolean,
            default: false,
        },
		contentPadding: {
			type: [Number, String],
            default: 15,
        },
        boxPadding: {
            type: [Number, String],
            default: 15,
        }
	};

	export default defineComponent({
		name: 'q-container',
		props: props,

		setup(props) {

            const cardShadow = computed({
              get() {
              	return props.cardMode ? "hover" : "never";
              },
              set(val) { }
            });

            return {
              cardShadow,
            }
		},
	});
</script>


<style scoped lang="scss">
    .system-custom-container {
        :deep(.el-card__body) {
            display: flex;
            flex-direction: column;
            flex: 1;
            overflow: auto;
            padding: 15px 0!important;
            .el-table {
                flex: 1;
            }
            .el-form {
                padding: 0px 15px;
            }
        }
    }
</style>