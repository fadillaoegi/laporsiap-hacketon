import type { Language } from '../../domain/models'
import type { LanguagePreferences } from '../../application/ports'
export class BrowserLanguagePreferences implements LanguagePreferences {
 read(): Language { try { return localStorage.getItem('laporsiap-language') === 'id' ? 'id' : 'en' } catch { return 'en' } }
 write(language: Language) { try { localStorage.setItem('laporsiap-language', language) } catch { /* Storage may be disabled. */ } }
}
