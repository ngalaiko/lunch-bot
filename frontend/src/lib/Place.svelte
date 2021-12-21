<script lang="ts">
  import type { Place } from '../api'
  import BoostButton from './BoostButton.svelte'
  export let place: Place
</script>

<div class="flex">
  <p>
    {place.name} - {place.chance * 100}%
  </p>
  <BoostButton {place} />
</div>
