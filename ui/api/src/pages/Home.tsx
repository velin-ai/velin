import { useEffect, useState } from 'react'
import { useI18n } from '../i18n/context'
import '../App.css'

function prefersReducedMotion() {
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

type TypedAnswerProps = {
  answer: string
  hiddenAnswer: string
}

function TypedAnswer({ answer, hiddenAnswer }: TypedAnswerProps) {
  const [typedAnswer, setTypedAnswer] = useState(() =>
    prefersReducedMotion() ? answer : '',
  )

  useEffect(() => {
    if (prefersReducedMotion()) {
      return
    }

    let index = 0

    const timer = window.setInterval(() => {
      index += 1
      setTypedAnswer(answer.slice(0, index))

      if (index >= answer.length) {
        window.clearInterval(timer)
      }
    }, 42)

    return () => window.clearInterval(timer)
  }, [answer])

  const isTyping = typedAnswer.length < answer.length

  return (
    <>
      <p className="response-text" aria-hidden="true">
        {typedAnswer}
        <span
          className={isTyping ? 'caret' : 'caret is-hidden'}
          aria-hidden="true"
        />
      </p>
      <p className="visually-hidden">{hiddenAnswer}</p>
    </>
  )
}

function Home() {
  const { t, language } = useI18n()
  const answer = t('home.example.answer')

  return (
    <div className="page">
      <a className="skip-link" href="#main">
        {t('home.skipToContent')}
      </a>

      <header className="site-header">
        <div className="wrap">
          <a className="brand" href="/" aria-label={t('common.homeLink')}>
            <img
              className="brand-mark"
              src="/logo.svg"
              alt=""
              width="46"
              height="46"
            />
            <span className="brand-name">velin</span>
            <span className="brand-divider" aria-hidden="true" />
            <span className="brand-tagline">{t('common.brandTagline')}</span>
          </a>

          <nav className="site-nav" aria-label={t('home.nav.label')}>
            <a href="#nasil-calisir">{t('home.nav.howItWorks')}</a>
            <a href="#altyapi">{t('home.nav.infrastructure')}</a>
            <a className="nav-login" href="/login">
              {t('home.nav.logIn')}
            </a>
          </nav>
        </div>
      </header>

      <main id="main">
        <section className="hero" aria-labelledby="hero-title">
          <div className="wrap">
            <div className="hero-copy">
              <h1 id="hero-title">{t('home.hero.title')}</h1>

              <p className="hero-lede">{t('home.hero.lede')}</p>

              <div className="hero-actions">
                <a className="button" href="/register">
                  {t('home.hero.createAccount')}
                </a>
                <a className="text-link" href="#nasil-calisir">
                  {t('home.hero.howItWorks')}
                </a>
              </div>
            </div>

            <figure className="request" aria-label={t('home.example.label')}>
              <div className="request-bar">
                <span className="request-method">POST</span>
                <span className="request-path">
                  /api/v1/chats/42/messages
                </span>
              </div>

              <div className="request-body">
                <p>
                  <span className="token-key">Authorization</span>
                  <span>Bearer vln_live_8f3c…</span>
                </p>
                <p>
                  <span className="token-key">Content-Type</span>
                  <span>application/json</span>
                </p>

                <pre className="request-json">
                  <code>
                    {'{\n  '}
                    <span className="json-key">"model"</span>
                    {': '}
                    <span className="json-string">"llama3.1"</span>
                    {',\n  '}
                    <span className="json-key">"content"</span>
                    {': "'}
                    {t('home.example.message')}
                    {'"'}
                    {'\n}'}
                  </code>
                </pre>
              </div>

              <div className="response">
                <div className="response-head">
                  <span>{t('home.example.response')}</span>
                  <span className="status">200</span>
                </div>

                <TypedAnswer
                  key={language}
                  answer={answer}
                  hiddenAnswer={t('home.example.hiddenAnswer', { answer })}
                />
              </div>
            </figure>
          </div>
        </section>

        <section
          className="steps"
          id="nasil-calisir"
          aria-labelledby="steps-title"
        >
          <div className="wrap">
            <div className="section-head">
              <h2 id="steps-title">{t('home.steps.title')}</h2>
              <p>{t('home.steps.lede')}</p>
            </div>

            <ol className="step-list">
              <li>
                <span className="step-number">01</span>
                <h3>{t('home.steps.oneTitle')}</h3>
                <p>{t('home.steps.oneBody')}</p>
              </li>
              <li>
                <span className="step-number">02</span>
                <h3>{t('home.steps.twoTitle')}</h3>
                <p>{t('home.steps.twoBody')}</p>
              </li>
              <li>
                <span className="step-number">03</span>
                <h3>{t('home.steps.threeTitle')}</h3>
                <p>{t('home.steps.threeBody')}</p>
              </li>
            </ol>
          </div>
        </section>

        <section
          className="infrastructure"
          id="altyapi"
          aria-labelledby="infra-title"
        >
          <div className="wrap">
            <div className="infra-copy">
              <h2 id="infra-title">{t('home.infrastructure.title')}</h2>
              <p>{t('home.infrastructure.body')}</p>
            </div>

            <dl className="spec-list">
              <div>
                <dt>{t('home.infrastructure.specs.server')}</dt>
                <dd>Go</dd>
              </div>
              <div>
                <dt>{t('home.infrastructure.specs.data')}</dt>
                <dd>PostgreSQL</dd>
              </div>
              <div>
                <dt>{t('home.infrastructure.specs.model')}</dt>
                <dd>Ollama</dd>
              </div>
              <div>
                <dt>{t('home.infrastructure.specs.identity')}</dt>
                <dd>{t('home.infrastructure.specs.identityValue')}</dd>
              </div>
            </dl>
          </div>
        </section>
      </main>

      <footer className="site-footer">
        <div className="wrap">
          <div className="footer-brand">
            <img src="/logo.svg" alt="" width="26" height="26" />
            <span>velin</span>
          </div>
          <p className="footer-note">{t('home.footerNote')}</p>
          <p className="footer-copy">{t('common.copyright')}</p>
        </div>
      </footer>
    </div>
  )
}

export default Home
