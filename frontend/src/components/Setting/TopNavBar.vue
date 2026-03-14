<template>
  <div class="top-nav w-[294px] bg-white/90  dark:bg-gray-500/80  " >
    <!-- 顶部进度条 -->
    <!-- <div v-if="showProgress" class="progress-bar"></div> -->
    <div class="action-row ">
       <div class="pending  bg-cyan-300/10 dark:bg-cyan-600/20 text-sky-800 dark:text-neutral-200 [--pulse-color:14,165,233] 
         dark:[--pulse-color:56,189,248]" 
           :class="{ active: props.pendingCount > 0 }"
           @click="filterPending">
        <span class="dot" :class="props.pendingCount===0?'bg-slate-400  dark:bg-slate-300':'bg-orange-400  dark:bg-orange-300'"></span>
       {{ props.pendingCount }} {{ $t('pages.home.TopNavBar.pending') }}
      </div>
      <div class="text-xs text-gray-900 dark:text-gray-200"  v-if="saved">
        {{ props.altertitle }}
      </div>

      <div class="actions">
        <button @click="onFilter" ><svg t="1772461059462" class="icon " :class="filterType==='all'?'fill-gray-900 dark:fill-gray-200':'fill-orange-900 dark:fill-orange-200'" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="12696" width="18" height="18"><path d="M516.973714 950.857143a67.584 67.584 0 0 1-30.427428-7.606857l-111.616-53.248a64.073143 64.073143 0 0 1-38.034286-58.368V499.346286v-2.56L100.937143 179.712a65.609143 65.609143 0 0 1-5.046857-71.021714c10.166857-22.893714 32.914286-35.547429 58.368-35.547429h611.328c25.380571 0 48.274286 12.653714 58.368 35.547429 10.166857 22.820571 10.166857 48.201143-5.12 68.461714L582.948571 494.226286v388.169143c0 22.820571-12.653714 43.154286-30.427428 55.808a55.442286 55.442286 0 0 1-35.547429 12.653714zM154.185143 133.997714c-2.486857 0-5.046857 0-5.046857 2.56v5.12l235.958857 317.074286a70.217143 70.217143 0 0 1 12.653714 40.594286v332.288c0 2.56 0 5.12 2.56 5.12l114.102857 53.248h5.12c2.56 0 2.56-2.56 2.56-5.12V499.346286c0-15.213714 5.046857-27.940571 12.653715-40.594286l235.885714-317.147429c2.56-2.486857 2.56-5.046857 0-5.046857 0-2.56-2.486857-2.56-5.046857-2.56H154.331429z m740.790857 482.011429H717.385143a29.257143 29.257143 0 0 1-30.427429-30.427429 29.257143 29.257143 0 0 1 30.427429-30.427428h177.590857a29.257143 29.257143 0 0 1 30.427429 30.427428 29.257143 29.257143 0 0 1-30.427429 30.427429z m0 126.829714H717.385143a29.257143 29.257143 0 0 1-30.427429-30.427428 29.257143 29.257143 0 0 1 30.427429-30.427429h177.590857a29.257143 29.257143 0 0 1 30.427429 30.427429 29.257143 29.257143 0 0 1-30.427429 30.427428z m0 126.829714H717.385143a29.257143 29.257143 0 0 1-30.427429-30.427428 29.257143 29.257143 0 0 1 30.427429-30.427429h177.590857a29.257143 29.257143 0 0 1 30.427429 30.427429 29.257143 29.257143 0 0 1-30.427429 30.427428z"  fill-opacity=".65" p-id="12697"></path><path d="M483.254857 949.76q16.822857 8.411429 33.718857 8.411429 22.162286 0 39.936-14.189715 33.353143-23.917714 33.353143-61.586285V496.64l234.422857-315.099429q24.868571-33.206857 5.924572-75.849142Q812.909714 65.828571 765.586286 65.828571H154.331429q-47.177143 0-64.950858 39.497143-22.528 40.96 5.851429 78.701715l234.422857 315.172571v332.434286q0 46.08 42.276572 64.950857l111.396571 53.174857z m33.718857 1.097143a55.442286 55.442286 0 0 0 35.547429-12.653714 69.705143 69.705143 0 0 0 30.427428-55.808V494.226286l235.885715-317.074286c15.286857-20.260571 15.286857-45.641143 5.12-68.461714-10.166857-22.893714-32.987429-35.547429-58.368-35.547429H154.331429c-25.380571 0-48.274286 12.653714-58.368 35.547429a65.609143 65.609143 0 0 0 5.12 70.948571l235.885714 317.147429v334.848c0 25.380571 15.213714 48.201143 38.034286 58.368l111.616 53.248c10.166857 5.12 20.333714 7.606857 30.427428 7.606857z m248.685715-816.859429c2.486857 0 4.973714 0 4.973714 2.56 2.56 0 2.56 2.56 0 5.12l-1.462857 2.048-0.146286 0.146286L534.674286 458.752a76.653714 76.653714 0 0 0-12.653715 40.594286v385.536c0 2.56 0 5.12-2.56 5.12h-5.12l-114.102857-53.321143c-2.56 0-2.56-2.486857-2.56-5.12V499.419429a70.217143 70.217143 0 0 0-12.653714-40.594286L149.138286 141.604571v-5.046857c0-2.56 2.56-2.56 5.12-2.56h611.328z m-607.597715 7.314286h603.721143L528.676571 454.656l-0.219428 0.292571q-13.677714 22.893714-13.677714 44.397715v382.683428l-109.714286-51.2V499.346286q0-26.112-14.116572-44.982857l-232.886857-313.051429z m559.323429 482.011429h177.590857q16.310857 0 27.062857-10.678858 10.678857-10.752 10.678857-27.062857 0-37.741714-37.741714-37.741714H717.385143q-37.741714 0-37.741714 37.741714t37.741714 37.741715z m0-7.314286h177.590857a29.257143 29.257143 0 0 0 30.427429-30.427429 29.257143 29.257143 0 0 0-30.427429-30.427428H717.385143a29.257143 29.257143 0 0 0-30.427429 30.427428 29.257143 29.257143 0 0 0 30.427429 30.427429z m0 134.144h177.590857q37.741714 0 37.741714-37.741714t-37.741714-37.741715H717.385143q-16.310857 0-27.062857 10.678857-10.678857 10.678857-10.678857 27.062858 0 37.741714 37.741714 37.741714z m0-7.314286h177.590857a29.257143 29.257143 0 0 0 30.427429-30.427428 29.257143 29.257143 0 0 0-30.427429-30.427429H717.385143a29.257143 29.257143 0 0 0-30.427429 30.427429 29.257143 29.257143 0 0 0 30.427429 30.427428z m0 134.144h177.590857q37.741714 0 37.741714-37.741714 0-16.310857-10.678857-27.062857-10.752-10.678857-27.062857-10.678857H717.385143q-16.310857 0-27.062857 10.678857-10.678857 10.752-10.678857 27.062857 0 37.741714 37.741714 37.741714z m0-7.314286h177.590857a29.257143 29.257143 0 0 0 30.427429-30.427428 29.257143 29.257143 0 0 0-30.427429-30.427429H717.385143a29.257143 29.257143 0 0 0-30.427429 30.427429 29.257143 29.257143 0 0 0 30.427429 30.427428z"  fill-opacity=".65" p-id="12698"></path></svg></button>
        <!-- <button @click="toggleView"><svg t="1772458939405" class="icon fill-gray-900 dark:fill-gray-200" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="10344" width="14" height="14"><path d="M379.854 528.516c63.85 0 115.608 51.76 115.608 115.614v264.258c0 63.85-51.76 115.612-115.608 115.612H115.608C51.76 1024 0 972.24 0 908.388v-264.26c0-63.85 51.76-115.612 115.608-115.612z m594.6 330.322c27.364 0 49.546 22.184 49.546 49.55 0 25.084-18.64 45.814-42.824 49.096l-6.722 0.452h-330.72c-27.364 0-49.548-22.184-49.548-49.548 0-25.086 18.64-45.816 42.824-49.096l6.722-0.454h330.72z m-594.6-231.226H115.608a16.516 16.516 0 0 0-16.516 16.518v264.258a16.516 16.516 0 0 0 16.516 16.516h264.246a16.516 16.516 0 0 0 16.516-16.516v-264.26a16.516 16.516 0 0 0-16.516-16.516z m594.6-165.16C1001.818 462.452 1024 484.636 1024 512c0 25.084-18.64 45.816-42.824 49.096l-6.722 0.452h-330.72c-27.364 0-49.548-22.184-49.548-49.548 0-25.084 18.64-45.816 42.824-49.096l6.722-0.452h330.72zM379.854 0c63.85 0 115.608 51.76 115.608 115.612v264.26c0 63.85-51.76 115.612-115.608 115.612H115.608C51.76 495.484 0 443.724 0 379.87V115.612C0 51.762 51.76 0 115.608 0z m0 99.096H115.608a16.516 16.516 0 0 0-16.516 16.516v264.26a16.516 16.516 0 0 0 16.516 16.516h264.246a16.516 16.516 0 0 0 16.516-16.518V115.612a16.516 16.516 0 0 0-16.516-16.516z m594.6-33.032c27.364 0 49.546 22.184 49.546 49.548 0 25.086-18.64 45.816-42.824 49.096l-6.722 0.454h-330.72c-27.364 0-49.548-22.184-49.548-49.55 0-25.084 18.64-45.814 42.824-49.096l6.722-0.452h330.72z"   p-id="10345"></path></svg></button> -->
        <!-- <button><svg t="1772608567745" class="icon fill-gray-900 dark:fill-gray-200" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="20743" width="18" height="18"><path d="M256 96A160 160 0 0 0 96 256v512A160 160 0 0 0 256 928h512a160 160 0 1 0 0-320h-139.477333l143.189333-143.189333a160 160 0 1 0-226.304-226.261334L416 368V256A160 160 0 0 0 256 96z m145.194667 739.370667l163.328-163.370667H768a96 96 0 0 1 0 192H384c6.698667-8.874667 12.458667-18.474667 17.194667-28.629333z m325.248-415.786667L416 730.026667v-271.530667l174.677333-174.677333a96 96 0 1 1 135.765334 135.765333zM256 864A96 96 0 0 1 160 768V256a96 96 0 0 1 192 0v512A96 96 0 0 1 256 864zM256 810.666667a42.666667 42.666667 0 1 0 0-85.333334 42.666667 42.666667 0 0 0 0 85.333334z"  p-id="20744"></path></svg></button> -->
        <button  @click="OpenTipsWindow" ><svg t="1772461138954" class="icon fill-gray-900 dark:fill-gray-200" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="14007" width="24" height="24"><path d="M503 839.1H219.6c-42.3 0-76.6-34.3-76.6-76.6V197.1c0-43.1 34.9-78 78-78h591.7c46.7 0 78.1 32.7 78.1 81.5v30.2c0 17.5-14.2 31.7-31.6 31.7-17.5 0-31.7-14.2-31.7-31.7v-30.3c0-18.1-6.8-18.1-14.7-18.1H221c-7.9 0-14.6 6.7-14.6 14.7v565.3c0 7.3 6 13.3 13.3 13.3h282.7c17.3 0 31.8 13.5 32.3 30.8 0.4 17.9-13.9 32.6-31.7 32.6z"   p-id="14008"></path><path d="M861.3 659.1H507.2c-17.5 0-31.7-14.2-31.7-31.7s14.2-31.7 31.7-31.7h354.1c17.5 0 31.7 14.2 31.7 31.7-0.1 17.5-14.2 31.7-31.7 31.7zM596.1 344.3H320.2c-17.5 0-31.7-14.2-31.7-31.7s14.2-31.7 31.7-31.7h275.9c17.5 0 31.7 14.2 31.7 31.7s-14.2 31.7-31.7 31.7zM440.8 486H320.2c-17.5 0-31.7-14.2-31.7-31.7s14.2-31.7 31.7-31.7h120.6c17.5 0 31.7 14.2 31.7 31.7 0 17.6-14.2 31.7-31.7 31.7z"   p-id="14009"></path><path d="M684.2 836.1c-17.5 0-31.7-14.2-31.7-31.7v-354c0-17.5 14.2-31.7 31.7-31.7s31.7 14.2 31.7 31.7v354c0 17.5-14.2 31.7-31.7 31.7z"   p-id="14010"></path></svg></button>
      </div>
    </div>

    <a-drawer
    title=""
    :placement="placement"
    :closable="false"
    width="120"
    height="80"
    :open="open"
    @close="onClose"
    :bodyStyle="{ padding: '4px' }" 
    class="bg-cyan-100  dark:bg-gray-400"
  >
   <div class="px-1 py-1">
      <!-- 标题 -->
      <div class="text-xs font-semibold text-gray-400 mb-4 tracking-wide dark:text-gray-100 ">
        筛选类型
      </div>
      <!-- 单选 -->
      <a-radio-group
        v-model:value="filterType"
        class="flex flex-col gap-1"
        @change="onChange"
      >
        <a-radio value="all">
          <span class="text-xs text-neutral-900 dark:text-white ">默认</span>
        </a-radio>
        <a-radio value="scheduled">
          <span class="text-xs text-neutral-900 dark:text-white ">定时提示</span>
        </a-radio>
        <a-radio value="immediate">
          <span class="text-xs text-neutral-900 dark:text-white">即时提示</span>
        </a-radio>
      </a-radio-group>
       <a-divider style="border-color: #ccc" dashed />
    </div>
  </a-drawer>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted} from "vue"
