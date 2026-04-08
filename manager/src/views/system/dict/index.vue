<template>
	<q-container :cardMode="getCardMode()" class="back-container">
		<el-container>
			<el-aside width="250px">
				<el-container>
					<el-header>
						<el-input placeholder="输入关键字进行过滤" v-model="filterNode" clearable></el-input>
						<el-button icon="ele-Plus" style="height: 30px; width: 30px; padding-left: 18px;" class="ml10" @click="addNode"></el-button>
					</el-header>
					<el-main class="nopadding">
						<el-scrollbar>
							<el-tree
									ref="treeRef"
									class="folder-tree"
									:data="treeData"
									:props="dictProps"
									:expand-on-click-node="false"
									:highlight-current="true"
									:filter-node-method="doFilterNode"
									:current-node-key="dictType"
									node-key="dictType"
									@node-click="nodeClick"
							>
								<template #default="{node, data}">
									<span class="custom-tree-node">
										<el-text class="label" truncated>{{ node.label }}</el-text>
										<span class="code">{{ data.dictType }}</span>
										<span class="do">
											<el-button-group>
												<el-button icon="ele-EditPen" size="small" @click.stop="editNode(data)"></el-button>
												<el-popconfirm
														width="auto"
														confirm-button-text="确认"
														cancel-button-text="取消"
														:title="'确认要删除字典【'+$t(data.dictName)+'】吗?'"
														@confirm="delNode(node, data)"
												>
													<template #reference>
													  <el-button icon="ele-Delete" size="small"></el-button>
													</template>
												  </el-popconfirm>
											</el-button-group>
										</span>
									</span>
								</template>
							</el-tree>
						</el-scrollbar>
					</el-main>
				</el-container>
			</el-aside>
			<el-container>
				<el-header class="back-container-header">
					<div>
						<el-button icon="ele-Plus" size="default" round type="success" class="ml10" @click="addDict" :disabled="dictType == ''">
							新增
						</el-button>
						<el-button icon="ele-Delete" size="default" round type="danger" class="ml10" @click="delDict(null)">
							删除
						</el-button>
					</div>
					<el-input v-model="filterDict" placeholder="请输入查询条件" clearable style="width:200px"></el-input>
				</el-header>
				<el-main>
					<q-table
							ref="tableRef"
							class="back-container-table"
							:data="tableData"
							:recordCount="total"
							:isSelection="true"
							:paginationVisible="false"
							row-key="dictCode"
					>
						<el-table-column label="" width="60">
							<template #default>
								<svg class="move" style="cursor: move;" viewBox="64 64 896 896" focusable="false" data-icon="holder" width="1em" height="1em" fill="currentColor" aria-hidden="true">
									<path d="M300 276.5a56 56 0 1056-97 56 56 0 00-56 97zm0 284a56 56 0 1056-97 56 56 0 00-56 97zM640 228a56 56 0 10112 0 56 56 0 00-112 0zm0 284a56 56 0 10112 0 56 56 0 00-112 0zM300 844.5a56 56 0 1056-97 56 56 0 00-56 97zM640 796a56 56 0 10112 0 56 56 0 00-112 0z"></path>
								</svg>
							</template>
						</el-table-column>
						<el-table-column label="字典标签" prop="dictLabel" width="200" show-overflow-tooltip></el-table-column>
						<el-table-column label="字典键值" prop="dictValue" width="100"></el-table-column>
						<el-table-column label="颜色" prop="dictColor" width="100">
							<template #default="{row}">
								<el-tag style="width: 100%;" :color="row.dictColor" v-if="row.dictColor && row.dictColor !== ''"/>
							</template>
						</el-table-column>
<!--						<el-table-column label="字典排序" prop="dictSort" width="100"></el-table-column>-->
						<el-table-column label="字典状态" prop="status" width="100">
							<template #default="{row}">
								<el-tag type="success" v-if="row.status">启用</el-tag>
								<el-tag type="info" v-else>禁用</el-tag>
							</template>
						</el-table-column>
<!--						<el-table-column label="创建时间" prop="createdAt" width="180" show-overflow-tooltip></el-table-column>-->
						<el-table-column label="备注" prop="remark" width="200" show-overflow-tooltip></el-table-column>
						<el-table-column label="操作" prop="opt" fixed="right" width="120">
							<template #default="{row}">
								<el-button size="small" text type="primary" @click="editDict(row)">修改</el-button>
								<el-divider direction="vertical" />
								<el-button size="small" text type="danger" @click="delDict(row)">删除</el-button>
							</template>
						</el-table-column>
					</q-table>
				</el-main>
			</el-container>
		</el-container>

		<EditNode ref="nodeDlgRef" @typeList="loadNode"></EditNode>
		<EidtDict ref="dictDlgRef" @dataList="getDictData" :dict-type="dictType"></EidtDict>
	</q-container>
</template>

