<script lang="ts">
  import { navigate } from 'svelte-routing'
  import { users } from '../../../api'

  const params = new URLSearchParams(window.location.search)

  const code = params.get('code') as string
  const next = params.get('next') as string
  const redirectUri = `${location.origin}/oauth/slack?next=${encodeURIComponent(next)}`

  users
    .slackOAuth(code, redirectUri)
    .then(() => {
      navigate(next)
    })
    .catch(e => {
      alert(`Error: ${e}`)
    })
</script>

<div>Please wait...</div>
