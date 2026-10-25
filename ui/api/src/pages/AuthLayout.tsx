import type { ReactNode } from 'react'
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
  return (
    <div className="auth-page">
      <aside className="auth-intro">
        <a className="auth-brand" href="/" aria-label="Velin ana sayfa">
          <img src="/logo.svg" alt="" width="36" height="36" />
          <span>velin</span>
        </a>

        <div className="auth-intro-body">
          <p className="auth-tagline">Kendi AI Homelab'in</p>
          <p className="auth-note">
            Hesaplar, API anahtarları ve sohbet geçmişi tek yerde. Veriniz kendi
            sunucunuzda kalır.
          </p>

          <dl className="auth-specs">
            <div>
              <dt>Oturum</dt>
              <dd>JWT, 30 dakika</dd>
            </div>
            <div>
              <dt>Yenileme</dt>
              <dd>30 gün</dd>
            </div>
            <div>
              <dt>API anahtarı</dt>
              <dd>SHA-256</dd>
            </div>
          </dl>
        </div>

        <p className="auth-copyright">© 2026 Velin</p>
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
