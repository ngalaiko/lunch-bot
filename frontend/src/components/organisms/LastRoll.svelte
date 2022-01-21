<script lang="ts" context="module">
  import { rolls } from '../../api'
  const fetching = rolls.list()
</script>

<script lang="ts">
  import { RollButton, TimeSince } from '../molecules'

  $: lastRoll = $rolls.sort((a, b) => b.time.getTime() - a.time.getTime()).at(0)

  const onRoll = () => rolls.create().catch(alert)
</script>

<div class="flex flex-col items-center">
  {#await fetching}
    loading...
  {:then}
    <div class="flex flex-col items-center m-3">
      {#if lastRoll}
        <h2 class="text-2xl">{lastRoll.place.name}</h2>
        <p class="text-sm text-slate-500"><TimeSince date={lastRoll.time} /></p>
      {:else}
        <h2 class="text-2xl">No rolls yet</h2>
      {/if}
    </div>
    <RollButton on:roll={onRoll} />
  {:catch e}
    <p>Error: {e.message}</p>
  {/await}
</div>
