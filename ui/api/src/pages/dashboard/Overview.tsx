import { useI18n } from '../../i18n/context'

type OverviewProps = {
  onOpenApiKeys: () => void
}

function Overview({ onOpenApiKeys }: OverviewProps) {
  const { t } = useI18n()

  return (
    <section className="panel" aria-labelledby="overview-title">
      <header className="panel-head">
        <h1 id="overview-title">{t('dashboard.overview.title')}</h1>
        <p className="panel-lede">{t('dashboard.overview.lede')}</p>
      </header>

      <dl className="overview-session">
        <div>
          <dt>{t('dashboard.overview.specs.session')}</dt>
          <dd>{t('dashboard.overview.specs.sessionValue')}</dd>
        </div>
        <div>
          <dt>{t('dashboard.overview.specs.refresh')}</dt>
          <dd>{t('dashboard.overview.specs.refreshValue')}</dd>
        </div>
        <div>
          <dt>{t('dashboard.overview.specs.storage')}</dt>
          <dd>{t('dashboard.overview.specs.storageValue')}</dd>
        </div>
      </dl>

      <button className="overview-action" type="button" onClick={onOpenApiKeys}>
        {t('dashboard.overview.action')}
      </button>
    </section>
  )
}

export default Overview
