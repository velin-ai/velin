import { useEffect, useState } from 'react'
import '../App.css'

const fullAnswer = 'Merhaba! Size nasıl yardımcı olabilirim?'

function prefersReducedMotion() {
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

function Home() {
  const [typedAnswer, setTypedAnswer] = useState(() =>
    prefersReducedMotion() ? fullAnswer : '',
  )

  useEffect(() => {
    if (prefersReducedMotion()) {
      return
    }

    let index = 0

    const timer = window.setInterval(() => {
      index += 1
      setTypedAnswer(fullAnswer.slice(0, index))

      if (index >= fullAnswer.length) {
        window.clearInterval(timer)
      }
    }, 42)

    return () => window.clearInterval(timer)
  }, [])

  const isTyping = typedAnswer.length < fullAnswer.length

  return (
    <div className="page">
      <a className="skip-link" href="#main">
        İçeriğe geç
      </a>

      <header className="site-header">
        <div className="wrap">
          <a className="brand" href="/" aria-label="Velin ana sayfa">
            <img
              className="brand-mark"
              src="/logo.svg"
              alt=""
              width="46"
              height="46"
            />
            <span className="brand-name">velin</span>
            <span className="brand-divider" aria-hidden="true" />
            <span className="brand-tagline">Kendi AI Homelab'in</span>
          </a>

          <nav className="site-nav" aria-label="Ana menü">
            <a href="#nasil-calisir">Nasıl çalışır</a>
            <a href="#altyapi">Altyapı</a>
            <a className="nav-login" href="/login">
              Giriş
            </a>
          </nav>
        </div>
      </header>

      <main id="main">
        <section className="hero" aria-labelledby="hero-title">
          <div className="wrap">
            <div className="hero-copy">
              <h1 id="hero-title">
                Yapay zekâ sohbetini kendi altyapınızda çalıştırın.
              </h1>

              <p className="hero-lede">
                Velin; hesapları, API anahtarlarını ve sohbet geçmişini tek
                yerde toplar. Model çağrıları Ollama üzerinden gider, veriniz
                sunucunuzdan çıkmaz.
              </p>

              <div className="hero-actions">
                <a className="button" href="/register">
                  Hesap oluştur
                </a>
                <a className="text-link" href="#nasil-calisir">
                  Nasıl çalışır
                </a>
              </div>
            </div>

            <figure className="request" aria-label="Örnek API isteği ve yanıtı">
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
                    {': '}
                    <span className="json-string">"Merhaba"</span>
                    {'\n}'}
                  </code>
                </pre>
              </div>

              <div className="response">
                <div className="response-head">
                  <span>yanıt</span>
                  <span className="status">200</span>
                </div>

                <p className="response-text" aria-hidden="true">
                  {typedAnswer}
                  <span
                    className={isTyping ? 'caret' : 'caret is-hidden'}
                    aria-hidden="true"
                  />
                </p>
                <p className="visually-hidden">Örnek yanıt: {fullAnswer}</p>
              </div>
            </figure>
          </div>
        </section>

        <section className="steps" id="nasil-calisir" aria-labelledby="steps-title">
          <div className="wrap">
            <div className="section-head">
              <h2 id="steps-title">Üç adımda ilk yanıt</h2>
              <p>
                Hesabınızı açın, bir anahtar üretin ve ilk mesajınızı gönderin.
              </p>
            </div>

            <ol className="step-list">
              <li>
                <span className="step-number">01</span>
                <h3>Hesap oluşturun</h3>
                <p>
                  Kullanıcı adı ve parolanızı belirleyin. Oturumlar kısa ömürlü
                  JWT ile, yenileme jetonları ayrı olarak yönetilir.
                </p>
              </li>
              <li>
                <span className="step-number">02</span>
                <h3>API anahtarı üretin</h3>
                <p>
                  Canlı veya test ortamı için bir anahtar oluşturun. Anahtar
                  yalnızca bir kez gösterilir; veritabanında SHA-256 özeti
                  saklanır.
                </p>
              </li>
              <li>
                <span className="step-number">03</span>
                <h3>Mesaj gönderin</h3>
                <p>
                  Sohbeti açın ve mesajınızı gönderin. Son 20 mesaj bağlam olarak
                  modele iletilir, yanıt sohbet geçmişine kaydedilir.
                </p>
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
              <h2 id="infra-title">Veriniz sizin sunucunuzda kalır.</h2>
              <p>
                Mesajlar PostgreSQL'de saklanır, model çağrıları doğrudan
                sunucunuzdaki Ollama'ya gider. Her istek yalnızca kendi
                sohbetinizin geçmişini taşır; başka bir hesabın verisine
                erişilemez.
              </p>
            </div>

            <dl className="spec-list">
              <div>
                <dt>Sunucu</dt>
                <dd>Go</dd>
              </div>
              <div>
                <dt>Veri</dt>
                <dd>PostgreSQL</dd>
              </div>
              <div>
                <dt>Model</dt>
                <dd>Ollama</dd>
              </div>
              <div>
                <dt>Kimlik</dt>
                <dd>JWT ve API anahtarı</dd>
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
          <p className="footer-note">Kendi AI Homelab'in.</p>
          <p className="footer-copy">© 2026 Velin</p>
        </div>
      </footer>
    </div>
  )
}

export default Home
