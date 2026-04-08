
/**
 * @Descripttion: table数据字典
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-07-05 11:30
 * @LastEditors:
 * @LastEditTime:
 *
 *	title: 标题，用于导出时的文件名
 *	tableName: 表名，用于缓存自定义表格配置
 *	rowKey: 记录主键key
 *	search: 数据筛选配置，根据配置自动创建过滤条件页面
 *  filter: 数据过滤配置，配置固定过滤条件
 *	columns: table列配置，根据要显示列信息，常用配置说明如下：
 *		key: 字段名，唯一
 *		width: 显示宽度，单位px，未配置时auto
 *		title: 列标题
 *		sortable: 是否排序
 *		fixed: 固定类型 true / left / right, 默认不固定
 *		show: 是否显示
 *	 	filters: 列筛选配置
 *	 	formatter: 数据格式化回调函数
 *	 	align: 对齐方式，left / center / right 默认left
 *	 	dicts: 配置字词, 如: [{text: '启用', value: 1}, {text: '禁用', value: 0}] 根据value值替换成text, 也可加入类型type（包含success/info/warning/danger）
 *	 	summary: 合计类型 static/total/sum
 *	 	showOverflowTooltip: 当列宽度不够显示不完整内容时，鼠标经过时是否提示内容
 *
*/

export const formatTest = (val: any) => {
		return val;
};

