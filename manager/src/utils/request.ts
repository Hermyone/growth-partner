import axios, { AxiosInstance } from 'axios';
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus';
import { Session } from '/@/utils/storage';
import sysConfig from '/@/config';
// import qs from 'qs';

// 配置新建一个 axios 实例
const service : AxiosInstance = axios.create({
	// baseURL: sysConfig.API_URL as any,
	baseURL: '/',
	timeout: sysConfig.TIMEOUT,
	withCredentials: true,
	headers: { 'Content-Type': 'application/json' },
	// paramsSerializer: {
	// 	// 序列化
	// 	serialize(params) {
	// 		return qs.stringify(params, { allowDots: true });
	// 	},
	// },
});

// 添加请求拦截器
service.interceptors.request.use(
	(config) => {
		// 在发送请求之前做些什么 token
		let token = Session.get('token');
		if (token) {
			(<any>config.headers).common[sysConfig.TOKEN_NAME] = sysConfig.TOKEN_PREFIX + token;
		}
		return config;
	},
	(error) => {
		// 对请求错误做些什么
		return Promise.reject(error);
	}
);

// 添加响应拦截器
service.interceptors.response.use(
	(response) => {
		// 对响应数据做点什么
		const res = response.data;
		const code = response.data.code;
		if(code){
			if (code === 401) {
				ElMessageBox.alert('登录状态已过期，请重新登录', '提示', {confirmButtonText:'确认'})
					.then(() => {
						// 清除浏览器全部临时缓存
						Session.clear();
						// 去登录页
						window.location.href = '/';
					})
					.catch(() => {});
			} else if (code !== 0) {
				// ElMessage.error(res.message);
				if(res.message.includes("invalid")){
					console.log(res.message)
					Session.clear();
					window.location.href = '/';
				}else{
					ElNotification({
						title: '系统错误',
						message: res.message,
						type: 'error',
					});
				}
				return Promise.reject(new Error(res.message));
			} else {
				return res;
			}
		}else{
			return res;
		}
	},
	(error) => {
		// 对响应错误做点什么
		if (error.message.indexOf('timeout') != -1) {
			ElMessage.error('网络超时');
		} else if (error.message == 'Network Error') {
			ElMessage.error('网络连接错误');
		} else {
			if (error.response.data) ElMessage.error(error.response.statusText);
			else ElMessage.error('接口路径找不到');
		}
		return Promise.reject(error);
	}
);

export const http = {

	/** get 请求
	 * @param  {string} url 接口地址
	 * @param  {object} params 请求参数
	 * @param  {object} config 参数
	 */
	get: function(url: string, params={}, config={}) {
		return new Promise((resolve, reject) => {
			service({
				method: 'get',
				url: url,
				params: params,
				...config
			}).then((response) => {
				resolve(response);
			}).catch((error) => {
				reject(error);
			})
		})
	},

	/** post 请求
	 * @param  {string} url 接口地址
	 * @param  {object} data 请求参数
	 * @param  {object} config 参数
	 */
	post: function(url: string, data={}, config={}) {
		return new Promise((resolve, reject) => {
			service({
				method: 'post',
				url: url,
				data: data,
				...config
			}).then((response) => {
				resolve(response);
			}).catch((error) => {
				reject(error);
			})
		})
	},

	/** put 请求
	 * @param  {string} url 接口地址
	 * @param  {object} data 请求参数
	 * @param  {object} config 参数
	 */
	put: function(url: string, data={}, config={}) {
		return new Promise((resolve, reject) => {
			service({
				method: 'put',
				url: url,
				data: data,
				...config
			}).then((response) => {
				resolve(response);
			}).catch((error) => {
				reject(error);
			})
		})
	},

	/** patch 请求
	 * @param  {string} url 接口地址
	 * @param  {object} data 请求参数
	 * @param  {object} config 参数
	 */
	patch: function(url: string, data={}, config={}) {
		return new Promise((resolve, reject) => {
			service({
				method: 'patch',
				url: url,
				data: data,
				...config
			}).then((response) => {
				resolve(response);
			}).catch((error) => {
				reject(error);
			})
		})
	},

	/** delete 请求
	 * @param  {string} url 接口地址
	 * @param  {object} data 请求参数
	 * @param  {object} config 参数
	 */
	delete: function(url: string, data={}, config={}) {
		return new Promise((resolve, reject) => {
			service({
				method: 'delete',
				url: url,
				data: data,
				...config
			}).then((response) => {
				resolve(response);
			}).catch((error) => {
				reject(error);
			})
		})
	},

	/** jsonp 请求
	 * @param  {string} url 接口地址
	 * @param  {string} name JSONP回调函数名称
	 */
	jsonp: function(url: string, name ='jsonp'){
		return new Promise((resolve) => {
			var script = document.createElement('script')
			var _id = `jsonp${Math.ceil(Math.random() * 1000000)}`
			script.id = _id
			script.type = 'text/javascript'
			script.src = url
			// @ts-ignore
			window[name] =(response: any) => {
				resolve(response)
				document.getElementsByTagName('head')[0].removeChild(script)
				try {
					// @ts-ignore
					delete window[name];
				}catch(e){
					// @ts-ignore
					window[name] = undefined;
				}
			}
			document.getElementsByTagName('head')[0].appendChild(script)
		})
	}
};

// 导出 axios 实例
export default service;
