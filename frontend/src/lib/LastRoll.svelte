<script lang="ts">
  import { rolls, places } from '../api'
  import { Button } from '../atoms'

  $: lastRoll = $rolls.shift()
  $: text = lastRoll
    ? `last rolled ${lastRoll.place.name} at ${lastRoll.time}`
    : 'no rolls yet'

  const onButtonClick = () => rolls.create().then(places.list).catch(alert)
</script>

<div>
  {#await rolls.list()}
    loading...
  {:then}
    <p>{text}</p>
    <Button on:click={onButtonClick}>Roll</Button>
  {:catch e}
    <p>Error: {e.message}</p>
  {/await}
</div>
