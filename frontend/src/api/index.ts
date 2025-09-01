import request from '../utils/request'
// 基础接口
/**
 * 打招呼
 * @returns Hello, AsyncLab
 */
export function sayHello() {
    return request({
        url: '/hello',
        method: 'GET'
    })
}
