import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    id("org.springframework.boot") version "3.0.2"
    id("io.spring.dependency-management") version "1.1.0"
    kotlin("jvm") version "1.7.22"
    kotlin("plugin.spring") version "1.7.22"
    kotlin("plugin.serialization") version "1.7.20"
}

group = "fi.oplab.findyagency"

version = "0.0.1-SNAPSHOT"

java.sourceCompatibility = JavaVersion.VERSION_17

repositories {
    mavenCentral()
    maven {
        name = "GitHubPackages"
        url = uri("https://maven.pkg.github.com/findy-network/findy-common-kt")
        credentials {
            username = System.getenv("USERNAME")
            password = System.getenv("TOKEN")
        }
    }
}

dependencies {
    implementation("org.findy_network:findy-common-kt-client:0.0.13")
    implementation("org.findy_network:findy-common-kt-stub:0.0.13")
    implementation("com.google.zxing:core:3.5.1")
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.4.1")
    implementation("org.jetbrains.kotlin:kotlin-reflect")
    implementation("org.jetbrains.kotlin:kotlin-stdlib-jdk8")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-reactor")
    implementation("org.reactivestreams:reactive-streams")
    implementation("org.springframework.boot:spring-boot-starter-web")
    testImplementation("org.springframework.boot:spring-boot-starter-test")
}

tasks.withType<KotlinCompile> {
    kotlinOptions {
        freeCompilerArgs = listOf("-Xjsr305=strict")
        jvmTarget = "17"
    }
}

tasks.withType<Test> { useJUnitPlatform() }
