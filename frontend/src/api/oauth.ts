import { restUri } from './api'

const slack = async (code: string, redirectUri: string): Promise<void> => {
  const response = await fetch(`${restUri}/oauth/slack`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      code: code,
      redirect_uri: redirectUri
    }),
    credentials: 'include'
  })
  if (response.status != 200) throw new Error('failed to exchange slack code')
}

export default {
  slack
}
