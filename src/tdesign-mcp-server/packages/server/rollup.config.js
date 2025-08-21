import commonjs from "@rollup/plugin-commonjs";
import json from "@rollup/plugin-json";
import resolve from "@rollup/plugin-node-resolve";
import replace from "@rollup/plugin-replace";
import terser from "@rollup/plugin-terser";
import typescript from "rollup-plugin-typescript2";

export default [
  {
    input: "./src/index.ts",
    output: [
      {
        dir: "dist",
        format: "esm",
        entryFileNames: "server.js",
        banner: "#!/usr/bin/env node",
        sourcemap: false,
        inlineDynamicImports: true
      }
    ],
    plugins: [
      replace({
        preventAssignment: true,
        "process.env.NODE_ENV": JSON.stringify(process.env.NODE_ENV)
      }),
      resolve(),
      commonjs(),
      json(),
      typescript({
        tsconfigOverride: {
          compilerOptions: {
            removeComments: true
          }
        }
      }),
      terser()
    ]
  }
];
