/**
 * 时间解析工具（完整版）
 * 返回 Unix 时间戳（秒）或 undefined
 */

export type ParseResult = number | undefined

const MINUTE = 60
const HOUR = 60 * MINUTE
const DAY = 24 * HOUR

/* =========================
 * 基础工具
 * ========================= */

function nowSec() {
  return Math.floor(Date.now() / 1000)
}

function todayAt(hour: number, minute = 0) {
  const d = new Date()
  d.setHours(hour, minute, 0, 0)
  return Math.floor(d.getTime() / 1000)
}

function tomorrowAt(hour: number, minute = 0) {
  const d = new Date()
  d.setDate(d.getDate() + 1)
  d.setHours(hour, minute, 0, 0)
  return Math.floor(d.getTime() / 1000)
}

function nextTime(hour: number, minute: number, now: number) {
  const d = new Date()
  d.setHours(hour, minute, 0, 0)
  let ts = Math.floor(d.getTime() / 1000)
  if (ts <= now) {
    d.setDate(d.getDate() + 1)
    ts = Math.floor(d.getTime() / 1000)
  }
  return ts
}

/* =========================
 * 中文数字处理（轻量但够用）
 * ========================= */

const CN_NUM: Record<string, number> = {
  一: 1,
  二: 2,
  三: 3,
  四: 4,
  五: 5,
  六: 6,
  七: 7,
  八: 8,
  九: 9,
  十: 10
}

function chineseNumberToInt(text: string): number | undefined {
  if (text === '十') return 10

  if (text.length === 2 && text[0] === '十') {
    return 10 + (CN_NUM[text[1]] ?? 0)
  }

  if (text.length === 2 && text[1] === '十') {
    return (CN_NUM[text[0]] ?? 0) * 10
  }

  if (text.length === 3 && text[1] === '十') {
    return (CN_NUM[text[0]] ?? 0) * 10 + (CN_NUM[text[2]] ?? 0)
  }

  return CN_NUM[text]
}

/**
 * 文本归一化
 * 把中文时间语义转成数字分钟
 */
function normalizeText(text: string): string {
  return text
    .replace(/半小时/g, '30分钟')
    .replace(/一刻钟/g, '15分钟')
    .replace(/三刻钟/g, '45分钟')
}

/* =========================
 * 时间规则
 * ========================= */

type TimeRule = {
  name: string
  regex: RegExp
  handler: (m: RegExpMatchArray, now: number) => number
}

const RULES: TimeRule[] = [
  // 🔴 中文数字分钟（十分钟）
  {
    name: '中文分钟',
    regex: /([一二三四五六七八九十]+)\s*分钟/,
    handler: (m, now) => {
      const v = chineseNumberToInt(m[1])
      if (!v) throw new Error('invalid cn number')
      return now + v * MINUTE
    }
  },

  // 🔴 数字分钟
  {
    name: 'X分钟后',
    regex: /(\d+)\s*(分钟|分|min|m)(后)?/,
    handler: (m, now) => now + Number(m[1]) * MINUTE
  },

  // 🔴 小时
  {
    name: 'X小时后',
    regex: /(\d+)\s*(小时|时|h)(后)?/,
    handler: (m, now) => now + Number(m[1]) * HOUR
  },

  // 🔴 天
  {
    name: 'X天后',
    regex: /(\d+)\s*(天|day|d)(后)?/,
    handler: (m, now) => now + Number(m[1]) * DAY
  },

  // 🟠 快捷语义
  {
    name: '马上',
    regex: /(马上|立即)/,
    handler: (_, now) => now + 1 * MINUTE
  },
  {
    name: '稍后',
    regex: /(稍后|等下)/,
    handler: (_, now) => now + 10 * MINUTE
  },

  // 🟡 模糊日期
  {
    name: '今天',
    regex: /(今天)(早上|上午|中午|下午|晚上)?/,
    handler: (m) => {
      const p = m[2]
      if (p === '中午') return todayAt(12)
      if (p === '下午') return todayAt(15)
      if (p === '晚上') return todayAt(20)
      return todayAt(9)
    }
  },
  {
    name: '明天',
    regex: /(明天)(早上|上午|中午|下午|晚上)?/,
    handler: (m) => {
      const p = m[2]
      if (p === '中午') return tomorrowAt(12)
      if (p === '下午') return tomorrowAt(15)
      if (p === '晚上') return tomorrowAt(20)
      return tomorrowAt(9)
    }
  },

  // 🟢 HH:mm
  {
    name: 'HH:mm',
    regex: /(\d{1,2}):(\d{2})/,
    handler: (m, now) => nextTime(Number(m[1]), Number(m[2]), now)
  },

  // 🟢 X点半
  {
    name: 'X点半',
    regex: /(\d{1,2})点半/,
    handler: (m, now) => nextTime(Number(m[1]), 30, now)
  },

  // 🟢 X点
  {
    name: 'X点',
    regex: /(\d{1,2})点/,
    handler: (m, now) => nextTime(Number(m[1]), 0, now)
  },

  // 🔵 兜底：30m / 2h
  {
    name: '数字单位',
    regex: /(\d+)(m|h|d)/,
    handler: (m, now) => {
      const v = Number(m[1])
      if (m[2] === 'm') return now + v * MINUTE
      if (m[2] === 'h') return now + v * HOUR
      return now + v * DAY
    }
  }
]

/* =========================
 * 主函数（对外 API）
 * ========================= */

export function parseTime(text: string): ParseResult {
  if (!text) return undefined

  const now = nowSec()
  const input = normalizeText(text.replace(/\s+/g, '').toLowerCase())

  for (const rule of RULES) {
    const m = input.match(rule.regex)
    if (m) {
      try {
        return rule.handler(m, now)
      } catch {
        return undefined
      }
    }
  }

  return undefined
}