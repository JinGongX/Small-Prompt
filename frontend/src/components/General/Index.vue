<script setup lang="ts">
import { ref,onMounted,watch } from 'vue'
import ListItem from '../Setting/ListRow.vue' 
import LanguageSwitcher from '../Setting/LanguageSwitcher.vue'
import ThemeSetting from '../Setting/ThemeSetting.vue'
import {  useSettingsStore } from '../../utils/settings'

const modelValue = ref(false)
const isSoundEnabled = ref(false)
const duration = ref(10) //const durationOptions = [5, 10, 15, 20, 30, 60]
const completion = ref(10) //const completionOptions = [5, 10, 15, 20, 30, 60]
const settings = useSettingsStore()
onMounted(async() => { 
  const saved = localStorage.getItem('sound-enabled')
  isSoundEnabled.value = saved === 'true'
  
  const savedDuration = localStorage.getItem('duration')
  if (savedDuration) {
    console.log('Loaded duration from localStorage:', savedDuration)
    duration.value = parseInt(savedDuration, 10)
  }
  const savedCompletion = localStorage.getItem('completion')
  if (savedCompletion) {
    completion.value = parseInt(savedCompletion, 10)
  }
})
watch(isSoundEnabled, (newVal) => {
  //localStorage.setItem('sound-enabled', String(newVal))
  settings.toggleSound(newVal)
})
watch(duration, (newVal) => {
  //localStorage.setItem('duration', String(newVal))
  settings.toggleDuration(newVal)
    // 通知子窗口也可以加上 BroadcastChannel
   const bc = new BroadcastChannel('settings')
   bc.postMessage({type:'duration',value:newVal})
})
watch(completion, (newVal) => {
  //localStorage.setItem('completion', String(newVal))
  settings.toggleCompletion(newVal)
})

</script>

<template > 
<div >
  <section >
    <!-- 应用设置 -->
  <h2 class="text-lg font-bold text-gray-800 dark:text-white" >{{ $t('components.general.title.tips') }}</h2>
      <div class="bg-white rounded-lg shadow divide-y divide-gray-200 dark:bg-gray-800 dark:divide-gray-700">
          <ListItem :label="$t('components.general.label.prompt_time')" :subLabel="$t('components.general.subLabel.sb_prompt_time')">
            <a-select  class="dark:text-white"
              v-model:value="duration"
              style="width: 60px" 
            >
              <a-select-option :value="5" >5</a-select-option>
              <a-select-option :value="10">10</a-select-option>
              <a-select-option :value="15">15</a-select-option>
              <a-select-option :value="20">20</a-select-option>
              <a-select-option :value="30">30</a-select-option>
              <a-select-option :value="60">60</a-select-option>
            </a-select>
          </ListItem>

          <ListItem :label="$t('components.general.label.completion_time')" :subLabel="$t('components.general.subLabel.sb_completion_time')">
            <a-select  class="dark:text-white"
              ref="select"
              v-model:value="completion"
              style="width: 60px" 
            >
              <a-select-option :value="5" >5</a-select-option>
              <a-select-option :value="10">10</a-select-option>
              <a-select-option :value="15">15</a-select-option>
              <a-select-option :value="20">20</a-select-option>
              <a-select-option :value="30">30</a-select-option>
              <a-select-option :value="60">60</a-select-option>
            </a-select>
          </ListItem>
      </div>
      <!-- 应用设置 -->
  <!-- <h2 class="text-lg font-bold text-gray-800 dark:text-white" >{{ $t('components.general.title.application') }}</h2>
      <div class="bg-white rounded-lg shadow divide-y divide-gray-200 dark:bg-gray-800 dark:divide-gray-700">
          <ListItem :label="$t('components.general.label.startup')" subLabel="">
            <input type="checkbox" class="sr-only peer" v-model="modelValue" />
            <div
              class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-500 rounded-full peer dark:bg-gray-300 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-500 relative"
            ></div>
          </ListItem>
      </div> -->
      <!-- 音效设置 -->
      <h2 class="text-lg font-bold text-gray-800 dark:text-white" >{{ $t('components.general.title.notify') }}</h2>
      <div class="bg-white rounded-lg shadow divide-y divide-gray-200 dark:bg-gray-800 dark:divide-gray-700">
         <ListItem :label="$t('components.general.label.notify_copy')" subLabel="">
            <input type="checkbox" class="sr-only peer" v-model="isSoundEnabled" />
            <div
              class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-500 rounded-full peer dark:bg-gray-300 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-500 relative"
            ></div>
          </ListItem>
      </div> 
    <!-- 外观设置 -->
        <h2 class="text-lg font-bold text-gray-800 dark:text-white"  >{{ $t('components.general.title.exterior') }}</h2>
      <div class="bg-white rounded-lg shadow divide-y divide-gray-200 dark:bg-gray-800 dark:divide-gray-700">
         <ListItem :label="$t('components.general.label.language')" subLabel="">
            <LanguageSwitcher />
          </ListItem>
          <ListItem :label="$t('components.general.label.theme')" subLabel="">
            <ThemeSetting />
          </ListItem>
      </div>
<!-- 应用更新 -->
     <!-- <h2 class="text-lg font-bold text-gray-800 dark:text-white" >{{ $t('components.general.title.update') }}</h2>
      <div class="mt-4 bg-white rounded-lg shadow divide-y divide-gray-200 dark:bg-gray-800 dark:divide-gray-700">
         <ListItem :label="$t('components.general.label.automatic_up')" subLabel="">
            <input type="checkbox" class="sr-only peer" v-model="modelValue" />
            <div
              class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-500 rounded-full peer dark:bg-gray-300 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-500 relative"
            ></div>
          </ListItem>
          <ListItem :label="$t('components.general.label.next_up')" subLabel="">
            <input type="checkbox" class="sr-only peer" v-model="modelValue" />
            <div
              class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-500 rounded-full peer dark:bg-gray-300 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-500 relative"
            ></div>
          </ListItem>
      </div>  -->


  </section>
  </div>
</template>
<style scoped>
 h2{
    text-align: left;
    margin:15px 15px 4px 15px;
    font-size:15px;
 }
 .rg_desc{
    text-align: left;
    margin:4px 0px 4px 8px;
    font-size:13px;
    color:#999;
 }
</style>