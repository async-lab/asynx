import request from '../utils/request'
import type { LoginRequest } from './types'

/**
 * 创建访问令牌（登录）
 * @param {Object} reqData 登录请求数据
 * @param {string} reqData.username 用户名
 * @param {string} reqData.password 密码
 * @returns 访问令牌
 */
export function createToken(reqData: LoginRequest) {
    return request({
        url: '/tokens',
        method: 'POST',
        data: reqData
    })
} 