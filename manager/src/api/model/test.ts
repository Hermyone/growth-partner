import sysConfig from '/@/config';
import { requestToObject_ex } from '/@/api/items';

const requestList = [
	{ key: 'dict', path: '/system/dict/data/list', method: 'get', name: '数据字典查询' },
	{ key: 'params', path: '/system/config/list', method: 'get', name: '配置参数查询' },
]

const mockList = [
	{ key: 'mock', path: 'mock', method: 'post', name: '测试' }
]

export default {
	...requestToObject_ex(requestList, `${sysConfig.API_URL}api/v1`),
	...requestToObject_ex(mockList, `${sysConfig.VITE_URL}`),
}
