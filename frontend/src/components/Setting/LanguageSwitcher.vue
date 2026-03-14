<template>
  <div >
    <a-select  class="dark:text-white"
              ref="select"
              v-model:value="locale"
              style="width: 100px"
            >
              <a-select-option value="zh" >简体中文</a-select-option>
              <a-select-option value="en">English</a-select-option>
              <a-select-option value="zh-HK">繁體中文</a-select-option>
            </a-select>
  </div>
</template>

<script setup lang="ts">
import { ref ,computed} from 'vue'
import { i18n, setupI18n } from '../../utils/i18n'
import { Locale } from '../../locales'
import { SetLanguage } from '../../../bindings/changeme/services/appservice'
import { langChannel } from '../../utils/langChannel'
// const locale = ref(i18n.global.locale.value)
const locale = computed({
  get: () => i18n.global.locale.value,
  set: (val) => {
    switchLang(val as Locale)
  }
})
// const bc = new BroadcastChannel('language')
const switchLang =async  (lang: Locale) => {
  await setupI18n(lang)//(locale.value as Locale)
  await SetLanguage(lang)// 更新配置项中的语言设置 and 更新菜单语言
   // 通知子窗口也可以加上 BroadcastChannel
   langChannel.postMessage(lang)//(locale.value)

  //await LoadNewAppMenu(lang)//(locale.value) // 更新菜单语言
};
</script>