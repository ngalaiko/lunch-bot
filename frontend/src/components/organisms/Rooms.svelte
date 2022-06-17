<script lang="ts">
  import { rooms } from '../../api'
  import { Room } from '../molecules'

  $: sortedList = $rooms.sort((a, b) => b.time.getTime() - a.time.getTime())
</script>

<div class="flex flex-col items-center">
  <div class="grid grid-cols-3">
    {#each sortedList as room}
      <Room {room} />
    {/each}
  </div>
</div>
