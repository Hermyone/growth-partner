import sysConfig from '/@/config';
import { requestToObject_ex, requestToObject_m } from '/@/api/items';

const system = [
	{ key: 'test', path: 'post/list', method: 'get', name: '测试' },

	{ key: 'user_list', path: 'user/list', method: 'get', name: '用户列表查询' },
]

const imAPI = {
	files: [
		{ key: 'query', path: 'files/fileList', method: 'post', name: '文件列表查询' },
		{ key: 'add', path: 'files/addFile', method: 'post', name: '增加文档' },
		{ key: 'edit', path: 'files/editFile', method: 'post', name: '修改文档' },
		{ key: 'delete', path: 'files/deleteFile', method: 'post', name: '删除文档' },
	],
	group: [
		{ key: 'query', path: 'group/groupList', method: 'post', name: '群组列表查询' },
		{ key: 'add', path: 'group/addGroup', method: 'post', name: '增加群组' },
		{ key: 'edit', path: 'group/editGroup', method: 'post', name: '修改群组' },
		{ key: 'delete', path: 'group/deleteGroup', method: 'post', name: '删除群组' },
		{ key: 'owner', path: 'group/owner', method: 'post', name: '设置群主' },
	],
	groupUser: [
		{ key: 'query', path: 'group/user/groupUserList', method: 'post', name: '群用户列表查询' },
		{ key: 'add', path: 'group/user/addGroupUser', method: 'post', name: '增加群用户' },
		{ key: 'edit', path: 'group/user/editGroupUser', method: 'post', name: '修改群用户' },
		{ key: 'delete', path: 'group/user/deleteGroupUser', method: 'post', name: '删除群用户' },
		{ key: 'role', path: 'group/user/setManager', method: 'post', name: '设置群用户身份' },
	],
	msg: [
		{ key: 'query', path: 'msg/msgList', method: 'post', name: '消息列表查询' },
		{ key: 'notice', path: 'msg/msgNotice', method: 'post', name: '广播消息列表查询' },
		{ key: 'add', path: 'msg/addMsg', method: 'post', name: '增加消息' },
		{ key: 'edit', path: 'msg/editMsg', method: 'post', name: '修改消息' },
		{ key: 'delete', path: 'msg/deleteMsg', method: 'post', name: '删除消息' },
		{ key: 'user', path: 'msg/taMsg', method: 'post', name: 'TA的会话' },
		{ key: 'sendNotice', path: 'notice/msg', method: 'post', name: '发送广播' },
	],
	friend: [
		{ key: 'query', path: 'friend/friendList', method: 'post', name: '好友列表查询' },
		{ key: 'friends', path: 'friend/friends', method: 'post', name: '联系人列表查询' },
	],
	words: [
		{ key: 'query', path: 'words/wordsList', method: 'post', name: '敏感词列表查询' },
		{ key: 'add', path: 'words/addWords', method: 'post', name: '增加敏感词' },
		{ key: 'edit', path: 'words/editWords', method: 'post', name: '修改敏感词' },
		{ key: 'status', path: 'words/editStatus', method: 'post', name: '修改敏感词状态' },
		{ key: 'delete', path: 'words/deleteWords', method: 'post', name: '删除敏感词' },
	],
}

const dataTotal = {
	home: [
		{ key: 'nums', path: 'total/nums', method: 'get', name: '首页数量合计' },
	],
	statistics: [
	]
}

const execProcedure = {
	store: [
		{ key: 'test', path: 'store/test', method: 'post', name: '测试' },
		{ key: 'common', path: 'store/common', method: 'post', name: '执行业务' },
	]
}

// @ts-ignore
export default {
	...requestToObject_ex(system, `${sysConfig.VITE_URL}api/v1/system/`),
	...requestToObject_m(imAPI, `${sysConfig.VITE_URL}api/v1/im/`),
	...requestToObject_m(dataTotal, `${sysConfig.VITE_URL}api/v1/im/`),
	...requestToObject_m(execProcedure, `${sysConfig.VITE_URL}api/v1/im/`),
}
