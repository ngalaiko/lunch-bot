<script lang="ts">
  import { users } from '../../api'
  import { Login } from '../organisms'
</script>

{#await users.getMe()}
  <p>Loading...</p>
{:then}
  <slot />
{:catch}
  <Login />
{/await}
