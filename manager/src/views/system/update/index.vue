<template>
	<q-container :cardMode="getCardMode()" class="back-container">
		<div class="back-container-body">
			<el-timeline v-if="state.logs.length > 0">
				<el-timeline-item
						v-for="(item, index) in state.logs"
						placement="top"
						:hollow="index !== 0"
						:timestamp="item.Title"
						:size="index === 0 || index === state.logs.length-1 ? 'large':'normal'"
						:type="index === 0 || index === state.logs.length-1 ? 'primary':''"
						:icon="Top"
				>
					<div style="display: flex; align-items: center; height: 32px;">
						<span>更新时间： {{ item.Time }}</span>
						<el-button type="primary" size="small" round text @click="onSwitchConext(item)" style="margin-left: 10px;">
							<el-icon v-if="item.more"><ele-ArrowUp /></el-icon>
							<el-icon v-else><ele-ArrowDown /></el-icon>
							{{ item.more ? '关闭' : '查看' }}
						</el-button>
					</div>
					<el-card shadow="hover" :body-style="{ padding: '0!important' }" v-if="item.more">
<!--						<div class="back-container-card">-->
<!--							<div style="padding: 15px">-->
<!--								<h4>{{ item.Title }}</h4>-->
<!--								<div class="update_context" v-html="item.Memo"></div>-->
<!--							</div>-->
<!--						</div>-->
						<v-md-preview :text="item.Memo"></v-md-preview>
					</el-card>
				</el-timeline-item>
			</el-timeline>

			<div v-else>
				<el-empty description="暂无升级日志"></el-empty>
			</div>
		</div>
	</q-container>
</template>

<script setup lang="ts" name="updatelogs">
	import { reactive, ref, onMounted } from 'vue';
	import { getCardMode } from '/@/utils/common';
	import { getUpdaeLogsList } from '/@/api/system/update';
	import { Top } from '@element-plus/icons-vue';

	const state = reactive({
		logs: [] as any,
	});

	const onSwitchConext = (item: any) => {
		if(item.more)
			item.more = false
		else
			item.more = true
	}

	onMounted(() => {
		getUpdaeLogsList({}).then((resp: any) => {
			state.logs = resp.data.list || [];
			if(state.logs.length > 0) state.logs[0].more = true;
		})
	});
</script>

<style lang="scss" scoped>
	.back-container {
		&-body {
			padding: 32px 100px;
		}
		h4 {
			padding: 10px 0!important;
		}
		&-card {
			min-width: 500px;
			display: flex;
			justify-content: space-between;
		}
	}
	.update_context {
		font-size: 14px;
		line-height: 24px;
		padding: 15px;
		font-family: Arial, Tahoma, Helvetica,sans-serif,"微软雅黑";
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
</style>