<script lang="ts">
  import { rolls, places } from '../api'

  const onButtonClick = () => rolls.create().then(places.list).catch(alert)
</script>

<button
  type="button"
  class="text-lg p-2 px-4 rounded-md drop-shadow-md bg-blue-500 text-white"
  on:click={onButtonClick}>Roll</button
>
