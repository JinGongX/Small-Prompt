export function debounce<T extends (...args: any[]) => void>(
  fn: T,
  delay = 400
) {
  let timer: number | undefined

  return (...args: Parameters<T>) => {
    clearTimeout(timer)
    timer = window.setTimeout(() => {
      fn(...args)
    }, delay)
  }
}



import { parseTime } from './timeParser'

const MINUTE = 60
const HOUR = 60 * MINUTE
const DAY = 24 * HOUR

export type TimePreview =
  | { status: 'empty' }
  | { status: 'invalid' }
  | { status: 'ok'; at: number; text: string }

function nowSec() {
  return Math.floor(Date.now() / 1000)
}

function formatTime(ts: number) {
  const d = new Date(ts * 1000)
  return d.toLocaleString('zh-CN', {
    month: 'numeric',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

export function parseTimePreview(text: string,tipType: string): TimePreview {
  if (tipType==='immediate'){
    return {
      status: 'ok',
      at: nowSec(),
      text: `将立即提醒`//`将在 ${m} 分钟后提醒（${formatTime(at)}）`
    }
  }
  if (!text || !text.trim()) {
    return { status: 'empty' }
  }
  

  const at = parseTime(text)
  if (!at) {
    return { status: 'invalid' }
  }

  const now = nowSec()
  const diff = at - now

  if (diff > 0 && diff < DAY) {
    if (diff < HOUR) {
      const m = Math.max(1, Math.round(diff / MINUTE))
      return {
        status: 'ok',
        at,
        text: `将在 ${m} 分钟后提醒`//`将在 ${m} 分钟后提醒（${formatTime(at)}）`
      }
    }

    const h = Math.round(diff / HOUR)
    return {
      status: 'ok',
      at,
      text: `将在 ${h} 小时后提醒`//`将在 ${h} 小时后提醒（${formatTime(at)}）`
    }
  }

  return {
    status: 'ok',
    at,
    text: `识别为 ${formatTime(at)}`
  }
}