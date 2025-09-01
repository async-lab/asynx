/**
 * 实现token设置和获取等的工具类
 */

import Cookies from 'js-cookie'

// Token 相关常量
const TOKEN_KEY = 'asynx_token'
const TOKEN_EXPIRATION_DAY = 1 // 存Cookie的token的过期时间 => 1天

// 用户名相关常量
const USERNAME_KEY = 'async_is_remember_username'

/**
 * 获取Cookie中的token
 * @returns token字符串或undefined
 */
export function getToken(): string | undefined {
    return Cookies.get(TOKEN_KEY)
}

/**
 * 存储Cookie中的token
 * @param token token字符串
 * @returns 设置结果
 */
export function setToken(token: string): string | undefined {
    return Cookies.set(TOKEN_KEY, token, { expires: TOKEN_EXPIRATION_DAY })
}

/**
 * 移除token（从sessionStorage和Cookie中）
 */
export function removeToken(): void {
    sessionStorage.removeItem(TOKEN_KEY)
    Cookies.remove(TOKEN_KEY)
}

/**
 * 设置用户名到localStorage
 * @param username 用户名
 */
export function setUsername(username: string = ''): void {
    localStorage.setItem(USERNAME_KEY, username)
}

/**
 * 从localStorage获取用户名
 * @returns 用户名或null
 */
export function getUsername(): string | null {
    return localStorage.getItem(USERNAME_KEY)
}