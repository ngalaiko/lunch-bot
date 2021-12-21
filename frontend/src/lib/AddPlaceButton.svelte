<script lang="ts">
  import { places } from '../api'

  let name = ''
  const onFormSubmit = () => places.create(name).catch(alert)
</script>

<form on:submit|preventDefault={onFormSubmit}>
  <label for="name">Add new place</label>
  <input required class="border-2" id="name" bind:value={name} type="text" />

  <input type="submit" value="Create" />
</form>
