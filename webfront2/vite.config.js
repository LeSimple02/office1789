import { fileURLToPath, URL } from 'node:url'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'


// https://vite.dev/config/
export default defineConfig((mode)=>{
	const env = loadEnv(mode, process.cwd(), '');
	return{
  // Use relative base so built assets resolve under Cordova file:// scheme
  base: './',
 
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
	
  },
}
})
