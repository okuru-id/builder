import { createApp } from 'vue'
import './style.css'
// Force Tailwind to scan inspector arbitrary-value safelist (w/h/radius px).
// Side-effect import; the module exports a string of class literals.
import '@/safelist-arbitrary'
import App from './App.vue'
import router from './router'

createApp(App).use(router).mount('#app')
