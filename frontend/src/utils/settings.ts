//settings.ts
import { defineStore } from 'pinia'

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    soundEnabled: true, // 默认开启音效
    //OCREnabled:true,//默认开启图片ocr
    //AutoStartEnabled:false,//默认不开机启动
    duration: 10, //默认提示时间为10秒 仅定时提示使用
    completion: 10, //默认完成时间为10秒  仅定时提示使用
  }),
  actions: {
    toggleSound(enabled: boolean) {
      this.soundEnabled = enabled
      localStorage.setItem('sound-enabled', String(enabled))
    },
    loadSettings() {
      const saved = localStorage.getItem('sound-enabled')
      this.soundEnabled = saved === 'true'

      const savedDuration = localStorage.getItem('duration')
      if (savedDuration) {
        this.duration = parseInt(savedDuration, 10)
      }

      const savedCompletion = localStorage.getItem('completion')
      if (savedCompletion) {
        this.completion = parseInt(savedCompletion, 10)
      }
      //ocr
    //   const savocred = localStorage.getItem('ocr-enabled')
    //   this.OCREnabled = savocred === 'true'
    //   //AutoStart
    //   const savedautostart = localStorage.getItem('autostart-enabled')
    //   this.AutoStartEnabled = savedautostart === 'true'
    },
      toggleDuration(value: number) {
        this.duration = value
        localStorage.setItem('duration', String(value))
      },
      toggleCompletion(value: number) {
        this.completion = value
        localStorage.setItem('completion', String(value))
      },
    //  toggleOCR(enabled: boolean) {
    //   this.OCREnabled = enabled
    //   localStorage.setItem('ocr-enabled', String(enabled))
    // },
    // toggleAutoStart(enabled: boolean) {
    //   this.AutoStartEnabled = enabled
    //   localStorage.setItem('autostart-enabled', String(enabled))
    // },
  }
})