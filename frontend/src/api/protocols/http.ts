const apiUri = import.meta.env.DEV
  ? 'https://localhost:8000/api/'
  : `${location.origin}/api/`

const get = async (path: string): Promise<any> => {
  const response = await fetch(apiUri + path, {
    method: 'GET',
    credentials: 'include'
  })
  if (response.status !== 200)
    throw new Error(`${response.status}: ${response.statusText}`)
  return await response.json()
}

const post = async (path: string, body?: any): Promise<any> => {
  const opts: RequestInit = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include'
  }
  if (body) opts.body = JSON.stringify(body)
  const response = await fetch(apiUri + path, opts)
  if (response.status !== 200)
    throw new Error(`${response.status}: ${response.statusText}`)
  const responseBody = await response.text()
  if (responseBody.length == 0) return {}
  return JSON.parse(responseBody)
}

export default {
  get,
  post
}
