<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, RouterLink } from 'vue-router'

const router = useRouter()
const code = ref('')
const query = ref('')
const sport = ref('')
const results = ref([])
const searching = ref(false)

function goCode() { if (code.value.trim()) router.push(`/live/${code.value.trim()}`) }

async function search() {
  searching.value = true
  try {
    const qs = new URLSearchParams()
    if (query.value) qs.set('q', query.value)
    if (sport.value) qs.set('sport', sport.value)
    const r = await fetch(`/api/v1/competitions/public?${qs}`)
    if (!r.ok) throw new Error()
    results.value = await r.json()
  } catch (e) { results.value = [] }
  finally { searching.value = false }
}
onMounted(search)
</script>

<template>
  <div class="landing">
    <header class="topbar">
      <div class="brand">🏆 Kampiun</div>
      <div class="topbar-links">
        <a href="#fitur" class="muted">Fitur</a>
        <a href="#bracket" class="muted">Bracket</a>
        <a href="#harga" class="muted">Harga</a>
        <RouterLink :to="{ path: '/dashboard' }" class="btn ghost">Masuk</RouterLink>
      </div>
    </header>

    <section class="hero">
      <h1>Skor pertandingan, milik semua.</h1>
      <p class="lede">
        Voli, badminton, sepak bola, mancing — catat skor, rakit bracket, tunjukkan
        live score. Dari kampung sampai turnamen besar. Offline atau online.
      </p>
      <div class="code-box">
        <input
          v-model="code"
          placeholder="Masukkan kode live score (mis. A7K2)"
          @keyup.enter="goCode"
        />
        <button @click="goCode" class="btn primary">Lihat Live</button>
      </div>
    </section>

    <section id="explore" class="explore">
      <h2>Jelajahi turnamen</h2>
      <div class="toolbar">
        <input v-model="query" placeholder="Cari nama turnamen…" @keyup.enter="search" />
        <select v-model="sport" @change="search">
          <option value="">Semua cabang</option>
          <option value="volleyball">Voli</option>
          <option value="badminton">Badminton</option>
          <option value="football">Sepak bola</option>
          <option value="fishing">Mancing</option>
        </select>
        <button class="btn primary" @click="search">{{ searching ? '…' : 'Cari' }}</button>
      </div>
      <div v-if="results.length" class="res">
        <div v-for="c in results" :key="c.id" class="resrow">
          <div>
            <strong>{{ c.name }}</strong>
            <span class="muted"> • {{ c.sport }} • kode {{ c.access_code }}</span>
          </div>
          <RouterLink :to="`/bracket/${c.id}`" class="btn ghost sm">Bracket</RouterLink>
        </div>
      </div>
      <p v-else-if="searching" class="muted">Mencari…</p>
      <p v-else class="muted">Belum ada turnamen publik.</p>
    </section>

    <section id="fitur" class="grid">
      <div class="card"><h3>⚡ Tap-tap skor</h3><p>Wasit ketuk blok tim, skor langsung. Nol latensi, simpan lokal.</p></div>
      <div class="card"><h3>🪜 Bracket cup</h3><p>Knock-out otomatis, custom jumlah peserta & seeding.</p></div>
      <div class="card"><h3>🔒 Live private</h3><p>Kode unik utk undangan — atau publik lewat pencarian.</p></div>
      <div class="card"><h3>📴 Offline dulu</h3><p>Jalan tanpa internet di kampung. Sinkron kalau ada jaringan.</p></div>
    </section>

    <footer class="muted">© {{ new Date().getFullYear() }} Kampiun — Egaliter, untuk semua pertandingan.</footer>
  </div>
</template>

<style scoped>
.landing {
  font-family: 'Source Sans 3', system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
  font-feature-settings: 'ss01';
  color: #061b31;
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 24px;
}
.topbar {
  display: flex; align-items: center; justify-content: space-between;
  position: sticky; top: 0; background: rgba(255,255,255,0.85);
  backdrop-filter: blur(12px); padding: 20px 0; z-index: 10;
}
.brand { font-weight: 600; font-size: 20px; letter-spacing: -0.3px; }
.topbar-links { display: flex; gap: 24px; align-items: center; }
.muted { color: #64748d; text-decoration: none; font-weight: 400; }
.hero { text-align: center; padding: 96px 0 72px; }
.hero h1 { font-size: 48px; font-weight: 300; line-height: 1.15; letter-spacing: -0.96px; max-width: 700px; margin: 0 auto 16px; }
.lede { font-size: 18px; font-weight: 300; color: #64748d; max-width: 560px; margin: 0 auto 32px; line-height: 1.4; }
.code-box { display: flex; gap: 12px; justify-content: center; max-width: 460px; margin: 0 auto; }
input {
  flex: 1; padding: 12px 16px; border: 1px solid #e5edf5; border-radius: 4px;
  font-size: 16px; font-family: inherit; outline: none; color: #061b31;
}
input:focus { border-color: #533afd; box-shadow: 0 0 0 2px rgba(83,58,253,0.15); }
.btn { padding: 12px 20px; border-radius: 4px; font-weight: 400; cursor: pointer; text-decoration: none; font-size: 15px; display: inline-block; }
.btn.primary { background: #533afd; color: #fff; }
.btn.primary:hover { background: #4434d4; }
.btn.ghost { border: 1px solid #b9b9f9; background: transparent; color: #533afd; }
.btn.ghost:hover { background: rgba(83,58,253,0.05); }
.grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 16px; padding: 32px 0 64px; }
.card {
  border: 1px solid #e5edf5; border-radius: 6px; padding: 24px; background: #fff;
  box-shadow: rgba(50,50,93,0.25) 0px 30px 45px -30px, rgba(0,0,0,0.1) 0px 18px 36px -18px;
}
.card h3 { margin: 8px 0 8px; font-size: 22px; font-weight: 300; letter-spacing: -0.22px; color: #061b31; }
.card p { color: #64748d; font-size: 16px; font-weight: 300; line-height: 1.4; margin: 0; }
footer { text-align: center; padding: 40px 0; font-size: 14px; }
.explore { padding: 16px 0 48px; }
.explore h2 { font-weight: 300; letter-spacing: -0.5px; margin-bottom: 8px; }
.toolbar { display: flex; gap: 10px; margin: 12px 0 20px; flex-wrap: wrap; }
.toolbar input { flex: 1; min-width: 180px; }
.toolbar select { padding: 12px 14px; border: 1px solid #e5edf5; border-radius: 4px; font-family: inherit; font-size: 15px; }
.res { display: flex; flex-direction: column; gap: 10px; }
.resrow { display: flex; justify-content: space-between; align-items: center; border: 1px solid #e5edf5; border-radius: 6px; padding: 14px 16px; background: #fff; box-shadow: rgba(50,50,93,0.08) 0px 4px 14px; }
.btn.sm { font-size: 13px; padding: 6px 12px; }
</style>