import { OpenTipsWindow } from  '../../../bindings/changeme/services/appservice'
import type { DrawerProps } from 'ant-design-vue';

const placement = ref<DrawerProps['placement']>('right');
const open = ref<boolean>(false);
const filterType = ref<'all' | 'scheduled' | 'immediate'>('all')
const syncing = ref(false)
const updateAvailable = ref(false)
const saved = ref(true)
const pendingCount = ref(0)
const showProgress = ref(false)
const emit = defineEmits<{
  (e: 'filter-change', value: 'all' | 'scheduled' | 'immediate'): void
}>()

function triggerSaved() {
  saved.value = true
  showProgress.value = true

  setTimeout(() => {
    saved.value = false
  }, 1200)

  setTimeout(() => {
    showProgress.value = false
  }, 800)
}
 
function onChange() {
   emit('filter-change', filterType.value)
  // 选择后自动关闭
  open.value = false
}
// function onFilter() {}
function toggleView() {}
function filterPending() {
    //emit('filter-pending', pendingCount.value) 
}

const onFilter = () => {
  open.value = true;
};
const onClose = () => {
  open.value = false;
};
const props = defineProps<{
  pendingCount: number,
  altertitle: string,
}>()


</script>

<style scoped>
.top-nav {
  /* width: 100%; */
  border-radius: 16px;
  /* background: rgba(255,255,255,0.85); */
  backdrop-filter: blur(20px);
  padding: 2px 4px;
  position: relative;
  box-shadow: 0 2px 12px rgba(0,0,0,0.04);
}


