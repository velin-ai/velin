import AuthForm from './AuthForm'
import AuthLayout from './AuthLayout'

function Login() {
  const registered =
    new URLSearchParams(window.location.search).get('registered') === '1'

  return (
    <AuthLayout
      title="Tekrar hoş geldiniz."
      lede="Hesabınıza giriş yapın; API anahtarlarınız ve sohbetleriniz sizi bekliyor."
      notice={
        registered
          ? 'Hesabınız oluşturuldu. Şimdi giriş yapabilirsiniz.'
          : undefined
      }
      footer={
        <p className="auth-switch">
          Hesabınız yok mu? <a href="/register">Hesap oluşturun</a>
        </p>
      }
    >
      <AuthForm mode="login" />
    </AuthLayout>
  )
}

export default Login
