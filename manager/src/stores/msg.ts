import { defineStore } from 'pinia';
import { MsgInfoStates, MsgInfoState } from '/@/types/pinia';
import { DemoNoticeList } from '/@/views/system/notices/mock';
import { getNoReadMessage } from '/@/api/system/message';

/**
 * 未读消息信息
 */

let timer: any = null;
let DemoMsgData = [
	{
		id: 1,
		title: '关于版本发布的通知',
		describe: '增加通知中心',
		time: '2023-08-25',
		status: 'Y',
	},
	{
		id: 2,
		title: '关于版本发布的通知',
		describe: '增加AI测试模块',
		time: '2023-08-08',
		status: 'Y',
	},
	{
		id: 3,
		title: '关于版本发布的通知',
		describe: '系统优化，修改已知的BUG',
		time: '2023-08-01',
		status: 'Y',
	},
];

export const msgInfo = defineStore('msg', {
	state: (): any => ({
			unreadCount: '',
			unreadList: [] as any,
	}),
	actions: {
		async setMsgInfos(data: any) {
			this.unreadList = data;
			this.unreadCount = data?.length;
		},

		async queryMsg() {
			// 测试数据
			let resp = await getNoReadMessage({});
			let data = resp.data?.list || [];
			this.unreadList = data.filter((item:any) => item.isread !== 1);
			if(this.unreadList.length > 0)
				this.unreadCount = this.unreadList.length;
			else
				this.unreadCount = '';
		},
		async setMsgReadStatus(id:any) {
			this.unreadList?.filter((item:any) => {
				if(item.id == id){
					item.isread = 1;
				}
			});
			this.queryMsg();
		},
		async clearMsg(){
			this.unreadList = [];
			this.queryMsg();
		},
		async startUpdate() {
			// if(timer === null){
			// 	timer = setInterval(() => {
			// 		this.queryMsg();
			// 	}, 5000);
			// }
		},

		async stopUpdate() {
			// if(timer){
			// 	timer && clearInterval(timer);
			// 	timer = null;
			// }
		},
	},
});
