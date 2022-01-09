<script lang="ts">
  import { rolls, places } from '../../api'
  import { RollButton } from '../molecules'

  const relativeTime = (date: Date) => {
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const day = 1000 * 60 * 60 * 24
    const days = Math.floor(diff / day)
    const weeks = Math.floor(days / 7)
    const months = Math.floor(days / 30)
    const years = Math.floor(days / 365)

    if (years > 0) {
      return `${years} year${years > 1 ? 's' : ''} ago`
    } else if (months > 0) {
      return `${months} month${months > 1 ? 's' : ''} ago`
    } else if (weeks > 0) {
      return `${weeks} week${weeks > 1 ? 's' : ''} ago`
    } else if (days > 0) {
      return `${days} day${days > 1 ? 's' : ''} ago`
    } else {
      return `${Math.floor(diff / 1000 / 60)} minute${
        Math.floor(diff / 1000 / 60) > 1 ? 's' : ''
      } ago`
    }
  }

  $: lastRoll = $rolls.shift()

  const onRoll = () => rolls.create().catch(alert)
</script>

<div class="flex flex-col items-center">
  {#await rolls.list()}
    loading...
  {:then}
    <div class="flex flex-col items-center m-3">
      {#if lastRoll}
        <h2 class="text-2xl">{lastRoll.place.name}</h2>
        <p class="text-sm text-slate-500">{relativeTime(lastRoll.time)}</p>
      {:else}
        <h2 class="text-2xl">No rolls yet</h2>
      {/if}
    </div>
    <RollButton on:roll={onRoll} />
  {:catch e}
    <p>Error: {e.message}</p>
  {/await}
</div>
