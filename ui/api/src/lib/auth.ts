type ApiResponse = {
  error?: string
  message?: string
  accessToken?: string
  refreshToken?: string
}

export type AuthResult = {
  status: number
  error?: string
  accessToken?: string
  refreshToken?: string
}

const ACCESS_TOKEN_KEY = 'velin.access_token'
const REFRESH_TOKEN_KEY = 'velin.refresh_token'

async function post(path: string, body: unknown): Promise<AuthResult> {
  const response = await fetch(path, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })

  let data: ApiResponse = {}

  try {
    data = (await response.json()) as ApiResponse
  } catch {
    // Gövde JSON değilse durum koduna göre devam et.
  }

  return {
    status: response.status,
    error: data.error,
    accessToken: data.accessToken,
    refreshToken: data.refreshToken,
  }
}

export function registerUser(username: string, password: string) {
  return post('/api/v1/auth/register', { username, password })
}

export async function loginUser(username: string, password: string) {
  const result = await post('/api/v1/auth/login', { username, password })

  if (result.status === 200 && result.accessToken && result.refreshToken) {
    window.localStorage.setItem(ACCESS_TOKEN_KEY, result.accessToken)
    window.localStorage.setItem(REFRESH_TOKEN_KEY, result.refreshToken)
  }

  return result
}
