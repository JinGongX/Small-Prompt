import { ref, onMounted, onUnmounted } from 'vue'

/**
 * 全局时间心跳（默认 1s）
 * 用于进度条 / 到期判断 / 自动完成
 */
export function useNowTick(interval = 1000) {
  const now = ref(Date.now())
  let timer: number | null = null

  onMounted(() => {
    timer = window.setInterval(() => {
      now.value = Date.now()
    }, interval)
  })

  onUnmounted(() => {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
  })

  return { now }
}