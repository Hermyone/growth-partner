<template>
    <el-drawer v-model="dialogVisible" :size="800">
        <template #header>
            <h2>数据导出</h2>
        </template>
        <el-container v-loading="loading">
            <el-main>
                <div class="dialog-header">
                    <el-form :model="params" :rules="rules" ref="formRef" label-width="80px">
                        <el-form-item label="表格名称" prop="tableName">
                            <el-input
                                    v-model="params.tableName"
                                    placeholder="请输入表格名称"
                                    clearable
                                    style="width:400px;"
                            />
                        </el-form-item>
                        <el-form-item label="导出标题" prop="tableCaption">
                            <el-input
                                    v-model="params.tableCaption"
                                    placeholder="请输入导出文件名和表格标题"
                                    clearable
                                    style="width:400px;"
                            />
                        </el-form-item>
                        <el-form-item label="导出类型" prop="fileType">
                            <el-radio-group v-model="params.fileType" size="small">
                                <el-radio-button value="xlsx">Excel</el-radio-button>
                                <el-radio-button value="csv" >CSV</el-radio-button>
                                <el-radio-button value="text" >Text</el-radio-button>
                            </el-radio-group>
                        </el-form-item>
                    </el-form>
                </div>

                <el-table
                        ref="tableRef"
                        size="small"
                        :stripe="true"
                        :data="data"
                >
                    <el-table-column type="index" label="ID" width="40" fixed="left" />
                    <el-table-column prop="flag" label="导出" width="50">
                        <template #default="{ row }">
                            <el-checkbox v-model="row.flag" true-label="Y" false-label="N" />
                        </template>
                    </el-table-column>
                    <el-table-column prop="primary" label="主键" width="40">
                        <template #default="{ row }">
                            <span :style="{color: row.primary ? 'red':'var(--el-text-color-disabled)'}">{{ row.primary ? '是':'否' }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="title" label="列名称" width="200" >
                        <template #default="{ row }">
                            <el-input v-model="row.title" size="small"></el-input>
                        </template>
                    </el-table-column>
                    <el-table-column prop="width" label="列宽" width="120">
                        <template #default="{ row }">
                            <el-input-number v-model="row.width" size="small" style="width: 100px"></el-input-number>
                        </template>
                    </el-table-column>
                    <el-table-column prop="type" label="列类型" width="80">
                        <template #default="{ row }">
                            <el-select v-model="row.type" size="small">
                                <el-option label="文本" value='text' />
                                <el-option label="图像" value='image' />
                            </el-select>
                        </template>
                    </el-table-column>
                    <el-table-column prop="dict" label="字典">
                        <template #default="{ row }">
                            <div style="display: flex; justify-content: space-between; align-items: center">
                                <el-text class="w-100px" truncated>{{ row.dict }}</el-text>
                                <el-button size="small" text type="primary" @click="editDict(row)">修改</el-button>
                            </div>
                        </template>
                    </el-table-column>
                </el-table>
            </el-main>

            <el-footer>
                <el-button type="primary" @click="doExport" >导出</el-button>
                <el-button type="primary" plain @click="doSave" >保存</el-button>
                <el-button @click="dialogVisible=false" >关闭</el-button>
            </el-footer>
        </el-container>

        <el-drawer title="字典" v-model="drawerVisible" :size="500" destroy-on-close >
            <el-main style="padding:5px">
                <q-form-table size="small" v-model="dicts" :addTemplate="addTemplate" drag-sort placeholder="暂无数据">
                    <el-table-column prop="text" label="名称">
                        <template #default="{ row }">
                            <el-input v-model="row.text" placeholder="请输入名称"></el-input>
                        </template>
                    </el-table-column>
                    <el-table-column prop="value" label="值" width="150">
                        <template #default="{ row }">
                            <el-input v-model="row.value" placeholder="请输入值"></el-input>
                        </template>
                    </el-table-column>
                </q-form-table>
            </el-main>
        </el-drawer>
    </el-drawer>
</template>

<script lang="ts">
	import { defineComponent, toRefs, onMounted, ref, reactive } from 'vue';
    import table2excel from 'js-table2excel';
    import Cookies from 'js-cookie';
    import { ElMessage } from 'element-plus';

	const props = {
        // 表格名称
        tableName: { type: String, default: '' },
        // 表格标题
        tableTitle: { type: String, default: '' },
        // 表格数据
        tableData: { type: Array, default: () => [] },
        // 数据行主键
		rowKey: { type: String, default: '' },
        // 定义列
        columns: { type: Array, default: () => [] },
    };

	export default defineComponent({
		name: 'q-export',
		props: props,
        setup(props, { emit }) {
            const dialogVisible = ref(false);
            const drawerVisible = ref(false);
            const formRef = ref();
            const state = reactive({
                loading: false,
                data: [] as any,
                params: {
                    fileType: 'xlsx',
                    tableName: props.tableName,
                    tableCaption: props.tableTitle,
                },
                addTemplate: {
                    text: '',
                    value: ''
                },
                dicts: [],
                rules: {
                    tableName: [ { required: true, message: '不能为空'} ],
                },
            });

            const doExport = () => {
            	// if(props.rowKey == ''){
            	// 	ElMessage.warning('导出失败, 数据未设置主键')
            	// 	return;
                // }

                let columns = state.data.filter((item: any) => item.flag === 'Y');
                if (columns.length <= 0) return ElMessage.warning('请先选择要导出的列');

                formRef.value?.validate((valid: boolean) => {
                    if (valid) {
                        switch (state.params.fileType) {
                            case 'xlsx':
                                toExcel(columns);
                                break;
                            case 'csv':
                                toText(columns, ',', '.csv');
                                break;
                            case 'text':
                                toText(columns, '\t', '.txt');
                                break;
                        }
                    }
                });
            };

            // 根据参数的设置格式化导出的数据
            const formatData = (cols: any[]) => {
                let formattedData = props.tableData.map((item:any)=> {
                    let newItem = {} as any;
                    cols.forEach(field => {
                        newItem[field.key] = item[field.key]??'';
                        if (field.dict && field.dict.length > 0) {
                            let dictItem = field.dict.find((dict:any) => dict.value == item[field.key]);
                            if (dictItem) {
                                newItem[field.key] = dictItem.text;
                            }
                        }
                    });
                    return newItem;
                });

                return formattedData;
            };

            const toExcel = (columns: any) => {
                let expData = formatData(columns);

                table2excel(columns, expData, state.params.tableCaption, state.params.tableCaption);
                ElMessage.warning('文件导出成功');
            };

            const toText = (columns: any, separator: string, extname: string) => {
                let expData = formatData(columns);
                let headers = {} as any;
                const csvRows = [] as any;
                // 添加表头
                columns.forEach((item:any) => {
                    headers[item.key] = item.title;
                });
                csvRows.push(Object.values(headers).join(separator));

                // 添加数据
                for (const row of expData) {
                    csvRows.push(Object.values(row).join(separator));
                }

                let csvContent = csvRows.join('\n');

                // 下载
                createDownload(csvContent, state.params.tableCaption + extname);

                ElMessage.warning('文件导出成功');
            };

            // 生成下载链接
            const createDownload = (content: any, filename: string) => {
                const blob = new Blob([content], { type: 'text/csv;charset=utf-8;' });
                const link = document.createElement('a');
                if (link.download !== undefined) {
                    const url = URL.createObjectURL(blob);
                    link.setAttribute('href', url);
                    link.setAttribute('download', filename);
                    link.style.visibility = 'hidden';
                    document.body.appendChild(link);
                    link.click();
                    document.body.removeChild(link);
                }
            };

            const doSave = () => {
                formRef.value?.validate((valid: boolean) => {
                    if(valid){
                        try {
                            new Promise(() => {
                                setTimeout(()=>{
                                    Cookies.set(state.params.tableName + '_export', JSON.stringify(state.data));
                                },500);
                            });
                        }catch(error){
                            ElMessage.error('保存失败');
                        }

                        ElMessage.success('保存成功');
                    }
                });
            };

            const open = () => {
                dialogVisible.value = true;
            };

            const editDict = (row: any) => {
                state.dicts = row.dict;
                drawerVisible.value = true;
            };

            onMounted(() => {
                if (props.columns) {
                    props.columns.forEach((item: any) => {
                        let r = {
                            primary: item.key == props.rowKey,
                            key: item.key,
                            title: item.title,
                            width: Number(item.width),
                            height: 32,
                            type: 'text',
                            flag: 'Y',
                            dict: item.dicts ? JSON.parse(JSON.stringify(item.dicts)) : [],
                        };
                        state.data.push(r);
                    });
                }
                const customData : any = Cookies.get(props.tableName + '_export');
                if(customData){
                    state.data = JSON.parse(customData);
                }
            });

            return {
                dialogVisible,
                drawerVisible,
                formRef,
                editDict,
                doExport,
                doSave,
                open,
                ...toRefs(state),
            };
        },
    });
</script>

<style lang="scss" scoped>
    .dialog-header {
        .el-form-item {
            margin-bottom: 15px!important;
        }
    }
</style>