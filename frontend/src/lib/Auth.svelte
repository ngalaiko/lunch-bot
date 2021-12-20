<script lang="ts">
  import { users } from '../api'
  import Login from './Login.svelte'
</script>

<div>
  {#await users.getMe()}
    Loading...
  {:then}
    <slot />
  {:catch}
    <Login />
  {/await}
</div>
