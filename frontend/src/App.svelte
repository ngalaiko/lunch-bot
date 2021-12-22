<script lang="ts">
  import { Router, Route } from 'svelte-routing'
  import Index from './pages/Index.svelte'
  import NotFound from './pages/NotFound.svelte'
  import { Auth } from './templates'
  import OAuthSlack from './pages/oauth/Slack.svelte'
</script>

<Router>
  <Route path="/oauth/slack"><OAuthSlack /></Route>
  <Route path="/"><Auth><Index /></Auth></Route>
  <Route path="*"><NotFound /></Route>
</Router>
