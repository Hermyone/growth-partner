<!--
 * @Descripttion: tag
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2024-01-02
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <el-tag :effect="custom ? 'dark' : 'light'" v-bind="$attrs" :color="getCustomColor()" :style="getCustomStyle()">
        <slot></slot>
    </el-tag>
</template>

<script lang="ts">
	import { toRefs, defineComponent, reactive, onMounted } from 'vue';

    const props = {
    	// 是否使用自定义颜色
        custom: {
            type: Boolean,
            default: false,
        },
        tagValue: {
        	type: String,
            default: ""
        }
    }

	export default defineComponent({
		name: 'q-tag',
        props: props,
		setup(props, { emit }) {

            const state = reactive({
            });

            // 页面加载时
            onMounted(() => {
            });

            const getCustomColor = () => {
            	if(props.custom) {
            		if(props.tagValue === "A")
            			return "#409EFF"
            		else if (props.tagValue === "B")
                        return "#ff8c00"
                    else if (props.tagValue === "C")
                        return "#F56C6C"
                    else if (props.tagValue === "D")
                        return "#67c23a"
                    else
						return "#409EFF"
                }
            }

            const getCustomStyle = () => {
            	if(props.custom) {
                    return {
                        "border-color" : getCustomColor(),
                        "border-radius": "4px"
                    }
                }
            }

			return {
                getCustomColor,
				getCustomStyle,
				...toRefs(state),
			};
		},
	});
</script>

<style scoped lang="scss">

</style>