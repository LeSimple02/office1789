import { createApp } from 'vue'
import { createPinia } from 'pinia'
import {gls} from '@/stores/global'
import App from './App.vue'
import router from './router'
import { createI18n } from 'vue-i18n'
import traduction from "./traduction.json"
import './assets/design-system.css'




const app = createApp(App)
app.use(createPinia())
app.use(router)

const currentLang = gls().lang || 'en'
const i18n = createI18n({
    locale: currentLang,
    fallbackLocale: 'en',
    messages: traduction,
    availableLocales: ['en', 'fr', 'es', 'zh', 'ja', 'it', 'de', 'pt', 'ru', 'ar', 'ko', 'nl'],
    warnHtmlMessage: false,
    escapeParameter: false,
})
console.log('[i18n] Loaded locale:', currentLang)
console.log('[i18n] Available locales:', i18n.global.availableLocales)

app.use(i18n)
app.mount('#app')
