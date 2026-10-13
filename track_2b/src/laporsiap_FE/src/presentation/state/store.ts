import { configureStore } from '@reduxjs/toolkit'
import type { CaseAccess, Language } from '../../domain/models'
import { workspaceSlice, initialWorkspace } from './workspaceSlice'
import { localeSlice } from './localeSlice'
export const createAppStore = (records: CaseAccess[] = [], language: Language = 'en') => configureStore({
 reducer: { workspace: workspaceSlice.reducer, locale: localeSlice.reducer },
 preloadedState: { workspace: initialWorkspace(records), locale: { language } },
 // Case text stays in memory. Never expose case state/actions via DevTools.
 devTools: false,
})
export type AppStore = ReturnType<typeof createAppStore>
export type RootState = ReturnType<AppStore['getState']>
export type AppDispatch = AppStore['dispatch']
