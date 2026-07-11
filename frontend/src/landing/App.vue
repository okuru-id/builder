<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { IconMail, IconBrandGithub, IconBrandLinkedin, IconBrandInstagram, IconBrandWhatsapp, IconBrandThreads, IconCode, IconBrandDocker, IconBrain, IconRoute, IconStar, IconLicense } from '@tabler/icons-vue'

const language = ref<'en' | 'id'>('en')
const copied = ref(false)
const template = ref('default')

const icomMap: Record<string, any> = { IconCode, IconBrain, IconRoute, IconBrandDocker }

const tmpl = {
  default: '',
  minimal: 'template-minimal',
  bold: 'template-bold',
}

const hero = ref({
  greeting_en: "Hi, I'm Kurob",
  greeting_id: 'Halo, Saya Kurob',
  description_en: 'Kurob is a fullstack programmer from Cilegon with 10+ years of experience building robust Laravel, React, DevOps, and automation solutions for modern businesses.',
  description_id: 'Kurob adalah fullstack programmer dari Cilegon dengan pengalaman 10+ tahun membangun solusi Laravel, React, DevOps, dan automation untuk kebutuhan bisnis modern.',
  profile_image: '/images/profile.webp',
  profile_image_mobile: '/images/profile-half.webp',
  cta_text: 'Contact Me',
  cta_link: '#contact',
})

const clientLogos = ref([
  { name: 'Danareksa', logo: '/images/danareksa.png' },
  { name: 'Krakatau Medika', logo: '/images/km.png' },
  { name: 'CIMB NIAGA', logo: '/images/cimb-niaga.png' },
  { name: 'IASH', logo: '/images/iash.png' },
  { name: 'Krakatau Steel', logo: '/images/ks.png' },
  { name: 'Krakatau Tirta Industri', logo: '/images/kti.png' },
  { name: 'Pemkot Cilegon', logo: '/images/pemkot-cilegon.png' },
  { name: 'Yayasan Pendidikan Telkom', logo: '/images/ypt.png' },
])

const services = ref([
  { icon: 'IconCode', title: 'Fullstack App', description_en: 'Building end-to-end web and mobile applications with modern frameworks and scalable architecture.', description_id: 'Membangun aplikasi web dan mobile end-to-end dengan framework modern dan arsitektur scalable.' },
  { icon: 'IconBrain', title: 'AI Integration Project', description_en: 'Integrating AI capabilities into existing systems to enhance automation and decision-making.', description_id: 'Mengintegrasikan kemampuan AI ke dalam sistem yang ada untuk meningkatkan otomatisasi dan pengambilan keputusan.' },
  { icon: 'IconRoute', title: 'Automation Workflow System', description_en: 'Designing workflow automation to streamline business processes and improve efficiency.', description_id: 'Merancang otomatisasi alur kerja untuk menyederhanakan proses bisnis dan meningkatkan efisiensi.' },
  { icon: 'IconBrandDocker', title: 'DevOps', description_en: 'Setting up CI/CD pipelines, containerization, and infrastructure for reliable deployments.', description_id: 'Menyiapkan pipeline CI/CD, containerization, dan infrastruktur untuk deployment yang andal.' },
])

