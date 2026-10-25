import { LANGUAGES, useI18n } from '../../i18n/context'

function Settings() {
  const { language, setLanguage, t } = useI18n()

  return (
    <section className="panel" aria-labelledby="settings-title">
      <header className="panel-head">
        <h1 id="settings-title">{t('dashboard.settings.title')}</h1>
        <p className="panel-lede">{t('dashboard.settings.lede')}</p>
      </header>

      <div className="settings-section">
        <h2>{t('dashboard.settings.language.title')}</h2>
        <p className="settings-description">
          {t('dashboard.settings.language.description')}
        </p>

        <fieldset className="language-options">
          <legend className="visually-hidden">
            {t('dashboard.settings.language.title')}
          </legend>

          {LANGUAGES.map((code) => (
            <label
              className="language-option"
              data-checked={language === code}
              key={code}
            >
              <input
                type="radio"
                name="language"
                value={code}
                checked={language === code}
                onChange={() => setLanguage(code)}
              />
              <span className="language-code">{code}</span>
              <span className="language-name">
                {t(`dashboard.settings.language.${code}`)}
              </span>
              <span className="language-marker" aria-hidden="true" />
            </label>
          ))}
        </fieldset>
      </div>
    </section>
  )
}

export default Settings
