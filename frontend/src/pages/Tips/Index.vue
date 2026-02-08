<!-- ================= Input Panel ================= -->
<script setup lang="ts">
import { ref,onMounted,onUnmounted,watch } from 'vue'
import { InsTips } from '../../../bindings/changeme/services/suistore' 
import{HideTipsWindow} from '../../../bindings/changeme/appservice'
import { applyTheme } from '../../utils/ThemeManager'
import { parseTime } from '../../utils/timeParser'
import { CheckCircleTwoTone, RocketOutlined, ThunderboltOutlined, FieldTimeOutlined } from '@ant-design/icons-vue';
import { debounce,parseTimePreview } from '../../utils/useDebounce'

const input = ref('')
const tipType = ref('scheduled') // 'scheduled'|'immediate'
export type TimePreview =
  | { status: 'empty' }
  | { status: 'invalid' }
  | { status: 'ok'; at: number; text: string }
const parseResult = ref<TimePreview>({ status: 'empty' })
const duration = ref(10)
const completion = ref(10)

const sendMessage = () => {
  if (!input.value.trim()) return
  // 发送消息逻辑
  console.log('发送消息:', input.value)
  const now = Math.floor(Date.now() / 1000)
  //const expireAt = now + 30 * 60 // 30 分钟后过期
  const parsed = parseTime(input.value)
  var snoozeAt = parsed || (now + duration.value * 60) // 默认 30 分钟后提醒 now+1*15 //
  if(tipType.value==='immediate'){
    snoozeAt=now
  }
  const expireAt = (snoozeAt + completion.value * 60)
  InsTips(tipType.value, detectType(input.value), input.value, expireAt, snoozeAt)
  input.value = ''
  tipType.value = 'scheduled'
  HideTipsWindow()
}
// 监听 Esc 键关闭面板
const onKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') {
    input.value = ''
    tipType.value = 'scheduled'
    HideTipsWindow()
  }
}
watch(input, (val) => {
  runParse(val)
})
const runParse = debounce((text: string) => {
  parseResult.value = parseTimePreview(text,tipType.value)
}, 400)

onMounted(() => {
  window.addEventListener('keydown', onKeydown)

  const savedCompletion = localStorage.getItem('completion')
  if (savedCompletion) {
    completion.value = parseInt(savedCompletion, 10)
  }
  const savedDuration = localStorage.getItem('duration')
  if (savedDuration) {
    duration.value = parseInt(savedDuration, 10)
  }

})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
})

//--- 识别类型和时间的简单函数 ---
type PromptType = 'security' | 'device' | 'life' | 'work'|'system' |'default'|'rest'
function detectType(text: string): PromptType {
  if (/密码|安全|验证/.test(text)) return 'security'
  if (/备份|存储|照片|电脑/.test(text)) return 'device'
  if (/吃饭|做饭|生活/.test(text)) return 'life'
  if (/休息|午休|咖啡|喝水|下班/.test(text)) return 'rest'
  if (/会议|任务|工作|开会/.test(text)) return 'work'
  return 'default'
}

// function parseTime(text: string): number | undefined {
//   const m = text.match(/(\d+)\s*(分钟|分|min|m)/)
//   if (m) return Date.now() + Number(m[1]) * 60 * 1000
//   return undefined
// }
// function parseTime(text: string): number | undefined {
//   const m = text.match(/(\d+)\s*(分钟|分|min|m)/)
//   if (!m) return undefined

//   const minutes = Number(m[1])
//   return Math.floor(Date.now() / 1000) + minutes * 60
// }

  const bc = new BroadcastChannel('theme')
  bc.onmessage = (e) => {
    applyTheme(e.data)
  }
  const setting = new BroadcastChannel('settings')
  setting.onmessage = (e) => {
    if(e.data.type==='duration'){ 
      duration.value=e.data.value
    }else if(e.data.type==='completion'){
      completion.value=e.data.value
    }
  }

  function sendtipType(item){
    if(item===tipType.value){
      return
    }else if(item!='immediate'){
      tipType.value='scheduled'
    }else{
      tipType.value='immediate'
    }
    runParse(input.value)
  }
</script>

<template>
  <div
    class="fixed inset-0 flex items-center justify-center bg-black/20 backdrop-blur-sm dark:bg-gray-400/20"
  >
    <div
      class="w-[420px] rounded-1xl bg-white/90 dark:bg-gray-500/90 dark:text-white backdrop-blur-md shadow-[0_16px_40px_rgba(0,0,0,0.18)] px-5 py-4"
    >
      <!-- Header -->
      <div class="mb-3 text-sm text-neutral-500 dark:text-white">
        轻提示 · 输入一句你不想忘的事
      </div>

      <!-- Input -->
      <textarea
        v-model="input"  @keydown.enter.prevent="sendMessage"
        rows="2"
        placeholder="例如：30 分钟后提醒我备份照片"
        class="w-full resize-none rounded-xl border border-neutral-200 bg-neutral-50 dark:bg-gray-100 px-4 py-3 text-sm text-neutral-900 placeholder:text-neutral-400 focus:border-neutral-300 focus:outline-none"
      />

      <!-- Hint -->
        <div class="mt-3 flex items-center justify-between  text-neutral-400 dark:text-white"> 
        <div class="text-xs">
          <!-- 自动识别类型 -  -->
          Enter 确认 · Esc 取消 · {{( parseResult.status === 'empty')
            ? '将默认提醒'
            : (parseResult.status === 'invalid')
            ? '将默认提醒'
            : `${parseResult.text}` }}
        </div>
         <div class="flex gap-3"> 
          <a-tooltip>
            <template #title>即时提示</template> 
            <ThunderboltOutlined  @click="sendtipType('immediate')" :class="tipType==='immediate'?'text-orange-300 dark:text-orange-400':'text-gray-400 dark:text-white'"   />
          </a-tooltip>
             <a-tooltip  placement="topRight">
            <template #title>定时提示</template> 
             <FieldTimeOutlined   @click="sendtipType('scheduled')" :class="tipType==='scheduled'?'text-orange-300 dark:text-orange-400':'text-gray-400 dark:text-white'" />
            </a-tooltip>
           </div>
      </div>
    </div>
  </div>
</template>