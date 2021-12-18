<script lang="ts">
  import rolls from '../stores/rolls'

  $: lastRoll = $rolls.shift()
</script>

<div>
  {#if lastRoll}
    last rolled {lastRoll.place.name} at {lastRoll.time}
  {:else}
    No rolls yet
  {/if}
</div>
