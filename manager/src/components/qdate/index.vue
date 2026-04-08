<!--
 * @Descripttion: 日期选择组件
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2024-07-10
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <el-date-picker
            v-bind="$attrs"
            unlink-panels
            class="w100"
            :type="type"
            :value-format="fmt"
            :range-separator="separator"
            :shortcuts="type === 'daterange' ? shortcuts:[]"
            placeholder="请选择日期"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
    >
    </el-date-picker>
</template>

<script lang="ts">
	import { defineComponent, onMounted, reactive, toRefs } from 'vue';

	const props = {
        type: { type: String, default: "date" },
        fmt: { type: String, default: "YYYY-MM-DD" },
		separator: { type: String, default: "-" },
	};

	export default defineComponent({
		name: 'q-date-picker',
		props: props,

		setup(props) {

			const state = reactive({
                shortcuts: [
                    {
                        text: '最近一周',
                        value: () => {
                            const end = new Date()
                            const start = new Date()
                            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
                            return [start, end]
                        },
                    },
                    {
                        text: '最近一月',
                        value: () => {
                            const end = new Date()
                            const start = new Date()
                            start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
                            return [start, end]
                        },
                    },
                    {
                        text: '最近3个月',
                        value: () => {
                            const end = new Date()
                            const start = new Date()
                            start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
                            return [start, end]
                        },
                    },
                ],
			});

            // 页面加载时
            onMounted(() => {

            });

			return {
				...toRefs(state),
			}
		},
	});
</script>

<style scoped lang="scss">
</style>
