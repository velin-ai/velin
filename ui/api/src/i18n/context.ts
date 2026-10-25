import { createContext, useContext } from 'react'

export type Language = 'en' | 'tr'

export const LANGUAGES: Language[] = ['en', 'tr']

export type TranslationParams = Record<string, string | number>

export type I18nContextValue = {
  language: Language
  setLanguage: (language: Language) => void
  t: (key: string, params?: TranslationParams) => string
  formatDate: (value: string) => string
}

export const I18nContext = createContext<I18nContextValue | null>(null)

export function useI18n() {
  const context = useContext(I18nContext)

  if (!context) {
    throw new Error('useI18n must be used within an I18nProvider')
  }

  return context
}
