export const uploadConfig = {
	filename: "file",								//form请求时文件的key
	successCode: 200,								//请求完成代码
	maxSize: 10,										//最大文件大小 默认10MB
	parseData: function (res:any) { //可根据服务修改该解析方法
		// return {
		// 	code: res.code,							//分析状态字段结构
		// 	fileName: res.data.fileName,//分析文件名称
		// 	src: res.data.src,					//分析图片远程地址结构
		// 	msg: res.message						//分析描述字段结构
		// }

		return {
			code: res.state == '0' ? 200 : 0,
			fileName: res.filename,
			src: res.url,
			msg: res.msg
		}
	}
}

// 格式化数组值
export const formatArr = (arr: any) => {
	var _arr : any = [];
	arr.forEach((item:any) => {
		if(item){
			_arr.push({
				name: item.name,
				url: item.url
			})
		}
	})
	return _arr;
};

// 数组转换为原始值
export const toStr = (arr:any) => {
	return arr.map((v:any) => v.url).join(",")
};

// 默认值转换为数组
export const toArr = (str: string) => {
	var _arr : any = [];
	var arr = str.split(",")
	arr.forEach(item => {
		if(item){
			var urlArr = item.split('/');
			var fileName = urlArr[urlArr.length - 1]
			_arr.push({
				name: fileName,
				url: item
			})
		}
	})
	return _arr;
};