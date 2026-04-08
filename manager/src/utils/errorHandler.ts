/* eslint-disable no-console */
/**
 * 全局代码错误捕捉
 */

export const errorMap : any = {
    InternalError: "Javascript引擎内部错误",
    ReferenceError: "未找到对象",
    TypeError: "使用了错误的类型或对象",
    RangeError: "使用内置对象时，参数超范围",
    SyntaxError: "语法错误",
    EvalError: "错误的使用了Eval",
    URIError: "URI错误"
}

export default {
	errorHandler(error: any, vm: any){
		//过滤HTTP请求错误
		if(error.status || error.status==0){
			return false
		}

		var errorName : any = errorMap[error.name] || "未知错误"

		console.warn(`[error]: ${error}`);
		console.error(error);

		vm.$nextTick(() => {
			vm.$notify.error({
				title: errorName,
				message: error
			});
		})
	},

	errorWin(message: any, source: string, lineno: number, colno: number, error: Error) {
		console.log('捕获到异常：',{message, source, lineno, colno, error});
	}
}