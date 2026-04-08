<!--
 * @Descripttion: 根据传入的API和字典使用table展示数据
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2024-07-10
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <q-table-v2
            ref="tableRef"
            :isSelection="false"
            :data="tableData"
            :recordCount="total"
            :row-key="dict.rowKey"
            :columns="dict.columns"
            :tableName="dict.tableName"
            :globalTitle="dict.title"
            @pageChange="onTablePageChange"
    >
        <slot name="customColumns"></slot>
    </q-table-v2>
</template>

<script lang="ts">
	import { defineComponent, reactive, toRefs } from 'vue';

	const props = {
        API: {
            type: Object,
            default: null,
        },
        dict: {
        	type: Object,
            default: () => {}
        },
        params: {
            type: Object,
            default: {},
        }
	};

	export default defineComponent({
		name: 'q-detail-table',
		props: props,
		setup(props, { emit }) {
			const state = reactive({
                loading: false,
                total: 0,
                tableData: [],
                params: {
                    pageNum: 1,
                    pageSize: 10,
                },
			});

            const onTablePageChange = (page: TablePageType) => {
                state.params.pageNum = page.pageNum;
                state.params.pageSize = page.pageSize;
                getTableData();
            };

            const getTableData = () => {
                if(!props.API)
                    return;

                state.loading = false;

                props.API.get(props.params).then((resp: any) => {
                    state.tableData = resp.data.list || [];
                    state.total = resp.data.total || state.tableData.length;
                })

                setTimeout(() => {
                    state.loading = false;
                }, 1000);
            };

			return {
                getTableData,
                onTablePageChange,
				...toRefs(state),
			}
		},
	});

</script>

<style scoped lang="scss">
</style>