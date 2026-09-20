<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const code = ref('')
function goCode() {
  if (code.value.trim()) router.push(`/live/${code.value.trim()}`)
}
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
  font-family: 'Inter', system-ui, sans-serif;
  color: #1a1a2e;
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 24px;
}
.topbar {
  display: flex; align-items: center; justify-content: space-between;
  padding: 20px 0;
}
.brand { font-weight: 800; font-size: 20px; letter-spacing: 0.5px; }
.topbar-links { display: flex; gap: 24px; align-items: center; }
.muted { color: #666; text-decoration: none; }
.hero { text-align: center; padding: 96px 0 72px; }
.hero h1 { font-size: 52px; line-height: 1.1; max-width: 700px; margin: 0 auto 16px; }
.lede { font-size: 18px; color: #555; max-width: 560px; margin: 0 auto 32px; line-height: 1.6; }
.code-box { display: flex; gap: 12px; justify-content: center; max-width: 460px; margin: 0 auto; }
input {
  flex: 1; padding: 14px 16px; border: 2px solid #e0e0e8; border-radius: 12px;
  font-size: 16px; outline: none;
}
input:focus { border-color: #4f46e5; }
.btn { padding: 14px 22px; border-radius: 12px; font-weight: 600; cursor: pointer; border: none; text-decoration: none; font-size: 15px; }
.btn.primary { background: #4f46e5; color: white; }
.btn.ghost { border: 2px solid #e0e0e8; background: white; color: #1a1a2e; }
.grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 16px; padding: 32px 0 64px; }
.card { border: 1px solid #eaeaf2; border-radius: 16px; padding: 24px; background: #fff; }
.card h3 { margin: 8px 0 8px; font-size: 17px; }
.card p { color: #667; font-size: 14px; line-height: 1.5; margin: 0; }
footer { text-align: center; padding: 40px 0; font-size: 14px; }
</style>
