<script lang="ts">
  import Places from '../lib/Places.svelte'
  import LastRoll from '../lib/LastRoll.svelte'
  import RollButton from '../lib/RollButton.svelte'
  import AddPlaceButton from '../lib/AddPlaceButton.svelte'
</script>

<main class="text-center">
  <LastRoll />
  <RollButton />
  <Places />
  <AddPlaceButton />
</main>
