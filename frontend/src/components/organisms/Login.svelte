<script lang="ts">
  import { LoginMethods } from '../molecules'

  const redirectUri = encodeURIComponent(
    `${location.origin}/oauth/slack?next=${encodeURIComponent(location.href)}`
  )
  const clientId = '1693172761239.2533308174103'
  const userScope = 'identity.basic'
  const slackUrl = `https://slack.com/oauth/v2/authorize?user_scope=${userScope}&redirect_uri=${redirectUri}&client_id=${clientId}`

  const handleOnSlack = () => {
    location.href = slackUrl
  }
</script>

<LoginMethods on:slack={handleOnSlack} />
