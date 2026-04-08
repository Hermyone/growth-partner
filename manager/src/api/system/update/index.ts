import request from '/@/utils/request';

export function getUpdaeLogsList(query:Object) {
    return request({
        url: '/api/v1/system/updateLog/list',
        method: 'get',
        params:query
    })
}

