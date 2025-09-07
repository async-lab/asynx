/**
 * 实现token设置和获取等的工具类
 */

import Cookies from 'js-cookie'
import type { User } from '@/api/types'

// Token 相关常量
const TOKEN_KEY = 'asynx_token'
const TOKEN_EXPIRATION_DAY = 1 // 存Cookie的token的过期时间 => 1天

// 用户名相关常量
const USERNAME_KEY = 'async_is_remember_username'
const USER_PROFILE_KEY = 'async_user_profile'

// 记住密码相关常量（仅本地使用）
const PWD_KEY = 'async_encrypted_password'
const PWD_IV_KEY = 'async_encrypted_password_iv'
const PWD_VER_KEY = 'async_encrypted_password_ver'
const PWD_VER = 'v1'
const APP_SECRET = 'ASYNC_APP_LOCAL_PWD_V1'

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

/**
 * 保存用户信息到本地
 */
export function setUserProfile(user: User): void {
    try {
        localStorage.setItem(USER_PROFILE_KEY, JSON.stringify(user))
    } catch {
        // ignore
    }
}

/**
 * 获取本地保存的用户信息
 */
export function getUserProfile(): User | null {
    const raw = localStorage.getItem(USER_PROFILE_KEY)
    if (!raw) return null
    try {
        return JSON.parse(raw) as User
    } catch {
        return null
    }
}

/**
 * 清除本地保存的用户信息
 */
export function clearUserProfile(): void {
    localStorage.removeItem(USER_PROFILE_KEY)
}

/**
 * 合并保存角色和类型到用户资料
 */
export function setUserRoleCategory(role?: string, category?: string): void {
    const currentRaw = localStorage.getItem(USER_PROFILE_KEY)
    let current: any = {}
    try { current = currentRaw ? JSON.parse(currentRaw) : {} } catch { current = {} }
    const next = {
        ...current,
        ...(role !== undefined ? { role } : {}),
        ...(category !== undefined ? { category } : {})
    }
    try { localStorage.setItem(USER_PROFILE_KEY, JSON.stringify(next)) } catch {}
}

export function getStoredRole(): string | undefined {
    const up: any = getUserProfile() as any
    return up?.role
}

export function getStoredCategory(): string | undefined {
    const up: any = getUserProfile() as any
    return up?.category
}

// ======== 本地加密记住密码 ========

async function deriveAesKey(username: string): Promise<CryptoKey> {
    const enc = new TextEncoder()
    const baseKey = await crypto.subtle.importKey(
        'raw',
        enc.encode(APP_SECRET),
        { name: 'PBKDF2' },
        false,
        ['deriveKey']
    )
    const salt = enc.encode('ASYNC_SALT_' + (username || 'anonymous'))
    return crypto.subtle.deriveKey(
        {
            name: 'PBKDF2',
            salt,
            iterations: 100_000,
            hash: 'SHA-256'
        },
        baseKey,
        { name: 'AES-GCM', length: 256 },
        false,
        ['encrypt', 'decrypt']
    )
}

export async function saveEncryptedPassword(username: string, password: string): Promise<void> {
    try {
        if (!password) {
            clearEncryptedPassword()
            return
        }
        const key = await deriveAesKey(username)
        const iv = crypto.getRandomValues(new Uint8Array(12))
        const enc = new TextEncoder()
        const cipher = await crypto.subtle.encrypt({ name: 'AES-GCM', iv }, key, enc.encode(password))
        const cipherBytes = new Uint8Array(cipher)
        const b64 = btoa(String.fromCharCode(...cipherBytes))
        const ivB64 = btoa(String.fromCharCode(...iv))
        localStorage.setItem(PWD_KEY, b64)
        localStorage.setItem(PWD_IV_KEY, ivB64)
        localStorage.setItem(PWD_VER_KEY, PWD_VER)
    } catch (e) {
        // 失败则不保存
        clearEncryptedPassword()
        console.error('Failed to save encrypted password', e)
    }
}

export async function loadEncryptedPassword(username: string): Promise<string | null> {
    try {
        const ver = localStorage.getItem(PWD_VER_KEY)
        if (ver !== PWD_VER) return null
        const b64 = localStorage.getItem(PWD_KEY)
        const ivB64 = localStorage.getItem(PWD_IV_KEY)
        if (!b64 || !ivB64) return null
        const cipher = Uint8Array.from(atob(b64), c => c.charCodeAt(0))
        const iv = Uint8Array.from(atob(ivB64), c => c.charCodeAt(0))
        const key = await deriveAesKey(username)
        const plainBuf = await crypto.subtle.decrypt({ name: 'AES-GCM', iv }, key, cipher)
        const dec = new TextDecoder()
        return dec.decode(plainBuf)
    } catch (e) {
        console.warn('Failed to load encrypted password', e)
        return null
    }
}

export function clearEncryptedPassword(): void {
    localStorage.removeItem(PWD_KEY)
    localStorage.removeItem(PWD_IV_KEY)
    localStorage.removeItem(PWD_VER_KEY)
}