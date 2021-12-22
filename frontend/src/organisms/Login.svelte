<script lang="ts">
  import { SlackIcon } from '../atoms'

  const redirectUri = encodeURIComponent(
    `${location.origin}/oauth/slack?next=${encodeURIComponent(location.href)}`
  )
  const clientId = '1693172761239.2533308174103'
  const userScope = 'identity.basic'
  const slackUrl = `https://slack.com/oauth/v2/authorize?user_scope=${userScope}&redirect_uri=${redirectUri}&client_id=${clientId}`
</script>

<div class="flex flex-col items-center">
  <h1 class="text-xl">Login</h1>
  <ul class="mt-2">
    <li>
      <a
        href={slackUrl}
        class="inline-flex items-center px-3 py-2 border border-transparent shadow-sm text-sm leading-4 font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        ><SlackIcon class="-ml-0.5 mr-2 h-4 w-4" />Slack</a
      >
    </li>
  </ul>
</div>
