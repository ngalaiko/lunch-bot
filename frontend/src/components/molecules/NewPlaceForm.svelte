<script lang="ts">
  import { Button, Input } from '../atoms'
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  export let name = ''
  const notEmpty = (value: string) => value.length > 0
  const onFormSubmit = () => dispatch('submit', { name })
</script>

<form on:submit|preventDefault={onFormSubmit} class="flex">
  <Input name="name" placeholder="Something new?" bind:value={name} valid={notEmpty} />
  <Button size="sm" on:click={onFormSubmit}>Add</Button>

  <input type="submit" class="hidden" />
</form>