export const tableDict : any = {
	tableDict: {
		title: '字典管理',
		tableName: 'dictTable',
		rowKey: 'dictId',
		search: [
			{ label: '字典类型', prop: 'dictType', placeholder: '请输入字典类型', required: false, type: 'input' },
			{ label: '字典名称', prop: 'dictName', placeholder: '请输入字典名称', required: false, type: 'input' },
			{
				label: '字典状态',
				prop: 'status',
				placeholder: '请选择',
				required: false,
				type: 'select',
				options: [
					{ label: '启用', value: '1' },
					{ label: '禁用', value: '0' },
				],
			},
		],
		columns: [
			{ key: 'dictId', width: '70', title: '字典ID', fixed: 'left', show: false },
			{ key: 'dictName', width: '200', title: '字典名称', fixed: 'left', show: true, showOverflowTooltip: true, formatter: (row:any) => { return formatTest(row.dictName) } },
			{ key: 'dictType', width: '260', title: '字典类型', show: true },
			{ key: 'status', width: '100', title: '字典状态', align: "center", show: true,
				// filters: [{text: '启用', value: 1}, {text: '禁用', value: 0}],
				dicts: [{text: '启用', value: 1, type: 'success'}, {text: '禁用', value: 0, type: 'info'}]
			},
			{ key: 'remark', title: '字典描述', show: true, showOverflowTooltip: true },
			{ key: 'createdAt', width: '180', title: '创建时间', align: "center", show: true },
			{ key: 'opt', width: '150', title: '操作', align: "center", fixed: 'right', show: true },
		],
	},

	tableDictList: {
		title: '字典列表',
		tableName: 'dictTableList',
		rowKey: 'dictCode',
		search: [
			// { label: '字典类型', prop: 'dictType', placeholder: '请输入字典类型', required: true, type: 'input', id: 'dictType' },
			{ label: '字典标签', prop: 'dictLabel', placeholder: '请输入字典标签', required: false, type: 'input' },
			{
				label: '字典状态',
				prop: 'status',
				placeholder: '请选择',
				required: false,
				type: 'select',
				options: [
					{ label: '启用', value: '1' },
					{ label: '禁用', value: '0' },
				],
			},
		],
		columns: [
			{ key: 'dictCode', width: '80', title: '字典编码', fixed: 'left', show: false },
			{ key: 'dictLabel', width: '200', title: '字典标签', fixed: 'left', show: true },
			{ key: 'dictValue', width: '120', title: '字典键值', show: true },
			{ key: 'dictSort', width: '100', title: '字典排序', show: true },
			{ key: 'status', width: '100', title: '字典状态', align: "center", show: true },
			{ key: 'createdAt', width: '180', title: '创建时间', align: "center", show: true },
			{ key: 'remark', title: '备注', show: true, showOverflowTooltip: true },
			{ key: 'opt', width: '150', title: '操作', align: "center", fixed: 'right', show: true },
		],
	},

	tableMockTest: {
		title: 'Mock测试',
		tableName: 'mock-test-data',
		rowKey: 'date',
		columns: [
			{ key: 'date', width: '120', title: '日期', sortable: true, fixed: 'left', show: true },
		],
	},

	tableOperDict: {
		title: '操作日志列表',
		tableName: 'operTableList',
		filter: [
			{
				title: '系统模块',
				isMore: false,
				isShowMore: false,
				prop: 'title',
				children: [
					{ label: '全部', 		value: '', active: true, },
					{ label: '操作日志',  active: false, },
					{ label: '岗位管理',  active: false, },
					{ label: '用户管理',  active: false, },
					{ label: '角色管理',  active: false, },
					{ label: '参数管理',  active: false, },
					{ label: '字典管理',  active: false, },
					{ label: '部门管理',  active: false, },
					{ label: '在线用户',  active: false, },
					{ label: '服务监控',  active: false, },
					{ label: '添加岗位',  active: false, },
					{ label: '添加菜单',  active: false, },
					{ label: '添加角色',  active: false, },
					{ label: '添加部门',  active: false, },
				],
			},
			{
				title: '请求方式',
				isMore: false,
				isShowMore: false,
				prop: 'requestMethod',
				children: [
					{ label: '全部', value: '', active: true, },
					{ label: '读取', value: 'GET', active: false, },
					{ label: '新增', value: 'POST', active: false, },
					{ label: '修改', value: 'PUT', active: false, },
					{ label: '删除', value: 'DELETE', active: false, },
				],
			},
			{
				title: '操作人员',
				isMore: false,
				isShowMore: false,
				prop: 'operName',
				children: [
					{ label: '全部', 		value: '', active: true, },
					{ label: 'admin', 	value: 'admin', active: false, },
					{ label: 'debug', 	value: 'debug', active: false, },
					{ label: 'commmon', value: 'commmon', active: false, },
					{ label: 'guest', 	value: 'guest', active: false, },
				],
			},
		],
	},

	table_im_message: {
		title: '消息列表',
		tableName: 'table_im_message',
		rowKey: 'id',
		search: [
			{ label: '关键字', prop: 'keyWords', placeholder: '请输入查询关键字', required: false, type: 'input' },
			{ label: '是否群消息', prop: 'is_group', placeholder: '', required: false, type: 'select',
				options: [
					{ label: '否', value: '0' },
					{ label: '是', value: '1' },
				],
			},
			{ label: '是否已读', prop: 'is_read', placeholder: '', required: false, type: 'select',
				options: [
					{ label: '未读', value: '0' },
					{ label: '已读', value: '1' },
				],
			},
		],
		treeOptions: [
			{ name: '所有消息', id: 'all' },
			{ name: '文本', id: 'text' },
			{ name: '图片', id: 'image' },
			{ name: '音频', id: 'voice' },
			{ name: '视频', id: 'video' },
			{ name: '文件', id: 'file' },
		],
		treeProps: {
			id: 'id',
			children: 'children',
			label: 'name',
		},
		columns: [
			{ key: 'msg_id', width: '100', title: '消息编号', fixed: 'left', keylink: false, show: true },
			{ key: 'type', width: '100', title: '消息类型', show: true },
			{ key: 'is_group', width: '100', title: '是否群消息', show: true },
			{ key: 'is_read', width: '100', title: '是否已读', show: true },
			{ key: 'created_at', width: '160', title: '发送时间', align: "center", show: true, showOverflowTooltip: true },
			{ key: 'from_user_name', width: '100', title: '发送方', show: true, keylink: true, showOverflowTooltip: true },
			{ key: 'to_user_name', width: '100', title: '接收方', show: true, keylink: true, showOverflowTooltip: true },
			{ key: 'content', width: '400', title: '消息内容', show: true, showOverflowTooltip: true },
		],
	},

	table_im_notice: {
		title: '消息广播',
		tableName: 'table_im_notice',
		rowKey: 'id',
		search: [
			{ label: '关键字', prop: 'keyWords', placeholder: '请输入查询关键字', required: false, type: 'input' },
		],
		columns: [
			{ key: 'msg_id', width: '100', title: '消息编号', fixed: 'left', keylink: false, show: true },
			{ key: 'type', width: '100', title: '消息类型', show: true },
			{ key: 'created_at', width: '160', title: '发送时间', align: "center", show: true, showOverflowTooltip: true },
			{ key: 'from_user_name', width: '100', title: '发送方', show: true, showOverflowTooltip: true },
			{ key: 'to_user_name', width: '100', title: '接收方', show: true, keylink: true, showOverflowTooltip: true },
			{ key: 'content', width: '400', title: '消息内容', show: true, showOverflowTooltip: true },
		],
	},

	table_im_message_ex: {
		title: '消息列表',
		tableName: 'table_im_message_ex',
		rowKey: 'id',
		search: [],
		columns: [
			{ key: 'msg_id', width: '100', title: '消息编号', fixed: 'left', keylink: false, show: true },
			{ key: 'type', width: '100', title: '消息类型', show: true },
			{ key: 'is_read', width: '100', title: '是否已读', show: true },
			{ key: 'created_at', width: '160', title: '发送时间', align: "center", show: true, showOverflowTooltip: true },
			{ key: 'from_user_name', width: '100', title: '发送方', show: true, showOverflowTooltip: true },
			{ key: 'to_user_name', width: '100', title: '接收方', show: true, showOverflowTooltip: true },
			{ key: 'content', width: '400', title: '消息内容', show: true, showOverflowTooltip: true },
		],
	},

	table_im_group: {
		title: '群组列表',
		tableName: 'table_im_group',
		rowKey: 'id',
		search: [
			{ label: '关键字', prop: 'keyWords', placeholder: '请输入查询关键字', required: false, type: 'input' },
			{ label: '类型', prop: 'is_public', placeholder: '', required: false, type: 'select',
				options: [
					{ label: '公开', value: '1' },
					{ label: '私密', value: '0' },
				],
			},
		],
		columns: [
			{ key: 'id', width: '100', title: '编号', fixed: 'left', keylink: false, show: false },
			{ key: 'avatar', width: '80', title: '头像', fixed: 'left', show: true },
			{ key: 'name', width: '300', title: '群组名称', fixed: 'left', show: true, showOverflowTooltip: true },
			{ key: 'owner_name', width: '100', title: '群主', show: true, keylink: true, showOverflowTooltip: true },
			{ key: 'created_name', width: '100', title: '创建人', show: true, keylink: true, showOverflowTooltip: true },
			{ key: 'created_at', width: '160', title: '创建时间', align: "center", show: true, showOverflowTooltip: true },
			{ key: 'is_public', width: '100', title: '是否公开', align: "center", show: false },
			{ key: 'opt', width: '280', title: '操作', align: "center", fixed: 'right', show: true },
		],
	},

	table_im_user: {
		title: '用户列表',
		tableName: 'table_im_user',
		rowKey: 'id',
		search: [
			{ label: '关键字', prop: 'keyWords', placeholder: '请输入查询关键字', required: false, type: 'input' },
		],
		columns: [
			{ key: 'avatar', width: '100', title: '头像', fixed: 'left', show: true, summary: 'static' },
			{ key: 'userName', width: '100', title: '账号', fixed: 'left', show: true, showOverflowTooltip: true, summary: 'total' },
			{ key: 'userNickname', width: '200', title: '昵称', fixed: 'left', show: true, showOverflowTooltip: true },
			{ key: 'sex', width: '80', title: '性别', show: true },
			{ key: 'userStatus', width: '80', title: '状态', show: true },
			{ key: 'isAdmin', width: '150', title: '类型', show: true, showOverflowTooltip: true },
			{ key: 'createdAt', width: '160', title: '注册时间', align: "center", show: true, showOverflowTooltip: true },
			{ key: 'mobile', width: '150', title: '手机号', show: true, showOverflowTooltip: true },
			{ key: 'userEmail', width: '150', title: '邮箱', show: true, showOverflowTooltip: true },
			{ key: 'opt', width: '200', title: '操作', align: "right", fixed: 'right', show: true },
		],
	},
	table_im_group_user: {
		title: '用户列表',
		tableName: 'table_im_group_user',
		rowKey: 'id',
		search: [],
		columns: [
			{ key: 'avatar', width: '80', title: '头像', fixed: 'left', show: true },
			{ key: 'account', width: '100', title: '账号', fixed: 'left', show: true, showOverflowTooltip: true },
			{ key: 'displayName', width: '200', title: '昵称', fixed: 'left', show: true, showOverflowTooltip: true },
			{ key: 'sex', width: '80', title: '性别', show: true },
			{ key: 'email', width: '150', title: '邮箱', show: true, showOverflowTooltip: true },
			{ key: 'opt', width: '80', title: '操作', align: "center", fixed: 'right', show: true },
		],
	},

	table_business_files: {
		title: '文件管理',
		tableName: 'table_business_files',
		rowKey: 'id',
		search: [
			{ label: '文件名称', prop: 'KeyWords', placeholder: '请输入条件', required: false, type: 'input' },
			{ label: '文件类型', prop: 'cate', placeholder: '', required: false, type: 'select',
				options: [
					{ label: '所有文件', value: '0' },
					{ label: '文档', value: '1' },
					{ label: '图片', value: '2' },
					{ label: '音频', value: '3' },
					{ label: '视频', value: '4' },
					{ label: '其它', value: '5' },
				],
			},
			{ label: '方向', prop: 'role', placeholder: '', required: false, type: 'select',
				options: [
					{ label: '所有文件', value: '0' },
					{ label: '我发送的', value: '1' },
					{ label: '我接收的', value: '2' },
				],
			},
		],
		treeOptions: [
			{ name: '所有文件', id: '0' },
			{ name: '文档', id: '1' },
			{ name: '图片', id: '2' },
			{ name: '音频', id: '3' },
			{ name: '视频', id: '4' },
			{ name: '其它', id: '5' },
		],
		treeProps: {
			id: 'id',
			children: 'children',
			label: 'name',
		},
		columns: [
			{ key: 'file_name', width: '300', title: '文件名称', fixed: 'left', keylink: false, show: true, showOverflowTooltip: true },
			{ key: 'file_size', width: '120', title: '文件大小', show: true, showOverflowTooltip: true },
			{ key: 'kind', width: '100', title: '类型', show: true, showOverflowTooltip: true },
			{ key: 'created_name', width: '120', title: '上传人', show: true, keylink: true, showOverflowTooltip: true },
			{ key: 'create_time', width: '160', title: '上传时间', align: "center", show: true, showOverflowTooltip: true },
			{ key: 'remark', title: '备注', show: true, showOverflowTooltip: true },
			{ key: 'opt', width: '200', title: '操作', align: "left", fixed: 'right', show: true },
		],
	},

	table_im_words: {
		title: '敏感词列表',
		tableName: 'table_im_words',
		rowKey: 'id',
		search: [
			{ label: '关键字', prop: 'keyWords', placeholder: '请输入查询关键字', required: false, type: 'input' },
		],
		columns: [
			{ key: 'id', width: '80', title: '编号', fixed: 'left', show: true },
			{ key: 'lang', width: '100', title: '词言', show: true, showOverflowTooltip: true },
			{ key: 'word', width: '200', title: '敏感词', show: true, showOverflowTooltip: true },
			{ key: 'status', width: '80', title: '状态', show: true },
			{ key: 'created_name', width: '100', title: '创建人', show: true, keylink: true, showOverflowTooltip: true },
			{ key: 'createdAt', width: '160', title: '创建时间', align: "center", show: true, showOverflowTooltip: true },
			{ key: 'opt', width: '150', title: '操作', align: "center", fixed: 'right', show: true },
		],
	},
}
