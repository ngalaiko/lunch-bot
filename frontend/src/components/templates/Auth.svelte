<script context="module" lang="ts">
  const loading = users.getMe()
</script>

<script lang="ts">
  import { users } from '../../api'
  import { Login, Loading } from '../organisms'
</script>

<Loading {loading}>
  {#if $users}
    <slot />
  {:else}
    <Login />
  {/if}
</Loading>
