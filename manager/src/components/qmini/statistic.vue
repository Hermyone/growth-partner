<template>
    <div class="statistic-style">
        <div class="statistic-style-title">
            {{ title }}
            <el-tooltip v-if="tips" effect="light">
                <template #content>
                    <div style="width: 200px;line-height: 2;">
                        {{ tips }}
                    </div>
                </template>
                <SvgIcon class="statistic-style-tips" name="fa fa-question-circle-o" :size="14" color="#909399" />
            </el-tooltip>
        </div>
        <div class="statistic-style-content">
            <span v-if="prefix" class="statistic-style-content-prefix">{{ prefix }}</span>
            <span class="statistic-style-content-value">{{ cmtValue }}</span>
            <span v-if="suffix" class="statistic-style-content-suffix">{{ suffix }}</span>
        </div>
        <div v-if="description || $slots.default" class="statistic-style-description">
            <slot>
                {{ description }}
            </slot>
        </div>
    </div>
</template>

<script lang="ts">
	import { computed, defineComponent, reactive, toRefs } from 'vue';

	const props = {
		// 标题
        title: { type: String, required: true, default: "" },
        // 数值
        value: { type: String, required: true, default: "" },
        // 数值前缀
        prefix: { type: String, default: "" },
        // 数值后缀
        suffix: { type: String, default: "" },
        // 描述信息
        description: { type: String, default: "" },
        // 提示信息
        tips: { type: String, default: "" },
        // 是否数值格式化
        groupSeparator: { type: Boolean, default: false }
	};

	export default defineComponent({
		name: 'q-statistic',
		props: props,

		setup(props) {

			const state = reactive({
			});

			const cmtValue = computed(() => {
				var num = props.value;
                if(!num.includes('.')){
                    num += '.'
                }
                return num.replace(/(\d)(?=(\d{3})+\.)/g, function ($0, $1) {
                    return $1 + ',';
                }).replace(/\.$/, '');
            });

			return {
                cmtValue,
				...toRefs(state),
			}
		},
	});
</script>

<style scoped lang="scss">
    .statistic-style {
        &-title {
            font-size: 12px;
            color: var(--el-text-color-secondary);
            margin-bottom: 10px;
            display: flex;
            align-items: center;
        }

        &-content {
            font-size: 20px;
            color: var(--el-color-primary);
            &-prefix {
                margin-right: 5px;
            }
            &-value {
                font-weight: bold;
            }
            &-suffix {
                margin-left: 5px;
                font-size: 12px;
            }
        }

        &-tips {
            margin-left: 5px;
        }

        &-description {
            margin-top: 10px;
            color: var(--el-text-color-secondary);
        }
    }
</style>
