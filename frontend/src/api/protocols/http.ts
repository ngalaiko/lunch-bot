const apiUri = 'https://localhost:8000/'

const get = async (path: string): Promise<any> => {
  const response = await fetch(apiUri + path, {
    method: 'GET',
    credentials: 'include'
  })
  if (response.status !== 200)
    throw new Error(`${response.status}: ${response.statusText}`)
  return await response.json()
}

const post = async (path: string, body: any): Promise<any> => {
  const response = await fetch(apiUri + path, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(body),
    credentials: 'include'
  })
  if (response.status !== 200)
    throw new Error(`${response.status}: ${response.statusText}`)
  return await response.json()
}

export default {
  get,
  post
}
