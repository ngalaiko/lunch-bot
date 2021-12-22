<script lang="ts">
  import { places } from '../api'
  import { Button, Input } from '../atoms'

  let name = ''
  const notEmpty = (value: string) => value.length > 0

  const onFormSubmit = () => {
    if (notEmpty(name))
      places
        .create(name)
        .then(() => (name = ''))
        .catch(alert)
  }
</script>

<form on:submit|preventDefault={onFormSubmit}>
  <Input name="name" placeholder="Weidao" bind:value={name} valid={notEmpty} />
  <Button on:click={onFormSubmit}>Add</Button>

  <input type="submit" class="hidden" />
</form>
