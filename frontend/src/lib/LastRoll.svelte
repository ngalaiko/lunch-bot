<script lang="ts">
  import { rolls } from '../api'

  $: lastRoll = $rolls.shift()
</script>

<div>
  {#await rolls.list()}
    loading...
  {:then}
    {#if lastRoll}
      last rolled {lastRoll.place.name} at {lastRoll.time}
    {:else}
      No rolls yet
    {/if}
  {:catch e}
    <p>Error: {e.message}</p>
  {/await}
</div>
