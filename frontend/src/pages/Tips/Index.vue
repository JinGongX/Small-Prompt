<!-- ================= Input Panel ================= -->
<script setup lang="ts">
import { ref,onMounted,onUnmounted,watch } from 'vue'
import { InsTips } from '../../../bindings/changeme/services/suistore' 
import{HideTipsWindow} from '../../../bindings/changeme/services/appservice'
import { applyTheme } from '../../utils/ThemeManager'
import { parseTime } from '../../utils/timeParser'
import { CheckCircleTwoTone, RocketOutlined, ThunderboltOutlined, FieldTimeOutlined,SendOutlined } from '@ant-design/icons-vue';
import { debounce,parseTimePreview } from '../../utils/useDebounce'
import { themeChannel,settingChannel } from '../../utils/langChannel'
import {  IsmacOS,OS_READY } from '../../utils/osinfo'
const ismacos=ref(false)
const inputRef = ref<HTMLInputElement | null>(null)
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
   // input.value = ''
   // tipType.value = 'scheduled'
    HideTipsWindow()
  }
}
watch(input, (val) => {
  runParse(val)
})
const runParse = debounce((text: string) => {
  parseResult.value = parseTimePreview(text,tipType.value)
}, 400)
let readyToAutoHide = false // 加入这个变量是为了避免一打开窗口就触发 blur 导致窗口关闭的情况，设置一个短暂的延迟后才允许自动隐藏
onMounted(async() => {
  window.addEventListener('keydown', onKeydown)

  inputRef.value?.focus()
  document.addEventListener("mouseup", () => {
    inputRef.value?.focus()
  })
  await OS_READY
  ismacos.value=IsmacOS();

  const savedCompletion = localStorage.getItem('completion')
  if (savedCompletion) {
    completion.value = parseInt(savedCompletion, 10)
  }
  const savedDuration = localStorage.getItem('duration')
  if (savedDuration) {
    duration.value = parseInt(savedDuration, 10)
  }

  setTimeout(() => {
    readyToAutoHide = true
  }, 1000)
   window.addEventListener('blur', onBlur)

})
let blurTimeout: any
onUnmounted(() => {
  clearTimeout(blurTimeout)
  window.removeEventListener('blur', onBlur)
  window.removeEventListener('keydown', onKeydown)
})
 

function onBlur() {
  if (readyToAutoHide) {
    blurTimeout = setTimeout(() => {
      // 如果此时已经重新获取焦点就不关闭了
      if (document.hasFocus()) return
      if(input.value.trim()) {
        // 如果输入框有内容，说明用户可能在输入，暂不关闭
        return
      }
      HideTipsWindow()
    }, 100)
  }
}

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

  // const bc = new BroadcastChannel('theme')
  themeChannel.onmessage = (e) => {
    applyTheme(e.data)
  }
  // const setting = new BroadcastChannel('settings')
  settingChannel.onmessage = (e) => {
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
  <!-- <div 
    class="drag-region fixed   flex  items-center justify-center  z-50"
  > -->
    <div :class="ismacos?'':'border border-gray-500 dark:border-gray-500'"
      class="w-[360px] inset-0 drag-region rounded-3xl flex flex-col bg-neutral-50 dark:bg-gray-500 dark:text-white  px-4 pt-4"
    >
      <!-- Header -->
      <!-- <div   class="mb-4 text-sm  text-neutral-500 dark:text-white">
        新提示
      </div> -->

    
       <!-- Input -->
      <textarea
        v-model="input"  @keydown.enter.prevent="sendMessage" ref="inputRef"
        rows="2"
        placeholder="有思路，就写下来" 
        class="w-full  no-drag indent-1 resize-none  font-medium  bg-neutral-50 dark:bg-gray-500  py-1 text-sm text-neutral-900  dark:text-gray-100 placeholder:text-neutral-400 dark:placeholder:text-gray-100 focus:outline-none caret-blue-600 dark:caret-blue-200"
      />
      <!-- Hint -->
        <div class="drag-region py-3 flex items-center justify-between  text-neutral-400 dark:text-white "> 
           <div class="flex gap-3"> 
          <a-tooltip  placement="topRight">
            <template #title>即时提示</template> 
            <ThunderboltOutlined  @click="sendtipType('immediate')"  :class="tipType==='immediate'?'text-orange-300 dark:text-orange-400':'text-gray-500 dark:text-white'"    />
          </a-tooltip>
             <a-tooltip  placement="topRight">
            <template #title>定时提示</template> 
             <FieldTimeOutlined   @click="sendtipType('scheduled')" :class="tipType==='scheduled'?'text-orange-300 dark:text-orange-400':'text-gray-500 dark:text-white'" />
            </a-tooltip>
           </div>

        <div class="text-xs border-dotted border-1 border-gray-400 dark:border-white rounded">
          <!-- 自动识别类型 -   Enter 确认 · Esc 取消 · -->
         {{( parseResult.status === 'empty')
            ? '将默认提醒'
            : (parseResult.status === 'invalid')
            ? '将默认提醒'
            : `${parseResult.text}` }}
        </div> 
           <div class="flex gap-3">
            <div @click="sendMessage" class="w-[30px] h-[30px]  text-neutral-50 rounded-full   items-center justify-center flex pl-1" :class="input.trim().length > 0 ? 'bg-neutral-800':'bg-neutral-300 dark:bg-neutral-400'">
              <SendOutlined  />
            </div>
            
           </div>
      </div>

       
    </div>
  <!-- </div> -->
</template>
<style>
.drag-region {
  -webkit-app-region: drag;
  --wails-draggable: drag;
}
.no-drag {
  --wails-draggable: no-drag;
}
</style>