package seeders

import (
	"okuru/app/facades"
	"okuru/app/models"
)

type LandingTemplateSeeder struct{}

func (s *LandingTemplateSeeder) Signature() string {
	return "LandingTemplateSeeder"
}

func (s *LandingTemplateSeeder) Run() error {
	templates := []models.LandingTemplate{
		{
			Name:        "Sistem",
			Description: "Single-platform operations landing page. Tailwind + Inter/Space Grotesk.",
			HTML:        sistemHTML(),
		},
	}

	// Remove obsolete templates from prior seeds.
	if _, err := facades.Orm().Query().
		Where("name IN ?", []string{"Professional", "Minimal", "Showcase"}).
		Delete(&models.LandingTemplate{}); err != nil {
		return err
	}

	for _, t := range templates {
		var existing models.LandingTemplate
		if err := facades.Orm().Query().UpdateOrCreate(
			&existing,
			models.LandingTemplate{Name: t.Name},
			t,
		); err != nil {
			return err
		}
	}
	return nil
}

func sistemHTML() string {
	return `<!DOCTYPE html>
<html lang="id">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Sistem — Landing Page</title>
<script src="https://cdn.tailwindcss.com"></script>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@500;600;700&family=Inter:wght@400;500;600&family=IBM+Plex+Mono:wght@400;500&display=swap" rel="stylesheet">
<style>
  :root{
    --bg:#F7F7F6;
    --surface:#FFFFFF;
    --ink:#14181F;
    --ink-soft:#5B6472;
    --line:#E4E4E1;
    --accent:#2F6F5E;
    --accent-ink:#F4FAF8;
  }
  body{ background:var(--bg); color:var(--ink); font-family:'Inter',sans-serif; }
  .font-display{ font-family:'Space Grotesk',sans-serif; }
  .font-mono{ font-family:'IBM Plex Mono',monospace; }
  .btn-accent{ background:var(--accent); color:#fff; }
  .btn-accent:hover{ background:#265d4f; }
  .dot{ width:6px; height:6px; border-radius:9999px; background:var(--accent); }
  ::selection{ background:var(--accent); color:#fff; }
  :focus-visible{ outline:2px solid var(--accent); outline-offset:2px; }
</style>
</head>
<body class="antialiased">

<!-- NAV -->
<header class="border-b" style="border-color:var(--line)">
  <div class="max-w-6xl mx-auto px-6 h-16 flex items-center justify-between">
    <a href="#" class="font-display font-semibold text-lg tracking-tight" style="color:var(--ink)">
      Sistem<span style="color:var(--accent)">.</span>
    </a>
    <nav class="hidden md:flex items-center gap-8 text-sm" style="color:var(--ink-soft)">
      <a href="#fitur" class="hover:text-[--ink] transition-colors">Fitur</a>
      <a href="#alur" class="hover:text-[--ink] transition-colors">Cara kerja</a>
      <a href="#harga" class="hover:text-[--ink] transition-colors">Harga</a>
      <a href="#faq" class="hover:text-[--ink] transition-colors">FAQ</a>
    </nav>
    <div class="flex items-center gap-3">
      <a href="#" class="hidden sm:inline text-sm font-medium" style="color:var(--ink-soft)">Masuk</a>
      <a href="#mulai" class="btn-accent text-sm font-medium px-4 py-2 rounded-md transition-colors">
        Mulai gratis
      </a>
    </div>
  </div>
</header>

<!-- HERO -->
<section class="max-w-6xl mx-auto px-6 pt-20 pb-16 grid lg:grid-cols-2 gap-14 items-center">
  <div>
    <div class="inline-flex items-center gap-2 text-xs font-mono px-3 py-1 rounded-full border mb-6" style="border-color:var(--line); color:var(--ink-soft)">
      <span class="dot"></span> status: operasional
    </div>
    <h1 class="font-display font-semibold text-4xl sm:text-5xl leading-[1.1] tracking-tight" style="color:var(--ink)">
      Satu sistem, untuk seluruh operasional Anda.
    </h1>
    <p class="mt-5 text-base sm:text-lg leading-relaxed max-w-lg" style="color:var(--ink-soft)">
      Kelola proses, data, dan tim dalam satu platform yang rapi dan bisa diandalkan — tanpa tools yang tercerai-berai.
    </p>
    <div class="mt-8 flex flex-wrap items-center gap-4">
      <a href="#mulai" class="btn-accent px-6 py-3 rounded-md text-sm font-medium transition-colors">
        Coba sekarang
      </a>
      <a href="#alur" class="px-6 py-3 rounded-md text-sm font-medium border transition-colors hover:bg-white" style="border-color:var(--line); color:var(--ink)">
        Lihat cara kerja
      </a>
    </div>
    <p class="mt-6 text-xs font-mono" style="color:var(--ink-soft)">Tidak perlu kartu kredit · Aktif dalam 5 menit</p>
  </div>

  <!-- Signature visual: terminal/status mockup -->
  <div class="rounded-xl border shadow-sm overflow-hidden" style="border-color:var(--line); background:var(--surface)">
    <div class="flex items-center gap-2 px-4 py-3 border-b" style="border-color:var(--line)">
      <span class="w-2.5 h-2.5 rounded-full" style="background:#DDD"></span>
      <span class="w-2.5 h-2.5 rounded-full" style="background:#DDD"></span>
      <span class="w-2.5 h-2.5 rounded-full" style="background:#DDD"></span>
      <span class="ml-3 text-xs font-mono" style="color:var(--ink-soft)">sistem — ringkasan</span>
    </div>
    <div class="p-5 font-mono text-sm space-y-3" style="color:var(--ink)">
      <div class="flex justify-between">
        <span style="color:var(--ink-soft)">proses_aktif</span>
        <span>128</span>
      </div>
      <div class="flex justify-between">
        <span style="color:var(--ink-soft)">tugas_selesai_hari_ini</span>
        <span>1.204</span>
      </div>
      <div class="flex justify-between">
        <span style="color:var(--ink-soft)">rata_waktu_respon</span>
        <span>2,4 dtk</span>
      </div>
      <div class="pt-3 border-t" style="border-color:var(--line)">
        <div class="flex items-center gap-2" style="color:var(--accent)">
          <span class="dot"></span>
          <span>semua sistem berjalan normal</span>
        </div>
      </div>
    </div>
  </div>
</section>

<!-- LOGOS / TRUST -->
<section class="border-y" style="border-color:var(--line); background:var(--surface)">
  <div class="max-w-6xl mx-auto px-6 py-8 flex flex-wrap items-center justify-between gap-6">
    <p class="text-xs font-mono uppercase tracking-wide" style="color:var(--ink-soft)">Dipercaya oleh tim operasional dari berbagai industri</p>
    <div class="flex flex-wrap items-center gap-8 text-sm font-display font-semibold" style="color:var(--ink-soft)">
      <span>Nadira</span>
      <span>Kertas Baja</span>
      <span>Lintas Data</span>
      <span>Mitra Logistik</span>
      <span>Beringin Grup</span>
    </div>
  </div>
</section>

<!-- FITUR -->
<section id="fitur" class="max-w-6xl mx-auto px-6 py-20">
  <div class="max-w-xl mb-14">
    <p class="text-xs font-mono uppercase tracking-wide mb-3" style="color:var(--accent)">Fitur utama</p>
    <h2 class="font-display font-semibold text-3xl tracking-tight" style="color:var(--ink)">
      Dibangun untuk operasional yang jelas dan terukur
    </h2>
  </div>
  <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-px" style="background:var(--line)">
    <div class="p-7" style="background:var(--surface)">
      <div class="w-9 h-9 rounded-md flex items-center justify-center mb-5" style="background:var(--accent-ink); color:var(--accent)">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 12h18M3 6h18M3 18h18"/></svg>
      </div>
      <h3 class="font-display font-semibold text-base mb-2" style="color:var(--ink)">Alur kerja terpusat</h3>
      <p class="text-sm leading-relaxed" style="color:var(--ink-soft)">Semua proses tercatat di satu tempat, dari permintaan masuk sampai selesai dikerjakan.</p>
    </div>
    <div class="p-7" style="background:var(--surface)">
      <div class="w-9 h-9 rounded-md flex items-center justify-center mb-5" style="background:var(--accent-ink); color:var(--accent)">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 20V10M18 20V4M6 20v-4"/></svg>
      </div>
      <h3 class="font-display font-semibold text-base mb-2" style="color:var(--ink)">Laporan real-time</h3>
      <p class="text-sm leading-relaxed" style="color:var(--ink-soft)">Pantau kinerja tim dan sistem lewat data yang selalu diperbarui, bukan laporan basi.</p>
    </div>
    <div class="p-7" style="background:var(--surface)">
      <div class="w-9 h-9 rounded-md flex items-center justify-center mb-5" style="background:var(--accent-ink); color:var(--accent)">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="10" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
      </div>
      <h3 class="font-display font-semibold text-base mb-2" style="color:var(--ink)">Akses & keamanan</h3>
      <p class="text-sm leading-relaxed" style="color:var(--ink-soft)">Atur peran dan hak akses per anggota tim, dengan jejak aktivitas yang bisa ditelusuri.</p>
    </div>
    <div class="p-7" style="background:var(--surface)">
      <div class="w-9 h-9 rounded-md flex items-center justify-center mb-5" style="background:var(--accent-ink); color:var(--accent)">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 4h16v16H4z"/><path d="M4 9h16M9 20V9"/></svg>
      </div>
      <h3 class="font-display font-semibold text-base mb-2" style="color:var(--ink)">Integrasi terbuka</h3>
      <p class="text-sm leading-relaxed" style="color:var(--ink-soft)">Hubungkan dengan sistem yang sudah Anda pakai lewat API, tanpa perlu migrasi total.</p>
    </div>
    <div class="p-7" style="background:var(--surface)">
      <div class="w-9 h-9 rounded-md flex items-center justify-center mb-5" style="background:var(--accent-ink); color:var(--accent)">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 3"/></svg>
      </div>
      <h3 class="font-display font-semibold text-base mb-2" style="color:var(--ink)">Notifikasi tepat waktu</h3>
      <p class="text-sm leading-relaxed" style="color:var(--ink-soft)">Dapat kabar begitu ada perubahan penting, tanpa harus bolak-balik memeriksa.</p>
    </div>
    <div class="p-7" style="background:var(--surface)">
      <div class="w-9 h-9 rounded-md flex items-center justify-center mb-5" style="background:var(--accent-ink); color:var(--accent)">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2l9 4.5v9L12 22l-9-6.5v-9z"/></svg>
      </div>
      <h3 class="font-display font-semibold text-base mb-2" style="color:var(--ink)">Skalabel dari awal</h3>
      <p class="text-sm leading-relaxed" style="color:var(--ink-soft)">Cocok untuk tim kecil, tetap stabil saat operasional Anda tumbuh lebih besar.</p>
    </div>
  </div>
</section>

<!-- ALUR KERJA -->
<section id="alur" class="border-t" style="border-color:var(--line); background:var(--surface)">
  <div class="max-w-6xl mx-auto px-6 py-20">
    <div class="max-w-xl mb-14">
      <p class="text-xs font-mono uppercase tracking-wide mb-3" style="color:var(--accent)">Cara kerja</p>
      <h2 class="font-display font-semibold text-3xl tracking-tight" style="color:var(--ink)">Tiga langkah, langsung berjalan</h2>
    </div>
    <div class="grid md:grid-cols-3 gap-10">
      <div>
        <p class="font-mono text-sm mb-3" style="color:var(--accent)">01</p>
        <h3 class="font-display font-semibold mb-2" style="color:var(--ink)">Hubungkan data Anda</h3>
        <p class="text-sm leading-relaxed" style="color:var(--ink-soft)">Impor data yang ada atau sambungkan lewat integrasi yang tersedia.</p>
      </div>
      <div>
        <p class="font-mono text-sm mb-3" style="color:var(--accent)">02</p>
        <h3 class="font-display font-semibold mb-2" style="color:var(--ink)">Atur alur & tim</h3>
        <p class="text-sm leading-relaxed" style="color:var(--ink-soft)">Sesuaikan proses dan hak akses sesuai struktur organisasi Anda.</p>
      </div>
      <div>
        <p class="font-mono text-sm mb-3" style="color:var(--accent)">03</p>
        <h3 class="font-display font-semibold mb-2" style="color:var(--ink)">Pantau & tindak lanjuti</h3>
        <p class="text-sm leading-relaxed" style="color:var(--ink-soft)">Lihat perkembangan lewat dashboard, ambil keputusan berdasarkan data.</p>
      </div>
    </div>
  </div>
</section>

<!-- CTA -->
<section id="mulai" class="max-w-6xl mx-auto px-6 py-20">
  <div class="rounded-xl border p-10 sm:p-14 text-center" style="border-color:var(--line); background:var(--surface)">
    <h2 class="font-display font-semibold text-2xl sm:text-3xl tracking-tight" style="color:var(--ink)">
      Siap merapikan sistem operasional Anda?
    </h2>
    <p class="mt-3 text-sm sm:text-base max-w-md mx-auto" style="color:var(--ink-soft)">
      Mulai gratis hari ini, tanpa komitmen jangka panjang.
    </p>
    <div class="mt-7 flex flex-wrap items-center justify-center gap-4">
      <a href="#" class="btn-accent px-6 py-3 rounded-md text-sm font-medium transition-colors">Buat akun</a>
      <a href="#" class="px-6 py-3 rounded-md text-sm font-medium border transition-colors" style="border-color:var(--line); color:var(--ink)">Hubungi kami</a>
    </div>
  </div>
</section>

<!-- FOOTER -->
<footer class="border-t" style="border-color:var(--line)">
  <div class="max-w-6xl mx-auto px-6 py-10 flex flex-col sm:flex-row items-center justify-between gap-4">
    <p class="font-display font-semibold" style="color:var(--ink)">Sistem<span style="color:var(--accent)">.</span></p>
    <p class="text-xs font-mono" style="color:var(--ink-soft)">© 2026 Sistem. Seluruh hak cipta dilindungi.</p>
    <div class="flex gap-6 text-sm" style="color:var(--ink-soft)">
      <a href="#" class="hover:text-[--ink]">Privasi</a>
      <a href="#" class="hover:text-[--ink]">Ketentuan</a>
      <a href="#" class="hover:text-[--ink]">Kontak</a>
    </div>
  </div>
</footer>

</body>
</html>
`
}
