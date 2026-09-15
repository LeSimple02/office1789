<script setup>
import {gls} from "@/stores/global"
import { RouterLink, RouterView} from 'vue-router'
import {ref, computed} from "vue"
import router from "@/router/index"
import { useDark, useToggle } from '@vueuse/core'
import verifback from "./views/backwork.js"

const isDark = useDark()
const toggleDark = useToggle(isDark)
const menuOpen = ref(false)

const languageNames = {
  'en': 'English', 'fr': 'Français', 'es': 'Español', 'zh': '中文',
  'ja': '日本語', 'it': 'Italiano', 'de': 'Deutsch', 'pt': 'Português',
  'ru': 'Русский', 'ar': 'العربية', 'ko': '한국어', 'nl': 'Nederlands'
}

const toggleMenu = () => {
  menuOpen.value = !menuOpen.value
  document.body.classList.toggle('menu-open', menuOpen.value)
}
const closeMenu = () => {
  menuOpen.value = false
  document.body.classList.remove('menu-open')
}

let menu1 = [{name : "Home", link : "/"},{name : "Login", link : "/login"},{name : "About", link : "/about"}]
let menu2 = [{name : "Mail", link : "/mail"}, {name : "Drive", link : "/drive"}, {name : "Chat", link : "/chat"}, {name : "Calendar", link : "/calendar"}, {name : "Account", link : "/account"}]
verifback()

let menu = computed(() => gls().log == 1 ? menu2 : menu1)
let popup = ref(0)
function decof(){
  localStorage.setItem("log", 0)
  gls().log = localStorage.getItem("log")
  popup.value = 0
  router.push('/login')
}
</script>

<template>
  <header>
    <div class="header-inner">
      <img class="logo" v-if="!isDark" src="@/assets/logo.png" alt="Office1789" />
      <img class="logo" v-if="isDark" src="@/assets/logol.png" alt="Office1789" />
      <button class="burger" @click="toggleMenu" aria-label="Menu">☰</button>
      <nav :class="{ open: menuOpen }">
        <RouterLink v-for="key in menu" :key="key.name" :to="key.link" @click="closeMenu">
          {{ $t(key.name) }}
        </RouterLink>
        <div class="nav-settings">
          <select v-model="$i18n.locale" @change="save()" aria-label="Language">
            <option v-for="locale in $i18n.availableLocales" :key="`locale-${locale}`" :value="locale">
              {{ languageNames[locale] || locale }}
            </option>
          </select>
          <button @click="toggleDark()" id="toggle" aria-label="Toggle theme">
            {{ isDark ? "☾" : "☀" }}
          </button>
        </div>
        <button id="bdeco" v-if="gls().log == 1" @click="popup = 1; closeMenu()" :title="$t('logout')">
          {{ $t('logout') }}
        </button>
      </nav>
    </div>
  </header>

  <div v-if="popup==1" id="sb" @click="popup=0"></div>
  <div v-if="popup==1" id="pop">
    <p>{{ $t('confirmLogout') }}</p>
    <div class="pop-actions">
      <button @click="popup=0" id="Cancel">{{ $t('cancel') }}</button>
      <button @click="decof" id="Yes">{{ $t('yes') }}</button>
    </div>
  </div>

  <main><RouterView /></main>

  <footer>
    <RouterLink to="/legalesmentions">{{ $t('legal') }}</RouterLink>
    <RouterLink to="/contact">Contact</RouterLink>
  </footer>
</template>

<script>
export default { methods: { save() { localStorage.setItem("lang", this.$i18n.locale) } } }
</script>

<style>
/* ===== Header ===== */
header {
  position: sticky; top: 0; z-index: 10000;
  background: var(--bg-glass);
  backdrop-filter: blur(16px) saturate(180%);
  border-bottom: 1px solid var(--border);
  transition: border-color var(--tb), box-shadow var(--tb);
}
header:hover { border-bottom-color: var(--border-strong); box-shadow: var(--sh-sm); }

.header-inner {
  max-width: var(--content-max); margin: 0 auto;
  display: flex; align-items: center; gap: 20px;
  padding: 0 var(--sp-5); height: var(--header-h);
}

.logo {
  width: 34px; height: 34px;
  border-radius: var(--r-sm);
  flex-shrink: 0;
  transition: transform var(--tb);
  animation: fadeInDown .4s var(--ease) both;
}
.logo:hover { transform: scale(1.08) rotate(-3deg); }

/* ===== Nav ===== */
nav {
  display: flex; align-items: center; gap: 4px; flex: 1;
  animation: fadeInDown .4s var(--ease) .05s both;
}
nav a {
  padding: 7px 14px; font-size: 14px; font-weight: 500;
  color: var(--text-secondary); text-decoration: none;
  border-radius: var(--r-sm);
  transition: color var(--tb), background var(--tb), transform var(--tb);
  white-space: nowrap; position: relative;
}
nav a::after {
  content: ''; position: absolute; bottom: 2px; left: 50%; width: 0; height: 2px;
  background: var(--ink); border-radius: 2px; transform: translateX(-50%);
  transition: width var(--tb);
}
nav a:hover { color: var(--text-primary); background: var(--bg-tertiary); }
nav a:hover::after { width: 60%; }
nav a.router-link-active {
  color: var(--paper); background: var(--ink);
}
nav a.router-link-active::after { display: none; }
nav a.router-link-active:hover { background: var(--ink-soft); }

.nav-settings { display: flex; align-items: center; gap: 8px; margin-left: auto; }

