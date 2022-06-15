<script lang="ts" context="module">
  import { rooms } from '../../api'
  const fetching = rooms.list()
</script>

<script lang="ts">
  import { Room, Loading } from '../molecules'

  $: sortedList = $rooms.sort((a, b) => b.time.getTime() - a.time.getTime())
</script>

<div class="flex flex-col items-center">
  {#await fetching}
    <Loading />
  {:then}
    <div class="grid grid-cols-3">
      {#each sortedList as room}
        <Room {room} />
      {/each}
    </div>
  {:catch e}
    <p>Error: {e.message}</p>
  {/await}
</div>
