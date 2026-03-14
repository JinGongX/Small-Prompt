<template>
  <div class="flex items-center justify-center h-full px-6">
    <div class="max-w-[420px] w-full text-center">

      <!-- 插图 -->
      <!-- <div class="flex justify-center mb-7">
        <div
          class="w-12 h-12 rounded-full flex items-center justify-center
                 bg-blue-500/5 dark:bg-white/5
                 animate-floatSoft"
        >
          
          <svg
            class="w-12 h-12 text-blue-500 dark:text-blue-400"
            fill="none"
            stroke="currentColor"
            stroke-width="1.6"
            viewBox="0 0 24 24"
          >
            <circle cx="12" cy="13" r="7" />
            <path stroke-linecap="round" d="M12 13V9" />
            <path stroke-linecap="round" d="M12 13L15 15" />
            <path stroke-linecap="round" d="M5 4L3 6" />
            <path stroke-linecap="round" d="M19 4L21 6" />
          </svg>
        </div>
      </div> -->

      <!-- 主标题 -->
      <h2 class="text-[20px] font-semibold text-gray-800 dark:text-white mb-2">
        从第一个轻提示开始 ✨
      </h2>

      <!-- 副标题 -->
      <p class="text-sm text-gray-500 dark:text-white/50 mb-6 leading-relaxed">
        让未来的你，收到现在的善意提醒
      </p>

      <!-- 按钮 -->
      <button 
        @click="OpenTipsWindow"
        class="w-40 inline-flex items-center gap-2 px-5 h-9
               bg-[#007AFF] hover:bg-[#0066E0] active:bg-[#0055C5]
               text-white text-sm font-medium
               rounded-[10px]
               shadow-sm hover:-translate-y-[1px]
               transition-all duration-200"
      >
        <span class="text-base leading-none">＋</span>
        创建提示
      </button>

      <!-- 辅助说明 -->
      <!-- <div class="mt-4 text-xs text-gray-400 dark:text-white/30">
        支持定时 · 分类
      </div> -->

    </div>
  </div>
</template>

<script setup lang="ts">
import { OpenTipsWindow } from  '../../../bindings/changeme/services/appservice'
// defineEmits(['create'])
</script>

<style scoped>
@keyframes floatSoft {
  0% { transform: scale(1); }
  50% { transform: scale(1.04); }
  100% { transform: scale(1); }
}

.animate-floatSoft {
  animation: floatSoft 4s ease-in-out infinite;
}
</style>