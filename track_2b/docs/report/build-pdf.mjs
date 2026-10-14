// Renders docs/report/report.html to the submission PDF with local Chrome.
// Usage (from track_2b): node docs/report/build-pdf.mjs [TeamName]
import { createRequire } from "node:module";
import { fileURLToPath, pathToFileURL } from "node:url";
import path from "node:path";

const here = path.dirname(fileURLToPath(import.meta.url));
const require = createRequire(path.join(here, "../../src/laporsiap_FE/package.json"));
const { chromium } = require("@playwright/test");

const team = (process.argv[2] || "LaporSiap").replace(/[^A-Za-z0-9_-]/g, "");
const out = path.join(here, `${team}_Report.pdf`);
const browser = await chromium.launch({ channel: "chrome" });
const page = await browser.newPage();
await page.goto(pathToFileURL(path.join(here, "report.html")).href, { waitUntil: "networkidle" });
await page.evaluate(() => document.fonts.ready);
await page.pdf({
  path: out,
  format: "A4",
  preferCSSPageSize: true,
  printBackground: true,
  displayHeaderFooter: true,
  headerTemplate: "<span></span>",
  footerTemplate: `<div style="width:100%;font-size:7px;color:#52656e;padding:0 15mm;display:flex;justify-content:space-between;font-family:Helvetica,Arial,sans-serif">
    <span>LaporSiap · Hack Apertus 2026 · Track 2B</span><span><span class="pageNumber"></span> / <span class="totalPages"></span></span></div>`,
});
await browser.close();
console.log(out);
