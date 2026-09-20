<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'

const email = ref(''), pass = ref('')
const token = ref(localStorage.getItem('kmp_token') || '')
const orgs = ref([])
const comps = ref([])
const msg = ref('')

async function api(path, opts = {}) {
  const h = { 'Content-Type': 'application/json', ...(opts.headers || {}) }
  if (token.value) h['Authorization'] = 'Bearer ' + token.value
  const r = await fetch(`/api/v1${path}`, { ...opts, headers: h })
  if (!r.ok) { const e = await r.json().catch(() => ({})); throw new Error(e.error || r.status) }
  return r.status === 204 ? null : r.json()
}

async function login() {
  try {
    const r = await api('/auth/login', { method: 'POST', body: JSON.stringify({ email: email.value, password: pass.value }) })
    token.value = r.token
    localStorage.setItem('kmp_token', r.token)
    msg.value = 'Masuk berhasil'
    await load()
  } catch (e) { msg.value = e.message }
}

async function load() {
  if (!token.value) return
  try { orgs.value = await api('/orgs') } catch (e) { console.log(e) }
  if (orgs.value.length) {
    comps.value = await api(`/orgs/${orgs.value[0].id}/competitions`).catch(() => [])
  }
}

onMounted(load)
</script>

<template>
  <div class="dash">
    <header class="bar">
      <span class="brand">🏆 Kampiun</span>
      <button v-if="token" class="btn ghost" @click="token='';localStorage.removeItem('kmp_token');orgs=[];comps=[]">Keluar</button>
    </header>

    <section v-if="!token" class="auth">
      <h2>Masuk penyelenggara</h2>
      <input v-model="email" placeholder="email" type="email" />
      <input v-model="pass" placeholder="password" type="password" @keyup.enter="login" />
      <button class="btn primary" @click="login">Masuk</button>
      <p v-if="msg" class="msg">{{ msg }}</p>
    </section>

    <section v-else class="content">
      <h2>Organisasi saya</h2>
      <ul v-if="orgs.length">
        <li v-for="o in orgs" :key="o.id" class="row">
          <span><strong>{{ o.name }}</strong> <span class="muted">({{ o.slug }})</span></span>
        </li>
      </ul>
      <p v-else class="muted">Belum ada organisasi. Buat lewat API dulu.</p>

      <h3>Kompetisi</h3>
      <ul v-if="comps.length">
        <li v-for="c in comps" :key="c.id" class="row">
          <span><strong>{{ c.name }}</strong> <span class="muted">• {{ c.sport }}</span></span>
          <RouterLink :to="`/live/${c.id}/`" class="btn ghost sm">Lihat</RouterLink>
        </li>
      </ul>
      <p v-else class="muted">Belum ada kompetisi.</p>
    </section>
  </div>
</template>

<style scoped>
.dash { font-family: 'Inter', system-ui, sans-serif; max-width: 760px; margin: 0 auto; padding: 24px; color: #1a1a2e; }
.bar { display: flex; justify-content: space-between; align-items: center; }
.brand { font-weight: 800; font-size: 18px; }
.auth { max-width: 320px; display: flex; flex-direction: column; gap: 10px; margin-top: 40px; }
input { padding: 12px 14px; border: 2px solid #e0e0e8; border-radius: 10px; font-size: 15px; }
.btn { padding: 12px 18px; border-radius: 10px; border: none; font-weight: 600; cursor: pointer; }
.btn.primary { background: #4f46e5; color: white; }
.btn.ghost { background: white; border: 2px solid #e0e0e8; }
.btn.sm { font-size: 13px; padding: 6px 12px; }
.row { background: #f7f7fb; border-radius: 12px; padding: 14px 18px; margin: 8px 0; display: flex; justify-content: space-between; align-items: center; }
.muted { color: #888; }
.msg { color: #4f46e5; }
h2, h3 { margin-top: 24px; }
</style>
