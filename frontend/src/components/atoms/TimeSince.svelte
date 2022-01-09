<script lang="ts">
  export let from: Date
  export let to: Date

  const humanDiff = (from: Date, to: Date): string => {
    const diff = to.getTime() - from.getTime()
    const seconds = Math.floor(diff / 1000)
    const minutes = Math.floor(diff / 1000 / 60)
    const days = Math.floor(minutes / 60 / 24)
    const weeks = Math.floor(days / 7)
    const months = Math.floor(days / 30)
    const years = Math.floor(days / 365)

    if (years > 0) {
      return years === 1 ? 'last year' : `${years} years ago`
    } else if (months > 0) {
      return months === 1 ? 'last month' : `${months} months ago`
    } else if (weeks > 0) {
      return weeks === 1 ? 'last week' : `${weeks} weeks ago`
    } else if (days > 0) {
      return days === 1 ? 'yesterday' : `${days} days ago`
    } else if (minutes > 0) {
      return minutes === 1 ? 'a minute ago' : `${minutes} minutes ago`
    } else if (seconds > 0) {
      return seconds === 1 ? 'a second ago' : `${seconds} seconds ago`
    } else {
      return 'just now'
    }
  }
</script>

<span>{humanDiff(from, to)}</span>
