package seeders

import (
	"okuru/app/facades"
	"okuru/app/models"
)

type DatabaseSeeder struct{}

func (s *DatabaseSeeder) Signature() string {
	return "DatabaseSeeder"
}

func (s *DatabaseSeeder) Run() error {
	if err := s.seedAdmin(); err != nil {
		return err
	}
	if err := s.seedCategories(); err != nil {
		return err
	}
	if err := s.seedSettings(); err != nil {
		return err
	}
	if err := s.seedLandingSections(); err != nil {
		return err
	}
	if err := s.seedLandingTemplates(); err != nil {
		return err
	}
	return nil
}

func (s *DatabaseSeeder) seedLandingSections() error {
	seeder := &LandingSectionSeeder{}
	return seeder.Run()
}

func (s *DatabaseSeeder) seedLandingTemplates() error {
	seeder := &LandingTemplateSeeder{}
	return seeder.Run()
}

func (s *DatabaseSeeder) seedAdmin() error {
	hashed, err := facades.Hash().Make("initial.1")
	if err != nil {
		return err
	}
	var user models.User
	return facades.Orm().Query().UpdateOrCreate(
		&user,
		models.User{Email: "admin@okuru.id"},
		models.User{Password: hashed},
	)
}

func (s *DatabaseSeeder) seedCategories() error {
	categories := []models.Category{
		{Slug: "web-development", NameEn: "Web Development", NameId: "Pengembangan Web"},
		{Slug: "devops", NameEn: "DevOps", NameId: "DevOps"},
		{Slug: "javascript", NameEn: "JavaScript", NameId: "JavaScript"},
		{Slug: "laravel", NameEn: "Laravel", NameId: "Laravel"},
		{Slug: "react", NameEn: "React", NameId: "React"},
		{Slug: "docker", NameEn: "Docker", NameId: "Docker"},
	}
	for _, c := range categories {
		var cat models.Category
		if err := facades.Orm().Query().UpdateOrCreate(
			&cat,
			models.Category{Slug: c.Slug},
			c,
		); err != nil {
			return err
		}
	}
	return nil
}

func (s *DatabaseSeeder) seedSettings() error {
	settings := []models.Setting{
		{Key: "hero_title_en", Value: "Hi, I'm Kurob"},
		{Key: "hero_title_id", Value: "Halo, Saya Kurob"},
		{Key: "hero_desc_en", Value: "I'm from Indonesia with 8+ years of experience as a full-stack web developer."},
		{Key: "hero_desc_id", Value: "Saya dari Indonesia dengan pengalaman 8+ tahun sebagai full-stack web developer."},
		{Key: "start_year", Value: "2016"},
		{Key: "landing_mode", Value: "custom"},
		{Key: "landing_template", Value: "default"},
		{Key: "landing_template_html", Value: sistemHTML()},
	}
	for _, st := range settings {
		var setting models.Setting
		if err := facades.Orm().Query().UpdateOrCreate(
			&setting,
			models.Setting{Key: st.Key},
			st,
		); err != nil {
			return err
		}
	}
	return nil
}

