import { getAccessToken } from './auth'

export type APIKey = {
  id: number
  title: string
  expires_at: string | null
  created_at: string
}

type ListResponse = {
  api_keys?: APIKey[]
}

type CreateResponse = {
  api_key?: string
}

export type ListAPIKeysResult = {
  status: number
  apiKeys: APIKey[]
}

export type CreateAPIKeyResult = {
  status: number
  apiKey?: string
}

function authHeaders(): Record<string, string> {
  const token = getAccessToken()

  return token ? { Authorization: `Bearer ${token}` } : {}
}

export async function listAPIKeys(): Promise<ListAPIKeysResult> {
  const response = await fetch('/api/v1/api-keys', {
    headers: authHeaders(),
  })

  let data: ListResponse = {}

  try {
    data = (await response.json()) as ListResponse
  } catch {
    // JSON olmayan yanıtlarda durum koduna göre ilerle.
  }

  return {
    status: response.status,
    apiKeys: data.api_keys ?? [],
  }
}

export async function createAPIKey(
  title: string,
  expiresAt: string | null,
): Promise<CreateAPIKeyResult> {
  const response = await fetch('/api/v1/api-key', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...authHeaders(),
    },
    body: JSON.stringify({
      title,
      expires_at: expiresAt,
    }),
  })

  let data: CreateResponse = {}

  try {
    data = (await response.json()) as CreateResponse
  } catch {
    // JSON olmayan yanıtlarda durum koduna göre ilerle.
  }

  return {
    status: response.status,
    apiKey: data.api_key,
  }
}
