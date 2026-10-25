import { useCallback, useEffect, useMemo, useState, type ReactNode } from 'react'
import en from './en.json'
import tr from './tr.json'
import {
  I18nContext,
  type I18nContextValue,
  type Language,
  type TranslationParams,
} from './context'

const dictionaries: Record<Language, unknown> = { en, tr }

const STORAGE_KEY = 'velin.language'
const DEFAULT_LANGUAGE: Language = 'en'

function lookup(dictionary: unknown, key: string) {
  let current: unknown = dictionary

  for (const part of key.split('.')) {
    if (typeof current !== 'object' || current === null) {
      return undefined
    }

    current = (current as Record<string, unknown>)[part]
  }

  return typeof current === 'string' ? current : undefined
}

function interpolate(template: string, params?: TranslationParams) {
  if (!params) {
    return template
  }

  return template.replace(/\{(\w+)\}/g, (match, name: string) => {
    const value = params[name]

    return value === undefined ? match : String(value)
  })
}

function readStoredLanguage(): Language {
  const stored = window.localStorage.getItem(STORAGE_KEY)

  return stored === 'tr' ? 'tr' : DEFAULT_LANGUAGE
}

function I18nProvider({ children }: { children: ReactNode }) {
  const [language, setLanguage] = useState<Language>(readStoredLanguage)

  const t = useCallback(
    (key: string, params?: TranslationParams) => {
      const template = lookup(dictionaries[language], key)

      return template ? interpolate(template, params) : key
    },
    [language],
  )

  const formatDate = useCallback(
    (value: string) =>
      new Intl.DateTimeFormat(language === 'tr' ? 'tr-TR' : 'en-US', {
        day: 'numeric',
        month: 'short',
        year: 'numeric',
      }).format(new Date(value)),
    [language],
  )

  useEffect(() => {
    document.documentElement.lang = language
    document.title = t('meta.title')
    window.localStorage.setItem(STORAGE_KEY, language)
  }, [language, t])

  const value = useMemo<I18nContextValue>(
    () => ({ language, setLanguage, t, formatDate }),
    [language, t, formatDate],
  )

  return <I18nContext.Provider value={value}>{children}</I18nContext.Provider>
}

export { I18nProvider }
