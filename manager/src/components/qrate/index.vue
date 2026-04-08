<!--
 * @Descripttion: rate 支持可绑定字符串类型数值
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2024-07-19
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <el-rate :v-bind="$attrs" v-model="defaultModel" disabled>
    </el-rate>
</template>

<script lang="ts">
	import { toRefs, defineComponent, reactive, onMounted, watch } from 'vue';

    const props = {
        rateValue: {
        	type: [Number, String],
            default: 1.0
        }
    }

	export default defineComponent({
		name: 'q-rate',
        props: props,
		setup(props, { emit }) {

            const state = reactive({
              defaultModel: 1.0
            });

            // 页面加载时
            onMounted(() => {
                state.defaultModel = formatValue(props.rateValue);
            });

            const formatValue = (val: any) => {
                if(typeof val === 'number')
                    return parseFloat(val.toFixed(1));
                else {
                    const num = parseFloat(val);
                    const factor = Math.pow(10, 1);
                    return Math.round(num * factor) / factor;
                }
            }

            watch(
                () => props.rateValue,
                (val) => {
                    state.defaultModel = formatValue(val);
                },{
                    deep: true
                }
            );

			return {
				...toRefs(state),
			};
		},
	});
</script>

<style scoped lang="scss">

</style>