<template>
    <span class="status-style" :class="[{'status-style-processing':pulse}, 'status-bg--'+type]"></span>
</template>

<script lang="ts">
	import { defineComponent, reactive, toRefs } from 'vue';

	const props = {
		// 类型
        type: { type: String, default: "primary" },
        // 脉冲
        pulse: { type: Boolean, default: false }
	};

	export default defineComponent({
		name: 'q-status-indicator',
		props: props,

		setup(props) {

			const state = reactive({
			});

			return {
				...toRefs(state),
			}
		},
	});
</script>

<style scoped lang="scss">
    .status-style{
        display: inline-block;
        background: #000;
        width: 8px;
        height: 8px;
        border-radius: 50%;
        vertical-align: middle;

        &-processing {
            position: relative;
        }
        &-processing:after {
            position: absolute;
            top:0px;
            left:0px;
            width: 100%;
            height: 100%;
            border-radius: 50%;
            background: inherit;
            content: '';
            animation: warn 1.2s ease-in-out infinite;
        }
    }
    .status-bg {
        &--primary{
            background: var(--el-color-primary);
        }
        &--success{
            background: var(--el-color-success);
        }
        &--warning{
            background: var(--el-color-warning);
        }
        &--danger{
            background: var(--el-color-danger);
        }
        &--info{
            background: var(--el-color-info);
        }
    }

    @keyframes warn {
        0% {
            transform: scale(0.5);
            opacity: 1;
        }
        30% {
            opacity: 0.7;
        }
        100% {
            transform: scale(2.5);
            opacity: 0;
        }
    }
</style>
