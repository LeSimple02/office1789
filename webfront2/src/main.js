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
    warnHtmlMessage: false,
    escapeParameter: false,
})
console.log('[i18n] Loaded locale:', currentLang)

app.use(i18n)
app.mount('#app')
