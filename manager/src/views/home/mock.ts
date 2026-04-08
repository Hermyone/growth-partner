import * as echarts from 'echarts';

export const colorList = ['#51A3FC', '#36C78B', '#FEC279', '#968AF5', '#E790E8'];

const HOME_URL = '/' // http://localhost:8888/#/home'

export const homeNums = [
	{
		value: '125,12',
		field: 'project_num',
		title: '数量',
		icon: 'iconfont icon-diannao-shuju',
		path: '/project/manager',
		valueColor: '#FF6462',
		iconBg: '--next-color-primary-lighter',
		iconColor: '--el-color-primary',
	},
	{
		value: '653,33',
		field: 'user_num',
		title: '占位',
		icon: 'iconfont icon-icon-',
		path: '/system/avatar/list',
		valueColor: '#6690F9',
		iconBg: '--next-color-success-lighter',
		iconColor: '--el-color-success',
	},
	{
		value: '125,65',
		field: 'project_user_num',
		title: '占位',
		icon: 'iconfont icon-shenqingkaiban',
		path: '/statistics/number',
		valueColor: '#6690F9',
		iconBg: '--next-color-warning-lighter',
		iconColor: '--el-color-warning',
	},
	{
		value: '520,43',
		field: 'sign_num',
		title: '占位',
		icon: 'iconfont icon-xuanzeqi',
		path: '/sign/in',
		valueColor: '#FF6462',
		iconBg: '--next-color-danger-lighter',
		iconColor: '--el-color-danger',
	},
];

export const quickEntrance = [
	{
		title: '用戶管理',
		icon: 'ele-User',
		path: '/system/avatar/list',
		iconColor: '--el-color-primary',
	},
	{
		title: '占位',
		icon: 'ele-Notebook',
		path: '/project/manager',
		iconColor: '--el-color-primary',
	},
	{
		title: '占位',
		icon: 'ele-SoldOut',
		path: '/project/check/in',
		iconColor: '--el-color-primary',
	},
	{
		title: '占位',
		icon: 'ele-Tickets',
		path: '/sign/list',
		iconColor: '--el-color-primary',
	},
	{
		title: '占位',
		icon: 'ele-Aim',
		path: '/attendance/work',
		iconColor: '--el-color-primary',
	},
	{
		title: '占位',
		icon: 'ele-Check',
		path: '/production/examine',
		iconColor: '--el-color-primary',
	},
	{
		title: '占位',
		icon: 'ele-Coordinate',
		path: '/attendance/examine',
		iconColor: '--el-color-primary',
	},
	{
		title: '占位',
		icon: 'ele-EditPen',
		path: '/production/workload',
		iconColor: '--el-color-primary',
	},
]

export const homeThree = [
	{
		icon: 'iconfont icon-yangan',
		label: '浅粉红',
		value: '2.1%OBS/M',
		iconColor: '#F72B3F',
	},
	{
		icon: 'iconfont icon-wendu',
		label: '深红(猩红)',
		value: '30℃',
		iconColor: '#91BFF8',
	},
	{
		icon: 'iconfont icon-shidu',
		label: '淡紫红',
		value: '57%RH',
		iconColor: '#88D565',
	},
	{
		icon: 'iconfont icon-shidu',
		label: '弱紫罗兰红',
		value: '107w',
		iconColor: '#88D565',
	},
	{
		icon: 'iconfont icon-zaosheng',
		label: '中紫罗兰红',
		value: '57DB',
		iconColor: '#FBD4A0',
	},
	{
		icon: 'iconfont icon-zaosheng',
		label: '紫罗兰',
		value: '57PV',
		iconColor: '#FBD4A0',
	},
	{
		icon: 'iconfont icon-zaosheng',
		label: '暗紫罗兰',
		value: '517Cpd',
		iconColor: '#FBD4A0',
	},
	{
		icon: 'iconfont icon-zaosheng',
		label: '幽灵白',
		value: '12kg',
		iconColor: '#FBD4A0',
	},
	{
		icon: 'iconfont icon-zaosheng',
		label: '海军蓝',
		value: '64fm',
		iconColor: '#FBD4A0',
	},
];

