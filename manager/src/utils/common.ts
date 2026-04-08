import { nextTick, defineAsyncComponent } from 'vue';
import type { App } from 'vue';
import * as svg from '@element-plus/icons-vue';
import router from '/@/router/index';
import pinia from '/@/stores/index';
import { storeToRefs } from 'pinia';
import { useThemeConfig } from '/@/stores/themeConfig';
import { i18n } from '/@/i18n/index';
import { Local } from '/@/utils/storage';
import { verifyUrl } from '/@/utils/toolsValidate';
import {Session} from "/@/utils/storage";
import sysConfig from "/@/config";
import { EmptyArrayType, EmptyObjectType, RouteItem } from '/@/types/global';

// 引入组件
const SvgIcon = defineAsyncComponent(() => import('/@/components/svgIcon/index.vue'));

// 内容显示是否使用卡片模式
export function getCardMode() : boolean {
	const stores = useThemeConfig(pinia);
	const { themeConfig } = storeToRefs(stores);
	let { layout } = themeConfig.value;
	return layout !== "columns";
}

export function getUpFileUrl(url:string){
    if(!url){
        return url
    }
    if (/^http|^blob/i.test(url)) {
        return url
    }
    let reg = new RegExp('^/*' + sysConfig.API_URL + "/*");
    return sysConfig.API_URL+url.replace(reg,'')
}

/**
 * 构造树型结构数据
 * @param {*} data 数据源
 * @param {*} id id字段 默认 'id'
 * @param {*} parentId 父节点字段 默认 'parentId'
 * @param {*} children 孩子节点字段 默认 'children'
 * @param {*} rootId 根Id 默认 0
 */
export function handleTree(data:any[], id:string, parentId:string, children:string, rootId:number):any[] {
    id = id || 'id'
    parentId = parentId || 'parentId'
    children = children || 'children'
    rootId = rootId || 0
    //对源数据深度克隆
    const cloneData = JSON.parse(JSON.stringify(data))
    //循环所有项
    const treeData =  cloneData.filter((father:any) => {
        let branchArr = cloneData.filter((child:any) => {
            //返回每一项的子级数组
            return father[id] === child[parentId]
        });
        branchArr.length > 0 ? father[children] = branchArr : '';
        //返回第一层
        return father[parentId] === rootId;
    });
    return treeData != '' ? treeData : data;
}

export function convertToTree(data: any[], id:string, parentId:string, rootId:string): any {
	id = id || 'id'
	parentId = parentId || 'parent'
	rootId = rootId || "0"

	// 创建一个哈希映射，以便快速查找父节点
	const map: Record<string, any> = {};
	data.forEach(item => {
		map[item[id]] = { ...item, children: [] };
	});

	// 构建树形结构
	const tree: any = [];
	data.forEach(item => {
		if (item[parentId] === rootId) {
			// 父节点为 "0" 的作为树的根节点
			tree.push(map[item[id]]);
		} else {
			// 找到父节点，并将当前项添加到其 children 中
			map[item[parentId]].children.push(map[item[id]]);
		}
	});

	return tree;
}

// 回显数据字典
export function selectDictLabel(data:any[], value:string):string {
	if(data == undefined){
		return value
	}

    let actions:string[]=[]
    data.map((item) => {
        if (item.value == value) {
            actions.push(item.label);
            return false;
        }
    })
    return actions.join('');
}

export function selectDictColor(data:any[], value:string):string {
	if(data == undefined){
		return value
	}
	let actions:string[]=[]
	data.map((item) => {
		if (item.value == value) {
			actions.push(item.color);
			return false;
		}
	})
	return actions.join('');
}

export function getToken():string{
    return Session.get('token')
}

// 日期格式化
export function parseTime(time:any, pattern:string) {
    if (arguments.length === 0 || !time) {
        return null
    }
    const format = pattern || '{y}-{m}-{d} {h}:{i}:{s}'
    let date
    if (typeof time === 'object') {
        date = time
    } else {
        if ((typeof time === 'string') && (/^[0-9]+$/.test(time))) {
            time = parseInt(time)
        } else if (typeof time === 'string') {
            time = time.replace(new RegExp(/-/gm), '/');
        }
        if ((typeof time === 'number') && (time.toString().length === 10)) {
            time = time * 1000
        }
        date = new Date(time)
    }
    const formatObj:any = {
        y: date.getFullYear(),
        m: date.getMonth() + 1,
        d: date.getDate(),
        h: date.getHours(),
        i: date.getMinutes(),
        s: date.getSeconds(),
        a: date.getDay()
    }
    const time_str = format.replace(/{(y|m|d|h|i|s|a)+}/g, (result, key) => {
        let value = formatObj[key]
        // Note: getDay() returns 0 on Sunday
        if (key === 'a') { return ['日', '一', '二', '三', '四', '五', '六'][value] }
        if (result.length > 0 && value < 10) {
            value = '0' + value
        }
        return value || 0
    })
    return time_str
}

