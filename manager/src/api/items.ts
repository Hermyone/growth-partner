import { http } from '/@/utils/request';

export function getItems(f: Function, query: Object) {
	query = query || { pageSize: 10000 };
	return f(query);
}

export function setItems(response: any, k: string, v: string):Array<ItemOptions> {
	const data: Array<ItemOptions> = [];
	k = k || 'id';
	v = v || 'name';
	if (response.data && response.data.list && response.data.list.length > 0) {
		response.data.list.forEach((e: any) => {
			data.push({
				key: e[k].toString(),
				value: e[v].toString(),
			});
		});
	}
	return data;
}

// 选项类型接口
export interface ItemOptions {
	key:string,
	value:string
}

/** 通过 options 数组获取 key 对应的 value */
export function getOptionValue(key: any, options: Array<any>,keyName:string,valName:string) {
	keyName = keyName??'key'
	valName = valName??'value'
	const option = options.find((value) => {
		return key + '' === value[keyName];
	});
	if (option !== undefined) {
		return option[valName];
	}
}

export function isEmpty(v: any) {
	if (v === '' || v === undefined || v === null) {
		return true;
	}

	if (typeof v === 'object') {
		if (Array.isArray(v)) {
			return v.length === 0;
		} else {
			return Object.keys(v).length === 0;
		}
	}
	return false;
}

// 将申请指令转成对象
export const requestToObject = (pathArr : any[], API_URL : string = '') => {
	let obj : any = {}
	pathArr.forEach((list:any) => {
		let [item, method = 'post' ] = list.split(':')
		obj[item] = {
			url: `${API_URL}/${item}`
		}
		obj[item][method] = async function (data = {}) {
			// @ts-ignore
			return await http[method](obj[item].url, data)
		}
	})
	return obj
}

export const requestToObject_ex = (requestList : any[], API_URL : string = '') => {
	let obj : any = {}
	requestList.forEach(({key, path, method, name, isFullPath}) => {
		// 默认根据最后一级path作为请求KEY
		key = key || path.replace(/.*\/(.+)$/, '$1')
		obj[key] = {
			url: isFullPath ? path : `${API_URL}${path}`,
			name
		}
		// 可以同时设置post和get 如："post,get"
		for(let md of method.replace(/(\s*[,， ]\s*)/g, ',').split(',')){
			if(!md) continue
			obj[key][md] = async function (data= {}) {
				// @ts-ignore
				return await http[md](obj[key].url, data)
			}
		}
	})
	return obj
}

export const requestToObject_m = (requestConfig: any, API_URL : string = '') => {
	if(Array.isArray(requestConfig)) {
		return requestToObject_ex(requestConfig, API_URL)
	}else{
		let obj : any ={}
		Object.keys(requestConfig).forEach((key:string) => {
			obj[key] = requestToObject_m(requestConfig[key], API_URL)
		})
		return obj
	}
}