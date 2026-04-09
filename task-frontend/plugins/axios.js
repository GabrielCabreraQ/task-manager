// plugins/axios.js
import axios from "axios";

export default defineNuxtPlugin((nuxtApp) => {
  const config = useRuntimeConfig();

  const api = axios.create({
    baseURL: config.public.apiBase,
    headers: {
      "Content-Type": "application/json",
    },
    timeout: 10000,
  });

  // Request interceptor
  api.interceptors.request.use(
    (config) => config,
    (error) => Promise.reject(error)
  );

  // Response interceptor
  api.interceptors.response.use(
    (response) => response,
    (error) => {
      const msg =
        error.response?.data?.error ||
        error.response?.data?.message ||
        error.message ||
        "Error desconocido";
      return Promise.reject(new Error(msg));
    }
  );

  nuxtApp.provide("api", api);
});