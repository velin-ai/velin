import {
  useCallback,
  useEffect,
  useState,
  type FormEvent,
} from 'react'
import { useI18n } from '../../i18n/context'
import {
  createAPIKey,
  listAPIKeys,
  type APIKey,
} from '../../lib/apiKeys'
import { clearSession } from '../../lib/auth'

function isExpired(apiKey: APIKey) {
  return apiKey.expires_at
    ? new Date(apiKey.expires_at).getTime() <= Date.now()
    : false
}

function redirectToLogin() {
  clearSession()
  window.location.assign('/login')
}

function ApiKeys() {
  const { t, formatDate } = useI18n()
  const [apiKeys, setApiKeys] = useState<APIKey[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [title, setTitle] = useState('')
  const [expiresAt, setExpiresAt] = useState('')
  const [creating, setCreating] = useState(false)
  const [createError, setCreateError] = useState('')
  const [newKey, setNewKey] = useState('')
  const [copied, setCopied] = useState(false)

  const loadKeys = useCallback(async () => {
    try {
      const result = await listAPIKeys()

      if (result.status === 401) {
        redirectToLogin()

        return
      }

      if (result.status !== 200) {
        setError('dashboard.apiKeys.errors.loadFailed')

        return
      }

      setApiKeys(result.apiKeys)
      setError('')
    } catch {
      setError('dashboard.apiKeys.errors.serverUnreachable')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    // eslint-disable-next-line react/set-state-in-effect -- initial data fetch
    void loadKeys()
  }, [loadKeys])

  async function handleCreate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const trimmedTitle = title.trim()

    if (!trimmedTitle) {
      setCreateError('dashboard.apiKeys.errors.titleRequired')

      return
    }

    setCreating(true)
    setCreateError('')
    setNewKey('')

    try {
      const result = await createAPIKey(
        trimmedTitle,
        expiresAt ? new Date(expiresAt).toISOString() : null,
      )

      if (result.status === 401) {
        redirectToLogin()

        return
      }

      if (result.status === 201 && result.apiKey) {
        setTitle('')
        setExpiresAt('')
        setNewKey(result.apiKey)
        await loadKeys()

        return
      }

      if (result.status === 409) {
        setCreateError('dashboard.apiKeys.errors.titleExists')

        return
      }

      if (result.status === 400) {
        setCreateError('dashboard.apiKeys.errors.invalidInput')

        return
      }

      setCreateError('dashboard.apiKeys.errors.createFailed')
    } catch {
      setCreateError('dashboard.apiKeys.errors.serverUnreachable')
    } finally {
      setCreating(false)
    }
  }

  async function handleCopy() {
    try {
      await navigator.clipboard.writeText(newKey)
      setCopied(true)
      window.setTimeout(() => setCopied(false), 2000)
    } catch {
      setCopied(false)
    }
  }

  return (
    <section className="panel" aria-labelledby="api-keys-title">
      <header className="panel-head">
        <h1 id="api-keys-title">{t('dashboard.apiKeys.title')}</h1>
        <p className="panel-lede">{t('dashboard.apiKeys.lede')}</p>
      </header>

      <form className="key-form" onSubmit={handleCreate}>
        <div className="key-field">
          <label htmlFor="key-title">
            {t('dashboard.apiKeys.form.title')}
          </label>
          <input
            id="key-title"
            name="title"
            type="text"
            maxLength={50}
            required
            placeholder={t('dashboard.apiKeys.form.titlePlaceholder')}
            value={title}
            onChange={(event) => setTitle(event.target.value)}
          />
        </div>

        <div className="key-field">
          <label htmlFor="key-expires">
            {t('dashboard.apiKeys.form.expires')}
          </label>
          <input
            id="key-expires"
            name="expires_at"
            type="datetime-local"
            value={expiresAt}
            onChange={(event) => setExpiresAt(event.target.value)}
          />
        </div>

        <button className="key-create" type="submit" disabled={creating}>
          {creating
            ? t('dashboard.apiKeys.form.submitting')
            : t('dashboard.apiKeys.form.submit')}
        </button>
      </form>

      {createError ? (
        <p className="form-error" role="alert">
          {t(createError)}
        </p>
      ) : null}

      {newKey ? (
        <div className="key-reveal" role="status">
          <div className="key-reveal-head">
            <div>
              <p className="key-reveal-title">
                {t('dashboard.apiKeys.reveal.title')}
              </p>
              <p className="key-reveal-note">
                {t('dashboard.apiKeys.reveal.note')}
              </p>
            </div>
            <button
              className="key-reveal-close"
              type="button"
              onClick={() => setNewKey('')}
            >
              {t('dashboard.apiKeys.reveal.close')}
            </button>
          </div>

          <div className="key-reveal-value">
            <code>{newKey}</code>
            <button className="key-copy" type="button" onClick={handleCopy}>
              {copied
                ? t('dashboard.apiKeys.reveal.copied')
                : t('dashboard.apiKeys.reveal.copy')}
            </button>
          </div>
        </div>
      ) : null}

      {loading ? (
        <p className="keys-message">{t('dashboard.apiKeys.list.loading')}</p>
      ) : error ? (
        <div className="keys-message">
          <p>{t(error)}</p>
          <button
            className="keys-retry"
            type="button"
            onClick={() => {
              setLoading(true)
              setError('')
              void loadKeys()
            }}
          >
            {t('dashboard.apiKeys.list.retry')}
          </button>
        </div>
      ) : apiKeys.length === 0 ? (
        <p className="keys-message">{t('dashboard.apiKeys.list.empty')}</p>
      ) : (
        <ul className="keys-list">
          {apiKeys.map((apiKey) => {
            const expired = isExpired(apiKey)
            const created = formatDate(apiKey.created_at)
            const meta = apiKey.expires_at
              ? t('dashboard.apiKeys.list.meta', {
                  created,
                  expires: formatDate(apiKey.expires_at),
                })
              : t('dashboard.apiKeys.list.metaNoExpiry', { created })

            return (
              <li className="key-row" key={apiKey.id}>
                <div>
                  <p className="key-title">{apiKey.title}</p>
                  <p className="key-meta">{meta}</p>
                </div>

                <span
                  className={expired ? 'key-state is-expired' : 'key-state'}
                >
                  {expired
                    ? t('dashboard.apiKeys.list.expired')
                    : t('dashboard.apiKeys.list.active')}
                </span>
              </li>
            )
          })}
        </ul>
      )}
    </section>
  )
}

export default ApiKeys
