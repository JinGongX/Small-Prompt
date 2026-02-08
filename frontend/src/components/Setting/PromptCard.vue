<!-- ---------------- PromptCard ---------------- -->
<script lang="ts">
import { defineComponent, computed } from 'vue'

export default defineComponent({
  name: 'PromptCard',
  props: {
    card: {
      type: Object as () => {
        id: number
        type: string  //'life' | 'work' | 'device' | 'security' | 'system'
        title: string
      },
      required: true,
    },
  },
  emits: ['close', 'click'],
  setup(props, { emit }) {
    const visual = computed(() => {
      const map = {
        life: { bg: 'bg-amber-50', image: '../../public/image.png' },
        work: { bg: 'bg-blue-50', image: '../../public/image2.png' },
        device: { bg: 'bg-teal-50', image: '../../public/image2.png' },
        security: { bg: 'bg-emerald-50', image: '../../public/image3.png' },
        system: { bg: 'bg-neutral-100', image: '../../public/image3.png' },
        default: { bg: 'bg-neutral-100', image: '../../public/image3.png' },
      }
      return map[props.card.type]
    })

    const handleClose = () => {
      emit('close', props.card.id)
    }

    const handleClick = () => {
      emit('click', props.card)
    }

    return { visual, handleClose, handleClick }
  },
})
 
</script>

<template #default>
  <div @click="handleClick"
    class="group flex w-[304px] gap-3 rounded-2xl bg-white/80 backdrop-blur-md px-1 py-1 shadow-[0_8px_24px_rgba(0,0,0,0.12)] transition"
  >
    <!-- Visual Block -->
    <div
      class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl"
      :class="visual.bg"
    >
      <img
        :src="visual.image"
        class="h-10 w-10 object-contain"
        draggable="false"
      />
    </div>

    <!-- Content -->
    <div class="flex-1 text-sm py-1  text-neutral-900 line-clamp-2 break-words" :title="card.title">
      {{ card.title }}
    </div>

    <!-- Close -->
     <div  class="flex h-10 w-8">
         <button    @click.stop="handleClose"
      class="opacity-0 group-hover:opacity-100 text-neutral-400 transition hover:text-neutral-600"
    >✓
    </button>
     </div>
  </div>
</template>
