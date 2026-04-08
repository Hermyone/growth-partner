
const WeatherRes = {
	// 天气图标
	Icon : [
		{ name:'晴', 			value:'0', 		icon:'q-icon-weather-clear' },
		{ name:'阴', 			value:'1', 		icon:'q-icon-weather-cloudy' },
		{ name:'多云', 		value:'2', 		icon:'q-icon-weather-cloudy-ex' },

		{ name:'阵雨', 		value:'3', 		icon:'q-icon-weather-shower' },
		{ name:'雷阵雨', 	value:'4', 		icon:'q-icon-weather-thunderstorm' },
		{ name:'小雨', 		value:'5', 		icon:'q-icon-weather-light-rain' },
		{ name:'中雨', 		value:'6', 		icon:'q-icon-weather-moderate-rain' },
		{ name:'大雨', 		value:'7', 		icon:'q-icon-weather-heavy-rain' },
		{ name:'暴雨', 		value:'8', 		icon:'q-icon-weather-heavy-rain' },
		{ name:'大暴雨', 	value:'9', 		icon:'q-icon-weather-heavy-rainstorm' },
		{ name:'特大暴雨', value:'10', 	icon:'q-icon-weather-extremely-heavy-rainstorm' },

		{ name:'雨夹雪', 	value:'11', 	icon:'q-icon-weather-sleet' },
		{ name:'阵雪', 		value:'12', 	icon:'q-icon-weather-snow-shower' },
		{ name:'小雪',			value:'13', 	icon:'q-icon-weather-spit' },
		{ name:'中雪', 		value:'14', 	icon:'q-icon-weather-snow' },
		{ name:'大雪', 		value:'15', 	icon:'q-icon-weather-major-snow' },
		{ name:'暴雪', 		value:'16', 	icon:'q-icon-weather-blizzard' },
		{ name:'冻雪', 		value:'17', 	icon:'q-icon-weather-frozen-snow' },

		{ name:'雾', 			value:'18', 	icon:'q-icon-weather-fog' },
		{ name:'霾', 			value:'19', 	icon:'q-icon-weather-haze' },
		{ name:'沙尘暴', 	value:'20', 	icon:'q-icon-weather-sandstorm' },
	],

	// Tag 颜色
	Tag : [
		{ name: '优',			value: '1',		color: '#02d403', AQI: '0-50', 		memo: '可以参加户外活动呼吸清新空气' },
		{ name: '良',			value: '2',		color: '#eded0b', AQI: '50-100', 	memo: '可以正常进行室外活动' },
		{ name: '轻度',		value: '3',		color: '#ff7e01', AQI: '101-150', memo: '敏感人群减少体力消耗大的户外活动' },
		{ name: '中度',		value: '4',		color: '#ff0100', AQI: '151-200', memo: '对敏感人群影响较大' },
		{ name: '重度',		value: '5',		color: '#7f0081', AQI: '201-300', memo: '所有人应适当减少室外活动' },
		{ name: '严重',		value: '6',		color: '#7d0000', AQI: '>300', 		memo: '尽量不要留在室外' },
	],
};

export default WeatherRes;