<script setup>
import { ref, onMounted } from 'vue'

const props = defineProps({ id: String })
const rounds = ref([])
const err = ref('')
const loading = ref(true)

async function load() {
  try {
    const r = await fetch(`/api/v1/competitions/${props.id}/bracket`)
    if (!r.ok) throw new Error('bracket tidak ditemukan')
    rounds.value = await r.json()
  } catch (e) { err.value = e.message }
  finally { loading.value = false }
}

onMounted(load)
</script>

<template>
  <div class="brk">
    <header class="head">
      <span class="brand">🏆 Kampiun</span>
      <h2>Bracket Turnamen</h2>
    </header>
    <p v-if="loading" class="muted">Memuat bracket…</p>
    <p v-else-if="err" class="err">{{ err }}</p>

    <div v-else class="bracket">
      <div v-for="r in rounds" :key="r.round" class="round">
        <div class="rhead" :class="{ final: r.round === rounds.length }">
          {{ r.round === rounds.length ? '🏆 FINAL' : 'Ronde ' + r.round }}
        </div>
        <div v-for="m in (r.matches?.length ? r.matches : [null])" :key="m ? m.match.id : r.round" class="match">
          <template v-if="m">
            <div class="slot" :class="{ win: m.match.winner_id === m.match.home_id }">
              <span class="name">{{ m.home_name || '—' }}</span>
              <span class="score">{{ m.match.status === 'finished' ? m.match.winner_id === m.match.home_id ? '✓' : '' : '' }}</span>
            </div>
            <div class="slot" :class="{ win: m.match.winner_id === m.match.away_id }">
              <span class="name">{{ m.away_name || '—' }}</span>
              <span class="score">{{ m.match.status === 'finished' && m.match.winner_id === m.match.away_id ? '✓' : '' }}</span>
            </div>
          </template>
          <template v-else>
            <div class="slot empty"><span class="name">—</span></div>
            <div class="slot empty"><span class="name">—</span></div>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.brk { font-family: 'Source Sans 3', system-ui, sans-serif; max-width: 1080px; margin: 0 auto; padding: 24px; color: #061b31; }
.head { display: flex; align-items: center; gap: 16px; margin-bottom: 24px; }
.brand { font-weight: 600; font-size: 18px; color: #533afd; }
h2 { font-weight: 300; letter-spacing: -0.5px; margin: 0; }
.bracket { display: flex; flex-wrap: wrap; gap: 28px; padding: 16px 0; }
.round { display: flex; flex-direction: column; gap: 16px; min-width: 200px; }
.rhead { font-size: 13px; font-weight: 600; letter-spacing: 1px; color: #64748d; text-transform: uppercase; margin-bottom: 6px; }
.rhead.final { color: #533afd; }
.match { display: flex; flex-direction: column; border: 1px solid #e5edf5; border-radius: 4px; overflow: hidden; box-shadow: rgba(50,50,93,0.12) 0px 4px 12px; }
.slot { display: flex; justify-content: space-between; align-items: center; padding: 8px 12px; font-size: 14px; }
.slot + .slot { border-top: 1px solid #eef2f7; }
.slot.win { background: rgba(83,58,253,0.06); font-weight: 600; }
.slot.empty { color: #bbb; }
.score { color: #533afd; font-weight: 700; }
.err { color: #e74c3c; }
.muted { color: #888; }
</style>
