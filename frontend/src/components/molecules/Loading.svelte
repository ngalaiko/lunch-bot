<script lang="ts">
  import { onMount } from 'svelte'

  let className: string = ''
  export { className as class }

  let modifier = -1
  $: count = 2
  onMount(() => {
    const interval = setInterval(() => {
      if (count === 3 || count === 0) modifier = -modifier
      count += modifier
    }, 150)
    return () => {
      clearInterval(interval)
    }
  })
</script>

<p class={className}>
  <span>Loading{'.'.repeat(count)}</span>
</p>
