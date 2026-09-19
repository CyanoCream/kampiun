# Kampiun — Live Score Platform

Aplikasi pencatatan skor & klasemen pertandingan apa pun — sepak bola, voli, badminton,
mancing, turnamen kampung sampai tingkat besar. Egaliter: satu tool untuk semua, offline
dan online.

## Status

- **Ini planning + scaffold.** Build dicicil bertahap (lihat `docs/roadmap.md`).
- Belum ada deploy.

## Struktur (rencana)

| Folder | Isi | Stack |
|---|---|---|
| `backend/` | API (auth, billing, cup, score engine, TrueSkill) | Go (DDD) |
| `web/` | Landing + dashboard admin + penyelenggara + live viewer | Vue 3 + Vite |
| `mobile/` | Portal + input skor wasit (tap-tap) | Kotlin Multiplatform |

## Dokumentasi

- `docs/brainstorm.md` — hasil brainstroming produk & keputusan.
- `docs/roadmap.md` — urutan bangun MVP.

## Catatan

Brand ini menggantikan "Laga/Liga" yang sudah dipakai pihak lain. Nama: **Kampiun**.
