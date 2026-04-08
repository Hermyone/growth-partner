import request from '/@/utils/request';

export function getMessageList(query?:Object) {
    return request({
        url: '/api/v1/system/message/list',
        method: 'get',
        params:query
    })
}


export function addMessage(data:object) {
    return request({
        url: '/api/v1/system/message/add',
        method: 'post',
        data:data
    })
}


export function editMessage(data:object) {
    return request({
        url: '/api/v1/system/message/edit',
        method: 'put',
        data:data
    })
}

export function sendMessage(data:object) {
    return request({
        url: '/api/v1/system/message/send',
        method: 'put',
        data:data
    })
}

export function deleteMessage(id:number) {
    return request({
        url: '/api/v1/system/message/delete',
        method: 'delete',
        data:{id}
    })
}

export function getReadMessage(data:object) {
    return request({
        url: '/api/v1/system/message/read/list',
        method: 'get',
        data:data
    })
}

export function setReadMessage(data:object) {
    return request({
        url: '/api/v1/system/message/read',
        method: 'post',
        data:data
    })
}

export function getNoReadMessage(data:object) {
    return request({
        url: '/api/v1/system/message/noread',
        method: 'get',
        data:data
    })
}