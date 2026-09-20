<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({ id: String })

const data = ref(null)
const match = ref(null)
const err = ref('')
const loading = ref(true)
let ws = null

async function fetchMatch() {
  try {
    const r = await fetch(`/api/v1/matches/${props.id}`)
    if (!r.ok) throw new Error('pertandingan tidak ditemukan')
    match.value = await r.json()
  } catch (e) { err.value = e.message }
}

async function fetchScore() {
  try {
    const r = await fetch(`/api/v1/matches/${props.id}/score`)
    if (!r.ok) throw new Error('belum ada skor / pertandingan tidak ditemukan')
    data.value = await r.json()
    loading.value = false
  } catch (e) { if (!err.value) err.value = e.message }
}

function connectWS() {
  const proto = location.protocol === 'https:' ? 'wss' : 'ws'
  ws = new WebSocket(`${proto}://${location.host}/api/v1/matches/${props.id}/ws`)
  ws.onmessage = (ev) => {
    const sc = JSON.parse(ev.data)
    // trigger confetti hanya saat transisi ke finished
    if (sc.Finished && !data.value?.Finished) showConfetti()
    data.value = sc
  }
  ws.onclose = () => { if (!data.value?.Finished) setTimeout(connectWS, 2000) } // reconnect
}

// Confetti ringan: elemen CSS yang jatuh.
function showConfetti() {
  const colors = ['#533afd', '#ea2261', '#f96bee', '#15be53', '#f39c12']
  const wrap = document.createElement('div')
  wrap.className = 'confetti'
  document.body.appendChild(wrap)
  for (let i = 0; i < 60; i++) {
    const p = document.createElement('div')
    p.className = 'piece'
    p.style.left = Math.random() * 100 + '%'
    p.style.background = colors[i % colors.length]
    p.style.animationDelay = (Math.random() * 2) + 's'
    p.style.transform = `rotate(${Math.random() * 360}deg)`
    wrap.appendChild(p)
  }
  setTimeout(() => wrap.remove(), 4000)
}

onMounted(() => {
  fetchMatch()
  fetchScore()
  connectWS()
})
onUnmounted(() => { if (ws) ws.close() })
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
      <div class="side" :class="match?.home ? '' : 'empty'">
        <span class="name">{{ match?.home?.name ?? 'Tim A' }}</span>
        <span class="sets">Set: {{ data?.HomeWon }}</span>
        <span v-if="data?.Finished && data?.HomeWon > data?.AwayWon" class="winer">🏆 pemenang</span>
      </div>
      <div class="mid">
        <span class="big">{{ data?.HomeScore ?? '–' }}</span>
        <span class="sep">:</span>
        <span class="big">{{ data?.AwayScore ?? '–' }}</span>
        <span v-if="data?.Finished" class="done-label">SELESAI</span>
      </div>
      <div class="side" :class="match?.away ? '' : 'empty'">
        <span class="name">{{ match?.away?.name ?? 'Tim B' }}</span>
        <span class="sets">Set: {{ data?.AwayWon }}</span>
        <span v-if="data?.Finished && data?.AwayWon > data?.HomeWon" class="winer">🏆 pemenang</span>
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
.side { text-align: center; padding: 32px; border-radius: 16px; color: white; }
.side:first-child { background: #4f46e5; }
.side:last-child { background: #e74c3c; }
.side.empty { background: #888; }
.name { display: block; font-size: 20px; font-weight: 700; }
.sets { display: block; margin-top: 8px; opacity: 0.85; font-size: 14px; }
.mid { display: flex; align-items: center; gap: 8px; }
.big { font-size: 56px; font-weight: 800; }
.sep { font-size: 56px; color: #999; }
.winer { display: block; margin-top: 12px; font-weight: 700; }
.done-label { position: absolute; font-size: 13px; color: #2ecc71; font-weight: 700; margin-left: 12px; }
@keyframes pulse { 50% { opacity: 0.3; } }
.muted { color: #888; }

/* confetti */
.confetti { position: fixed; inset: 0; pointer-events: none; overflow: hidden; z-index: 999; }
.piece {
  position: absolute; top: -20px; width: 10px; height: 16px; border-radius: 2px;
  animation: fall 2.6s ease-in forwards;
}
@keyframes fall { to { transform: translateY(110vh) rotate(720deg); } }
</style>