// 返回当前时间
export function getCurrentDatetime(fmt: string = '{y}-{m}-{d} {h}:{i}:{s}') {
	return parseTime(new Date, fmt);
}

export function getCurrentDate(fmt : string = '{y}-{m}-{d}') {
	return parseTime(new Date, fmt);
}

export function getCurrentTime(fmt : string = '{h}:{i}:{s}') {
	return parseTime(new Date, fmt);
}

/**
 * 导出全局注册 element plus svg 图标
 * @param app vue 实例
 * @description 使用：https://element-plus.gitee.io/zh-CN/component/icon.html
 */
export function elSvg(app: App) {
	const icons = svg as any;
	for (const i in icons) {
		let key1 = `ele-${icons[i].name}`;
		app.component(key1, icons[i]);

		// 为兼容习惯，注册名称为el-icon-XX的图标名称
		let key2 = `el-icon${icons[i].name}`;
		key2 =  key2.replace(/[A-Z]/g, (match:string)=>'-'+match.toLowerCase());
		app.component(key2.toLowerCase(), icons[i]);
	}
	app.component('SvgIcon', SvgIcon);
}

/**
 * 设置浏览器标题国际化
 * @method const title = useTitle(); ==> title()
 */
export function useTitle() {
	const stores = useThemeConfig(pinia);
	const { themeConfig } = storeToRefs(stores);
	nextTick(() => {
		let webTitle = '';
		let globalTitle: string = themeConfig.value.globalTitle;
		const { path, meta } = router.currentRoute.value;
		if (path === '/login') {
			webTitle = <any>meta.title;
		} else {
			webTitle = setTagsViewNameI18n(router.currentRoute.value);
		}
		document.title = `${webTitle} - ${globalTitle}` || globalTitle;
	});
}

/**
 * 设置 自定义 tagsView 名称、 自定义 tagsView 名称国际化
 * @param params 路由 query、params 中的 tagsViewName
 * @returns 返回当前 tagsViewName 名称
 */
export function setTagsViewNameI18n(item: any) {
	let tagsViewName: any = '';
	const { query, params, meta } = item;
	if (query?.tagsViewName || params?.tagsViewName) {
		if (/\/zh-cn|en|zh-tw\//.test(query?.tagsViewName) || /\/(zh-cn|en|zh-tw)\//.test(params?.tagsViewName)) {
			// 国际化
			const urlTagsParams = (query?.tagsViewName && JSON.parse(query?.tagsViewName)) || (params?.tagsViewName && JSON.parse(params?.tagsViewName));
			tagsViewName = urlTagsParams[i18n.global.locale.value];
		} else {
			// 非国际化
			tagsViewName = query?.tagsViewName || params?.tagsViewName;
		}
	} else {
		// 非自定义 tagsView 名称
		tagsViewName = i18n.global.t(<any>meta.title);
	}
	return tagsViewName;
}

/**
 * 图片懒加载
 * @param el dom 目标元素
 * @param arr 列表数据
 * @description data-xxx 属性用于存储页面或应用程序的私有自定义数据
 */
export const lazyImg = (el: any, arr: any) => {
	const io = new IntersectionObserver((res) => {
		res.forEach((v: any) => {
			if (v.isIntersecting) {
				const { img, key } = v.target.dataset;
				v.target.src = img;
				v.target.onload = () => {
					io.unobserve(v.target);
					arr[key]['loading'] = false;
				};
			}
		});
	});
	nextTick(() => {
		document.querySelectorAll(el).forEach((img) => io.observe(img));
	});
};

/**
 * 全局组件大小
 * @returns 返回 `window.localStorage` 中读取的缓存值 `globalComponentSize`
 */
