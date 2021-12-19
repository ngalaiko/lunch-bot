import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import fs from 'fs'

export default defineConfig({
  server: {
    https: {
      key: fs.readFileSync('./.cert/key.pem'),
      cert: fs.readFileSync('./.cert/cert.pem')
    }
  },
  plugins: [svelte()]
})
