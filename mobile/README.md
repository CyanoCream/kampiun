# Kampiun Mobile — Kotlin Multiplatform

Scaffold KMP: layar **tap-tap skor wasit** + client API ke backend Kampiun.
Target awal Android; iOS menyusul.

## Struktur

```
mobile/
  settings.gradle.kts
  build.gradle.kts
  gradle/libs.versions.toml   (version catalog)
  composeApp/                 (shared Compose Multiplatform + Android app)
    src/
      commonMain/kotlin/...   (kode shared: api, viewmodel, UI tap-score)
      androidMain/kotlin/...  (Android entry)
```

## Build

Di dev machine (kamu, ada Android Studio):

```bash
gradlew :composeApp:assembleDebug
```

Catatan: VPS ini hanya bisa compile target shared (metadata) utk cek error source.
Android full build butuh Android SDK — jalankan lokal.

## Alur

- Wasit login → pilih match → layar **dua blok besar** (home/away), ketuk = skor +1,
  long-press = menu (minus/undo, selesai set), simpan lokal-optimistik, sync naik ke API.
