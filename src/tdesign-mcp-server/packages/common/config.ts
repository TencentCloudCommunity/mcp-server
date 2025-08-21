import path from "path";
import { getDirname } from "./toolkit";

const __dirname = getDirname(import.meta.url);

if (process.env.NODE_ENV !== "production") {
  import("dotenv").then((dotenv) => {
    dotenv.config({ path: path.join(__dirname, "../../.env"), quiet: true });
  });
}

interface EnvConfig {
  DOCS_SOURCE_MODE?: "online" | "local";
}

type EnvKey = keyof EnvConfig;

export const getEnvConfig = <K extends EnvKey>(key: K) => {
  const value = process.env[key];
  return value as EnvConfig[K]; // 类型断言确保匹配
};
