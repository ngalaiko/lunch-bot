<script lang="ts" context="module">
  import { rolls, places, boosts } from '../../api'
  const fetching = Promise.all([rolls.list(), places.list(), boosts.list()])
</script>

<script lang="ts">
  import { TimeSince } from '../molecules'
  $: fullHistory = [...$places, ...$rolls, ...$boosts].sort(
    (a, b) => b.time.getTime() - a.time.getTime()
  )
</script>

<div class="flex flex-col items-center">
  {#await fetching}
    <p>Loading...</p>
  {:then}
    <ul class="flex flex-col items-stretch space-y-2">
      {#each fullHistory as item}
        <li class="flex row items-center justify-between">
          <p>
            {#if item.__typename === 'Roll'}
              {item.user.name} rolled {item.place.name}
            {:else if item.__typename === 'Boost'}
              {item.user.name} boosted {item.place.name}
            {:else if item.__typename === 'Place' && item.user}
              {item.user.name} added {item.name}
            {/if}
          </p>
          <TimeSince date={item.time} />
        </li>
      {/each}
    </ul>
  {:catch e}
    <p>Error: {e.message}</p>
  {/await}
</div>
