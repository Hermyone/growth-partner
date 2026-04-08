<!--
 * @Descripttion: 根据传入的API和字典使用table展示数据的对话框
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2024-07-23
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <q-dialog :v-bind="$attrs" draggable show-fullscreen>
        <slot name="filter"></slot>

        <q-table-v2
                ref="tableRef"
                :isSelection="selection"
                :highlightCurrentRow="true"
                :data="tableData"
                :recordCount="total"
                :row-key="dict.rowKey"
                :columns="dict.columns"
                :tableName="dict.tableName"
                :globalTitle="dict.title"
                :pageSize="params.pageSize"
                @current-change="onHandleCurrentChange"
                @pageChange="onTablePageChange"
        >
            <slot name="columns"></slot>
        </q-table-v2>

        <template #footer>
			<slot></slot>
        </template>
    </q-dialog>
</template>

<script lang="ts">
	import { defineComponent, reactive, toRefs, ref } from 'vue';

	const props = {
        API: {
            type: Object,
            default: null,
        },
        dict: {
        	type: Object,
            default: () => {}
        },
        extendParams: {
            type: Object,
            default: {},
        },
        respDataNode: {
            type: String,
            default: "list",
        },
        selection: {
        	type: Boolean,
            default: false,
        }
	};

	export default defineComponent({
		name: 'q-dialog-table',
		props: props,
		setup(props, { emit }) {
			const tableRef = ref();
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

            const onSearch = (data: any) => {
                state.params = Object.assign({}, state.params, { ...data });
                tableRef.value?.pageReset();
            };

            const getTableData = () => {
                if(!props.API)
                    return;

                state.loading = false;

                state.params = Object.assign({}, state.params, { ...props.extendParams });
                props.API.get(state.params).then((resp: any) => {
                    state.tableData = resp.data[props.respDataNode] || [];
                    state.total = resp.data.total || state.tableData.length;
                })

                setTimeout(() => {
                    state.loading = false;
                }, 1000);
            };

            const getSelectlist = (field: string) => {
                return tableRef.value?.getSelectlist(field);
            }

            const onHandleCurrentChange = (val: any) => {
            	emit("onSelect", val);
            };

			return {
                tableRef,
                getTableData,
                getSelectlist,
                onSearch,
                onTablePageChange,
                onHandleCurrentChange,
				...toRefs(state),
			}
		},
	});

</script>

<style scoped lang="scss">
</style>