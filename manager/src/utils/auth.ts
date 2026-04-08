import { useUserInfo } from '/@/stores/userInfo';
import { judementSameArr } from '/@/utils/arrayOperation';

/**
 * 单个权限验证
 * @param value 权限值
 * @returns 有权限，返回 `true`，反之则反
 */
export function auth(value: string): boolean {
	const stores = useUserInfo();
	return stores.userInfos.authBtnList.some((v: string) => v === value);
}

/**
 * 多个权限验证，满足一个则为 true
 * @param value 权限值
 * @returns 有权限，返回 `true`，反之则反
 */
export function auths(value: Array<string>): boolean {
	let flag = false;
	const stores = useUserInfo();
	stores.userInfos.authBtnList.map((val: string) => {
		value.map((v: string) => {
			if (val === v) flag = true;
		});
	});
	return flag;
}

/**
 * 多个权限验证，全部满足则为 true
 * @param value 权限值
 * @returns 有权限，返回 `true`，反之则反
 */
export function authAll(value: Array<string>): boolean {
	const stores = useUserInfo();
	return judementSameArr(value, stores.userInfos.authBtnList);
}

/**
 *
 * @param data 具有权限的路径列表
 * @param targetPath 需要判断的路径
 * @returns 有权限，返回 `true`，反之则反
 */
export function authPath(data: Array<any>, targetPath: string): boolean {
	// 定义一个递归函数来检查数组中的每个元素
	function checkItem(item: any): boolean {
		// 如果当前项是一个对象，并且包含path属性
		if (typeof item === 'object' && item !== null && 'path' in item) {
			// 检查path是否匹配
			if (item.path === targetPath) {
				return true;
			}
		}

		// 如果当前项有children属性且是一个数组
		if (Array.isArray(item.children)) {
			// 递归检查children数组
			for (const child of item.children) {
				if (checkItem(child)) {
					return true;
				}
			}
		}

		// 如果没有找到匹配的path，则返回false
		return false;
	}

	// 遍历data数组，并对每个元素调用checkItem
	for (const item of data) {
		if (checkItem(item)) {
			return true;
		}
	}

	// 如果没有找到匹配的path，则返回false
	return false;
}