<script lang="ts">
  import type { Place } from './api'
  export let place: Place
</script>

<p>
  {place.name} - {place.chance * 100}%
</p>
