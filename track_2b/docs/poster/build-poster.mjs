// Renders docs/poster/poster.html to an A3 PDF and a high-resolution PNG with local Chrome.
// Usage (from track_2b): node docs/poster/build-poster.mjs
import { createRequire } from "node:module";
import { fileURLToPath, pathToFileURL } from "node:url";
import path from "node:path";

const here = path.dirname(fileURLToPath(import.meta.url));
const require = createRequire(path.join(here, "../../src/laporsiap_FE/package.json"));
const { chromium } = require("@playwright/test");

const browser = await chromium.launch({ channel: "chrome" });
// A3 at 96 dpi is 1123 × 1587 CSS px; scale 3 gives a ~290 dpi PNG.
const page = await browser.newPage({ viewport: { width: 1123, height: 1587 }, deviceScaleFactor: 3 });
await page.goto(pathToFileURL(path.join(here, "poster.html")).href, { waitUntil: "networkidle" });
await page.evaluate(() => document.fonts.ready);
await page.pdf({ path: path.join(here, "LaporSiap-Poster-A3.pdf"), preferCSSPageSize: true, printBackground: true });
await page.screenshot({ path: path.join(here, "LaporSiap-Poster.png") });
await browser.close();
console.log(path.join(here, "LaporSiap-Poster-A3.pdf"));
