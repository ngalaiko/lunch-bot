<script context="module" lang="ts">
  const loading = users.getMe()
</script>

<script lang="ts">
  import { users } from '../../api'
  import { Login } from '../organisms'
</script>

{#await loading}
  <p>Loading...</p>
{:then}
  {#if $users}
    <slot />
  {:else}
    <Login />
  {/if}
{/await}