const openSourceProjects = ref([
  {
    title: { en: 'Resident Fee Management System', id: 'Sistem Pengelolaan Iuran Warga' },
    description: { en: 'A modern solution for managing resident fee payments designed to facilitate transparent recording, automatic calculations, and guaranteed security with modern encryption.', id: 'Solusi modern untuk manajemen pembayaran iuran warga yang dirancang untuk memudahkan pencatatan, perhitungan, dan pelacakan iuran secara transparan dan efisien.' },
    githubUrl: 'https://github.com/kurob1993/kelola',
    technologies: [
      { name: 'Laravel', bgColor: 'bg-red-500/20', textColor: 'text-red-400' },
      { name: 'PHP', bgColor: 'bg-purple-500/20', textColor: 'text-purple-400' },
      { name: 'PostgreSQL', bgColor: 'bg-indigo-500/20', textColor: 'text-indigo-400' },
      { name: 'Tailwind CSS', bgColor: 'bg-cyan-500/20', textColor: 'text-cyan-400' },
    ],
    stats: { stars: 5, license: 'MIT License' },
  },
  {
    title: { en: 'SQL Server Docker Backup Script', id: 'Script Backup SQL Server Docker' },
    description: { en: 'A comprehensive bash script for automating SQL Server database backups in Docker environments.', id: 'Script bash komprehensif untuk mengotomatisasi backup database SQL Server dalam lingkungan Docker.' },
    githubUrl: 'https://github.com/kurob1993/sql-server-dockerize-backup',
    technologies: [
      { name: 'Docker', bgColor: 'bg-blue-500/20', textColor: 'text-blue-400' },
      { name: 'Bash', bgColor: 'bg-orange-500/20', textColor: 'text-orange-400' },
      { name: 'SQL Server', bgColor: 'bg-gray-500/20', textColor: 'text-gray-400' },
      { name: 'DevOps', bgColor: 'bg-green-500/20', textColor: 'text-green-400' },
    ],
    stats: { language: 'Shell 100%', type: 'Backup Automation' },
  },
  {
    title: { en: 'SonarQube Report Generator', id: 'Generator Laporan SonarQube' },
    description: { en: 'A Go-based web application for SonarQube report generation with Jenkins integration.', id: 'Aplikasi web berbasis Go untuk pembuatan laporan SonarQube dengan integrasi Jenkins.' },
    githubUrl: 'https://github.com/okuru-id/sonar-report',
    technologies: [
      { name: 'Go', bgColor: 'bg-cyan-500/20', textColor: 'text-cyan-400' },
      { name: 'Docker', bgColor: 'bg-blue-500/20', textColor: 'text-blue-400' },
      { name: 'SonarQube', bgColor: 'bg-green-500/20', textColor: 'text-green-400' },
      { name: 'Jenkins', bgColor: 'bg-red-500/20', textColor: 'text-red-400' },
    ],
    stats: { stars: 0, license: 'MIT License' },
  },
  {
    title: { en: 'Tactile Match Game', id: 'Permainan Tactile Match' },
    description: { en: 'A fun and interactive tactile matching game built with JavaScript.', id: 'Permainan mencocokkan tactile yang seru dan interaktif dibangun dengan JavaScript.' },
    githubUrl: 'https://github.com/okuru-id/tactile-match-game',
    technologies: [
      { name: 'JavaScript', bgColor: 'bg-yellow-500/20', textColor: 'text-yellow-400' },
      { name: 'HTML', bgColor: 'bg-orange-500/20', textColor: 'text-orange-400' },
      { name: 'CSS', bgColor: 'bg-blue-500/20', textColor: 'text-blue-400' },
    ],
    stats: { language: 'JavaScript 100%', type: 'Game' },
  },
])

const customSections = ref<Record<string, any>[]>([])

const knownTypes = ['hero', 'clients', 'services', 'projects', 'cta']

const cta = ref({
  heading: 'Tell me about your next project',
  email: 'kurob@okuru.id',
  whatsapp: '628999702143',
})

const techColors: Record<string, { bg: string; text: string }> = {
  Laravel: { bg: 'bg-red-500/20', text: 'text-red-400' },
  PHP: { bg: 'bg-purple-500/20', text: 'text-purple-400' },
  PostgreSQL: { bg: 'bg-indigo-500/20', text: 'text-indigo-400' },
  'Tailwind CSS': { bg: 'bg-cyan-500/20', text: 'text-cyan-400' },
  Docker: { bg: 'bg-blue-500/20', text: 'text-blue-400' },
  Bash: { bg: 'bg-orange-500/20', text: 'text-orange-400' },
  'SQL Server': { bg: 'bg-gray-500/20', text: 'text-gray-400' },
  DevOps: { bg: 'bg-green-500/20', text: 'text-green-400' },
  Go: { bg: 'bg-cyan-500/20', text: 'text-cyan-400' },
  SonarQube: { bg: 'bg-green-500/20', text: 'text-green-400' },
  Jenkins: { bg: 'bg-red-500/20', text: 'text-red-400' },
  JavaScript: { bg: 'bg-yellow-500/20', text: 'text-yellow-400' },
  HTML: { bg: 'bg-orange-500/20', text: 'text-orange-400' },
  CSS: { bg: 'bg-blue-500/20', text: 'text-blue-400' },
}

const techStyle = (name: string) => techColors[name] ?? { bg: 'bg-gray-500/20', text: 'text-gray-400' }

const serviceIcon = (name: string) => icomMap[name] ?? IconCode

