import { existsSync, promises as fs } from "fs";
import path from "path";

import { addTdPrefix, cleanText, DOM_RULE_MAP, getDirname } from "@tdesign-mcp-server/common";
import { TD_DOCS_OUTPUT_DIR } from "../utils/path";

const __dirname = getDirname(import.meta.url);
const SNAPSHOT_PATH = path.join(__dirname, "./snapshots/web/__snapshots__/setup.test.ts.snap");

(async function main() {
  console.log(`--------(web)--------`);
  if (!existsSync(SNAPSHOT_PATH)) {
    console.warn(`✗ Snapshot file not found.`);
    return;
  }
  await extractAllDom();
  console.log(`✓ Extract DOM successfully`);
})();

async function extractAllDom() {
  const snapshotContent = await fs.readFile(SNAPSHOT_PATH, "utf-8");
  const regex = /exports\[`[^`]+ \[([^\]]+)\] [^\]]+\`] = `([\s\S]+?)`;/g;

  let match;
  while ((match = regex.exec(snapshotContent)) !== null) {
    const compName = match[1];
    let compDom = match[2];
    compDom = cleanText(compDom, DOM_RULE_MAP);

    const compDomPath = path.join(TD_DOCS_OUTPUT_DIR, addTdPrefix("web-dom"), `${compName}.html`);
    if (!existsSync(compDomPath)) {
      await fs.mkdir(path.dirname(compDomPath), { recursive: true });
    }
    await fs.writeFile(compDomPath, compDom);
  }
}
