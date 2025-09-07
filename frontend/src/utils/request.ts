import axios from 'axios'
import type { InternalAxiosRequestConfig, AxiosResponse, AxiosError } from 'axios'
import {
  getToken,
  removeToken
} from './auth'
import { 
  useWarningConfirm,
  useFailedTip
} from './msgTip'

/**
 * API响应数据接口
 */
interface ApiResponse<T = any> {
  code?: number
  msg?: string
  data?: T
}

/**
 * 封装axios
 */

axios.defaults.withCredentials = true

// 获取API基础URL，优先使用环境变量，否则使用当前域名的/api路径
const getApiBaseUrl = () => {
    if (import.meta.env.VITE_API_BASE_URL) {
        return import.meta.env.VITE_API_BASE_URL
    }
    // 默认使用当前域名的/api路径
    return `${window.location.origin}/api`
}

const request = axios.create({
    baseURL: getApiBaseUrl()
})

// 添加请求拦截-token处理
request.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
        // 不需要认证的接口直接放行
        const publicApis = ['/login', '/tokens']
        if (publicApis.includes(config.url || '')) {
            return config
        }

        const token = getToken()
        const tokenString = 'Bearer ' + token
        // 需要认证的接口但没有token就取消请求，并提示
        if (!token) {
            console.log('Token 未找到，取消请求')
            return Promise.reject(new Error('登录异常，Token 未找到，请求被取消'))
        }
        // 在请求头设置token
        config.headers.set('Authorization', tokenString)

        return config
    },
    (error: AxiosError) => {
        useFailedTip('请求错误：' + error.message)
        return Promise.reject(error)
    }
)

// 添加响应拦截
request.interceptors.response.use(
    (response: AxiosResponse<ApiResponse>) => {
        const { data } = response
        
        // 检查是否有code字段，如果没有则直接返回数据
        if (!data.hasOwnProperty('code')) {
            return Promise.resolve(data)
        }
        
        return Promise.resolve(data.data)
    },
    (error: AxiosError) => {
        // 检查是否是HTML响应（通常表示重定向到前端页面）
        if (error.response?.data && typeof error.response.data === 'string' && error.response.data.includes('<!doctype html>')) {
            console.error('API请求被重定向到前端页面，请检查代理配置或后端服务状态')
            useFailedTip('API服务不可用，请检查后端服务是否正常运行')
        } else {
            useFailedTip('响应错误：' + error.message)
        }
        
        if (!getToken()) {
            useWarningConfirm('登录异常，即将跳转登录页，重新登录').then(() => {
                // router.push('/login')
                window.location.reload()
            })
        }

        const status = error.response?.status

        if (status === 401) { // authorized，token过期或token异常等的返回码
            useWarningConfirm('登录过期或异常，即将跳转登录页，重新登录').then(() => {
                removeToken()
                window.location.reload()
            })
            return Promise.reject(error)
        }

        const payload: any = error.response?.data
        const text = typeof payload === 'string'
            ? payload
            : (typeof payload?.data === 'string' ? payload.data : (payload?.msg || payload?.message || ''))
        useFailedTip(text)
        return Promise.reject(text)
    }
)

export default request