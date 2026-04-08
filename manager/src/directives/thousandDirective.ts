import type { App } from 'vue';
import { formatFloat } from '/@/utils/formatNumber';

/**
 * 格式化显示el-input, el-input-number
 * @directive 默认方式：v-thousands
 * @directive 参数方式：v-thousands="{fmt: ',0.00', focus:'center'}"
 * 参数fmt, 格式化参数，默认0.##
 *  0.00 最多保留两位小数据，不够时补0
 *	0.## 最多保留两位小数据
 *  ,0.00 强制保留两位小数据，不够时补0，千分位显示
 *  ,0.## 最多保留两位小数据，千分位显示
 *  0.00#### 最小保留两位小数，最多6位小数
 * 参数focus,获得焦点时光标位置 left|center|right, 默认center
 */

const dealThousands = (el: any, binding: any) => {
	const { fmt = '0.##' } =
		binding.value == null || typeof binding.value == 'string' ? { fmt: binding.value } : binding.value

	const value = el.value == '' ? 0 : String(el.value).replace(/,/g, '')
	el.type = 'text'
	el.value = formatFloat(value, fmt)
};

export function thousandDirective(app: App) {
	app.directive('thousands', {
		mounted(el, binding) {
			// 获取input 节点
			const suffixDom = el.querySelector('.el-input__suffix')
			if(el.tagName.toLocaleUpperCase() !== 'INPUT'){
				el = el.getElementsByTagName('input')[0]
			}
			// 千分位格式化
			setTimeout(() => {
				dealThousands(el, binding)
				el.classList.add('fmt-input-number', 'focus-' + (binding.value?.focus || 'center'))
				suffixDom && suffixDom.classList.add('m15')
			}, 8);
			// 获得焦点时
			el.onfocus = (e: any) => {
				const s = el.value.replace(/,/g, '')
				el.value = parseFloat(s) == 0 ? '' : parseFloat(s)
				el.type = 'number'
			}
			// 失去焦点
			el.onblur = (e: any) => {
				dealThousands(el, binding)
			}
		},
		updated(el, binding) {
			// 获取input 节点
			if(el.tagName.toLocaleUpperCase() !== 'INPUT') {
				el = el.getElementsByTagName('input')[0]
			}
			// 被动修改时重新格式化
			if(!el.isEqualNode(document.activeElement)) {
				dealThousands(el, binding)
			}
		},
	});
}
