import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import { ViteImageOptimizer } from "vite-plugin-image-optimizer";
import AutoImport from "unplugin-auto-import/vite";
import Components from "unplugin-vue-components/vite";
import { ElementPlusResolver } from "unplugin-vue-components/resolvers";
import ElementPlus from "unplugin-element-plus/vite";

export default defineConfig({
  plugins: [
    vue(),
    ViteImageOptimizer({
      // 图片压缩
      png: { quality: 100 },
      jpeg: { quality: 100 },
      webp: { quality: 100 },
    }),
    AutoImport({
      resolvers: [ElementPlusResolver()],
    }),
    Components({
      resolvers: [ElementPlusResolver({ importStyle: "sass" })],
    }),
    ElementPlus({
      defaultLocale: "zh-cn",
      useSource: true,
    }),
  ],
  resolve: {
    alias: {
      "@": "/src",
    },
  },
  server: {
    proxy: {
      "/api": {
        target:
          process.env.VITE_API_BASE_URL ||
          "https://asynx.internal.asynclab.club",
        changeOrigin: true,
        secure: false,
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
    },
  },
  build: {
    minify: "terser"
  },
});
