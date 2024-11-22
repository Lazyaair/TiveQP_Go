/// <reference types="vite/client" />
/// <reference types="element-plus/global" />

declare module '*.vue' {
  import { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
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