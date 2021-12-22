import { http } from './protocols'

const slack = async (code: string, redirectUri: string): Promise<void> => {
  await http.post('oauth/slack', { code, redirectUri })
}

export default {
  slack
}
