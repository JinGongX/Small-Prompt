<template>
  <div >
  <section >
  <div class="flex items-center justify-between">
  <h2 class="text-lg font-bold text-gray-800 dark:text-white" >提示记录管理</h2> 
  <div>
    <input
    v-model="keyword"
    placeholder="搜索提示内容...Enther键搜索 "
    @keydown.enter.prevent="onSearch"
    class="w-64 px-3 py-1 text-sm rounded-lg border 
           bg-white dark:bg-neutral-800
           focus:outline-none focus:ring-2 focus:ring-blue-500/40"
     />
  </div>
   <span class="text-sm text-neutral-400">
        共 {{ tips.length }} 条
  </span>
      
  </div>
  <div class="h-full w-full ">
    <!-- Table -->
    <div class="overflow-hidden rounded-xl border border-neutral-200 dark:border-neutral-700 bg-white  dark:bg-gray-600 dark:divide-gray-700">
      <table class="w-full text-sm">
        <thead class="bg-neutral-100 dark:bg-gray-800 text-neutral-600 dark:text-neutral-300">
          <tr>
            <th class="px-4 py-3 text-center">标题</th>
            <th class="px-3">状态</th>
            <!-- <th class="px-3">置顶</th> -->
            <th class="px-3">创建时间</th>
            <th class="px-3 text-center">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="tip in tips"
            :key="tip.id"
            class="border-t border-neutral-100 dark:border-neutral-700 hover:bg-neutral-50 dark:hover:bg-neutral-700/50 transition"
          >
            <!-- Title -->
            <td class="px-4 py-3">
              <div v-if="editingId !== tip.id">
                <div class="font-medium text-neutral-800 dark:text-neutral-100">
                  {{ tip.title }}
                </div>
                <div class="text-xs text-neutral-400 truncate max-w-xs">
                  <!-- {{ tip.desc }} -->
                </div>
              </div>
              <div v-else class="space-y-2">
                <input
                  v-model="editCache.title"
                  class="w-full rounded-md border px-2 py-1 text-sm dark:bg-neutral-700"
                />
                <textarea
                  v-model="editCache.desc"
                  class="w-full rounded-md border px-2 py-1 text-sm dark:bg-neutral-700"
                  rows="2"
                />
              </div>
            </td>
            <!-- State -->
            <td class="text-center">
              <button
                @click="toggleState(tip)"
                class="text-xs px-2 py-1 rounded-full"
                :class="tip.state === 8
                  ? 'bg-green-100 text-green-700'
                  : 'bg-yellow-100 text-yellow-700'"
              >
                {{ tip.state === 8 ? '已完成' : '进行中' }}
              </button>
            </td>
            <!-- Pinned -->
            <!-- <td class="text-center">
              <button
                @click="togglePinned(tip)"
                class="text-lg transition"
              >
                {{ tip.pinned ? '📌' : '📍' }}
              </button>
            </td> -->
            <!-- Time -->
            <td class="text-center text-xs text-neutral-400">
              {{ formatTime(tip.createdAt || 0) }}
            </td>
            <!-- Actions -->
            <td class="px-3 text-right">
              <div class="flex items-center justify-end gap-2">
                <template v-if="editingId === tip.id">
                  <button @click="saveEdit(tip)" class="text-green-600 text-xs">
                    保存
                  </button>
                  <button @click="cancelEdit" class="text-neutral-400 text-xs">
                    取消
                  </button>
                </template>
                <template class="" v-else>
                  <!-- <button @click="startEdit(tip)" class="text-blue-600 text-xs">
                    编辑
                  </button> -->
                  <button
                    @click="removeTip(tip.id)"
                    class="w-8 text-red-500 text-xs border-2"
                  >
                    删除
                  </button>
                </template>
              </div>
            </td>
          </tr>

          <tr v-if="tips.length === 0">
            <td colspan="5" class="py-10 text-center text-neutral-400">
              暂无提示记录
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>

  </section></div>
</template>

<script setup lang="ts">
import { reactive, ref,onMounted } from 'vue'
import{GetTips,UpTipsState,UpTipsPinned } from  '../../../bindings/changeme/services/suistore'
import { InputSearch } from 'ant-design-vue'

const keyword = ref('')
interface PromptCard {
  id: number
  type: string //PromptType
  category: string
  title: string
  state: number
  pinned: number
  pinnedat?: number | null
  snoozeat?: number | null
  createdAt?: number
  _isNew?: boolean // 是否为新添加的卡片（用于触发动画）
}
const tips = ref<PromptCard[]>([])

// 编辑状态
const editingId = ref<number | null>(null)
const editCache = reactive({
  title: '',
  desc: '',
})

 
const fetchTips = async () => {
  tips.value = await GetTips();
};

const onSearch = () => {
  if (!keyword.value.trim()) {
    fetchTips()
    return
  }
  console.log('搜索关键词:', keyword.value);
  tips.value = tips.value.filter(tip => tip.title.includes(keyword.value));
}
onMounted(async() => {
  await fetchTips();  
 // Events.On('tipEvent', onTipEvent);
});

// ------- methods -------
function startEdit(tip: PromptCard) {
  editingId.value = tip.id
  editCache.title = tip.title
}

function cancelEdit() {
  editingId.value = null
}

function saveEdit(tip: PromptCard) {
  tip.title = editCache.title
 // tip.desc = editCache.desc
  editingId.value = null

  // TODO: 调用后端 update 接口
}

function removeTip(id: number) {
  tips.value = tips.value.filter(t => t.id !== id)
  // TODO: 调用后端 delete 接口
  UpTipsState(id, 0);
}

function toggleState(tip: PromptCard) {
  tip.state = tip.state === 1 ? 8 : 1
  // TODO: update state
}

function togglePinned(tip: PromptCard) {
  tip.pinned = tip.pinned === 1 ? 0 : 1
  tip.pinnedat = tip.pinned ? Date.now() / 1000 : undefined
  // TODO: update pinned
}

function formatTime(ts: number) {
  return new Date(ts * 1000).toLocaleString()
}
</script>

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