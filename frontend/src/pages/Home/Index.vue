<!-- frontend/src/pages/Home/Index.vue -->
<script setup lang="ts">
import { computed,ref,onMounted,onUnmounted,watch } from 'vue'
import PromptCard from '../../components/Setting/PromptCard.vue';
import{GetTips,UpTipsState,UpTipsPinned,UpTipsDelayed } from  '../../../bindings/changeme/services/suistore'
import {Events} from "@wailsio/runtime";
import { applyTheme } from '../../utils/ThemeManager'
import Icon,{ CheckCircleTwoTone, PushpinOutlined , HistoryOutlined } from '@ant-design/icons-vue';
import { useNowTick } from '../../utils/useNowTick'

const MAX_SNOOZE = 3 // 最大延迟次数
// ---------- 类型定义 ----------
type PromptType = 'life' | 'work' | 'device' | 'security' | 'system' |'default'

interface PromptCard {
  id: number
  type: string //PromptType
  category: string
  title: string
  state: number
  pinned: number
  pinnedat?: number | null
  snoozeat?: number | null
  expireat: number
  snoozecount?: number // 延迟次数
  _isNew?: boolean // 是否为新添加的卡片（用于触发动画）
  activatedAt?: number  // ⭐ 实际弹出时间（秒）废弃
  autoCompleteAt?: number // activatedAt + 15min 废弃

}

const { now } = useNowTick(1000)
const AUTO_COMPLETE_DURATION =ref(15 * 60) 
function calcProgress(card: PromptCard) {  
  if (!card.snoozeat) return 0
  const totalDuration = card.expireat - (card.snoozeat ||  0)
  const elapsed = now.value / 1000 - (card.snoozeat ||  0)
  return Math.min(elapsed / totalDuration, 1) 
}
 
function formatTime(ts: number) {
  return new Date(ts * 1000).toLocaleString()
} 

// 根据类型获取样式
function getVisual(type: string) {
  const map: Record<string, { bg: string; image: string }> = {
    life: { bg: 'bg-amber-50', image: '/eat.png' },
    work: { bg: 'bg-blue-50', image: '/image2.png' },
    device: { bg: 'bg-teal-50', image: '/image2.png' },
    security: { bg: 'bg-emerald-50', image: '/image3.png' },
    system: { bg: 'bg-neutral-100', image: '/image3.png' },
    rest: { bg: 'bg-purple-50', image: '/coffee.png' },
    default: { bg: 'bg-neutral-100', image: '/default.png' },
  }
  return map[type] || map.default
} 
const tipsHistory = ref<PromptCard[]>([]);
const fetchTips = async () => {
  tipsHistory.value = await GetTips();
};
const cards = computed(() =>
  tipsHistory.value.filter(tip => tip.type==='immediate' && tip.state === 1)
)

const cardscheduleds = computed(() => 
  tipsHistory.value.filter(tip => tip.type==='scheduled' && tip.state === 2)
)
// 添加新卡片并触发动画
function addCard(card) {
  tipsHistory.value.unshift({
    ...card, 
    _isNew: true,
    activatedAt:  now.value / 1000,
    autoCompleteAt: now.value / 1000 + AUTO_COMPLETE_DURATION.value, 
  })
  // 下一帧移除标记，避免后续重排再触发动画
  requestAnimationFrame(() => {
    card._isNew = false
  })
}
const isReady = ref(false)
const onTipEvent = (item) => { 
  const raw = Array.isArray(item.data)
    ? item.data[0]
    : item.data

  addCard(raw)  
}
onMounted(async() => {
  await fetchTips();
  isReady.value = true
  const savedCompletion = localStorage.getItem('completion')
  if (savedCompletion) {
    AUTO_COMPLETE_DURATION.value = parseInt(savedCompletion, 10)*60
  } 
  Events.On('tipEvent', onTipEvent);
});
onUnmounted(() => {
  Events.Off('tipEvent'); 
})
computed(() => {
  watch(
  () => tipsHistory.value,
  (list) => {
    list.forEach(card => {
      if (card._isNew) {
        requestAnimationFrame(() => {
          card._isNew = false
        })
      }
    })
  },
  { immediate: true }
)
})

watch(now, () => {
  cardscheduleds.value.forEach(card => {
    if (
      card.expireat &&
      now.value / 1000 >= card.expireat
    ) { 
      onCardClose(card.id)
    }
  })
})
//
const onCardClick = (card) => {
  console.log('点击卡片', card)
} 

