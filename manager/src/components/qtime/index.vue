<!--
 * @Descripttion: 时间显示组件，支持按时段问候语
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-07-05
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <div class="item-background">
        <h2 style="height: 45px; padding-bottom: 15px">{{ title }}</h2>
        <div class="time">
            <h2>{{ time }}</h2>
            <p>{{ day }}</p>
        </div>
        <slot></slot>
    </div>
</template>

<script lang="ts">
	import { defineComponent, onMounted, reactive, toRefs } from 'vue';

	const props = {
        dateFmt: { type: String, default: "YYYY 年 MM 月 DD 日" },
        timeFmt: { type: String, default: "HH:mm:ss" },
	};

	export default defineComponent({
		name: 'q-time',
		props: props,

		setup(props) {

            const weeks = ['星期一','星期二','星期三','星期四','星期五','星期六','星期日']
			const state = reactive({
                title: '',
				time: '',
                day: '',
			});

			const dateFormat = (timestamp: number|string|Date, format: string, week?: boolean) => {
                var date = new Date(timestamp)
                function fixedTwo (value: number): string {
                    return value < 10 ? '0' + value : String(value)
                }
                var res = format
                if (res.includes('SSS')) {
                    const S = date.getMilliseconds()
                    res = res.replace('SSS', '0'.repeat(3 - String(S).length) + S)
                }
                if (res.includes('YY')) {
                    const Y = date.getFullYear()
                    res = res.includes('YYYY') ? res.replace('YYYY', String(Y)) : res.replace('YY', String(Y).slice(2, 4))
                }
                if (res.includes('M')) {
                    const M = date.getMonth() + 1
                    res = res.includes('MM') ? res.replace('MM', fixedTwo(M)) : res.replace('M', String(M))
                }
                if (res.includes('D')) {
                    const D = date.getDate()
                    res = res.includes('DD') ? res.replace('DD', fixedTwo(D)) : res.replace('D', String(D))
                }
                if (res.includes('H')) {
                    const H = date.getHours()
                    res = res.includes('HH') ? res.replace('HH', fixedTwo(H)) : res.replace('H', String(H))
                }
                if (res.includes('m')) {
                    var m = date.getMinutes()
                    res = res.includes('mm') ? res.replace('mm', fixedTwo(m)) : res.replace('m', String(m))
                }
                if (res.includes('s')) {
                    var s = date.getSeconds()
                    res = res.includes('ss') ? res.replace('ss', fixedTwo(s)) : res.replace('s', String(s))
                }

                if(week)
                    res = res + '  ' + weeks[date.getDay()-1] ;

                return res;
            };

			const showTime = () => {
                state.time = dateFormat(new Date(), props.timeFmt);
                state.day = dateFormat(new Date(), props.dateFmt, true);
            };

            // 页面加载时
            onMounted(() => {
            	showTime();
                setInterval(()=>{
                    showTime();
                    let hour: number = new Date().getHours();
                    if (hour < 6) state.title = '凌晨好';
                    else if (hour < 9) state.title = '早上好';
                    else if (hour < 12) state.title = '上午好';
                    else if (hour < 14) state.title = '中午好';
                    else if (hour < 17) state.title = '下午好';
                    else if (hour < 19) state.title = '傍晚好';
                    else if (hour < 22) state.title = '晚上好';
                    else state.title = '夜里好';
                },1000);
            });

			return {
                showTime,
				...toRefs(state),
			}
		},
	});
</script>

<style scoped lang="scss">
    .item-background {
        padding: 15px;
        background: linear-gradient(to right, var(--el-color-primary), var(--el-color-primary-light-4));color: #fff;
    }

    .time h2 {font-size: 40px;}
    .time p {font-size: 14px;margin-top: 13px;opacity: 0.7;}
</style>
