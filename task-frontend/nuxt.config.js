// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2026-04-08",

  devtools: { enabled: true },

  modules: ["@pinia/nuxt"],

  css: [
    "~/assets/css/main.scss",
  ],

  runtimeConfig: {
    public: {
      apiBase: process.env.API_BASE_URL || "http://localhost:8080",
    },
  },

  vite: {
    css: {
      preprocessorOptions: {
        scss: {
          quietDeps: true,
          silenceDeprecations: ["legacy-js-api", "import", "if-function", "global-builtin"],
        },
      },
    },
  },

  app: {
    head: {
      title: "TaskFlow — Gestión de Tareas",
      meta: [
        { charset: "utf-8" },
        { name: "viewport", content: "width=device-width, initial-scale=1" },
        { name: "description", content: "Aplicación profesional para gestión de tareas" },
      ],
      link: [
        { rel: "preconnect", href: "https://fonts.googleapis.com" },
        {
          rel: "stylesheet",
          href: "https://fonts.googleapis.com/css2?family=Space+Mono:wght@400;700&family=DM+Sans:wght@300;400;500;600&display=swap",
        },
        // Font Awesome desde CDN — evita el 404 de assets locales
        {
          rel: "stylesheet",
          href: "https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css",
          integrity: "sha512-DTOQO9RWCH3ppGqcWaEA1BIZOC6xxalwEsw9c2QQeAIftl+Vegovlnee1c9QX4TctnWMn13TZye+giMm8e2LwA==",
          crossorigin: "anonymous",
        },
      ],
    },
  },
});