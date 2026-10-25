import { useEffect, useState } from 'react'
import { useI18n } from '../i18n/context'
import { clearSession, getAccessToken } from '../lib/auth'
import ApiKeys from './dashboard/ApiKeys'
import Overview from './dashboard/Overview'
import Settings from './dashboard/Settings'
import '../App.css'
import './dashboard.css'

type View = 'overview' | 'api-keys' | 'settings'

function OverviewIcon() {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.6"
      aria-hidden="true"
    >
      <rect x="3.5" y="3.5" width="7" height="7" />
      <rect x="13.5" y="3.5" width="7" height="7" />
      <rect x="3.5" y="13.5" width="7" height="7" />
      <rect x="13.5" y="13.5" width="7" height="7" />
    </svg>
  )
}

function KeyIcon() {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.6"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
    >
      <circle cx="8" cy="15" r="4" />
      <path d="M11 12l8.5-8.5M17 6l2.5 2.5" />
    </svg>
  )
}

function SettingsIcon() {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.6"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
    >
      <circle cx="12" cy="12" r="3" />
      <path d="M19.4 15a1.7 1.7 0 0 0 .34 1.87l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.7 1.7 0 0 0-1.87-.34 1.7 1.7 0 0 0-1.03 1.56V21a2 2 0 1 1-4 0v-.09A1.7 1.7 0 0 0 8.9 19.4a1.7 1.7 0 0 0-1.87.34l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06A1.7 1.7 0 0 0 4.6 15a1.7 1.7 0 0 0-1.56-1.03H3a2 2 0 1 1 0-4h.09A1.7 1.7 0 0 0 4.6 8.9a1.7 1.7 0 0 0-.34-1.87l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06A1.7 1.7 0 0 0 8.9 4.6a1.7 1.7 0 0 0 1.03-1.56V3a2 2 0 1 1 4 0v.09a1.7 1.7 0 0 0 1.03 1.56 1.7 1.7 0 0 0 1.87-.34l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.7 1.7 0 0 0-.34 1.87V9c.2.6.66 1.1 1.56 1.03H21a2 2 0 1 1 0 4h-.09a1.7 1.7 0 0 0-1.56 1.03z" />
    </svg>
  )
}

function Dashboard() {
  const { t } = useI18n()
  const authenticated = Boolean(getAccessToken())
  const [view, setView] = useState<View>('overview')

  useEffect(() => {
    if (!authenticated) {
      window.location.replace('/login')
    }
  }, [authenticated])

  if (!authenticated) {
    return null
  }

  function handleLogout() {
    clearSession()
    window.location.assign('/login')
  }

  return (
    <div className="dash-shell">
      <aside className="sidebar">
        <a
          className="brand sidebar-brand"
          href="/"
          aria-label={t('common.homeLink')}
        >
          <img
            className="brand-mark"
            src="/logo.svg"
            alt=""
            width="38"
            height="38"
          />
          <span className="brand-name">velin</span>
        </a>

        <nav className="sidebar-nav" aria-label={t('dashboard.nav.label')}>
          <button
            className="sidebar-link"
            type="button"
            aria-current={view === 'overview' ? 'page' : undefined}
            onClick={() => setView('overview')}
          >
            <OverviewIcon />
            {t('dashboard.nav.overview')}
          </button>

          <button
            className="sidebar-link"
            type="button"
            aria-current={view === 'api-keys' ? 'page' : undefined}
            onClick={() => setView('api-keys')}
          >
            <KeyIcon />
            {t('dashboard.nav.apiKeys')}
          </button>
        </nav>

        <div className="sidebar-foot">
          <div className="sidebar-meta">
            <a className="sidebar-home" href="/">
              {t('common.home')}
            </a>
            <button
              className="sidebar-logout"
              type="button"
              onClick={handleLogout}
            >
              {t('common.logOut')}
            </button>
          </div>

          <button
            className="sidebar-link sidebar-settings"
            type="button"
            aria-current={view === 'settings' ? 'page' : undefined}
            onClick={() => setView('settings')}
          >
            <SettingsIcon />
            {t('dashboard.nav.settings')}
          </button>
        </div>
      </aside>

      <main className="dash-content">
        {view === 'overview' ? (
          <Overview onOpenApiKeys={() => setView('api-keys')} />
        ) : view === 'api-keys' ? (
          <ApiKeys />
        ) : (
          <Settings />
        )}
      </main>
    </div>
  )
}

export default Dashboard
