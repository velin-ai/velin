import { useI18n } from '../i18n/context'
import AuthForm from './AuthForm'
import AuthLayout from './AuthLayout'

function Login() {
  const { t } = useI18n()

  const registered =
    new URLSearchParams(window.location.search).get('registered') === '1'

  return (
    <AuthLayout
      title={t('auth.login.title')}
      lede={t('auth.login.lede')}
      notice={registered ? t('auth.login.notice') : undefined}
      footer={
        <p className="auth-switch">
          {t('auth.login.switchText')}{' '}
          <a href="/register">{t('auth.login.switchLink')}</a>
        </p>
      }
    >
      <AuthForm mode="login" />
    </AuthLayout>
  )
}

export default Login
