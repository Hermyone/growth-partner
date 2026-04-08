/**
 * 锚点
 * @param id 要跳转到的锚点id
 */
export function goAnchor(id: string) {
	const el = document.querySelector('#'+id);
	if(el){
			el.scrollIntoView({
					behavior: 'smooth'
			})
	}
}

/**
 * 文件下载
 * @param url 要下载的url地址
 * @param filename 保存的文件名
 */
export function download(url: string, filename: string) {
	const link = document.createElement('a');
	if (link.download !== undefined) {
		link.setAttribute('href', url);
		link.setAttribute('download', filename);
		link.style.visibility = 'hidden';
		document.body.appendChild(link);
		link.click();
		document.body.removeChild(link);
	}
}