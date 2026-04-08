/**
 * 格式化浮点数
 * @param v 需要需式化的数值
 * @param format 格式字符串
 * @description format 支持 0.00 0.## ,0.00等
 * @returns 返回格式化后的数值字符串
 */
export function formatFloat(v: any, format: string) : string {
	if(typeof(v) != 'number'){
		const x = Number(v)
		if(isNaN(x)){
			return v
		}
		v = x
	}

	if(format == '') return v.toString()

	let fv = v.toString().split(".")[1] || ''
	let fs = format.split(".")
	if(fs.length == 2) {
		let sIdx = fs[1].indexOf('#');
		let nIdx = fs[1].indexOf('0');
		if(sIdx == -1 && nIdx == -1) {
			v = v.toFixed(0)
		}else{
			if(sIdx == -1){
				v = v.toFixed(fs[1].length)
			}else{
				let len = sIdx
				if(fv.length > len){
					if(fv.length < fs[1].length){
						len = fv.length
					} else {
						len = fs[1].length
					}
				}
				v = v.toFixed(len)
			}
		}
	}

	let ds = fs[0].indexOf(',')
	if(ds == -1){
		return v.toString()
	}

	return v.toString().replace(/\d+/, function(n: string) {
		return n.replace(/(\d)(?=(\d{3})+$)/g, function($1) {
			return $1 + ","
		})
	})
}
