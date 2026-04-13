import { defineConfig } from "vite";

export default defineConfig({
  build: {
    outDir: "dist",
    target: "es2020",
    minify: "esbuild",
    cssCodeSplit: true,
    reportCompressedSize: false,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes("node_modules/echarts")) return "echarts";
          if (id.includes("node_modules/zrender")) return "echarts";
          if (id.includes("node_modules/gsap")) return "gsap";
        },
      },
    },
  },
  esbuild: {
    legalComments: "none",
  },
});