func defaultCustomTemplate() string {
	return `<!DOCTYPE html>
<html lang="id">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Kurob - Fullstack Programmer Cilegon</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet">
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:'Inter',system-ui,sans-serif;line-height:1.6;color:#1a1a1a;background:#fafafa;-webkit-font-smoothing:antialiased}
.max-w{max-width:1024px;margin:0 auto;padding:0 1rem}
.fade-in{animation:fadeIn .6s ease-out}
@keyframes fadeIn{from{opacity:0;transform:translateY(20px)}to{opacity:1;transform:translateY(0)}}
@keyframes marquee{0%{transform:translateX(0)}100%{transform:translateX(-50%)}}
.marquee{display:flex;animation:marquee 30s linear infinite}
.marquee:hover{animation-play-state:paused}
header{position:fixed;top:0;left:0;right:0;z-index:50;background:rgba(250,250,250,0.8);-webkit-backdrop-filter:blur(12px);backdrop-filter:blur(12px);border-bottom:1px solid rgba(229,229,229,0.5)}
header .inner{max-width:1024px;margin:0 auto;padding:.75rem 1rem;display:flex;align-items:center;justify-content:space-between}
header .logo{font-size:.875rem;font-weight:600;color:#1a1a1a}
header nav a{font-size:.8125rem;color:#6b7280;text-decoration:none;cursor:pointer;transition:color .2s}
header nav a:hover{color:#1a1a1a}
main{padding-top:8rem;padding-bottom:4rem}
@media(max-width:768px){main{padding-top:6rem}}
section{margin-bottom:4rem}
@media(max-width:640px){section{margin-bottom:3rem}}
.hero{text-align:center}
.hero-desktop{display:flex;flex-direction:row;align-items:start;gap:2rem;margin-bottom:3rem}
.hero-desktop .img-col{width:50%}
.hero-desktop .img-col img{width:100%;height:auto;object-fit:cover;border-radius:.75rem}
.hero-desktop .text-col{width:50%;align-self:center;text-align:left}
.hero-desktop .text-col .greeting{font-size:1.25rem;font-weight:600;color:#1a1a1a;margin-bottom:.5rem}
@media(min-width:1024px){.hero-desktop .text-col .greeting{font-size:1.5rem}}
.hero-desktop .text-col .desc{font-size:1.125rem;font-weight:500;color:#6b7280;line-height:1.2;margin-bottom:1.5rem}
@media(min-width:1024px){.hero-desktop .text-col .desc{font-size:1.25rem}}
.hero-mobile{display:none}
@media(max-width:767px){.hero-desktop{display:none}.hero-mobile{display:block;margin-bottom:2rem}}
.hero-mobile .img-wrap{width:100%;position:relative}
.hero-mobile .img-wrap img{width:100%;height:auto;object-fit:cover}
.hero-mobile .img-overlay{position:absolute;bottom:0;left:0;right:0;background:rgba(255,255,255,0.8);-webkit-backdrop-filter:blur(12px);backdrop-filter:blur(12px);padding:1rem;text-align:left}
.hero-mobile .img-overlay .greeting{font-size:1.25rem;font-weight:600;color:#1a1a1a;margin-bottom:.25rem}
@media(min-width:640px){.hero-mobile .img-overlay .greeting{font-size:1.5rem}}
.hero-mobile .img-overlay .desc{font-size:1rem;font-weight:500;color:#6b7280;line-height:1.2}
.btn-primary{display:inline-flex;align-items:center;gap:.5rem;background:#1a1a1a;color:#fff;padding:.625rem 1.5rem;border-radius:999px;font-size:.875rem;font-weight:500;border:none;cursor:pointer;text-decoration:none;transition:opacity .2s}
.btn-primary:hover{opacity:.9}
.lang-toggle{display:flex;align-items:center;gap:.25rem;background:#fff;border-radius:999px;padding:.25rem;border:1px solid #e5e5e5;width:fit-content;margin:2rem auto 0}
.lang-toggle button{padding:.25rem .75rem;font-size:.8125rem;font-weight:500;border-radius:999px;border:none;cursor:pointer;transition:all .2s;background:transparent;color:#6b7280}
.lang-toggle button.active{background:#1a1a1a;color:#fff}
.lang-toggle button:not(.active):hover{color:#1a1a1a}
.logos-wrap{background:#fff;border-radius:1.5rem;padding:1.5rem 0;overflow:hidden;box-shadow:0 4px 40px rgba(0,0,0,0.06)}
@media(min-width:640px){.logos-wrap{border-radius:1.5rem;padding:2.5rem 0}}
.logos-wrap .inner{display:flex;animation:marquee 30s linear infinite}
.logos-wrap .inner:hover{animation-play-state:paused}
.logos-wrap .inner .item{flex-shrink:0;display:flex;align-items:center;justify-content:center;padding:0 2rem;opacity:.6;transition:opacity .2s}
@media(min-width:640px){.logos-wrap .inner .item{padding:0 2.5rem}}
.logos-wrap .inner .item:hover{opacity:1}
.logos-wrap .inner .item img{height:2rem;width:auto;object-fit:contain;filter:grayscale(1);transition:filter .2s}
@media(min-width:640px){.logos-wrap .inner .item img{height:2.5rem}}
@media(min-width:768px){.logos-wrap .inner .item img{height:3rem}}
.logos-wrap .inner .item img:hover{filter:grayscale(0)}
.colab-text{text-align:center}
.colab-text p{font-size:1rem;color:#6b7280}
@media(min-width:640px){.colab-text p{font-size:1.25rem}}
@media(min-width:768px){.colab-text p{font-size:1.5rem}}
.services-grid{display:grid;grid-template-columns:repeat(2,1fr);gap:1rem}
@media(min-width:1024px){.services-grid{grid-template-columns:repeat(4,1fr)}}
@media(min-width:640px){.services-grid{gap:1.5rem}}
.service-card{background:#fff;border-radius:1.5rem;padding:1.5rem;text-align:center;box-shadow:0 4px 40px rgba(0,0,0,0.06);transition:all .3s;cursor:default}
@media(min-width:640px){.service-card{padding:2rem}}
.service-card:hover{transform:translateY(-4px);box-shadow:0 8px 60px rgba(0,0,0,0.1)}
.service-icon{width:2.5rem;height:2.5rem;border-radius:.75rem;background:#fafafa;display:flex;align-items:center;justify-content:center;margin:0 auto .75rem}
@media(min-width:640px){.service-icon{width:3.5rem;height:3.5rem;border-radius:1rem;margin-bottom:1.5rem}}
.service-icon svg{width:1.25rem;height:1.25rem;color:#1a1a1a}
@media(min-width:640px){.service-icon svg{width:1.75rem;height:1.75rem}}
.service-card h3{font-size:.875rem;font-weight:600;color:#1a1a1a;margin-bottom:.5rem}
@media(min-width:640px){.service-card h3{font-size:1.125rem;margin-bottom:.75rem}}
.service-card p{font-size:.75rem;color:#6b7280;line-height:1.5;display:none}
@media(min-width:640px){.service-card p{display:block;font-size:.875rem}}
h2.section-title{font-size:1.25rem;font-weight:700;color:#1a1a1a;margin-bottom:2rem}
@media(min-width:640px){h2.section-title{font-size:1.5rem}}
@media(min-width:768px){h2.section-title{font-size:1.875rem}}
.projects-grid{display:grid;grid-template-columns:1fr;gap:1rem}
@media(min-width:768px){.projects-grid{grid-template-columns:repeat(2,1fr)}}
@media(min-width:640px){.projects-grid{gap:1.5rem}}
.project-card{background:#fff;border-radius:1.5rem;padding:1.5rem;box-shadow:0 4px 40px rgba(0,0,0,0.06);transition:all .3s;display:flex;flex-direction:column;height:100%;text-decoration:none;color:inherit}
.project-card:hover{transform:translateY(-4px);box-shadow:0 8px 60px rgba(0,0,0,0.1)}
.project-card .title-row{display:flex;align-items:center;gap:.5rem;margin-bottom:.75rem}
.project-card .title-row svg{width:1.25rem;height:1.25rem;color:#1a1a1a;flex-shrink:0}
.project-card .title-row h3{font-size:.875rem;font-weight:600;color:#1a1a1a;transition:color .2s}
@media(min-width:640px){.project-card .title-row h3{font-size:1rem}}
.project-card:hover .title-row h3{color:#4b5563}
.project-card .desc{font-size:.8125rem;color:#6b7280;line-height:1.5;flex-grow:1;margin-bottom:1rem}
.project-card .techs{display:flex;flex-wrap:wrap;gap:.5rem;margin-bottom:1rem}
.project-card .techs span{padding:.25rem .5rem;font-size:.75rem;border-radius:999px}
.project-card .meta{display:flex;align-items:center;gap:1rem;font-size:.75rem;color:#9ca3af}
.project-card .meta span{display:inline-flex;align-items:center;gap:.25rem}
.cta-section .cta-card{background:#1a1a1a;border-radius:1.5rem;padding:2rem 1.5rem;text-align:center;color:#fff}
@media(min-width:640px){.cta-section .cta-card{border-radius:1.5rem;padding:3rem}}
@media(min-width:768px){.cta-section .cta-card{padding:4rem}}
.cta-section .cta-card .mail-icon{width:3rem;height:3rem;border-radius:999px;background:rgba(255,255,255,0.1);display:flex;align-items:center;justify-content:center;margin:0 auto 1.5rem}
@media(min-width:640px){.cta-section .cta-card .mail-icon{width:4rem;height:4rem;margin-bottom:2rem}}
.cta-section .cta-card .mail-icon svg{width:1.5rem;height:1.5rem}
@media(min-width:640px){.cta-section .cta-card .mail-icon svg{width:2rem;height:2rem}}
.cta-section .cta-card h2{font-size:1.25rem;font-weight:700;margin-bottom:1.5rem}
@media(min-width:640px){.cta-section .cta-card h2{font-size:1.5rem}}
@media(min-width:768px){.cta-section .cta-card h2{font-size:2.25rem}}
@media(min-width:1024px){.cta-section .cta-card h2{font-size:3rem}}
.cta-section .cta-card .btns{display:flex;flex-direction:column;gap:.75rem;justify-content:center}
@media(min-width:640px){.cta-section .cta-card .btns{flex-direction:row;gap:1rem}}
.cta-section .cta-card .btns .btn-light{background:#fff;color:#1a1a1a;padding:.75rem 2rem;border-radius:999px;font-size:.875rem;font-weight:500;border:none;cursor:pointer;display:inline-flex;align-items:center;gap:.5rem;justify-content:center;text-decoration:none;transition:background .2s}
.cta-section .cta-card .btns .btn-light:hover{background:#f3f4f6}
.cta-section .cta-card .btns .btn-outline{background:transparent;color:#fff;padding:.75rem 2rem;border-radius:999px;font-size:.875rem;font-weight:500;border:1px solid rgba(255,255,255,0.3);cursor:pointer;display:inline-flex;align-items:center;gap:.5rem;justify-content:center;text-decoration:none;transition:border-color .2s}
.cta-section .cta-card .btns .btn-outline:hover{border-color:#fff}
footer{border-top:1px solid #e5e5e5}
footer .inner{max-width:1024px;margin:0 auto;padding:2rem 1rem;display:flex;flex-direction:column;align-items:center;gap:1rem}
@media(min-width:768px){footer .inner{flex-direction:row;justify-content:space-between}}
footer .inner p{font-size:.875rem;color:#6b7280}
footer .inner .social{display:flex;align-items:center;gap:1.5rem}
footer .inner .social a{color:#6b7280;transition:color .2s}
footer .inner .social a:hover{color:#1a1a1a}
.toast{position:fixed;bottom:1rem;left:50%;transform:translateX(-50%);background:#1a1a1a;color:#fff;font-size:.875rem;padding:.75rem 1rem;border-radius:.75rem;box-shadow:0 4px 24px rgba(0,0,0,0.2);z-index:100;display:none;animation:fadeIn .3s ease-out}
.toast.show{display:block}
.custom-section{margin-bottom:4rem}
.custom-section h2{font-size:1.25rem;font-weight:700;color:#1a1a1a;margin-bottom:2rem;text-transform:capitalize}
@media(min-width:640px){.custom-section h2{font-size:1.5rem}}
@media(min-width:768px){.custom-section h2{font-size:1.875rem}}
.custom-section .items-grid{display:grid;grid-template-columns:1fr;gap:1rem}
@media(min-width:768px){.custom-section .items-grid{grid-template-columns:repeat(2,1fr)}}
@media(min-width:640px){.custom-section .items-grid{gap:1.5rem}}
.custom-section .item-card{background:#fff;border-radius:1.5rem;padding:1.25rem 1.5rem;box-shadow:0 4px 40px rgba(0,0,0,0.06)}
.custom-section .item-card .field{margin-bottom:.5rem}
.custom-section .item-card .field .k{font-size:.75rem;color:#9ca3af;text-transform:capitalize}
.custom-section .item-card .field .v{font-size:.875rem;color:#1a1a1a;margin-left:.25rem}
.custom-section .plain-card{background:#fff;border-radius:1.5rem;padding:1.25rem 1.5rem;box-shadow:0 4px 40px rgba(0,0,0,0.06);max-width:640px}
.custom-section .plain-card .field{margin-bottom:.5rem}
.custom-section .plain-card .field .k{font-size:.75rem;color:#9ca3af;text-transform:capitalize}
.custom-section .plain-card .field .v{font-size:.875rem;color:#1a1a1a;margin-left:.25rem}
</style>
</head>
<body>

<div id="app">
  <header>
    <div class="inner">
      <span class="logo">OKURU.ID</span>
      <nav>
        <a onclick="scrollToEl('open-source')">Open Source</a>
      </nav>
    </div>
  </header>

  <main>
    <!-- Hero -->
    <section class="max-w hero fade-in" id="hero-section">
      <div id="hero-desktop" class="hero-desktop"></div>
      <div id="hero-mobile" class="hero-mobile"></div>
      <div id="lang-toggle" class="lang-toggle"></div>
    </section>

    <!-- Client Logos -->
    <section class="max-w" id="clients-section"></section>

    <!-- Collaborate Text -->
    <section class="max-w colab-text" id="colab-section"></section>

    <!-- Services -->
    <section class="max-w" id="services-section"></section>

    <!-- Open Source Projects -->
    <section class="max-w" id="projects-section"></section>

    <!-- Custom Sections -->
    <div id="custom-sections"></div>

    <!-- CTA -->
    <section class="max-w cta-section" id="cta-section"></section>
  </main>

  <footer>
    <div class="inner">
      <p>&copy; <span id="year"></span> All rights reserved.</p>
      <div class="social">
        <a href="https://linkedin.com/in/kurob1993" target="_blank" rel="noopener" aria-label="LinkedIn">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z"/><rect width="4" height="12" x="2" y="9"/><circle cx="4" cy="4" r="2"/></svg>
        </a>
        <a href="https://www.threads.com/@_okuru.id" target="_blank" rel="noopener" aria-label="Threads">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19.25 8.5c-1.02-2.47-3.15-4-5.75-4-3.72 0-6.75 3.36-6.75 8s3.03 8 6.75 8c2.6 0 4.73-1.53 5.75-4m-5.75-10c-2.21 0-4 2.69-4 6s1.79 6 4 6c1.66 0 3.06-1.29 3.64-3.5M10 13h4"/></svg>
        </a>
        <a href="https://instagram.com/kurob1993" target="_blank" rel="noopener" aria-label="Instagram">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="20" height="20" x="2" y="2" rx="5" ry="5"/><path d="M16 11.37A4 4 0 1 1 12.63 8 4 4 0 0 1 16 11.37z"/><line x1="17.5" x2="17.51" y1="6.5" y2="6.5"/></svg>
        </a>
      </div>
    </div>
  </footer>

  <div class="toast" id="toast"></div>
</div>

<script>
(function(){
var lang = 'en';
var copied = false;
var DATA = {{DATA}};

function render(){
  var hero = DATA.hero || {};
  var clients = DATA.clients || {};
  var services = DATA.services || {};
  var projects = DATA.projects || {};
  var cta = DATA.cta || {};

  var greeting = lang === 'en' ? (hero.greeting_en || '') : (hero.greeting_id || '');
  var desc = lang === 'en' ? (hero.description_en || '') : (hero.description_id || '');

  // Hero Desktop
  document.getElementById('hero-desktop').innerHTML =
    '<div class="img-col"><img src="' + (hero.profile_image || '/images/profile.webp') + '" alt="Kurob" loading="eager"></div>' +
    '<div class="text-col">' +
      '<p class="greeting">' + greeting + '</p>' +
      '<p class="desc">' + desc + '</p>' +
      '<a href="#contact" class="btn-primary">' + (lang === 'en' ? (hero.cta_text || 'Contact Me') : 'Hubungi Saya') +
        '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 8l4 4m0 0l-4 4m4-4H3"/></svg>' +
      '</a>' +
    '</div>';

  // Hero Mobile
  document.getElementById('hero-mobile').innerHTML =
    '<div class="img-wrap">' +
      '<img src="' + (hero.profile_image_mobile || hero.profile_image || '/images/profile-half.webp') + '" alt="Kurob" loading="eager">' +
      '<div class="img-overlay">' +
        '<p class="greeting">' + greeting + '</p>' +
        '<p class="desc">' + desc + '</p>' +
      '</div>' +
    '</div>';

  // Language toggle
  var toggleHTML =
    '<button onclick="setLang(\'en\')" class="' + (lang === 'en' ? 'active' : '') + '">EN</button>' +
    '<button onclick="setLang(\'id\')" class="' + (lang === 'id' ? 'active' : '') + '">ID</button>';
  document.getElementById('lang-toggle').innerHTML = toggleHTML;

  // CTA button in mobile hero
  var ctaMobile = document.querySelector('.hero-mobile + .lang-toggle + .btn-primary');
  // Client logos
  var items = clients.items || [];
  if(items.length){
    var logosHTML = '<div class="logos-wrap"><div class="inner">';
    var doubled = items.concat(items);
    for(var i=0;i<doubled.length;i++){
      logosHTML += '<div class="item"><img src="' + doubled[i].logo + '" alt="' + (doubled[i].name||'') + '" loading="lazy"></div>';
    }
    logosHTML += '</div></div>';
    document.getElementById('clients-section').innerHTML = logosHTML;
  } else {
    document.getElementById('clients-section').innerHTML = '';
  }

  // Colab text
  document.getElementById('colab-section').innerHTML =
    '<p>' + (lang === 'en'
      ? 'Collaborate with a fullstack programmer from Cilegon to build impactful web products, automation, and digital experiences.'
      : 'Berkolaborasi dengan fullstack programmer dari Cilegon untuk membangun produk web, automation, dan pengalaman digital yang berdampak.') +
    '</p>';

  // Services
  var svcItems = services.items || [];
  if(svcItems.length){
    var svcHTML = '<div class="services-grid">';
    var iconMap = {IconCode:'<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M16 18l6-6-6-6M8 6l-6 6 6 6"/></svg>',IconBrain:'<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 4a4 4 0 0 0-4 4c0 2 1 3 2 4l-1 3h6l-1-3c1-1 2-2 2-4a4 4 0 0 0-4-4z"/><path d="M9 15v3m6-3v3"/></svg>',IconRoute:'<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 17a2 2 0 1 0 4 0 2 2 0 0 0-4 0z"/><path d="M17 17a2 2 0 1 0 4 0 2 2 0 0 0-4 0z"/><path d="M7 17h10"/><path d="M14 8h4l2 4"/></svg>',IconBrandDocker:'<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="4" y="6" width="4" height="4" rx="1"/><rect x="9" y="6" width="4" height="4" rx="1"/><rect x="14" y="6" width="4" height="4" rx="1"/><rect x="9" y="11" width="4" height="4" rx="1"/><rect x="14" y="11" width="4" height="4" rx="1"/><path d="M4 11h5"/><path d="M18 11h2a2 2 0 0 1 2 2v1"/></svg>'};
    for(var si=0;si<svcItems.length;si++){
      var s = svcItems[si];
      var icon = iconMap[s.icon] || iconMap.IconCode;
      svcHTML += '<div class="service-card"><div class="service-icon">' + icon + '</div><h3>' + s.title + '</h3><p>' + (lang === 'en' ? (s.description_en||'') : (s.description_id||'')) + '</p></div>';
    }
    svcHTML += '</div>';
    document.getElementById('services-section').innerHTML = svcHTML;
  }

  // Projects
  var projItems = projects.items || [];
  if(projItems.length){
    var projHTML = '<h2 class="section-title">' + (lang === 'en' ? 'Open Source Projects' : 'Proyek Open Source') + '</h2><div class="projects-grid">';
    for(var pi=0;pi<projItems.length;pi++){
      var p = projItems[pi];
      var techs = p.technologies || [];
      var techHTML = '';
      var techColors = {Laravel:'bg:#fecaca;color:#f87171;',PHP:'bg:#e9d5ff;color:#a78bfa;',PostgreSQL:'bg:#c7d2fe;color:#818cf8;','Tailwind CSS':'bg:#cffafe;color:#22d3ee;',Docker:'bg:#bfdbfe;color:#60a5fa;',Bash:'bg:#fed7aa;color:#fb923c;','SQL Server':'bg:#d1d5db;color:#9ca3af;',Go:'bg:#cffafe;color:#22d3ee;',SonarQube:'bg:#bbf7d0;color:#4ade80;',Jenkins:'bg:#fecaca;color:#f87171;',JavaScript:'bg:#fef08a;color:#eab308;',HTML:'bg:#fed7aa;color:#fb923c;',CSS:'bg:#bfdbfe;color:#60a5fa;'};
      for(var ti=0;ti<techs.length;ti++){
        var t = techs[ti];
        var tc = techColors[t] || 'bg:#d1d5db;color:#9ca3af;';
        techHTML += '<span style="padding:2px 8px;font-size:12px;border-radius:999px;' + tc + 'background-opacity:0.2;">' + t + '</span>';
      }
      projHTML +=
        '<a href="' + (p.github_url||'#') + '" target="_blank" rel="noopener" class="project-card">' +
          '<div class="title-row">' +
            '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"/></svg>' +
            '<h3>' + (lang === 'en' ? (p.title_en||'') : (p.title_id||'')) + '</h3>' +
          '</div>' +
          '<p class="desc">' + (lang === 'en' ? (p.description_en||'') : (p.description_id||'')) + '</p>' +
          '<div class="techs">' + techHTML + '</div>' +
          '<div class="meta"><span><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg> 0</span><span>MIT License</span></div>' +
        '</a>';
    }
    projHTML += '</div>';
    document.getElementById('projects-section').innerHTML = projHTML;
  }

  // Custom sections
  var customHTML = '';
  var knownTypes = ['hero','clients','services','projects','cta'];
  for(var key in DATA){
    if(knownTypes.indexOf(key) === -1 && DATA.hasOwnProperty(key)){
      var cs = DATA[key];
      customHTML += '<section class="max-w custom-section"><h2>' + key + '</h2>';
      if(cs.items && cs.items.length){
        customHTML += '<div class="items-grid">';
        for(var ci=0;ci<cs.items.length;ci++){
          customHTML += '<div class="item-card">';
          for(var ck in cs.items[ci]){
            if(cs.items[ci].hasOwnProperty(ck)){
              customHTML += '<div class="field"><span class="k">' + ck + ':</span><span class="v">' + cs.items[ci][ck] + '</span></div>';
            }
          }
          customHTML += '</div>';
        }
        customHTML += '</div>';
      } else {
        customHTML += '<div class="plain-card">';
        for(var ck2 in cs){
          if(cs.hasOwnProperty(ck2)){
            customHTML += '<div class="field"><span class="k">' + ck2 + ':</span><span class="v">' + cs[ck2] + '</span></div>';
          }
        }
        customHTML += '</div>';
      }
      customHTML += '</section>';
    }
  }
  document.getElementById('custom-sections').innerHTML = customHTML;

  // CTA
  var ctaHeading = cta.heading || 'Tell me about your next project';
  var ctaEmail = cta.email || 'kurob@okuru.id';
  var ctaWhatsapp = cta.whatsapp || '628999702143';
  document.getElementById('cta-section').innerHTML =
    '<div class="cta-card">' +
      '<div class="mail-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M3 8l7.89 5.26a2 2 0 0 0 2.22 0L21 8M5 19h14a2 2 0 0 0 2-2V7a2 2 0 0 0-2-2H5a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2z"/></svg></div>' +
      '<h2>' + ctaHeading + '</h2>' +
      '<div class="btns">' +
        '<button class="btn-light" onclick="copyEmail()">' +
          '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 8l7.89 5.26a2 2 0 0 0 2.22 0L21 8M5 19h14a2 2 0 0 0 2-2V7a2 2 0 0 0-2-2H5a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2z"/></svg> Contact Me' +
        '</button>' +
        '<a href="https://wa.me/' + ctaWhatsapp + '" target="_blank" rel="noopener" class="btn-outline">' +
          '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/></svg> WhatsApp' +
        '</a>' +
      '</div>' +
    '</div>';

  document.getElementById('year').textContent = new Date().getFullYear();
}

window.setLang = function(l){
  lang = l;
  render();
};

window.scrollToEl = function(id){
  var el = document.getElementById(id);
  if(el) el.scrollIntoView({behavior:'smooth'});
};

window.copyEmail = function(){
  var email = (DATA.cta && DATA.cta.email) || 'kurob@okuru.id';
  navigator.clipboard.writeText(email).then(function(){
    var t = document.getElementById('toast');
    t.textContent = 'Email copied to clipboard!';
    t.style.display = 'block';
    setTimeout(function(){t.style.display = 'none';}, 2000);
  });
};

render();
})();
</script>
</body>
</html>`
}