onMounted(async () => {
  try {
    const [landingRes, settingsRes] = await Promise.all([
      fetch('/api/v1/landing'),
      fetch('/api/v1/settings/public'),
    ])
    const json = await landingRes.json()
    const data = json.data ?? {}
    try {
      const settings = await settingsRes.json()
      if (settings.data?.landing_template) template.value = settings.data.landing_template
    } catch { /* ignore */ }    if (data.hero) Object.assign(hero.value, data.hero)
    if (data.clients?.items) clientLogos.value = data.clients.items
    if (data.services?.items) services.value = data.services.items.map((s: any) => ({ ...s, icon: s.icon || 'IconCode' }))
    if (data.projects?.items) {
      openSourceProjects.value = data.projects.items.map((p: any) => ({
        title: { en: p.title_en, id: p.title_id },
        description: { en: p.description_en, id: p.description_id },
        githubUrl: p.github_url,
        technologies: (p.technologies ?? []).map((t: string) => ({ name: t, ...techStyle(t) })),
        stats: { stars: 0, license: 'MIT License' },
      }))
    }
    if (data.cta) Object.assign(cta.value, data.cta)
    customSections.value = Object.entries(data)
      .filter(([type]) => !knownTypes.includes(type))
      .map(([type, content]) => ({ type, content: content as any }))
  } catch { /* fallback to hardcoded defaults */ }
})

function handleCopyEmail() {
  navigator.clipboard.writeText(cta.value.email)
  copied.value = true
  setTimeout(() => (copied.value = false), 2000)
}

function scrollTo(id: string) {
  document.getElementById(id)?.scrollIntoView({ behavior: 'smooth' })
}
</script>

