import { ElMessage, ElMessageBox } from 'element-plus'
import type { ElMessageBoxOptions, MessageBoxData } from 'element-plus'

/**
 * 消息类型枚举
 */
type MessageType = 'success' | 'error' | 'warning' | 'info'

/**
 * 快捷生成消息提示以及消息对话框等的工具文件
 */

/**
 * 内部使用的消息提示函数
 * @param msg 消息内容
 * @param type 消息类型
 */
function useTip(msg: string = '', type: MessageType = 'success'): void {
    ElMessage({
        showClose: true,
        message: msg,
        type,
    })
}

/**
 * 返回成功的消息提示
 * @param msg 成功信息
 */
export function useSuccessTip(msg: string = '操作成功'): void {
    useTip(msg, 'success')
}

/**
 * 返回失败的消息提示
 * @param msg 错误信息
 */
export function useFailedTip(msg: string = '操作失败'): void {
    useTip(msg, 'error')
}

/**
 * 返回警告的消息提示
 * @param msg 警告信息
 */
export function useWarnTip(msg: string = '操作警告'): void {
    useTip(msg, 'warning')
}

/**
 * 返回info的消息提示
 * @param msg 信息
 */
export function useInfoTip(msg: string = '操作完成'): void {
    useTip(msg, 'info')
}

/**
 * 全自定义的操作提示框
 * @param message 提示信息
 * @param title 提示标题
 * @param options 提示框结构 
 * @returns 提示框的promise对象
 */
export function useConfirm(
    message: string, 
    title: string, 
    options: ElMessageBoxOptions
): Promise<MessageBoxData> {
    return ElMessageBox.confirm(message, title, options)
}

/**
 * 一般时候的操作提示框
 * @param message 提示信息 
 * @returns 提示框的promise对象
 */
export function useSimpleConfirm(message: string): Promise<MessageBoxData> {
    return useConfirm(message, '温馨提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
        showClose: true
    })
}

/**
 * 警示框（只能点确定）
 * @param message 提示信息
 * @returns 警示框的promise对象
 */
export function useWarningConfirm(message: string): Promise<MessageBoxData> {
    return useConfirm(message, '警告提醒', {
        confirmButtonText: '我知道了',
        showCancelButton: false,
        showClose: false,
        type: 'warning',
    })
}