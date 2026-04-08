<!--
 * @Descripttion: 封装echarts
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-07-10
 * @LastEditors:
 * @LastEditTime:
-->

<template>
	<div ref="echartRef" class="echarts-style" :style="{height:height, width:width}"></div>
	<object ref="objectRef"
			tabindex="-1"
			type="text/html"
			aria-hidden="true"
			data="about:blank"
			style="display: block;
            position: absolute;
            top: 0px;
            left: 0px;
            width: 100%;
            height: 100%;
            border: none;
            padding: 0px;
            margin: 0px;
            opacity: 0;
            z-index: -1000;
            pointer-events: none;">
	</object>
</template>

<script>
	import * as echarts from 'echarts';
	import T from './echarts-theme-T';
	echarts.registerTheme('T', T);
	const unwarp = (obj) => obj && (obj.__v_raw || obj.valueOf() || obj);

	export default {
		...echarts,
		name: "q-echarts",
		props: {
			height: { type: String, default: "100%" },
			width: { type: String, default: "100%" },
			nodata: {type: Boolean, default: false },
			option: { type: Object, default: () => {} }
		},
		data() {
			return {
				isActivat: false,
				myChart: null
			}
		},
		watch: {
			option: {
				deep:true,
				handler (v) {
					unwarp(this.myChart).setOption(v);
				}
			}
		},
		computed: {
			myOptions: function() {
				return this.option || {};
			}
		},
		activated(){
			if(!this.isActivat){
				this.$nextTick(() => {
					this.myChart.resize()
				})
			}
		},
		deactivated(){
			this.isActivat = false;
		},
		mounted(){
			this.isActivat = true;
			this.$nextTick(() => {
				this.draw();
			})
		},
		methods: {
			draw(){
				var myChart = echarts.init(this.$refs.echartRef, 'T');
				myChart.setOption(this.myOptions);
				this.myChart = myChart;

				this.$refs["objectRef"].contentDocument.defaultView.addEventListener("resize", () => {
					this.myChart.resize()
				});
				// window.addEventListener('resize', () => this.myChart.resize());
			},
		}
	}
</script>

<style scoped lang="scss">
	.echarts-style {
		padding: 5px;
	}
</style>