plugins {
    id("org.jetbrains.kotlin.multiplatform") version "2.0.21"
}

repositories {
    mavenCentral()
}

kotlin {
    jvm()               // host-side utk menjalankan/test common tanpa Android SDK
    // Android target diaktifkan di dev machine (butuh Android SDK)
    // androidTarget()

    sourceSets {
        commonMain.dependencies {
            implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.9.0")
            implementation("io.ktor:ktor-client-core:3.0.1")
            implementation("io.ktor:ktor-client-content-negotiation:3.0.1")
            implementation("io.ktor:ktor-serialization-kotlinx-json:3.0.1")
            implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.7.3")
        }
        commonTest.dependencies {
            implementation(kotlin("test"))
        }
        jvmMain.dependencies {
            implementation("io.ktor:ktor-client-okhttp:3.0.1")
        }
    }
}

kotlin.sourceSets.all {
    languageSettings.optIn("kotlinx.coroutines.ExperimentalCoroutinesApi")
}
