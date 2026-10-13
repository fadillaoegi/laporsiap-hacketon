import { createSlice } from '@reduxjs/toolkit'
import type { PayloadAction } from '@reduxjs/toolkit'
import type { Language } from '../../domain/models'
export const localeSlice = createSlice({ name: 'locale', initialState: { language: 'en' as Language }, reducers: { setLanguage(state, action: PayloadAction<Language>) { state.language = action.payload } } })
