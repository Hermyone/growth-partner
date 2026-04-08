const modules : any = {}

/**
 *  将值放入名称的最后一级，多级目录每层为object
 * @param baseObj Object
 * @param names []string
 * @param value Object
 */
function setValue(baseObj: any, names: any, value: any) {
	let obj = {}
	if(names.length > 1) {
		if(baseObj.hasOwnProperty(names[0])) {
			obj = baseObj[names[0]]
		}else{
			baseObj[names[0]] = obj
		}

		setValue(obj, names.splice(1), value)
	}else if(names.length == 1){
		baseObj[names[0]] = value
	}
}

/**
 * 将ts文件解析到指定的对象中，按目录名放入
 * @param obj
 * @param files 文件列表 如：./model/items.ts
 * @param removePrefix 需移除的前辍，空是不移除
 */
function processFiles(obj: any, files: any, removePrefix: string) {
	for (let [key, value] of Object.entries(files)) {
		// @ts-ignore
		if(!value || !value.default){
			continue
		}

		// 去扩展名和路径前的./
		key = key.replace(/(\.\/|\.ts)/g, '')
		const names = key.split('/').filter((val) => val !== '' && (removePrefix === '' || val !== removePrefix))
		// @ts-ignore
		setValue(obj, names, value.default)
	}
}

const files : any = import.meta.globEager('./model/*.ts')
processFiles(modules, files, 'model')

export default modules
