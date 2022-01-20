<script lang="ts">
  import { Router, Route } from 'svelte-routing'
  import { Index, NotFound, Slack, History, List } from './pages'
</script>

<Router>
  <Route path="/oauth/slack"><Slack /></Route>
  <Route path="/history"><History /></Route>
  <Route path="/list"><List /></Route>
  <Route path="/"><Index /></Route>
  <Route path="*"><NotFound /></Route>
</Router>