/* ===== Select ===== */
select {
  padding: 6px 28px 6px 12px; font-size: 13px; font-weight: 500;
  color: var(--text-primary); background: var(--bg-primary);
  border: 1px solid var(--border-strong); border-radius: var(--r-sm);
  cursor: pointer; transition: border-color var(--tb);
  color-scheme: light; appearance: none; -webkit-appearance: none;
  background-image: url("data:image/svg+xml;charset=utf-8,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%23737373' stroke-width='2'%3E%3Cpolyline points='6 9 12 15 18 9'/%3E%3C/svg%3E");
  background-repeat: no-repeat; background-position: right 8px center;
}
select:hover { border-color: var(--ink); }
select:focus { outline: none; border-color: var(--ink); }
select option { background: var(--bg-primary); color: var(--text-primary); }
.dark select { color-scheme: dark; }

/* ===== Toggle ===== */
#toggle {
  width: 36px; height: 36px;
  display: flex; align-items: center; justify-content: center;
  font-size: 18px; line-height: 1;
  background: var(--bg-primary); border: 1px solid var(--border-strong);
  border-radius: var(--r-sm); color: var(--text-secondary);
  cursor: pointer; transition: all var(--tb);
}
#toggle:hover { color: var(--text-primary); border-color: var(--ink); transform: scale(1.08) rotate(8deg); }

/* ===== Logout ===== */
#bdeco {
  padding: 7px 18px; font-size: 13px; font-weight: 600;
  color: var(--paper); background: var(--ink); border: 1px solid var(--ink);
  border-radius: var(--r-sm); cursor: pointer;
  transition: all var(--tb); white-space: nowrap;
}
#bdeco:hover { background: var(--ink-soft); border-color: var(--ink-soft); transform: translateY(-2px); box-shadow: var(--sh-sm); }
#bdeco:active { transform: translateY(0); }

/* ===== Burger ===== */
.burger {
  display: none; background: none; border: none; font-size: 22px;
  color: var(--text-primary); cursor: pointer; padding: 4px 8px;
  margin-left: auto; z-index: 10001;
}

/* ===== Main ===== */
main { min-height: calc(100vh - var(--header-h) - 56px); }

/* ===== Modal ===== */
#sb {
  position: fixed; inset: 0; z-index: 999;
  background: var(--bg-overlay); backdrop-filter: blur(6px);
  animation: fadeIn .25s var(--ease) both;
}
#pop {
  position: fixed; left: 50%; top: 50%; transform: translate(-50%, -50%);
  width: calc(100% - 32px); max-width: 380px;
  padding: 32px 28px; background: var(--bg-primary);
  border: 1px solid var(--border); border-radius: var(--r-lg);
  box-shadow: var(--sh-xl); z-index: 1000; text-align: center;
  animation: popIn .3s var(--ease-spring) both;
}
#pop p { font-size: 1.0625rem; font-weight: 500; color: var(--text-primary); margin: 0 0 24px 0; line-height: 1.5; }
.pop-actions { display: flex; gap: 10px; justify-content: center; }
#pop .pop-actions #Cancel {
  flex: 1; padding: 11px 20px; font-size: .9375rem; font-weight: 600;
  color: var(--text-primary); background: var(--bg-primary);
  border: 1px solid var(--border-strong); border-radius: var(--r);
  cursor: pointer; transition: all var(--tb);
}
#pop .pop-actions #Cancel:hover { border-color: var(--ink); background: var(--bg-tertiary); transform: translateY(-1px); }
#pop .pop-actions #Yes {
  flex: 1; padding: 11px 20px; font-size: .9375rem; font-weight: 600;
  color: var(--paper); background: var(--ink); border: 1px solid var(--ink);
  border-radius: var(--r); cursor: pointer; transition: all var(--tb);
}
#pop .pop-actions #Yes:hover { background: var(--ink-soft); border-color: var(--ink-soft); transform: translateY(-1px); box-shadow: var(--sh-sm); }

/* ===== Footer ===== */
footer {
  display: flex; justify-content: center; gap: 4px;
  padding: 16px var(--sp-5);
  border-top: 1px solid var(--border); background: var(--bg-secondary);
  animation: fadeIn .5s var(--ease) both;
}
footer a {
  padding: 5px 12px; font-size: 13px; font-weight: 500;
  color: var(--text-secondary); border-radius: var(--r-sm);
  transition: color var(--tb), background var(--tb);
}
footer a:hover { color: var(--text-primary); background: var(--bg-tertiary); }

/* ===== Responsive ===== */
@media (max-width: 768px) {
  .header-inner { padding: 0 var(--sp-4); gap: 12px; }
  .burger { display: block; }
  nav {
    position: fixed; top: var(--header-h); left: 0; right: 0; bottom: 0;
    flex-direction: column; align-items: stretch; gap: 4px;
    padding: 20px 16px; background: var(--bg-primary);
    border-top: 1px solid var(--border);
    transform: translateX(-100%); transition: transform .3s var(--ease);
    z-index: 9999; overflow-y: auto; box-shadow: var(--sh-lg);
  }
  nav.open { transform: translateX(0); }
  body.menu-open { overflow: hidden; }
  nav a {
    padding: 14px 16px; font-size: 16px; border-radius: var(--r);
    border-bottom: 1px solid var(--border); animation: fadeInLeft .3s var(--ease) both;
  }
  nav a::after { display: none; }
  nav a.router-link-active { border-bottom-color: transparent; }
  .nav-settings {
    flex-direction: column; align-items: stretch; gap: 10px;
    margin: 12px 0 0; padding-top: 16px; border-top: 1px solid var(--border);
  }
  .nav-settings select, #toggle { width: 100%; height: 44px; padding: 6px 12px; }
  #toggle { justify-content: center; }
  #bdeco { width: 100%; padding: 14px; font-size: 16px; margin-top: 8px; }
}
@media (max-width: 480px) {
  .header-inner { padding: 0 16px; }
  .logo { width: 30px; height: 30px; }
}
</style>
