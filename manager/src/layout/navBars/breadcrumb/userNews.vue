<template>
	<div class="layout-navbars-breadcrumb-user-news">
		<div class="head-box">
			<div class="head-box-title">{{ $t('message.user.newTitle') }}</div>
			<div class="head-box-btn" v-if="newsList.length > 0" @click="onAllReadClick">{{ $t('message.user.newBtn') }}</div>
		</div>
		<div class="content-box">
			<template v-if="newsList.length > 0">
				<ul class="msg-list">
					<li v-for="item in newsList" v-bind:key="item.id">
						<a target="_blank" @click="readmsg(item)">
							<div class="msg-list__icon">
								<el-badge value="new" type="danger">
									<el-icon :size="32" class="table-box-item-icon">
										<ele-Message />
									</el-icon>
								</el-badge>
							</div>
							<div class="msg-list__main">
								<h2>{{item.title}}</h2>
							</div>
							<div class="msg-list__time">
								<p>{{item.time}}</p>
							</div>
						</a>
					</li>
				</ul>
			</template>
			<el-empty :image-size="100" :description="$t('message.user.newDesc')" v-else></el-empty>
		</div>
		<div class="foot-box" @click="onGoToGiteeClick" v-if="newsList.length > 0">{{ $t('message.user.newGo') }}</div>

		<q-dialog v-model="moreVisible" :title="msgTitle">
			<div class="markdown-body">
				<v-md-preview :text="msg"></v-md-preview>
			</div>
		</q-dialog>
	</div>
</template>

<script lang="ts">
import { onMounted, onUnmounted, reactive, ref, toRefs, defineComponent, getCurrentInstance } from 'vue';
import MQTT from '/@/api/mqtt/index';
import mqtt from 'mqtt';
import { OnMessageCallback } from 'mqtt';
import { useRouter } from 'vue-router';
import { msgInfo } from '/@/stores/msg';
import { storeToRefs } from 'pinia';
import { setReadMessage } from '/@/api/system/message';

export default defineComponent({
	name: 'layoutBreadcrumbUserNews',
	setup(props, { emit }) {
		const { proxy } = <any>getCurrentInstance();
		const router = useRouter();
		const msgStores = msgInfo();
		const { unreadCount, unreadList } = storeToRefs(msgStores);
		const mq = ref<MQTT | null>(null);

		// const initMQ = (url: string, callback: OnMessageCallback) => {
		// 	const subscription = ref({
		// 		topic: "topic/mqttx",
		// 		qos: 0 as mqtt.QoS,
		// 	});
		// 	// 设置订阅参数
		// 	mq.value = new MQTT(url, subscription);
		// 	// 初始化mqtt
		// 	mq.value.init();
		// 	// 链接mqtt
		// 	mq.value.link();
		// 	getMessage(callback);
		// };

		const getMessage = (callback: OnMessageCallback) => {
			mq.value?.get(callback);
		};
		
		const state = reactive({
			moreVisible: false,
			msgTitle: '',
			msg: '',
			newsList: unreadList,
		});

		onMounted(() => {
			proxy.mittBus.emit('messageChange', unreadCount);

			// 	initMQ('mqtt://localhost/', (topic, message) => {
			// 		console.log(topic, message);
			// 		// state.newsList.push(message);
			// 	});
		});

		onUnmounted(() => {
			if (mq.value) {
				mq.value.unsubscribes();
				mq.value.over();
			}
		});

		// 全部已读点击
		const onAllReadClick = () => {
			msgStores.clearMsg();
			emit('closed');
		};
		// 前往通知中心点击
		const onGoToGiteeClick = () => {
			router.push('/home/notices');
			emit('closed');
		};
		// 打开消息
		const readmsg = (item: any) => {
			state.msgTitle = item.title;
			state.msg = item.content;
			state.moreVisible = true;
			// 标记已读
			if(item.isread !== 1){
				setReadMessage({ id: [item.id] }).then(() => {
					msgStores.setMsgReadStatus(item.id);
					if(state.newsList!.length <= 0){
						emit('closed');
					}
				});
			}
		};

		return {
			onAllReadClick,
			onGoToGiteeClick,
			readmsg,
			...toRefs(state),
		};
	},
});
</script>

<style scoped lang="scss">
.layout-navbars-breadcrumb-user-news {
	.head-box {
		display: flex;
		border-bottom: 1px solid var(--el-border-color-lighter);
		box-sizing: border-box;
		color: var(--el-text-color-primary);
		justify-content: space-between;
		height: 35px;
		align-items: center;
		.head-box-btn {
			color: var(--el-color-primary);
			font-size: 13px;
			cursor: pointer;
			opacity: 0.8;
			&:hover {
				opacity: 1;
			}
		}
	}
	.content-box {
		font-size: 13px;
		.content-box-item {
			padding-top: 12px;
			&:last-of-type {
				padding-bottom: 12px;
			}
			.content-box-msg {
				color: var(--el-text-color-secondary);
				margin-top: 5px;
				margin-bottom: 5px;
			}
			.content-box-time {
				color: var(--el-text-color-secondary);
			}
		}
	}
	.foot-box {
		height: 32px;
		color: var(--el-color-primary);
		background-color: var(--el-color-primary-light-8);
		border-radius: 16px;
		font-size: 13px;
		cursor: pointer;
		opacity: 0.8;
		display: flex;
		align-items: center;
		justify-content: center;
		border-top: 1px solid var(--el-border-color-lighter);
		margin-top: 10px;
		&:hover {
			opacity: 1;
		}
	}
	:deep(.el-empty__description p) {
		font-size: 13px;
	}
}
.msg-list li {border-top:1px solid var(--el-border-color-lighter);}
.msg-list li a {display: flex;padding:5px;}
.msg-list li a:hover {background: var(--el-color-primary-light-9);}
.msg-list__icon {width: 40px;margin-right: 30px; padding-top: 10px; color: var(--el-color-primary-light-5);}
.msg-list__main {flex: 1;}
.msg-list__main h2 {padding-top: 15px; font-size: 15px;font-weight: normal;color: var(--el-text-color);}
.msg-list__time {width: 140px;text-align: left;color: var(--el-text-color-disabled);}
</style>
