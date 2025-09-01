import request from '../utils/request'
import type { 
    User, 
    RegisterRequest, 
    ChangePasswordRequest, 
    ModifyRoleRequest, 
    ModifyCategoryRequest 
} from './types'

/**
 * 获取用户列表
 * @returns 用户列表
 */
export function getUserList() {
    return request({
        url: '/users',
        method: 'GET'
    })
}

/**
 * 注册新用户
 * @param {Object} reqData 注册用户请求数据
 * @param {string} reqData.username 用户名
 * @param {string} reqData.givenName 名
 * @param {string} reqData.surName 姓
 * @param {string} reqData.mail 邮箱
 * @param {string} reqData.role 角色
 * @param {string} reqData.category 账号类型
 * @returns 注册结果
 */
export function registerUser(reqData: RegisterRequest) {
    return request({
        url: '/users',
        method: 'POST',
        data: reqData
    })
}

/**
 * 获取用户信息
 * @param {string} uid 用户ID，使用 'me' 可获取当前用户信息
 * @returns 用户信息
 */
export function getUserInfo(uid: string) {
    return request({
        url: `/users/${uid}`,
        method: 'GET'
    })
}

/**
 * 删除用户
 * @param {string} uid 用户ID，不能使用 'me'
 * @returns 删除结果
 */
export function deleteUser(uid: string) {
    return request({
        url: `/users/${uid}`,
        method: 'DELETE'
    })
}

/**
 * 获取账号类型
 * @param {string} uid 用户ID，使用 'me' 可获取当前用户类型
 * @returns 账号类型: system|member|external
 */
export function getUserCategory(uid: string) {
    return request({
        url: `/users/${uid}/category`,
        method: 'GET'
    })
}

/**
 * 更改账号类型
 * @param {string} uid 用户ID，不能使用 'me'
 * @param {Object} reqData 修改账号类型请求数据
 * @param {string} reqData.category 账号类型 system|member|external
 * @returns 修改结果
 */
export function modifyUserCategory(uid: string, reqData: ModifyCategoryRequest) {
    return request({
        url: `/users/${uid}/category`,
        method: 'PUT',
        data: reqData
    })
}

/**
 * 修改密码
 * @param {string} uid 用户ID，使用 'me' 可修改当前用户密码
 * @param {Object} reqData 修改密码请求数据
 * @param {string} reqData.password 新密码
 * @returns 修改结果
 */
export function changePassword(uid: string, reqData: ChangePasswordRequest) {
    return request({
        url: `/users/${uid}/password`,
        method: 'PUT',
        data: reqData
    })
}

/**
 * 获取账号角色
 * @param {string} uid 用户ID，使用 'me' 可获取当前用户角色
 * @returns 账号角色: admin|default|restricted
 */
export function getUserRole(uid: string) {
    return request({
        url: `/users/${uid}/role`,
        method: 'GET'
    })
}

/**
 * 更改账号角色
 * @param {string} uid 用户ID，不能使用 'me'
 * @param {Object} reqData 修改账号角色请求数据
 * @param {string} reqData.role 角色 admin|default|restricted
 * @returns 修改结果
 */
export function modifyUserRole(uid: string, reqData: ModifyRoleRequest) {
    return request({
        url: `/users/${uid}/role`,
        method: 'PUT',
        data: reqData
    })
} 