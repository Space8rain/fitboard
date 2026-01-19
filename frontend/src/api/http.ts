import axios, { type InternalAxiosRequestConfig } from "axios";

export const api = axios.create({
  baseURL: "http://localhost:3000/api",
  timeout: 5000,
});

// Добавляем Telegram initData в каждый запрос
api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const tg = window.Telegram?.WebApp;

  if (tg?.initData) {
    config.headers.set("X-Telegram-Init-Data", tg.initData);
  }

  return config;
});