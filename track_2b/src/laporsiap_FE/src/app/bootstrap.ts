import { CaseWorkflow } from '../application/caseWorkflow'
import { BrowserSessionVault } from '../infrastructure/browser/sessionVault'
import { BrowserFileSaver } from '../infrastructure/browser/fileSaver'
import { BrowserLanguagePreferences } from '../infrastructure/browser/languagePreferences'
import { HttpCaseRepository } from '../infrastructure/http/httpCaseRepository'
import { createAppStore } from '../presentation/state/store'
const sessions = new BrowserSessionVault()
export const preferences = new BrowserLanguagePreferences()
export const workflow = new CaseWorkflow(new HttpCaseRepository(sessions), sessions, new BrowserFileSaver())
export const store = createAppStore(sessions.records(), preferences.read())
