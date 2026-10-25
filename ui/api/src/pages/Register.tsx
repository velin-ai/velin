import AuthForm from './AuthForm'
import AuthLayout from './AuthLayout'

function Register() {
  return (
    <AuthLayout
      title="Hesabınızı oluşturun."
      lede="İlk API anahtarınızı üretmek için bir hesap açın. Kurulum gerekmez."
      footer={
        <p className="auth-switch">
          Zaten hesabınız var mı? <a href="/login">Giriş yapın</a>
        </p>
      }
    >
      <AuthForm mode="register" />
    </AuthLayout>
  )
}

export default Register
