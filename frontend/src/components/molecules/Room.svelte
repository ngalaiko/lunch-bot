<script lang="ts">
  import type { Room } from '../../api'

  export let room: Room
</script>

<div class="flex items-end w-24 h-24 bg-slate-100 shadow hover:shadow-lg">
  <p>{room.name}</p>
</div>
