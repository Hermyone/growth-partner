const CONFIG = {
	//标题
	APP_NAME: import.meta.env.VITE_APP_TITLE,

	//版本号
	APP_VER: "1.0.0",

	//内核版本号
	CORE_VER: "1.0.0",

	//API接口地址
	API_URL: import.meta.env.VITE_API_URL + ':' + import.meta.env.VITE_API_PORT + "/", // process.env.NODE_ENV === 'development' && process.env.VITE_APP_PROXY === 'true' ? "/api" : process.env.VITE_API_URL,

	VITE_URL: import.meta.env.VITE_API_URL + ':' + import.meta.env.VITE_PORT + "/",

	FILE_URL: import.meta.env.VITE_API_URL + ':' + import.meta.env.VITE_API_PORT + "/file/upload",

	//请求超时
	TIMEOUT: 20000,

	//TokenName
	TOKEN_NAME: "Authorization",

	//Token前缀，注意最后有个空格，如不需要需设置空字符串
	TOKEN_PREFIX: "Bearer ",

	//默认每个显示记录数
	DEF_PAGE_SIZE: 20,
}

function getAddress(str: string): string | null {
	const regex = /(https?:\/\/)?([\d.]+|localhost)/;
	const match = str.match(regex);
	return match ? match[2] : null;
}

export default CONFIG;