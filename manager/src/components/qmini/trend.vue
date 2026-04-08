<template>
    <span class="trend-style" :class="'trend-style--'+type">
		<el-icon v-if="iconType=='P'" class="trend-style-icon">
            <ele-Top />
        </el-icon>
		<el-icon v-if="iconType=='N'" class="trend-style-icon">
            <ele-Bottom />
        </el-icon>
		<el-icon v-if="iconType=='Z'" class="trend-style-icon">
            <ele-Right />
        </el-icon>
		<em class="trend-style-prefix">{{ prefix }}</em>
		<em class="trend-style-value">{{ modelValue }}</em>
		<em class="trend-style-suffix">{{ suffix }}</em>
	</span>
</template>

<script lang="ts">
	import { computed, defineComponent, reactive, toRefs } from 'vue';

	const props = {
        // 数值
        modelValue: { type: Number, default: 0},
        // 数值前缀
        prefix: { type: String, default: "" },
        // 数值后缀
        suffix: { type: String, default: "" },
        // 是否翻转
        reverse: { type: String, default: "" },
	};

	export default defineComponent({
		name: 'q-trend',
		props: props,

		setup(props) {

			const state = reactive({
			});

			const absValue = computed(() => {
				return Math.abs(props.modelValue);
            });

			const iconType = computed((v: string) => {
                if(props.modelValue == 0){
                    v = 'Z'
                }else if(props.modelValue < 0){
                    v = 'N'
                }else if(props.modelValue > 0){
                    v = 'P'
                }
			    return v;
            });

            const type = computed((v: string) => {
                if(props.modelValue == 0){
                    v = 'Z'
                }else if(props.modelValue < 0){
                    v = props.reverse?'P':'N'
                }else if(props.modelValue > 0){
                    v = props.reverse?'N':'P'
                }
                return v;
            });

			return {
                absValue,
                iconType,
                type,
				...toRefs(state),
			}
		},
	});
</script>

<style scoped lang="scss">
    .trend-style{
        display: flex;
        align-items: center;

        &-icon {
            margin-right: 2px;
         }
        &-prefix {
             margin-right: 2px;
         }
        &-suffix {
             margin-left: 2px;
         }
        &--P {
             color: #f56c6c;
         }
        &--N {
             color: #67c23a;
         }
        &--Z {
             color: #555;
         }
        :deep(em){
            font-style: normal;
        }
    }
</style>
