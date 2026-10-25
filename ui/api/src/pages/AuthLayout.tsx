import type { ReactNode } from 'react'
import { useI18n } from '../i18n/context'
import './auth.css'

type AuthLayoutProps = {
  title: string
  lede: string
  notice?: string
  children: ReactNode
  footer?: ReactNode
}

function AuthLayout({
  title,
  lede,
  notice,
  children,
  footer,
}: AuthLayoutProps) {
  const { t } = useI18n()

  return (
    <div className="auth-page">
      <aside className="auth-intro">
        <a className="auth-brand" href="/" aria-label={t('common.homeLink')}>
          <img src="/logo.svg" alt="" width="36" height="36" />
          <span>velin</span>
        </a>

        <div className="auth-intro-body">
          <p className="auth-tagline">{t('auth.tagline')}</p>
          <p className="auth-note">{t('auth.note')}</p>

          <dl className="auth-specs">
            <div>
              <dt>{t('auth.specs.session')}</dt>
              <dd>{t('auth.specs.sessionValue')}</dd>
            </div>
            <div>
              <dt>{t('auth.specs.refresh')}</dt>
              <dd>{t('auth.specs.refreshValue')}</dd>
            </div>
            <div>
              <dt>{t('auth.specs.key')}</dt>
              <dd>SHA-256</dd>
            </div>
          </dl>
        </div>

        <p className="auth-copyright">{t('common.copyright')}</p>
      </aside>

      <main className="auth-panel">
        <div className="auth-main">
          <h1>{title}</h1>
          <p className="auth-lede">{lede}</p>

          {notice ? (
            <p className="auth-notice" role="status">
              {notice}
            </p>
          ) : null}

          {children}

          {footer}
        </div>
      </main>
    </div>
  )
}

export default AuthLayout
