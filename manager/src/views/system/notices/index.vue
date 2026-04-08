<template>
	<q-container :cardMode="getCardMode()" class="back-container">
		<el-container>
			<el-tabs v-model="activePage" tab-position="left" class="back-container-tabs" @tab-change="doFilter">
				<el-tab-pane name="broadcast">
					<template #label>
						<el-badge :value="state.badge.broadcast">
							<span>系统广播</span>
						</el-badge>
					</template>
				</el-tab-pane>
				<el-tab-pane name="mail">
					<template #label>
						<el-badge :value="state.badge.mail">
							<span>站内信</span>
						</el-badge>
					</template>
				</el-tab-pane>
			</el-tabs>
			<el-container>
				<el-header>
					<el-radio-group v-model="state.msgState" size="small" @change="doFilter">
					<el-radio-button value="">全部</el-radio-button>
					<el-radio-button value="0" >未读</el-radio-button>
					<el-radio-button value="1" >已读</el-radio-button>
					</el-radio-group>

					<div>
						<el-input v-model="state.filterMsg" placeholder="请输入过滤条件" clearable style="width:250px"></el-input>
						<el-dropdown style="padding-left: 5px" @command="doCommand">
							<el-button type="primary">
								更多
								<el-icon class="el-icon--right">
									<ele-ArrowDown />
								</el-icon>
							</el-button>
							<template #dropdown>
								<el-dropdown-menu>
									<el-dropdown-item command="readall">标记全部已读</el-dropdown-item>
								</el-dropdown-menu>
							</template>
						</el-dropdown>
					</div>
				</el-header>
				<el-main style="padding: 0!important;">
					<q-table
							ref="tableRef"
							class="table-box"
							row-key="id"
							:recordCount="state.total"
							:data="state.tableData"
							:isSerialNo="false"
							:isSelection="false"
							:showHeader="false"
							:footerButton="false"
							:footerVisible="false"
					>
						<el-table-column label="消息" >
							<template #default="{row}">
								<div class="table-box-item">
									<el-icon :size="32" class="table-box-item-icon">
										<el-badge is-dot class="item" v-if="row.isread !== 1"><ele-Message /></el-badge>
										<div v-else><ele-Message /></div>
									</el-icon>

									<div class="table-box-item-content">
										<el-row type="flex" class="table-box-item-content-title">
											<el-space wrap>
												<span>{{ row.title }}</span>
												<el-tag size="small" v-if="row.level === 1">重要</el-tag>
												<el-tag size="small" type="danger" v-if="row.isread === 0" >未读</el-tag>
											</el-space>
										</el-row>
										<el-row class="table-box-item-content-text">
											<span class="table-box-item-content-desc">{{ row.time }}</span>
										</el-row>
									</div>
								</div>
							</template>
						</el-table-column>
						<el-table-column label="操作" align="center" width="100px" fixed="right">
							<template #default="{row}">
								<el-button size="small" type="primary" text @click="handleView(row)">详情</el-button>
							</template>
						</el-table-column>
					</q-table>
				</el-main>
			</el-container>
		</el-container>

		<q-dialog v-model="state.dialog.moreVisible" :title="state.dialog.msgTitle">
			<div class="markdown-body">
				<v-md-preview :text="state.dialog.msg"></v-md-preview>
			</div>
		</q-dialog>
	</q-container>
</template>

<script setup lang="ts" name="notices">
	import { reactive, ref, onMounted, watch } from 'vue';
	import { getCardMode } from '/@/utils/common';
	import { msgInfo } from '/@/stores/msg';
	import { getMessageList, setReadMessage, getReadMessage } from '/@/api/system/message';

	const activePage  = ref("broadcast");
	const state = reactive({
		loading: false,
		msgData: [] as any,
		tableData: [] as any,
		readData: [] as any,
		filterMsg: '',
		msgState: '',
		badge: {
			broadcast: '',
			mail: '',
		},
		dialog :{
			moreVisible: false,
			msgTitle: '',
			msg: '',
		}
	});

	const getNoReadNumber = (kind : number) => {
		var num = state.msgData.filter((item: any) => item.kind === kind && item.isread !== 1).length;
		return num > 0 ? num : '';
	};

	const doFilter = () => {
		state.tableData = state.msgData.filter((item:any) => {
			let flag = false;
			let kind = activePage.value==="broadcast" ? 0:1
			if(state.msgState === ''){
				flag = item.kind === kind;
			}else if(state.msgState === '0'){
				flag = item.kind === kind && item.isread !== 1;
			}else{
				flag = item.kind === kind && item.isread === 1;
			}

			if(state.filterMsg !== ''){
				let fields = ["title", "content"];
				flag = flag && fields.some((key) => {
					return String(item[key]).toLowerCase().indexOf(state.filterMsg) > -1;
				});
			}

			state.badge.broadcast = <string>getNoReadNumber(0);
			state.badge.mail = <string>getNoReadNumber(1);
			return flag;
		});
	};

	const handleView = (row: any) => {
		state.dialog.msgTitle = row.title;
		state.dialog.msg = row.content;
		state.dialog.moreVisible = true;
		// 标记已读
		if(row.isread !== 1){
			setReadMessage({ id: [row.id] }).then(() => {
				state.msgData.filter((item:any) => {
					if(item.id === row.id){
						item.isread = 1;
						doFilter();
					}
				});
			});
		}
	};

	const doCommand = (cmd: string) => {
		if(cmd === 'readall'){
			let ids : any = [];
			state.msgData.forEach((item:any) => {
				if(item.isread !== 1)
					ids.push(item.id)
			})
			if(ids.length > 0){
				setReadMessage({ id: ids }).then(() => {
					state.msgData.filter((item: any) => item.isread = 1);
					doFilter();
				});
			}
		}
	};

	watch(
			()=>state.filterMsg,
			(val)=>{
				doFilter();
			}
	);

	const loadData = () => {
		getReadMessage({}).then((resp:any) => {
			let data = resp.data.list || [];
			data.forEach((item: any) => {
				state.readData.push(item.msg_id);
			})

			let params = {
				pageNum: 1,
				pageSize: 100,
				user: 'Y',
			}
			getMessageList(params).then((resp:any) => {
				state.msgData = resp.data.list || [];
				if(state.readData.length > 0){
					state.msgData.forEach((item: any) => {
						item.isread = state.readData.includes(item.id) ? 1 : 0;
					})
				}

				doFilter();
			});
		})
	};

	onMounted(() => {
		loadData();
	});
</script>

<style lang="scss" scoped>
	.back-container {
		:deep(.el-card__body) {
			padding: 0 !important;
		}
		&-tabs {
			:deep(.el-tabs__header) {
				margin-right: 2px;
			}
			:deep(.el-tabs__item) {
				min-width: 120px;
				padding-right: 35px;
				padding-top: 15px;
			}
		}
	}
	.table-box {
		height: 100%;

		&-item {
			display: flex;
			height: 42px;
			&-icon {
				width: 60px;
				height: 100%;
				color: var(--el-color-primary-light-5);
			}
			&-content {
				padding-left: 5px;
				&-title {
					color: var(--el-text-color);
					font-weight: bold;
				}
				&-text {
					color: var(--el-text-color);
				}
				&-desc {
					padding-right: 15px;
					color: var(--el-text-color-placeholder);
				}
			}
		}
	}
	.markdown-body {
		padding: 0 15px;
		pre {
			background-color: var(--el-bg-color);
			color: var(--el-text-color-primary);
		}
	}
</style>