<script>
	import Sortable from 'sortablejs';
	import EditNode from './component/editDic.vue';
	import EidtDict from './component/editDicData.vue';
	import { ElMessage, ElMessageBox } from 'element-plus';
	import { getCardMode } from '/@/utils/common';
	import { deleteType, getTypeList } from '../../../api/system/dict/type';
	import { deleteData, getDataList, editData } from '../../../api/system/dict/data';

	export default {
		name: 'systemDic',
		components: {
			EditNode,
			EidtDict
		},
		data() {
			return {
				total: 0,
				loading: false,
				tableData: [],
				dictProps: {
					label: 'dictName'
				},
				treeData: [],
				filterNode: '',
				filterDict: '',
				dictType: '',
				getCardMode
			}
		},
		watch: {
			filterNode(val){
				this.$refs.treeRef.filter(val);
			}
		},
		mounted() {
			this.loadNode();
			this.rowDrop();
		},
		computed: {
		},
		methods: {
			//节点过滤
			doFilterNode(value, data){
				if (!value) return true;
				var targetText = String(data.dictName + data.dictType).toLocaleLowerCase();
				return targetText.indexOf(value.toLocaleLowerCase()) !== -1;
			},
			//行拖拽
			rowDrop(){
				const _this = this;
				const tbody = this.$refs.tableRef.$el.querySelector('.el-table__body-wrapper tbody');
				Sortable.create(tbody, {
					handle: ".move",
					animation: 300,
					ghostClass: "ghost",
					onEnd({ newIndex, oldIndex }) {
						const tableData = _this.tableData;
						const currRow = tableData.splice(oldIndex, 1)[0];
						tableData.splice(newIndex, 0, currRow);
						_this.doChangeOrder(tableData);
					}
				})
			},
			//加载树
			async loadNode(){
				let params = {
					pageNum: 1,
					pageSize: 9999
				}
				getTypeList(params).then((res) => {
					this.treeData = res.data.dictTypeList ?? [];

					if(this.treeData.length > 0) {
						let firstNode = this.treeData[0];
						this.nodeClick({ dictType: firstNode.dictType });
					}
				});
			},
			//树增加
			addNode(){
				this.$nextTick(() => {
					this.$refs.nodeDlgRef.openDialog();
				})
			},
			//编辑树
			editNode(data){
				this.$nextTick(() => {
					this.$refs.nodeDlgRef.openDialog(data);
				})
			},
			//树点击事件
			nodeClick(data){
				this.dictType = data.dictType;
				this.getDictData();
			},
			getDictData(){
				let params = {
					pageNum: 1,
					pageSize: 999,
					dictType: this.dictType
				}
				getDataList(params).then((res) => {
					this.tableData = res.data.list ?? [];
					this.total = res.data.total || this.tableData.length;
				});
			},
			//删除树
			delNode(node, data){
				var dicCurrentKey = this.$refs.treeRef.getCurrentKey();
				deleteType([data.dictId]).then(() => {
					this.$refs.treeRef.remove(data.dictType);
					if(dicCurrentKey === data.dictType){
						if(this.treeData.length > 0){
							var firstNode = this.treeData[0];
							if(firstNode){
								this.nodeClick({dictType: firstNode.dictType})
							}else{
								this.nodeClick({dictType: ''})
							}
						}
					}
					ElMessage.success('数据已更新');
				});
			},
			//添加明细
			addDict(){
				this.$nextTick(() => {
					this.$refs.dictDlgRef.openDialog();
				})
			},
			//编辑明细
			editDict(row){
				this.$nextTick(() => {
					this.$refs.dictDlgRef.openDialog(row);
				})
			},
			//删除明细
			async delDict(row){
				let msg = '确认要删除所选数据？';
				let ids = [];
				if (row) {
					msg = `确认要删除字典【${row.dictLabel}】吗?`;
					ids = [row.dictCode];
				} else {
					ids = this.$refs.tableRef.getSelectlist('dictCode');
				}
				if (ids.length === 0) {
					ElMessage.error('请选择要删除的数据。');
					return;
				}
				ElMessageBox.confirm(msg, '提示', {
					confirmButtonText: '确认',
					cancelButtonText: '取消',
					type: 'warning',
				})
						.then(() => {
							deleteData(ids).then(() => {
								ElMessage.success('删除成功');
								this.nodeClick({dictType: this.dictType})
							});
						})
						.catch(() => {});
			},
			doChangeOrder(data){
				var idx = 0;
				data.forEach(item => {
					item.dictSort = idx
					editData(item).then(() => {});
					idx++;
				})
			},
		},
	}
</script>

<style scoped lang="scss">
	.back-container{
		padding: 0 !important;
		margin: 0 !important;
		.el-main {
			padding: 0 0 15px 0;
		}
		&-table {
			height: 100%;
		}
		:deep(.el-card__body) {
			padding: 0 !important;
		}
		&-header{
			align-items: flex-start;
		}
	}
	.custom-tree-node {
		display: flex;
		flex: 1;
		align-items: center;
		justify-content: space-between;
		font-size: 12px;
		padding: 5px;
		height: 100%;
		.el-text {
			font-size: 13px;
		}
	}
	.custom-tree-node .label {
		max-width: 150px;
	}
	.custom-tree-node .code {
		font-size: 12px;
		color: var(--el-text-color-disabled);
	}
	.custom-tree-node .do {
		display: none;
	}
	.custom-tree-node:hover .code {
		display: none;
	}
	.custom-tree-node:hover .do {
		display: inline-block;
	}
</style>