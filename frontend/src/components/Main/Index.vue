<template>
  <div class="h-screen w-screen   flex flex-col  dark:text-black dark:text-white">
    <div class="flex flex-1    overflow-hidden">
      <!-- 左侧请求栏 -->
      <div class="w-44 bg-gray-100/20 p-3 space-y-2 font-bold text-base dark:bg-gray-800" style="padding-top:40px">
        <div v-for="(item, index) in requests" :key="index" @click="handleMenu(item)" 
             :class="[' cursor-pointer p-2 rounded text-left',selected === item.id ? 'bg-orange-300/90' : 'hover:bg-gray-300/70']">
            <component :is="item.icon" :style="['margin-right: 10px;vertical-align: middle;',item.id==='shortcut'?'font-size: 19px':'font-size: 18px']" />
            <span>{{ item.label }}</span>
        </div>
         <!-- <button @click="openSecond" style="width:90%;text-align: left;"  >  
           <SlackOutlined style="vertical-align: middle;margin-right: 10px;"/><span>轻提示</span></button> 
           <button @click="openTips" style="width:90%;text-align: left;"  >   
          <SendOutlined style="vertical-align: middle;margin-right: 10px;"  />
          <span>写提示</span></button>  -->
      </div>

      <!-- 主内容 -->
      <div class="flex-1 p-4 overflow-y-auto max-h-screen scroll-container bg-gray-50 dark:bg-gray-950 scrollbar-thin">
        <component :is="getComponent"  />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref ,computed,watch ,nextTick,onMounted,onBeforeUnmount} from 'vue'
import { useI18n } from 'vue-i18n'
import { HistoryOutlined,SettingOutlined,BulbOutlined ,FormOutlined,SlackOutlined ,SendOutlined } from '@ant-design/icons-vue';
import  ShortcutOutlinedicon   from '../Setting/ShortcutOutlined.vue'; 
import { useSettingsStore } from '../../utils/settings'
import {Events} from "@wailsio/runtime"; 

const settings = useSettingsStore()
const { t } = useI18n()

type MenuItem =
  | {
      id: string
      label: string
      icon: any
      type: 'component'
    }
  | {
      id: string
      label: string
      icon: any
      type: 'window'
      action: () => void
    }

const requests = computed<MenuItem[]>(() => [
  {id:'general',label:t('menus.general'),icon:SettingOutlined,type: 'component'},
  {id:'shortcut',label:t('menus.shortcut'),icon:ShortcutOutlinedicon,type: 'component'},
  {id:'history',label:t('menus.history'),icon:HistoryOutlined,type: 'component'}, 
   // ====== 窗口类 ======
  {
    id: 'light-tip',
    label: t('menus.wetips'),
    icon: SlackOutlined,
    type: 'window',
    action: OpenSecondWindow
  },
  {
    id: 'write-tip',
    label: t('menus.writeprompts'),
    icon: SendOutlined,
    type: 'window',
    action: OpenTipsWindow
  },

  {id:'about',label:t('menus.about'),icon:BulbOutlined,type: 'component'},
])
const selected = ref('general')
import General from '../General/Index.vue'
import Shortcut from '../Shortcut/Index.vue'
import About from '../About/Index.vue'
import History from '../History/Index.vue'

const components = {
  history: History,
  general: General,
  shortcut: Shortcut,
  about: About
}

const getComponent = computed(() => components[selected.value])
// 菜单变化时滚动到顶部
watch(selected, () => {
  nextTick(() => {
    const el = document.querySelector('.scroll-container')
    if (el) el.scrollTop = 0
  })
})

function handleMenu(item: MenuItem) {
  if (item.type === 'component') {
    selected.value = item.id
  }

  if (item.type === 'window') {
    item.action()
  }
}

let audio: HTMLAudioElement
audio = new Audio('/assets/notify.mp3') // 或写成 '/notify.mp3'
const handleClipboardChanged= (event: any)=>{
   // 判断是否开启音效
  if (settings.soundEnabled) {
    audio.currentTime = 0
    audio.play().catch(err => console.warn("音频播放失败:", err))
  } 
}
let unsubscribe: () => void
onMounted(async() => {  
  // const saved = localStorage.getItem('sound-enabled')
  // isSoundEnabled.value = saved === 'true'
  settings.loadSettings() // 载入持久化设置
  unsubscribe=Events.On('tipEvent', (event: any) => {
    handleClipboardChanged(event)
  }); 
})

onBeforeUnmount(() => {
  unsubscribe?.()
})


import { OpenSecondWindow, OpenTipsWindow } from  '../../../bindings/changeme/appservice'
//import { OpenTTTTSecondWindow } from '../../../bindings/changeme/appservice'
// function openSecond() {
//   OpenSecondWindow()
// }
// function openTips() {
//   OpenTipsWindow()
// }
</script>

<style scoped>
html, body {
  margin: 0;
  overflow-x: hidden; /* 禁止横向滚动 */
  overflow-y: hidden;
}
.drag-region {
  -webkit-app-region: drag;
}

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