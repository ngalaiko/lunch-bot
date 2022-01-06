<script lang="ts">
  import { Button } from '../atoms'
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  const onButtonClick = () => dispatch('roll')
</script>

<Button on:click={onButtonClick}>Roll</Button>
