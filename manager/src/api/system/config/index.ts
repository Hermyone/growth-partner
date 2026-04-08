import { http } from '/@/utils/request';

export function getConfigList(query:Object) {
    return http.get('/api/v1/system/config/list', query)
}

export function getConfig(id:number) {
    return http.get('/api/v1/system/config/get', {id})
}

export function addConfig(data:any) {
    return http.post('/api/v1/system/config/add', data)
}

export function editConfig(data:any) {
    return http.put('/api/v1/system/config/edit', data)
}

export function deleteConfig(ids:number[]) {
    return http.delete('/api/v1/system/config/delete', {ids})
}
