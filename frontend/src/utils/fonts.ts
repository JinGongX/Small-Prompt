export type FontMode = 'soft' | 'warm' | 'focus' | 'emotion'

export const fontMap: Record<FontMode, string> = {
  soft: '"Nunito", sans-serif',
  warm: '"Fredoka", sans-serif',
  focus: '"Inter", sans-serif',
  emotion: '"Patrick Hand", cursive'
}

export const fontMeta = [
  { key: 'soft', name: '温柔模式' },
  { key: 'warm', name: '轻松模式' },
  { key: 'focus', name: '专注模式' },
  { key: 'emotion', name: '情绪模式' }
]