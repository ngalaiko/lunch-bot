<script lang="ts">
  import { Button, Input } from '../atoms'
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  let name = ''
  const notEmpty = (value: string) => value.length > 0
  const onFormSubmit = () => {
    dispatch('submit', { name })
  }
</script>

<form on:submit|preventDefault={onFormSubmit}>
  <Input name="name" placeholder="Weidao" bind:value={name} valid={notEmpty} />
  <Button on:click={onFormSubmit}>Add</Button>

  <input type="submit" class="hidden" />
</form>