/* 顶部滑动条 */
/* .progress-bar {
  position: absolute;
  top: 0;
  left: 0;
  height: 2px;
  width: 100%;
  background: linear-gradient(to right,#3b82f6,#60a5fa);
  animation: slide 0.8s ease;
} */

/* @keyframes slide {
  from { transform: translateX(-100%); }
  to { transform: translateX(0); }
} */

.status-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  font-weight: 500;
}

.syncing {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #3b82f6;
}

.spinner {
  width: 12px;
  height: 12px;
  border: 2px solid #3b82f6;
  border-top: 2px solid transparent;
  border-radius: 50%;
  animation: rotate 1s linear infinite;
}

@keyframes rotate {
  to { transform: rotate(360deg); }
}

.update {
  color: #f59e0b;
}

.saved {
  color: #22c55e;
}

.pending {
  font-size: 11px;
  padding: 2px 6px 2px 6px;
  border-radius: 20px;
  /* background: rgba(59,130,246,0.1); */
  /* color: #3b82f6; */
  display: flex;
  /* position: relative; */
  align-items: center;
  gap: 4px;
  cursor: pointer;
  transition: all 0.3s;
}
.pending {
  --pulse-color: 59,130,246; /* 默认浅色模式 */
}

.dark .pending {
  --pulse-color: 56,189,248; /* 暗黑模式颜色 */
}
.pending.active {
  animation: pulse 2s infinite;
}