<template>
  <div :class="['min-h-screen bg-[#fafafa]', tmpl[template as keyof typeof tmpl]]">
    <!-- Toast -->
    <div v-if="copied" class="fixed bottom-4 left-1/2 -translate-x-1/2 bg-[#1a1a1a] text-white text-sm px-4 py-3 rounded-lg shadow-lg z-100 animate-fade-in">
      Email copied to clipboard!
    </div>

    <!-- Header -->
    <header class="fixed top-0 left-0 right-0 z-50 bg-[#fafafa]/80 backdrop-blur-md border-b border-[#e5e5e5]/50">
      <div class="max-width-container mx-auto px-4 py-3 flex items-center justify-between">
        <div class="flex items-center">
          <span class="text-xs sm:text-sm font-semibold text-[#1a1a1a]">OKURU.ID</span>
        </div>

        <nav aria-label="Navigasi utama desktop" class="hidden md:flex items-center gap-4 text-xs sm:text-sm text-[#6b7280]">
          <button @click="scrollTo('open-source')" class="hover:text-[#1a1a1a] transition-colors cursor-pointer">Open Source</button>
        </nav>

        <nav aria-label="Navigasi utama mobile" class="flex md:hidden items-center gap-3 text-[#6b7280]">
          <button @click="scrollTo('open-source')" class="hover:text-[#1a1a1a] transition-colors text-xs sm:text-sm cursor-pointer">Open Source</button>
        </nav>
      </div>
    </header>

    <main class="sm:pt-32 pb-12 sm:pb-16">
      <!-- Hero -->
      <section class="max-width-container mx-auto px-4 text-center mb-16 sm:mb-24" aria-labelledby="hero-title">
        <div class="animate-fade-in">
          <!-- Mobile -->
          <div class="md:hidden w-full mb-8 sm:mb-12">
            <div class="w-full relative">
              <img
                :src="hero.profile_image_mobile"
                alt="Kurob, fullstack programmer Cilegon"
                class="w-full h-auto object-cover"
                fetchpriority="high"
                decoding="async"
              />
              <div class="absolute bottom-0 left-0 right-0 bg-white/80 backdrop-blur-md p-4 text-left">
                <p class="text-xl sm:text-2xl font-semibold text-[#1a1a1a] mb-1">
                  {{ language === 'en' ? hero.greeting_en : hero.greeting_id }}
                </p>
                <h1 id="hero-title" class="text-base font-medium text-[#6b7280] leading-tight">
                  {{ language === 'en' ? hero.description_en : hero.description_id }}
                </h1>
              </div>
            </div>
          </div>

          <!-- Desktop -->
          <div class="hidden md:flex flex-row items-start gap-8 mb-12">
            <div class="w-1/2">
              <img
                :src="hero.profile_image"
                alt="Kurob, fullstack programmer Cilegon"
                class="w-full h-auto object-cover rounded-lg"
                fetchpriority="high"
                decoding="async"
              />
            </div>
            <div class="w-1/2 self-center">
              <p class="text-xl lg:text-2xl font-semibold text-[#1a1a1a] mb-2">
                {{ language === 'en' ? hero.greeting_en : hero.greeting_id }}
              </p>
              <p class="text-lg lg:text-xl font-medium text-[#6b7280] leading-tight mb-6">
                {{ language === 'en' ? hero.description_en : hero.description_id }}
              </p>
              <button @click="scrollTo('contact')" class="btn-primary inline-flex items-center gap-2">
                {{ language === 'en' ? hero.cta_text : 'Hubungi Saya' }}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Language Toggle Mobile -->
          <div class="mt-6 sm:mt-8 flex flex-col items-center gap-4 max-w-xl mx-auto md:hidden">
            <div class="flex items-center gap-1 bg-white rounded-full p-1 border border-[#e5e5e5]">
              <button
                @click="language = 'en'"
                :class="['px-3 py-1 text-xs font-medium rounded-full transition-all', language === 'en' ? 'bg-[#1a1a1a] text-white' : 'text-[#6b7280] hover:text-[#1a1a1a]']"
              >
                EN
              </button>
              <button
                @click="language = 'id'"
                :class="['px-3 py-1 text-xs font-medium rounded-full transition-all', language === 'id' ? 'bg-[#1a1a1a] text-white' : 'text-[#6b7280] hover:text-[#1a1a1a]']"
              >
                ID
              </button>
            </div>
            <button @click="scrollTo('contact')" class="btn-primary inline-flex items-center gap-2">
              {{ language === 'en' ? hero.cta_text : 'Hubungi Saya' }}
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
              </svg>
            </button>
          </div>

          <!-- Language Toggle Desktop -->
          <div class="hidden md:flex justify-center mt-8">
            <div class="flex items-center gap-1 bg-white rounded-full p-1 border border-[#e5e5e5]">
              <button
                @click="language = 'en'"
                :class="['px-3 py-1 text-xs font-medium rounded-full transition-all', language === 'en' ? 'bg-[#1a1a1a] text-white' : 'text-[#6b7280] hover:text-[#1a1a1a]']"
              >
                EN
              </button>
              <button
                @click="language = 'id'"
                :class="['px-3 py-1 text-xs font-medium rounded-full transition-all', language === 'id' ? 'bg-[#1a1a1a] text-white' : 'text-[#6b7280] hover:text-[#1a1a1a]']"
              >
                ID
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- Client Logos -->
      <section class="max-width-container mx-auto px-4 mb-12 sm:mb-16">
        <div class="bg-white rounded-2xl sm:rounded-3xl py-6 sm:py-10 shadow-[0_4px_40px_rgba(0,0,0,0.06)] overflow-hidden">
          <div class="flex animate-marquee">
            <template v-for="(client, index) in [...clientLogos, ...clientLogos]" :key="`${client.name}-${index}`">
              <div class="flex-shrink-0 flex items-center justify-center px-8 sm:px-10 opacity-60 hover:opacity-100 transition-opacity">
                <img
                  :src="client.logo"
                  :alt="client.name"
                  class="h-8 sm:h-10 md:h-12 w-auto object-contain grayscale hover:grayscale-0 transition-all"
                  loading="lazy"
                  decoding="async"
                />
              </div>
            </template>
          </div>
        </div>
      </section>

      <!-- Collaborate Text -->
      <section class="max-width-container mx-auto px-4 mb-12 sm:mb-16 text-center">
        <p class="text-base sm:text-xl md:text-2xl text-[#6b7280]">
          {{ language === 'en'
            ? 'Collaborate with a fullstack programmer from Cilegon to build impactful web products, automation, and digital experiences.'
            : 'Berkolaborasi dengan fullstack programmer dari Cilegon untuk membangun produk web, automation, dan pengalaman digital yang berdampak.' }}
        </p>
      </section>

      <!-- Services -->
      <section class="max-width-container mx-auto px-4 mb-12 sm:mb-16">
        <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6">
          <div
            v-for="service in services"
            :key="service.title"
            class="bg-white rounded-2xl sm:rounded-3xl p-4 sm:p-6 sm:p-8 shadow-[0_4px_40px_rgba(0,0,0,0.06)] text-center transition-all duration-300 hover:-translate-y-1 hover:shadow-[0_8px_60px_rgba(0,0,0,0.1)]"
          >
            <div class="w-10 h-10 sm:w-14 sm:h-14 rounded-xl sm:rounded-2xl bg-[#fafafa] flex items-center justify-center mx-auto mb-3 sm:mb-6">
              <component :is="serviceIcon(service.icon)" class="w-5 h-5 sm:w-7 sm:h-7 text-[#1a1a1a]" />
            </div>
            <h3 class="text-sm sm:text-lg font-semibold text-[#1a1a1a] mb-2 sm:mb-3">{{ service.title }}</h3>
            <p class="text-xs sm:text-sm text-[#6b7280] leading-relaxed hidden sm:block">{{ language === 'en' ? service.description_en : service.description_id }}</p>
          </div>
        </div>
      </section>

      <!-- Open Source -->
      <section id="open-source" class="max-width-container mx-auto px-4 mb-12 sm:mb-16">
        <h2 class="text-xl sm:text-2xl md:text-3xl font-bold text-[#1a1a1a] mb-8">
          {{ language === 'en' ? 'Open Source Projects' : 'Proyek Open Source' }}
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 sm:gap-6">
          <a
            v-for="(project, index) in openSourceProjects"
            :key="index"
            :href="project.githubUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="group"
          >
            <div class="bg-white rounded-2xl sm:rounded-3xl p-5 sm:p-6 shadow-[0_4px_40px_rgba(0,0,0,0.06)] transition-all duration-300 hover:-translate-y-1 hover:shadow-[0_8px_60px_rgba(0,0,0,0.1)] h-full flex flex-col">
              <div class="flex items-center gap-2 mb-3">
                <IconBrandGithub class="w-5 h-5 text-[#1a1a1a]" />
                <h3 class="text-sm sm:text-base font-semibold text-[#1a1a1a] group-hover:text-[#4b5563] transition-colors">
                  {{ language === 'en' ? project.title.en : project.title.id }}
                </h3>
              </div>
              <p class="text-xs sm:text-sm text-[#6b7280] leading-relaxed flex-grow mb-4">
                {{ language === 'en' ? project.description.en : project.description.id }}
              </p>
              <div class="flex flex-wrap gap-2 mb-4">
                <span
                  v-for="(tech, techIndex) in project.technologies"
                  :key="techIndex"
                  :class="`px-2 py-1 text-xs rounded-full ${tech.bgColor} ${tech.textColor}`"
                >
                  {{ tech.name }}
                </span>
              </div>
              <div class="flex items-center gap-4 text-xs text-[#9ca3af]">
                <span v-if="project.stats.stars !== undefined" class="inline-flex items-center gap-1">
                  <IconStar :size="14" /> {{ project.stats.stars }}
                </span>
                <span v-if="project.stats.language">{{ project.stats.language }}</span>
                <span v-if="project.stats.license" class="inline-flex items-center gap-1">
                  <IconLicense :size="14" /> {{ project.stats.license }}
                </span>
                <span v-if="project.stats.type">{{ project.stats.type }}</span>
              </div>
            </div>
          </a>
        </div>
      </section>

      <!-- Custom Sections -->
      <section v-for="cs in customSections" :key="cs.type" class="max-width-container mx-auto px-4 mb-12 sm:mb-16">
        <h2 class="text-xl sm:text-2xl md:text-3xl font-bold text-[#1a1a1a] mb-8 capitalize">{{ cs.type }}</h2>
        <div v-if="cs.content.items" class="grid grid-cols-1 md:grid-cols-2 gap-4 sm:gap-6">
          <div v-for="(item, ci) in cs.content.items" :key="ci"
            class="bg-white rounded-2xl sm:rounded-3xl p-5 sm:p-6 shadow-[0_4px_40px_rgba(0,0,0,0.06)]"
          >
            <div v-for="(v, k) in item" :key="k as string" class="mb-1">
              <span class="text-xs text-[#9ca3af] capitalize">{{ k }}:</span>
              <span class="text-sm text-[#1a1a1a] ml-1">{{ v }}</span>
            </div>
          </div>
        </div>
        <div v-else class="bg-white rounded-2xl sm:rounded-3xl p-5 sm:p-6 shadow-[0_4px_40px_rgba(0,0,0,0.06)]">
          <div v-for="(v, k) in cs.content" :key="k as string" class="mb-2">
            <span class="text-xs text-[#9ca3af] capitalize">{{ k }}:</span>
            <span class="text-sm text-[#1a1a1a] ml-1">{{ v }}</span>
          </div>
        </div>
      </section>

      <!-- CTA -->
      <section id="contact" class="max-width-container mx-auto px-4 mb-12 sm:mb-16">
        <div class="bg-[#1a1a1a] rounded-2xl sm:rounded-3xl p-8 sm:p-12 md:p-16 text-center text-white">
          <div class="w-12 h-12 sm:w-16 sm:h-16 rounded-full bg-white/10 flex items-center justify-center mx-auto mb-6 sm:mb-8">
            <svg class="w-6 h-6 sm:w-8 sm:h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
          </div>
          <h2 class="text-xl sm:text-2xl md:text-4xl lg:text-5xl font-bold mb-6 sm:mb-8">
            {{ cta.heading }}
          </h2>
          <div class="flex flex-col sm:flex-row gap-3 sm:gap-4 justify-center">
            <button
              @click="handleCopyEmail"
              class="bg-white text-[#1a1a1a] px-6 sm:px-8 py-3 sm:py-4 rounded-full font-medium text-sm hover:bg-gray-100 transition-colors inline-flex items-center gap-2 justify-center"
            >
              <IconMail :size="18" />
              Contact Me
            </button>
            <a
              :href="`https://wa.me/${cta.whatsapp}`"
              target="_blank"
              rel="noopener noreferrer"
              class="bg-transparent text-white px-6 sm:px-8 py-3 sm:py-4 rounded-full font-medium text-sm border border-white/30 hover:border-white transition-colors inline-flex items-center gap-2 justify-center"
              aria-label="Hubungi Kurob melalui WhatsApp"
            >
              <IconBrandWhatsapp :size="18" />
              WhatsApp
            </a>
          </div>
        </div>
      </section>
    </main>

    <!-- Footer -->
    <footer class="border-t border-[#e5e5e5]">
      <div class="max-width-container mx-auto px-4 md:px-8 py-8 flex flex-col md:flex-row items-center justify-between gap-4">
        <p class="text-sm text-[#6b7280]">© {{ new Date().getFullYear() }} All rights reserved.</p>
        <div class="flex items-center gap-6">
          <a href="https://linkedin.com/in/kurob1993" target="_blank" rel="noopener noreferrer" class="text-[#6b7280] hover:text-[#1a1a1a] transition-colors" aria-label="LinkedIn Kurob">
            <IconBrandLinkedin :size="20" />
          </a>
          <a href="https://www.threads.com/@_okuru.id" target="_blank" rel="noopener noreferrer" class="text-[#6b7280] hover:text-[#1a1a1a] transition-colors" aria-label="Threads Kurob">
            <IconBrandThreads :size="20" />
          </a>
          <a href="https://instagram.com/kurob1993" target="_blank" rel="noopener noreferrer" class="text-[#6b7280] hover:text-[#1a1a1a] transition-colors" aria-label="Instagram Kurob">
            <IconBrandInstagram :size="20" />
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>

<style>
/* Minimal template: more whitespace, no shadows */
.template-minimal .shadow-\[0_4px_40px_rgba\(0\,0\,0\,0\.06\)\] {
  box-shadow: none !important;
  border: 1px solid #e5e5e5;
}
.template-minimal .shadow-\[0_8px_60px_rgba\(0\,0\,0\,0\.1\)\] {
  box-shadow: none !important;
  border: 1px solid #d1d5db;
}
.template-minimal section > div > .rounded-2xl,
.template-minimal section > div > .rounded-3xl {
  border-radius: 0.75rem;
}

/* Bold template: brand accent, larger everything */
.template-bold {
  --brand: #2563eb;
}
.template-bold .btn-primary {
  background-color: var(--brand) !important;
  font-weight: 700;
  padding: 0.75rem 2rem;
}
.template-bold h2,
.template-bold h1 {
  letter-spacing: -0.025em;
}
.template-bold section > div > .bg-\[\#1a1a1a\] {
  background: linear-gradient(135deg, #1e3a8a, #111827) !important;
}
</style>
