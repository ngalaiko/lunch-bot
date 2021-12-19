<script lang="ts">
  import users from '../api/users'
  import Login from './Login.svelte'
</script>

<div>
  {#await users.getMe()}
    Loading...
  {:then user}
    {#if user}
      <slot />
    {:else}
      <Login />
    {/if}
  {/await}
</div>
