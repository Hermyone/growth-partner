<template>
	<div class="q-small-title">硬盘资源</div>
	<div class="el-table--enable-row-hover el-table--medium">
		<table cellspacing="0" style="width: 100%">
			<tbody>
			<tr>
				<td>
					<div class="cell">盘符</div>
				</td>
				<td>
					<div class="cell">总大小</div>
				</td>
				<td>
					<div class="cell">可用大小</div>
				</td>
				<td>
					<div class="cell">已用大小</div>
				</td>
				<td>
					<div class="cell">已用百分比</div>
				</td>
			</tr>
			<tr v-for="(sysFile, index) in sysInfo.diskList" :key="index">
				<td>
					<div class="cell">{{ sysFile.path }}</div>
				</td>
				<td>
					<div class="cell">{{ memorySizeFormat(sysFile.total) }}</div>
				</td>
				<td>
					<div class="cell">{{ memorySizeFormat(sysFile.free) }}</div>
				</td>
				<td>
					<div class="cell">{{ memorySizeFormat(sysFile.used) }}</div>
				</td>
				<td>
					<el-progress :text-inside="true" :stroke-width="20" :color="color" :percentage="sysFile.usedPercent" />
<!--					<div class="cell">{{ sysFile.usedPercent }}%</div>-->
				</td>
			</tr>
			</tbody>
		</table>
	</div>

	<div class="q-small-title">CPU / 内存</div>
	<el-row>
		<el-col :span="12">
			<div style="min-height: 200px; " ref="chartsWarningRef1"></div>
		</el-col>
		<el-col :span="12">
			<div style="min-height: 200px; " ref="chartsWarningRef2"></div>
		</el-col>
	</el-row>

</template>

<script lang="ts">
	import { defineComponent, reactive, onMounted, getCurrentInstance, toRefs } from 'vue';
	import * as echarts from 'echarts';
	import 'echarts-wordcloud';
	import { getSysInfo } from '/@/api/system/monitor/server';

	let interval: any = null;
	export default defineComponent({
		name: "personal-space",
		setup(){
			const { proxy } = <any>getCurrentInstance();
			const state = reactive({
				myCharts: [] as any[],
				sysInfo: {},
				color: [
					{ color: '#67C23A', percentage: 30 },
					{ color: '#E6A23C', percentage: 70 },
					{ color: '#F56C6C', percentage: 90 },
				],
			});

			let myChart1: any;
			let myChart2: any;
			function setOptChart1(value: number) {
				myChart1.setOption({
					series: [
						{
							data: [
								{
									value: value,
									name: 'CPU使用率',
								},
							],
						},
					],
				});
			}

			function setOptChart2(value: number) {
				myChart2.setOption({
					series: [
						{
							data: [
								{
									value: value,
									name: '内存使用率',
								},
							],
						},
					],
				});
			}

			//CPU
			const initChartCPU = () => {
				myChart1 = echarts.init(proxy.$refs.chartsWarningRef1);
				const option = {
					tooltip: {
						formatter: '{a} <br/>{b} : {c}%',
					},
					series: [
						{
							type: 'gauge',
							name: 'CPU',
							radius: '100%',
							title: {
								show: true,
								fontSize: 12,
								offsetCenter: [0, '95%'],
							},
							axisLine: {
								show: true,
								lineStyle: {
									// 属性lineStyle控制线条样式
									color: [
										[0.3, '#4dabf7'],
										[0.6, '#69db7c'],
										[0.8, '#ffa94d'],
										[1, '#ff6b6b'],
									],
								},
							},
							axisLabel: {
								color: 'inherit',
								distance: 15,
								fontSize: 10,
							},
							detail: {
								valueAnimation: true,
								formatter: '{value}%',
								textStyle: {
									fontSize: 24,
									color: 'red',
								},
								offsetCenter: ['0', '70%'],
							},
						},
					],
				};
				myChart1.setOption(option);
				state.myCharts.push(myChart1);
			};

			//内存
			const initChartRAM = () => {
				myChart2 = echarts.init(proxy.$refs.chartsWarningRef2);
				const option = {
					tooltip: {
						formatter: '{a} <br/>{b} : {c}%',
					},
					series: [
						{
							type: 'gauge',
							name: '内存',
							radius: '100%',
							title: {
								show: true,
								fontSize: 12,
								offsetCenter: [0, '95%'],
							},
							axisLine: {
								show: true,
								lineStyle: {
									// 属性lineStyle控制线条样式
									color: [
										[0.3, '#4dabf7'],
										[0.6, '#69db7c'],
										[0.8, '#ffa94d'],
										[1, '#ff6b6b'],
									],
								},
							},
							axisLabel: {
								color: 'inherit',
								distance: 15,
								fontSize: 10,
							},
							detail: {
								valueAnimation: true,
								formatter: '{value}%',
								textStyle: {
									fontSize: 24,
									color: 'red',
								},
								offsetCenter: ['0', '70%'], //表盘数据(30%)位置
							},
						},
					],
				};
				myChart2.setOption(option);
				state.myCharts.push(myChart2);
			};

			function getSystemInfo() {
				getSysInfo().then((res: any) => {
					const { code, data } = res;
					if (code === 0) {
						state.sysInfo = data;
						setOptChart1(data.cpuUsed);
						setOptChart2(data.memUsage);
					}
				});
			};

			// 页面加载时
			onMounted(() => {
				initChartCPU();
				initChartRAM();
			});

			return {
				...toRefs(state),
				getSystemInfo,
				setOptChart1,
				setOptChart2,
			};
		},
		created() {
			this.getSystemInfo();
			if (interval === null) {
				interval = setInterval(() => {
					this.getSystemInfo();
				}, 5000);
			}
		},
		unmounted() {
			if (interval !== null) {
				clearInterval(interval);
				interval = null;
			}
		},
		data() {
			return {};
		},
		methods: {
			memorySizeFormat(size: any) {
				size = parseFloat(size);
				let rank = 0;
				let rankchar = 'Bytes';
				while (size > 1024 && rankchar != 'TB') {
					size = size / 1024;
					rank++;
					if (rank == 1) {
						rankchar = 'KB';
					} else if (rank == 2) {
						rankchar = 'MB';
					} else if (rank == 3) {
						rankchar = 'GB';
					} else if (rank == 4) {
						rankchar = 'TB';
					}
				}
				return size.toFixed(2) + ' ' + rankchar;
			},
		},
	});

</script>

<style scoped lang="scss">
	.cell {
		box-sizing: border-box;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: normal;
		word-break: break-all;
		line-height: 36px;
		padding-left: 10px;
		padding-right: 10px;
	}
</style>

