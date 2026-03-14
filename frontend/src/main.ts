import { createApp } from 'vue'
import Antd from 'ant-design-vue';
import App from './App.vue'
import './main.css'  
import 'ant-design-vue/dist/reset.css';
import router from './router/index'
import {i18n,setupI18n } from './utils/i18n'
import { initTheme } from './utils/ThemeManager'//add
import { createPinia } from 'pinia'
import {GetLanguage} from  '../bindings/changeme/services/suistore' 

initTheme()
if (window.location.hash === '' || window.location.hash === '#/') {
  router.replace('/second')
}
async function bootstrap(){
    await setupI18n('zh')//默认语言或从后端读取
    createApp(App).use(Antd).use(router).use(i18n).use(createPinia()).mount('#app')
    router.isReady().then(async () => {
        const lang = await GetLanguage() //从后端读取语言设置 lang as any||
        console.log("当前语言:", lang)
        await setupI18n(lang as any|| 'zh')
    })
}
bootstrap()