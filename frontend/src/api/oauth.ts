import { rest } from './protocols'

const slack = async (code: string, redirectUri: string): Promise<void> => {
  await rest.post('oauth/slack', { code, redirectUri })
}

export default {
  slack
}
