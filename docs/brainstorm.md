# Brainstorming — Kampiun

> Status: keputusan produk, belum implementasi.

## Produk

Web + mobile pencatatan skor pertandingan umum (sepak bola, voli, badminton, mancing, dll).
Satu engine skor inti + konfigurasi ("template") per jenis pertandingan — bukan per-olahraga
masing-masing. Egaliter: sama-sama dipakai turnamen kampung sampai nasional, offline atau online.

## Deliverable (5)

1. **Landing page** — web (Vue UI / Go template).
2. **Dashboard admin (super admin)** — web. Akses penuh: override tim/orang (data bawaan),
   bypass pembayaran, laporan keuangan, dashboard pertandingan berlangsung.
3. **Dashboard penyelenggara** (komunitas/kampung) — web + mobile.
4. **Mobile** — portal admin + portal penyelenggara + **input skor wasit** (layar "tap-tap":
   blok besar Tim A / Tim B, tap = +1, long-press = menu minus/set/perhatian, tampilan
   optimistik, simpan lokal dulu, sinkron naik belakangan).
5. **Live score viewer** — web. Akses **private** via kode unik pendek (mis `SC-8K2M`),
   atau **publik** via search + filter (cup/venue/lokasi).

## Cup / Bracket

- Bracket knock-out gaya game: SVG garis nyambung pemenang, animasi pindah round, glow,
  confetti di juara. UI/UX interaktif seperti game.
- **Custom**: jumlah peserta bebas, nama cup, seeding/unggulan.
- **Beli-cup**: bayar sekali → durasi cup (mis 2 bulan), bisa diperpanjang dgn syarat.
  Akses tetap selama cup berlangsung.

## Billing

- **2 model**: langganan (subscription) & **beli-cup** (one-time, durasi).
- `Plan` punya `kind` = `subscription` | `one_time`. Order diproses → buat `Cup` + `CupAccess`
  dgn `expires_at`.
- **Admin bypass pembayaran** — aktivasi manual tanpa order (pola un-di, lifecycle).

## Referral

- User bawa kode referal ajak daftar / beli → reward/komisi. Detail reward & aturan
  BELUM digali (apa reward, kapan cair, syarat).

## Offline / Online

- Venue tertutup (kampung) jalan **tanpa internet**: satu device = satu wasit → device itu
  source of truth (keputusan kunci). Online = sinkron NAik (device → cloud), bukan multi-master.
- Bila perlu beberapa device edit live di 1 venue → host LAN (WebSocket ke host), BUKAN mesh P2P.
- Online memungkinkan: backup, akses dari luar, klasemen lintas-venue.

## Ranking — TrueSkill

- Rating naik/turun dari kekuatan lawan → ranking tak monoton (menang atas kuat naik banyak,
  kalah dari lemah turun drastis). Implementasi di **backend** (satu sumber rating), mobile
  hanya tampil.

## Stack

- **Backend**: Go — arsitektur DDD (pola repo `un-di`), JWT auth, notifikasi Telegram +
  approve/reject order dari bot, billing pola undangan.
- **Web**: Vue 3 + Vite (Tailwind).
- **Mobile**: Kotlin Multiplatform — Android-first dulu; iOS menyusul.
- **Rust**: OPSIONAL, batasi ke modul kecil (mis lib TrueSkill) — bukan backend, bukan full core.
- Semua pakai versi terbaru.

## Pertanyaan Belum Terjawab

1. iOS target sekarang atau nanti? (nentuin Android-first vs KMP penuh)
2. Rust sungguh perlu atau sekadar "keren"?
3. Satu wasit per pertandingan — konfirmasi akhir?
4. Reward referral: apa & kapan cair?
