/**
 * API 响应数据接口
 */
export interface ApiResponse<T = any> {
    code: number
    msg: string
    data: T
}

/**
 * 用户信息接口
 */
export interface User {
    username: string
    givenName: string
    surName: string
    mail: string
    role: string
    category: string
}

/**
 * 登录请求接口
 */
export interface LoginRequest {
    username: string
    password: string
}

/**
 * 注册用户请求接口
 */
export interface RegisterRequest {
    username: string
    givenName: string
    surName: string
    mail: string
    role: string
    category: string
}

/**
 * 修改密码请求接口
 */
export interface ChangePasswordRequest {
    password: string
}

/**
 * 修改角色请求接口
 */
export interface ModifyRoleRequest {
    role: string
}

/**
 * 修改账号类型请求接口
 */
export interface ModifyCategoryRequest {
    category: string
}

/**
 * 角色类型枚举
 */
export type RoleType = 'admin' | 'default' | 'restricted'

/**
 * 账号类型枚举
 */
export type CategoryType = 'system' | 'member' | 'external' 