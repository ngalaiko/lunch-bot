<script lang="ts">
  import { boosts, places } from '../../api'
  import { Place, NewPlaceForm } from '../molecules'

  const handleOnBoost = (e: CustomEvent) =>
    boosts.create(e.detail.id).then(places.list).catch(alert)

  let newPlaceName = ''
  const handleOnSubmit = (e: CustomEvent) => {
    if (newPlaceName.length > 0)
      places
        .create(e.detail.name)
        .then(places.list)
        .catch(alert)
        .finally(() => {
          newPlaceName = ''
        })
  }

  $: sortedList = $places
    .sort((a, b) => b.time.getTime() - a.time.getTime())
    .sort((a, b) => b.chance - a.chance)
</script>

<div class="flex flex-col items-center">
  <ul class="flex flex-col items-stretch space-y-2">
    <li><NewPlaceForm on:submit={handleOnSubmit} bind:name={newPlaceName} /></li>
    {#each sortedList as place}
      <li><Place on:boost={handleOnBoost} {place} /></li>
    {/each}
  </ul>
</div>
