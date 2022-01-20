<script lang="ts">
  import { TimeSince } from '../atoms'
  import { onMount } from 'svelte'

  export let date: Date
  const secondsInOneDay = 24 * 60 * 60 * 1000

  let now = new Date()
  onMount(() => {
    if (now.getTime() - date.getTime() > secondsInOneDay) return () => {}
    const interval = setInterval(() => {
      now = new Date()
    }, 1000)
    return () => {
      clearInterval(interval)
    }
  })
</script>

<TimeSince from={date} to={now} />
