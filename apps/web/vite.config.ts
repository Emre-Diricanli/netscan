import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// Vite 7 + React 19
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    open: false,
  },
  preview: { port: 5173 },
});