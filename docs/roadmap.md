# Roadmap — Kampiun

Urutan bangun (dicicil, satu repo). MVP = alur inti: skor + live viewer jalan dulu.

## Fase 1 — Fondasi (MVP inti)

- [ ] Scaffold backend Go (DDD): auth JWT, config, migrasi, seeding superadmin.
- [ ] Domain inti: `Organization` (penyelenggara), `Competition`, `Participant` (tim/orang),
      `Match`, `Scoreboard` + `ScoreEvent` (log tap-tap).
- [ ] Score engine inti + template per jenis pertandingan (voli/badminton/mancing/sepakbola).
- [ ] Tap-tap scorer (mobile) — offline-first, simpan lokal, sync naik.
- [ ] Live viewer web dgn kode unik (private) + search publik + filter.
- [ ] Dashboard penyelenggara (web): kelola peserta, pertandingan, venue.

## Fase 2 — Cup & Bracket

- [ ] Buat cup: jumlah peserta bebas, seeding, nama.
- [ ] Bracket knock-out (single elimination) + byes otomatis.
- [ ] Visualisasi bracket game-like (SVG animasi, confetti).

## Fase 3 — Billing & Admin

- [ ] Billing 2 model: langganan & beli-cup (durasi, perpanjang).
- [ ] Notifikasi Telegram + approve/reject order dari bot.
- [ ] Super admin: override data, bypass pembayaran, laporan keuangan, dashboard berlangsung.

## Fase 4 — Advanced

- [ ] TrueSkill ranking (backend).
- [ ] Referral (reward TBD).
- [ ] IOS release (setelah Android stabil).

## Prinsip

- Offline dulu baru online; satu source of truth (device wasit).
- Bukan multi-master P2P; kalau perlu multi-device → host LAN.
- Satu backend untuk semua UI (web + mobile).
- Semua tech versi terbaru.
