<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({ id: String })

const data = ref(null)
const err = ref('')
const loading = ref(true)
let timer = null

async function fetchScore() {
  try {
    const r = await fetch(`/api/v1/matches/${props.id}/score`)
    if (!r.ok) throw new Error('match tidak ditemukan atau belum ada skor')
    data.value = await r.json()
    err.value = ''
  } catch (e) {
    err.value = e.message
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchScore()
  timer = setInterval(fetchScore, 3000)
})
onUnmounted(() => clearInterval(timer))
</script>

<template>
  <div class="live">
    <header class="head">
      <span class="dot" :class="{ done: data?.Finished }"></span>
      <span class="lbl">{{ data?.Finished ? 'SELESAI' : 'LIVE' }}</span>
      <span class="code muted">kode: {{ id }}</span>
    </header>

    <p v-if="loading" class="muted">Mengambil skor…</p>
    <p v-else-if="err" class="err">{{ err }}</p>

    <div v-else class="board">
      <div class="side home">
        <span class="name">Tim A</span>
        <span class="sets">Set: {{ data?.HomeWon }}</span>
      </div>
      <div class="mid">
        <span class="big">{{ data?.HomeUnits?.at(-1) ?? '–' }}</span>
        <span class="sep">:</span>
        <span class="big">{{ data?.AwayUnits?.at(-1) ?? '–' }}</span>
      </div>
      <div class="side away">
        <span class="name">Tim B</span>
        <span class="sets">Set: {{ data?.AwayWon }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.live { font-family: 'Inter', system-ui, sans-serif; max-width: 720px; margin: 0 auto; padding: 24px; color: #1a1a2e; }
.head { display: flex; gap: 12px; align-items: center; margin-bottom: 16px; }
.dot { width: 12px; height: 12px; border-radius: 50%; background: #e74c3c; animation: pulse 1.2s infinite; }
.dot.done { background: #2ecc71; animation: none; }
.lbl { font-weight: 700; letter-spacing: 1px; }
.code { font-size: 13px; }
.err { color: #e74c3c; }
.board { display: grid; grid-template-columns: 1fr auto 1fr; align-items: center; gap: 16px; }
.side { text-align: center; padding: 32px; border-radius: 16px; color: white; }
.side.home { background: #4f46e5; }
.side.away { background: #e74c3c; }
.name { display: block; font-size: 20px; font-weight: 700; }
.sets { display: block; margin-top: 8px; opacity: 0.85; font-size: 14px; }
.mid { display: flex; align-items: center; gap: 8px; }
.big { font-size: 56px; font-weight: 800; }
.sep { font-size: 56px; color: #999; }
@keyframes pulse { 50% { opacity: 0.3; } }
.muted { color: #888; }
</style>
