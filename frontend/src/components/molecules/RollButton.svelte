<script lang="ts">
  import { CircleButton } from '../atoms'
  import { createEventDispatcher } from 'svelte'
  import Icon from 'svelte-awesome'
  import { refresh } from 'svelte-awesome/icons'

  const dispatch = createEventDispatcher()

  const onButtonClick = () => dispatch('roll')
</script>

<CircleButton on:click={onButtonClick}>
  <Icon data={refresh} scale={1.5} />
</CircleButton>
