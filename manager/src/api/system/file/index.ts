import request from '/@/utils/request';

/**
 * 文件上传接口 (Python写的一个文件服务器，根据不同文件服务器类型可再作调整)
 *
 * file = 待上传的文件
 * path = 文件保存位置
 * filename = 文件名
 *
 * 具体参数说明详见协议文档
 */

export function uplaodFile(file: File, path: string, filename: string){
    var data = new FormData();
    data.append('file', file);
    data.append('filename', filename);
    data.append('path', path);

    return request({
        url: '/file',
        method: 'post',
        data: data,
    });
}