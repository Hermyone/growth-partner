<!--
 * @Descripttion: 天气组件
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-07-19
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <div class="weather-style" :style="isSmall() ? 'width: 100%; height: 60px':'width: 100%; height: 100%'">
        <el-row class="weather-style-top" v-if="!isSmall()">
            <el-col :span="24">
                {{ city }}
            </el-col>
        </el-row>

        <el-row class="weather-style-content">
            <el-col :span="isSmall()?4:5">
                <SvgIcon :name="weatherIcon" class="weather-style-icon" :size="isSmall()?24:50" :color="iconColor" />
            </el-col>
            <el-col :span="isSmall()?16:13">
                <div :class="isSmall()?'weather-style-content-mainsmall':'weather-style-content-main'" :style="{color: textColor}">
                    <span>{{ weather }}</span>

                    <span class="weather-style-content-tag" :style="{background: tagColor}">{{ airquality }}</span>
                </div>
                <div class="weather-style-content-sub" v-if="!isSmall()">{{ winddirection }}</div>
            </el-col>
            <el-col :span="isSmall()?4:6" :class="isSmall()?'weather-style-rightsmall':'weather-style-right'" :style="{color: textColor}">
                <span>{{ temperature }}</span>
                <span :class="isSmall()?'weather-style-rightsmall-font':'weather-style-right-font'">°C</span>
            </el-col>
        </el-row>
    </div>
</template>

<script lang="ts">
	import { defineComponent, nextTick, onMounted, computed, reactive, toRefs, watch } from 'vue';
	import WeatherRes from './weather';
    import axios from 'axios';

	const props = {
		// 城市
        city: { type: String, default: '' },
        // URL 获取天气的URL
        url: { type: String, default: 'https://restapi.amap.com/v3/weather/weatherInfo?key=5d2d3e6c0d5188bec134fc4fc1b139e0' },
	    // 大小 small/large
        size: { type: String, default: 'large' },
        // 图标颜色
        iconColor: { type: String, default: null },
        // 字体颜色
        textColor: { type: String, default: null },
    };

	export default defineComponent({
		name: 'q-weather',
		props: props,

		setup(props, { emit }) {
            const state = reactive({
                city: props.city,
                // 天气
                weather: '暂无',
                weatherIcon: 'q-icon-weather-clear',
                // 温度 °C
                temperature: '0',
                // 空气质量
                airquality: '暂无',
                // 风向
                winddirection: '暂无',
                // tag 颜色
                tagColor: '#00e401',
            });

            onMounted(() => {
            	nextTick(() => {
                    updateWeather();
                });

                // 定时更新
                setInterval(()=>{
                    updateWeather();
                },3600000);
            });

            // 更新天气
            const updateWeather = () => {
                new Promise(() => {
                    axios({
                        method: 'get',
                        url: props.url,
                        params: {'city': props.city},
                    }).then((response) => {
                    	let js = response.data;
                    	if(js.status !== '0' && js.count > 0){
                    		let item = js.lives[0];

                            state.weather = item.weather;
                            state.temperature = item.temperature;
                            state.airquality = '轻度';
                            state.winddirection = item.winddirection + item.windpower;

                            updateIconName();
                            updateTagColor();
                        }else{
                            emit('error', js.info);
                        }
                    }).catch((error) => {
                        console.log("天气更新失败", error);
                        emit('error', error);
                    })
                })
            };

            const updateIconName = () => {
                WeatherRes.Icon.forEach((item : any) => {
                    if(item.name === state.weather){
                    	state.weatherIcon = item.icon;
                    	return;
                    }
                });
            };

            const updateTagColor = () => {
                WeatherRes.Tag.forEach((item : any) => {
                    if(item.name === state.airquality){
                        state.tagColor = item.color;
                        return;
                    }
                });
            }

            const updateValue = (v : any) => {
                state.city = v['city'];
                state.weather = v['weather'];
                state.temperature = v['temperature'];
                state.airquality = v['airquality'];
                state.winddirection = v['winddirection'];

                updateIconName();
                updateTagColor();
            };

            const isSmall = () => {
                return props.size == 'small';
            };

            watch(
            () => props.city,
            ()=>{
            	    state.city = props.city;
                    updateWeather();
                },
            );

            watch(
                () => props.url,
                ()=>{
                    updateWeather();
                },
            );

            return {
                updateWeather,
                updateValue,
                isSmall,
                ...toRefs(state),
            }
        },
	});
</script>

<style scoped lang="scss">
    .weather-style {
        padding: 5px 10px;
        &-top {
            height: 15%;
            align-items: center;
            text-align: right;
        }

        &-icon {
            padding-top: 5px;
        }
        &-content {
            height: 85%;
            align-items: center;

            &-main {
                font-size: 18px;
                font-weight: bold;
                color: var(--el-text-color-primary);
                span {
                    padding: 0 5px;
                }
            }
            &-mainsmall {
                font-size: 16px;
                color: var(--el-text-color-primary);
                span {
                    padding: 0 5px;
                }
            }
            &-tag {
                font-size: 12px;
                border-radius: 2px;
                color: white; // var(--el-color-white);
            }
            &-sub {
                font-size: 12px;
                padding: 10px 5px;
                color: var(--el-text-color-disabled);
            }
        }

        &-right {
            font-size: 40px;
            color: var(--el-text-color-primary);
            &-font {
                position:relative;
                top: -20px;
                left: 5px;
                font-size: 14px;
            }
        }
        &-rightsmall {
            font-size: 18px;
            color: var(--el-text-color-primary);
            padding-top: 5px;
        }
    }
</style>
