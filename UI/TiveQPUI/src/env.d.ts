/// <reference types="vite/client" />
/// <reference types="element-plus/global" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

// particles.js 类型声明
declare global {
  interface Window {
    particlesJS: (id: string, config: any) => void
  }
}

declare module 'leaflet' {
  export interface MapOptions {
    preferCanvas?: boolean
  }
}

interface ImportMetaEnv {
  readonly VITE_APP_TITLE: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
} 