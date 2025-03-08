import vuetify, { transformAssetUrls } from "vite-plugin-vuetify";
export default defineNuxtConfig({
  compatibilityDate: "2024-11-01",
  devtools: { enabled: true },
  build: {
    transpile: ["vuetify"],
  },
  modules: [
    (_options, nuxt) => {
      nuxt.hooks.hook("vite:extendConfig", (config) => {
        // @ts-expect-error I dont remember why this is needed
        config.plugins.push(vuetify({ autoImport: true }));
      });
    },
    "nuxt-echarts",
  ],
  vite: {
    vue: {
      template: {
        transformAssetUrls,
      },
    },
  },
  echarts: {
    charts: ["BarChart"],
    components: ["DatasetComponent", "GridComponent", "TooltipComponent"],
  },
  routeRules: {
    "/api/v1/**": { proxy: `${process.env.BASE_URL}/api/v1/**` },
  },
});
