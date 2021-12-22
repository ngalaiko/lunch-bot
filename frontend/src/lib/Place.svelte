<script lang="ts">
  import { boosts, places, Place } from '../api'
  import { Button } from '../atoms'
  export let place: Place

  const onButtonClick = () => boosts.create(place.id).then(places.list).catch(alert)
</script>

<div class="flex">
  <p>
    {place.name} - {place.chance * 100}%
  </p>
  <Button on:click={onButtonClick}>Boost</Button>
</div>
