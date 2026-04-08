
/**
 * @Descripttion: 数据字典
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-08-01 11:30
 * @LastEditors:
 * @LastEditTime:
 *
 *
*/
export const dataDict : any = {
	dictKind:[
		{
			"id": "1",
			"code": "dataDict",
			"name": "数据字典",
			"children": [
				{
					"id": "11",
					"code": "data_static",
					"name": "静态字典",
					"children": [
						{
							"id": "111",
							"code": "notice",
							"name": "通知类型"
						},
						{
							"id": "112",
							"code": "sex",
							"name": "性别"
						},
					]
				},
				{
					"id": "12",
					"code": "data_dynamics",
					"name": "动态字典"
				}
			]
		},
		{
			"id": "2",
			"code": "string",
			"name": "字符字典",
			"children": [
				{
					"id": "21",
					"code": "string_const",
					"name": "常量"
				},
				{
					"id": "22",
					"code": "string_code",
					"name": "错误码"
				},
				{
					"id": "23",
					"code": "string_msg",
					"name": "提示信息"
				}
			]
		},
		{
			"id": "3",
			"code": "function",
			"name": "功能码",
			"children": [
				{
					"id": "31",
					"code": "function_query",
					"name": "查询功能码"
				}
			]
		},
		{
			"id": "4",
			"code": "protocol",
			"name": "协议",
			"children": [
				{
					"id": "41",
					"code": "protocol_send",
					"name": "申请协议"
				},
				{
					"id": "42",
					"code": "protocol_recv",
					"name": "接收协议"
				}
			]
		}
	],
	dictItems: [
		{
			"id": "111",
			"key": "1",
			"name": "发布通知",
			"memo": "",
			"enabled": "1"
		},
		{
			"id": "111",
			"key": "2",
			"name": "转发通知",
			"memo": "",
			"enabled": "1"
		},
		{
			"id": "111",
			"key": "3",
			"name": "事务通知",
			"memo": "",
			"enabled": "0"
		},
		{
			"id": "112",
			"key": "0",
			"name": "女",
			"memo": "",
			"enabled": "1"
		},
		{
			"id": "112",
			"key": "1",
			"name": "男",
			"memo": "",
			"enabled": "1"
		},
		{
			"id": "112",
			"key": "2",
			"name": "保密",
			"memo": "",
			"enabled": "0"
		},
	]
}

export const dateRangeShortcuts : any = [
	{
		text: '最近一周',
		value: () => {
			const end = new Date()
			const start = new Date()
			start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
			return [start, end]
		},
	},
	{
		text: '最近一月',
		value: () => {
			const end = new Date()
			const start = new Date()
			start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
			return [start, end]
		},
	},
	{
		text: '最近3个月',
		value: () => {
			const end = new Date()
			const start = new Date()
			start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
			return [start, end]
		},
	},
]