export const globalComponentSize = (): string => {
	const stores = useThemeConfig(pinia);
	const { themeConfig } = storeToRefs(stores);
	return Local.get('themeConfig')?.globalComponentSize || themeConfig.value?.globalComponentSize;
};

/**
 * 对象深克隆
 * @param obj 源对象
 * @returns 克隆后的对象
 */
export function deepClone(obj: any) {
	let newObj: any;
	try {
		newObj = obj.push ? [] : {};
	} catch (error) {
		newObj = {};
	}
	for (let attr in obj) {
		if (obj[attr] && typeof obj[attr] === 'object') {
			newObj[attr] = deepClone(obj[attr]);
		} else {
			newObj[attr] = obj[attr];
		}
	}
	return newObj;
}

/**
 * 判断是否是移动端
 */
export function isMobile() {
	if (
		navigator.userAgent.match(
			/('phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone')/i
		)
	) {
		return true;
	} else {
		return false;
	}
}

/**
 * 判断数组对象中所有属性是否为空，为空则删除当前行对象
 * @description @感谢大黄
 * @param list 数组对象
 * @returns 删除空值后的数组对象
 */
export function handleEmpty(list: any) {
	const arr = [];
	for (const i in list) {
		const d = [];
		for (const j in list[i]) {
			d.push(list[i][j]);
		}
		const leng = d.filter((item) => item === '').length;
		if (leng !== d.length) {
			arr.push(list[i]);
		}
	}
	return arr;
}

/**
 * 打开外部链接
 * @param val 当前点击项菜单
 */
export function handleOpenLink(val: RouteItem) {
	const { origin, pathname } = window.location;
	router.push(val.path);
	if (verifyUrl(<string>val.meta?.isLink)) window.open(val.meta?.isLink);
	else window.open(`${origin}${pathname}#${val.meta?.isLink}`);
}

/**
 * 将字节数转带单位的辽符串
 * @param bytes
 */
export function convertBytesToUnit(bytes: number) {
	if (bytes === 0) {
		return "0";
	}

	const k = 1024;
	const units = ["B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

	let i = Math.floor(Math.log(bytes) / Math.log(k));

	if (i === 0) {
		return bytes + " B";
	} else if (i === 1) {
		return (bytes / k).toFixed(2) + " KB";
	} else if (i === 2) {
		return (bytes / (k * k)).toFixed(2) + " MB";
	} else {
		// 对于大于MB的情况，这里只是简单示例展示到MB以上，可根据实际需求继续完善
		return (bytes / (k * k)).toFixed(2) + " MB";
	}
}

/**
 * 从列表中根据关键字找到指定字段的值
 * @param source 要查找的源数据列表
 * @param ifField 需要判断的条件字段名
 * @param ifValue 需要判断的条件字段值
 * @param vField 需要返回的数据字段名
 */
export function findValueFromAarray(source: any, ifField: string, ifValue:string, vField: string) {
	if(!source)
		return [];

	let item : any= source.find((v:any) => v[ifField] === ifValue);
	if(item) return item[vField];
}

/**
 * 统一批量导出
 * @method elSvg 导出全局注册 element plus svg 图标
 * @method useTitle 设置浏览器标题国际化
 * @method setTagsViewNameI18n 设置 自定义 tagsView 名称、 自定义 tagsView 名称国际化
 * @method lazyImg 图片懒加载
 * @method globalComponentSize() element plus 全局组件大小
 * @method deepClone 对象深克隆
 * @method isMobile 判断是否是移动端
 * @method handleEmpty 判断数组对象中所有属性是否为空，为空则删除当前行对象
 */
const common = {
	elSvg: (app: App) => {
		elSvg(app);
	},
	useTitle: () => {
		useTitle();
	},
	setTagsViewNameI18n(route: any) {
		return setTagsViewNameI18n(route);
	},
	lazyImg: (el: any, arr: EmptyArrayType) => {
		lazyImg(el, arr);
	},
	globalComponentSize: () => {
		return globalComponentSize();
	},
	deepClone: (obj: EmptyObjectType) => {
		return deepClone(obj);
	},
	isMobile: () => {
		return isMobile();
	},
	handleEmpty: (list: EmptyArrayType) => {
		return handleEmpty(list);
	},
	handleOpenLink: (val: RouteItem) => {
		handleOpenLink(val);
	},
};

// 统一批量导出
export default common;
