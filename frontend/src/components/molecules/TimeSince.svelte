<script lang="ts">
  import { TimeSince } from '../atoms'
  import { onMount } from 'svelte'

  export let date: Date

  $: now = new Date()
  onMount(() => {
    const interval = setInterval(() => {
      now = new Date()
    }, 1000)
    return () => {
      clearInterval(interval)
    }
  })
</script>

<TimeSince from={date} to={now} />