export function getLineOption(bgColor:any, txtColor:any, lineColor:any, xData: any, yData1: any, yData2: any) {
	return {
		backgroundColor: bgColor,
		title: {
			text: '当日生产计划完成情况',
			left: 'center',
			textStyle: { fontSize: '15', color: txtColor },
		},
		grid: { top: 70, right: 20, bottom: 30, left: 30 },
		tooltip: { trigger: 'axis' },
		legend: { data: ['计划量', '实情完成量'], right: 0 },
		xAxis: {
			data: xData,
		},
		yAxis: [
			{
				type: 'value',
				name: '工作量',
				splitLine: { show: true, lineStyle: { type: 'dashed', color: lineColor } },
			},
		],
		series: [
			{
				name: '计划量',
				type: 'line',
				symbolSize: 4,
				symbol: 'circle',
				smooth: true,
				data: yData1,
				lineStyle: { color: '#409eff' },
				itemStyle: { color: '#409eff', borderColor: '#409eff' },
				areaStyle: {
					color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
						{ offset: 0, color: '#fe9a8bb3' },
						{ offset: 1, color: '#fe9a8b03' },
					]),
				},
			},
			{
				name: '实情完成量',
				type: 'line',
				symbolSize: 4,
				symbol: 'circle',
				smooth: true,
				data: yData2,
				lineStyle: { color: '#9373ee' },
				itemStyle: { color: '#9373ee', borderColor: '#9373ee' },
				areaStyle: {
					color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
						{ offset: 0, color: '#9373eeb3' },
						{ offset: 1, color: '#9373ee03' },
					]),
				},
				emphasis: {
					itemStyle: {
						color: {
							type: 'radial',
							x: 0.5,
							y: 0.5,
							r: 0.5,
							colorStops: [
								{ offset: 0, color: '#9E87FF' },
								{ offset: 0.4, color: '#9E87FF' },
								{ offset: 0.5, color: '#fff' },
								{ offset: 0.7, color: '#fff' },
								{ offset: 0.8, color: '#fff' },
								{ offset: 1, color: '#fff' },
							],
						},
						borderColor: '#9E87FF',
						borderWidth: 2,
					},
				},
			},
		],
	};
}

export function getPieOption(bgColor: any, txtColor: any, isIsDark: boolean, xData: any, yData: any) : any {
	var getname = xData;
	var getvalue = yData;
	var data = [];
	for (var i = 0; i < getname.length; i++) {
		data.push({ name: getname[i], value: getvalue[i] });
	}
	return {
		backgroundColor: bgColor,
		title: {
			text: '项目用工占比',
			left: 'center',
			textStyle: { fontSize: '15', color: txtColor },
		},
		tooltip: { trigger: 'item', formatter: '{b} <br/> {c} 人' },
		legend: {
			type: 'scroll',
			orient: 'vertical',
			right: '0%',
			left: '65%',
			top: 'center',
			itemWidth: 14,
			itemHeight: 14,
			data: getname,
			textStyle: {
				rich: {
					name: {
						fontSize: 16,
						width: 200,
						height: 35,
						padding: [0, 0, 0, 60],
						color: txtColor,
					},
					rate: {
						fontSize: 16,
						height: 35,
						width: 40,
						padding: [0, 0, 0, 30],
						color: txtColor,
					},
				},
			},
		},
		series: [
			{
				type: 'pie',
				radius: ['60', isIsDark ? '90' : '100'],
				center: ['32%', '50%'],
				padAngle: 5,
				itemStyle: {
					borderRadius: 10,
					borderColor: isIsDark ? '#000' : '#fff',
					color: function (params: any) {
						return colorList[params.dataIndex];
					},
				},
				label: { show: false },
				labelLine: { show: false },
				data: data,
			},
		],
	};
}