/* @keyframes pulse {
  0% { box-shadow: 0 0 0 0 rgba(59,130,246,0.4); }
  70% { box-shadow: 0 0 0 6px rgba(59,130,246,0); }
  100% { box-shadow: 0 0 0 0 rgba(59,130,246,0); }
} */

@keyframes pulse {
  0% { box-shadow: 0 0 0 0 rgba(var(--pulse-color), 0.4); }
  70% { box-shadow: 0 0 0 6px rgba(var(--pulse-color), 0); }
  100% { box-shadow: 0 0 0 0 rgba(var(--pulse-color), 0); }
}

/* .pending.active {
    animation: glow 2s ease-in-out infinite;
} */


/* @keyframes glow {
  0%,100% { 
    box-shadow: 0 0 0 0 rgba(56,189,248,0.2);
    transform: scale(1);
  }
  50% { 
    box-shadow: 0 0 10px 3px rgba(56,189,248,0.5);
    transform: scale(1.03);
  }
} */

.dot {
  width: 8px;
  height: 8px;
  /* background: #3b82f6; */
  border-radius: 50%;
}

.action-row {
  /* margin-top: 8px; */
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-right: 8px;
}

.actions {
  display: flex;
  gap: 12px;
}

.actions  button{
    width:20px;
    display: flex;
  align-items: center;    /* 垂直居中 */
  justify-content: center; /* 水平居中 */
}

.add-btn {
  margin-left: -4px;
  /* background: #3b82f6;
  color: white;
  border-radius: 12px;
  width: 26px;
  height: 26px;
  border: none;
  font-size: 18px; */
}

/** */

</style>