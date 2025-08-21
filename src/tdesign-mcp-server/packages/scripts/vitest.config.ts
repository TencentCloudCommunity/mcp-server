import path from "path";
import { fileURLToPath } from "url";

import react from "@vitejs/plugin-react";
import { defineConfig } from "vitest/config";

// Vitest 暂时不支持加载外部（common）的 TS 文件，所以重复维护一份 getDirName 写法
const __dirname = path.dirname(fileURLToPath(import.meta.url));

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@tdesign-mcp-server/docs-react": path.resolve(__dirname, "../docs/tdesign-react"),
      // 省得 Docs 目录下也安装依赖
      "tdesign-react": path.resolve(__dirname, "./node_modules/tdesign-react"),
      "tdesign-icons-react": path.resolve(__dirname, "./node_modules/tdesign-icons-react")
    }
  },
  test: {
    environment: "jsdom"
  }
});
