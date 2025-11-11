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

const i18n = createI18n({
    locale: gls().lang,
    fallbackLocale: 'en',
    messages: traduction
})

app.use(i18n)
app.mount('#app')