const onCardClose = (id: number) => { 
  const index = tipsHistory.value.findIndex(c => c.id === id)
  if (index !== -1) tipsHistory.value.splice(index, 1) // 删除卡片，避免留白
  UpTipsState(id,8); // 更新后端状态为已完成
} 

function canSnooze(card: PromptCard) {
  return (card.snoozecount ?? 0) < MAX_SNOOZE && card.state === 2
}
const onCardDelayed = (card: PromptCard) => {
  const newExpireAt = card.expireat + AUTO_COMPLETE_DURATION.value // snoozeAt + 自动完成时长
  card.expireat = newExpireAt
  card.snoozecount = (card.snoozecount ?? 0) + 1
  UpTipsDelayed(card.id, newExpireAt); // 更新后端状态为显示中（如果之前是已显示）
}

// const onCardPin = (id: number) => { 
//   const card = tipsHistory.value.find(c => c.id === id)
//   if (!card) return
//   const newPinned = card.pinned ? 0 : 1
//   //card.pinned = newPinned 
//   try  {
//      // 1️⃣ 先更新前端（响应式立刻生效）
//      card.pinned = newPinned
//      // 2️⃣ 再更新后端（持久化存储）
//      UpTipsPinned(id, newPinned);
//    } catch (error) {
//      // 如果后端更新失败，回滚前端状态
//      card.pinned = card.pinned
//   } 
// }

  const bc = new BroadcastChannel('theme')
  bc.onmessage = (e) => {
    applyTheme(e.data)
  }
  const setting = new BroadcastChannel('settings')
  setting.onmessage = (e) => {
    if(e.data.type==='duration'){ 
      AUTO_COMPLETE_DURATION.value=e.data.value*60
    }
  }
</script>

<template>
  <!-- Prompt Stack -->
