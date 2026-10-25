import { useI18n } from '../i18n/context'
import AuthForm from './AuthForm'
import AuthLayout from './AuthLayout'

function Register() {
  const { t } = useI18n()

  return (
    <AuthLayout
      title={t('auth.register.title')}
      lede={t('auth.register.lede')}
      footer={
        <p className="auth-switch">
          {t('auth.register.switchText')}{' '}
          <a href="/login">{t('auth.register.switchLink')}</a>
        </p>
      }
    >
      <AuthForm mode="register" />
    </AuthLayout>
  )
}

export default Register
