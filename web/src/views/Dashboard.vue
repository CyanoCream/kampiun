<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'

const email = ref(''), pass = ref('')
const token = ref(localStorage.getItem('kmp_token') || '')
const orgs = ref([])
const comps = ref([])
const msg = ref('')
const err = ref('')

// form org
const newOrgName = ref('')
// form kompetisi (utk org terpilih)
const selOrg = ref('')
const newComp = ref({ name: '', sport: 'volleyball', public: true })
// form peserta
const selComp = ref('')
const newParticipants = ref('')

async function api(path, opts = {}) {
  const h = { 'Content-Type': 'application/json', ...(opts.headers || {}) }
  if (token.value) h['Authorization'] = 'Bearer ' + token.value
  const r = await fetch(`/api/v1${path}`, { ...opts, headers: h })
  if (!r.ok) { const e = await r.json().catch(() => ({})); throw new Error(e.error || r.status) }
  return r.json()
}

async function login() {
  try {
    const r = await api('/auth/login', { method: 'POST', body: JSON.stringify({ email: email.value, password: pass.value }) })
    token.value = r.token; localStorage.setItem('kmp_token', r.token)
    msg.value = 'Masuk berhasil'; err.value = ''
    await load()
  } catch (e) { err.value = e.message }
}

async function load() {
  if (!token.value) return
  try { orgs.value = await api('/orgs') } catch (e) { orgs.value = [] }
  if (orgs.value.length && !selOrg.value) selOrg.value = orgs.value[0].id
  await loadComps()
}

async function loadComps() {
  if (!selOrg.value) { comps.value = []; return }
  try { comps.value = await api(`/orgs/${selOrg.value}/competitions`) } catch (e) { comps.value = [] }
  if (comps.value.length && !selComp.value) selComp.value = comps.value[0].id
}

async function createOrg() {
  try {
    const o = await api('/orgs', { method: 'POST', body: JSON.stringify({ name: newOrgName.value }) })
    newOrgName.value = ''; msg.value = `Organisasi '${o.name}' dibuat`; err.value = ''
    selOrg.value = o.ID
    await load()
  } catch (e) { err.value = e.message }
}

async function createComp() {
  try {
    const c = await api('/competitions', { method: 'POST', body: JSON.stringify({ ...newComp.value, org_id: selOrg.value }) })
    msg.value = `Kompetisi '${c.name}' dibuat — kode akses: ${c.access_code}`; err.value = ''
    await loadComps()
    selComp.value = c.ID
  } catch (e) { err.value = e.message }
}

async function addParticipants() {
  const names = newParticipants.value.split(',').map(s => s.trim()).filter(Boolean)
  if (names.length < 2) { err.value = 'Minimal 2 peserta'; return }
  try {
    await api(`/competitions/${selComp.value}/participants`, { method: 'POST', body: JSON.stringify({ names }) })
    newParticipants.value = ''; msg.value = `${names.length} peserta ditambahkan, bracket dibuat`
  } catch (e) { err.value = e.message }
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
      <p v-if="err" class="err">{{ err }}</p>
      <p v-if="msg" class="msg">{{ msg }}</p>
    </section>

    <section v-else class="content">
      <p v-if="msg" class="msg">{{ msg }}</p>
      <p v-if="err" class="err">{{ err }}</p>

      <h2>Organisasi</h2>
      <div class="form-inline">
        <input v-model="newOrgName" placeholder="nama organisasi (mis. Kampung Jaya)" @keyup.enter="createOrg" />
        <button class="btn primary" @click="createOrg">Buat</button>
      </div>
      <ul v-if="orgs.length" class="list">
        <li v-for="o in orgs" :key="o.id" class="row" :class="{ active: o.id === selOrg }" @click="selOrg=o.id;loadComps()">
          <strong>{{ o.name }}</strong> <span class="muted">/{{ o.slug }}</span>
        </li>
      </ul>
      <p v-else class="muted">Belum ada organisasi.</p>

      <template v-if="selOrg">
        <h2>Kompetisi</h2>
        <div class="form">
          <input v-model="newComp.name" placeholder="nama turnamen (mis. Voli Kampung 2026)" />
          <select v-model="newComp.sport">
            <option value="volleyball">Voli</option>
            <option value="badminton">Badminton</option>
            <option value="football">Sepak bola</option>
            <option value="fishing">Mancing</option>
            <option value="other">Lainnya</option>
          </select>
          <label class="chk"><input type="checkbox" v-model="newComp.public" /> Publik (bisa dicari)</label>
          <button class="btn primary" @click="createComp">Buat kompetisi</button>
        </div>

        <ul v-if="comps.length" class="list">
          <li v-for="c in comps" :key="c.id" class="row" :class="{ active: c.id === selComp }" @click="selComp=c.id">
            <strong>{{ c.name }}</strong> <span class="muted">• {{ c.sport }} • kode {{ c.access_code }}</span>
          </li>
        </ul>

        <template v-if="selComp">
          <h2>Peserta</h2>
          <div class="form-inline">
            <input v-model="newParticipants" placeholder="Pisahkan dgn koma: Garuda, Rajawali, Elang" @keyup.enter="addParticipants" />
            <button class="btn primary" @click="addParticipants">Tambah</button>
          </div>
          <p class="hint">Bracket knock-out dibuat otomatis dari jumlah peserta.</p>
          <RouterLink :to="`/live/${selComp}`" class="btn ghost lnk">Lihat live viewer</RouterLink>
        </template>
      </template>
    </section>
  </div>
</template>

<style scoped>
.dash { font-family: 'Inter', system-ui, sans-serif; max-width: 760px; margin: 0 auto; padding: 24px; color: #1a1a2e; }
.bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.brand { font-weight: 800; font-size: 18px; }
.auth { max-width: 320px; display: flex; flex-direction: column; gap: 10px; margin: 24px 0; }
input, select { padding: 12px 14px; border: 2px solid #e0e0e8; border-radius: 10px; font-size: 15px; }
.btn { padding: 12px 18px; border-radius: 10px; border: none; font-weight: 600; cursor: pointer; }
.btn.primary { background: #4f46e5; color: white; }
.btn.ghost { background: white; border: 2px solid #e0e0e8; color: #1a1a2e; text-decoration: none; }
.btn.lnk { margin-top: 10px; display: inline-block; }
.list { padding: 0; list-style: none; }
.row { background: #f7f7fb; border-radius: 12px; padding: 12px 16px; margin: 6px 0; display: flex; justify-content: space-between; align-items: center; cursor: pointer; border: 2px solid transparent; }
.row.active { border-color: #4f46e5; }
.form-inline { display: flex; gap: 10px; margin: 10px 0; }
.form-inline input { flex: 1; }
.form { display: flex; flex-direction: column; gap: 10px; margin: 10px 0; }
.chk { display: flex; align-items: center; gap: 6px; }
.hint { color: #888; font-size: 13px; }
.muted { color: #888; }
.err { color: #e74c3c; }
.msg { color: #2ecc71; }
h2 { margin-top: 28px; }
</style>
