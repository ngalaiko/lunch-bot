<script lang="ts">
  import type { Place } from '../api/places'
  export let place: Place
</script>

<p>
  {place.name} - {place.chance * 100}%
</p>
