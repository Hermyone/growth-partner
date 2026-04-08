import { nextTick } from 'vue'

// 渲染el-table树形样式，渲染完成/展开/收起时需要执行一次
export const useRenderTableTree = () => {
	// 参数为el-table的dom元素，即tableRef.value?.$el
	const renderTableTree = (tableDom: any) => {
		nextTick(() => {
			let rowDoms = tableDom.querySelectorAll('[class*=el-table__row--level]')
			let lvHeightArr : Array<number> = []
			rowDoms.forEach((rowDom: any, index: number) => {
				let rowHeight = rowDom.offsetHeight
				if(!rowHeight) return

				// 定位可展开的列
				let colDom = rowDom.querySelector(
					'.el-table__cell .el-table__expand-icon, .el-table__cell .el-table__indent'
				)
				while (colDom && !colDom.classList.contains('el-table__cell')) {
					colDom = colDom.parentNode
				}
				colDom?.classList.add('tree-column')
				// 计算高度
				// @ts-ignore
				const curRowHalf = parseInt(rowHeight / 2)
				let level = [...rowDom.classList].find((i) => i.indexOf('el-table__row--level-') == 0)
				level = level ? Number(level.split('el-table__row--level-')[1]) : 0
				// 顶级不加lineDom
				if(level == 0){
					lvHeightArr = [0, curRowHalf - 7]
					return
				}

				for(let i=1; i<level; i++){
					lvHeightArr[i] += rowHeight
				}
				// 给当前行增加lineDom，并重置当前级累计
				let indentDom = colDom.querySelector('.el-table__indent')
				if(indentDom != null) {
					let lineDom = indentDom.querySelector('.expand_line')
					if(lineDom == null){
						lineDom = document.createElement('span')
						lineDom.classList.add('expand_line')
						indentDom.appendChild(lineDom)
					}
					lineDom.style.height = (lvHeightArr[level] || 0) + curRowHalf + 'px'
				}
				lvHeightArr[level] = curRowHalf
				lvHeightArr[level + 1] = curRowHalf - 7
			})
		})
	}

	return { renderTableTree }
}


