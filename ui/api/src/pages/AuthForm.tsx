import { useState, type FormEvent } from 'react'
import { loginUser, registerUser } from '../lib/auth'

type AuthFormProps = {
  mode: 'login' | 'register'
}

function errorMessage(
  mode: AuthFormProps['mode'],
  status: number,
  code?: string,
) {
  if (mode === 'login') {
    if (status === 401) {
      return 'Kullanıcı adı veya parola hatalı.'
    }

    if (status === 400) {
      return 'Kullanıcı adı ve parola gerekli.'
    }

    return 'Giriş yapılamadı. Lütfen tekrar deneyin.'
  }

  if (status === 409) {
    return 'Bu kullanıcı adı zaten kullanılıyor.'
  }

  if (status === 400) {
    if (code === 'username_too_long') {
      return 'Kullanıcı adı en fazla 20 karakter olabilir.'
    }

    return 'Kullanıcı adı ve parola gerekli.'
  }

  return 'Kayıt tamamlanamadı. Lütfen tekrar deneyin.'
}

function AuthForm({ mode }: AuthFormProps) {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const isLogin = mode === 'login'

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    setError('')
    setLoading(true)

    try {
      const result = isLogin
        ? await loginUser(username.trim(), password)
        : await registerUser(username.trim(), password)

      if (isLogin && result.status === 200) {
        window.location.assign('/')

        return
      }

      if (!isLogin && result.status === 201) {
        window.location.assign('/login?registered=1')

        return
      }

      setError(errorMessage(mode, result.status, result.error))
    } catch {
      setError('Sunucuya ulaşılamadı. Bağlantınızı kontrol edin.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form className="auth-form" onSubmit={handleSubmit}>
      <div className="auth-field">
        <label htmlFor="username">Kullanıcı adı</label>
        <input
          id="username"
          name="username"
          type="text"
          autoComplete="username"
          autoCapitalize="none"
          autoFocus
          spellCheck={false}
          maxLength={20}
          required
          value={username}
          onChange={(event) => setUsername(event.target.value)}
          aria-invalid={error ? true : undefined}
        />
      </div>

      <div className="auth-field">
        <label htmlFor="password">Parola</label>

        <div className="auth-password">
          <input
            id="password"
            name="password"
            type={showPassword ? 'text' : 'password'}
            autoComplete={isLogin ? 'current-password' : 'new-password'}
            required
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            aria-invalid={error ? true : undefined}
          />

          <button
            className="auth-reveal"
            type="button"
            aria-label={
              showPassword ? 'Parolayı gizle' : 'Parolayı göster'
            }
            aria-pressed={showPassword}
            onClick={() => setShowPassword((visible) => !visible)}
          >
            {showPassword ? 'Gizle' : 'Göster'}
          </button>
        </div>
      </div>

      {error ? (
        <p className="auth-error" role="alert">
          {error}
        </p>
      ) : null}

      <button className="auth-submit" type="submit" disabled={loading}>
        {loading
          ? isLogin
            ? 'Giriş yapılıyor…'
            : 'Hesap oluşturuluyor…'
          : isLogin
            ? 'Giriş yap'
            : 'Hesap oluştur'}
      </button>

      <p className="auth-security">
        {isLogin
          ? 'Oturumlar kısa ömürlü JWT ile korunur.'
          : "Parolanız bcrypt ile hash'lenir; düz metin olarak saklanmaz."}
      </p>
    </form>
  )
}

export default AuthForm