<div class="drag-region">
     <div v-if="isReady &&cards.length > 0" class="fixed h-full top-2 right-4 pr-0.5 space-y-1 font-sans select-none   h-[86vh] overflow-y-auto  overflow-x-hidden scrollbar-thin pb-2"  > 
      <TransitionGroup name="card"> 
      <div
      v-for="card in cardscheduleds"
      :key="card.id"
      @click="onCardClick(card)"
      :data-new="card._isNew"
      class="group flex flex-col w-[294px] gap-1 rounded-2xl bg-white backdrop-blur-md px-1 py-1 shadow-[0_8px_24px_rgba(0,0,0,0.12)] transition. dark:bg-gray-400/80 dark:text-white hover:bg-white/90 dark:hover:bg-gray-500 cursor-pointer"
    >
     <div class="flex gap-1">
      <!-- Visual Block -->
      <div
         class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl"
        :class="getVisual(card.category).bg"
      >
        <img
          :src="getVisual(card.category).image"
          class="h-10 w-10 object-contain"
          draggable="false"
        />
      </div>
      <!-- Content -->
      <div
        class="flex-1 text-sm py-1 text-neutral-900 line-clamp-2 break-words  dark:text-white "
        :title="card.title"
      >
        {{ card.title }}
      </div>
      <!-- Close -->
      <div class="flex  flex-col h-10 w-8 gap-1 items-center justify-center mt-1">
        <!-- <PushpinOutlined   @click.stop="onCardPin(card.id)" :class="card.pinned===1?'text-orange-600 dark:text-orange-300':'text-gray-600 dark:text-white'" /> -->
        <button @click.stop="onCardDelayed(card)" :disabled="!canSnooze(card)" :title="`延迟提示 ${MAX_SNOOZE - (card.snoozecount ?? 0)} 次剩余`" >
          <!-- <HistoryOutlined  :rotate="60"  /> -->
           <icon>
              <template #component>
                <svg t="1770301650084" class="icon" viewBox="0 0 1088 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="10861" width="16" height="16"><path d="M955.935701 495.904165h-42.751986c-6.271998 0.032-11.999996-3.615999-14.559996-9.311997s-1.504-12.351996 2.815999-16.959995l78.431976-104.447967a15.839995 15.839995 0 0 1 25.343992 0l79.615975 105.375967a15.679995 15.679995 0 0 1-12.543996 25.183992h-42.719987l0.288 16.223995C1029.599678 794.624072 799.19975 1023.68 515.007839 1024 230.719928 1023.68 0.32 794.688072 0 512.00016 0.32 229.440248 230.655928 0.38432 515.007839 0.00032c100.799969-0.128 199.423938 29.247991 283.359912 84.447974a36.511989 36.511989 0 0 1-3.839999 63.26398 37.023988 37.023988 0 0 1-36.767989-2.207999 440.031862 440.031862 0 0 0-242.751924-72.415978c-243.679924 0.32-441.119862 196.703939-441.375862 438.975863 0.256 242.175924 197.727938 438.463863 441.375862 438.751863a443.103862 443.103862 0 0 0 317.535901-134.495958 437.887863 437.887863 0 0 0 123.391961-320.4159zM541.695831 286.080231v252.639921h177.919944a40.639987 40.639987 0 0 1 40.767987 40.543987c0.032 10.751997-4.255999 21.055993-11.871996 28.671991-7.679998 7.583998-18.047994 11.871996-28.895991 11.871996h-218.719932a40.383987 40.383987 0 0 1-40.831987-40.543987V286.080231c0-22.431993 18.271994-40.575987 40.767987-40.575988a40.639987 40.639987 0 0 1 40.799988 40.575988h0.064z m0 0" p-id="10862"></path></svg>             
              </template>
           </icon>
        </button>
        <button
          @click.stop="onCardClose(card.id)" title="完成"
          class="opacity-0 group-hover:opacity-100 text-neutral-400 transition hover:text-neutral-600"
        >
         <CheckCircleTwoTone  />
        </button> 
      </div>
      </div>
      <div 
       class="mt-[-3px] h-[3px] w-full overflow-hidden rounded-full bg-neutral-200/70 dark:bg-gray-600/60" :title="`自动完成时间：${formatTime(card.expireat||0)}`"
       >
      <div
          class="h-full  bg-emerald-500 transition-[width] duration-1000 linear" 
          :style="{ width: `${calcProgress(card) * 100}%` }"
        />
        </div>
    </div> 

       <!-- 即时类提示卡片 -->
    <div
      v-for="card in cards"
      :key="card.id"
      :data-new="card._isNew"
      @click="onCardClick(card)"
      class="group flex w-[294px] gap-1 rounded-2xl bg-white/80 backdrop-blur-md  px-1 py-1 shadow-[0_8px_24px_rgba(0,0,0,0.12)] transition. dark:bg-gray-400/80 dark:text-white hover:bg-white/90 dark:hover:bg-gray-500 cursor-pointer"
    >
      <!-- Visual Block -->
      <div
        class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl"
        :class="getVisual(card.category).bg"
      >
        <img
          :src="getVisual(card.category).image"
          class="h-10 w-10 object-contain"
          draggable="false"
        />
      </div>
      <!-- Content -->
      <div
        class="flex-1 text-sm py-1 text-neutral-900 line-clamp-2 break-words  dark:text-white "
        :title="card.title"
      >
        {{ card.title }}
      </div>
      <!-- Close -->
      <div class="flex  flex-col h-10 w-8 gap-1 items-center justify-center mt-1">
        <!-- <PushpinOutlined   @click.stop="onCardPin(card.id)" :class="card.pinned===1?'text-orange-600 dark:text-orange-300':'text-gray-600 dark:text-white'" /> -->
        <button
          @click.stop="onCardClose(card.id)"
          class="opacity-0 group-hover:opacity-100 text-neutral-400 transition hover:text-neutral-600"
        >
         <CheckCircleTwoTone  />
        </button>
        
      </div>
    </div> 
    </TransitionGroup>
  </div>
</div>
 
</template>
<style>
/* 只有 data-new="true" 的卡片才有 enter 动画 */
.card-enter-from[data-new="true"] {
  opacity: 0;
  transform: translateY(-6px) scale(0.85);
}

.card-enter-active[data-new="true"] {
  transition:
    transform 0.35s cubic-bezier(0.34, 1.56, 0.64, 1),
    opacity 0.2s ease-out;
}

.card-enter-to[data-new="true"] {
  opacity: 1;
  transform: translateY(0) scale(1);
}

/* 非新增卡片：禁用 enter 动画 */
.card-enter-active:not([data-new="true"]) {
  transition: none;
}

.card-leave-from {
  opacity: 1;
  transform: scale(1) translateX(0);
}

.card-leave-active {
  transition:
    transform 0.25s ease-in,
    opacity 0.2s ease-in;
}

.card-leave-to {
  opacity: 0;
  transform: scale(0.92) translateX(28px);
}
</style>
<style scoped> 
  /* 全局样式（可放在 main.css 或 tailwind.css） */
.scrollbar-thin::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

.scrollbar-thin::-webkit-scrollbar-thumb {
  background-color: #c1c1c1;
  border-radius: 4px;
}

.dark .scrollbar-thin::-webkit-scrollbar-thumb {
  background-color: #555;
}

.scrollbar-thin::-webkit-scrollbar-track {
  background-color: transparent;
}
</style>