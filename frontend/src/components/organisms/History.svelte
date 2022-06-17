<script lang="ts">
  import { rolls, places, boosts } from '../../api'
  import type { Place, Boost, Roll } from '../../api'
  import { getWeek, weekday } from '../../lib/time'

  type HistoryItem = Place | Boost | Roll

  const groupByWeek = (items: HistoryItem[]) => {
    const groups: Map<number, HistoryItem[]> = new Map()
    for (const item of items) {
      const week = getWeek(item.time)
      const group = groups.get(week) || []
      group.push(item)
      groups.set(week, group)
    }
    return groups
  }

  const groupByYear = (items: HistoryItem[]) => {
    const groups: Map<number, HistoryItem[]> = new Map()
    for (const item of items) {
      const year = item.time.getFullYear()
      const group = groups.get(year) || []
      group.push(item)
      groups.set(year, group)
    }
    return groups
  }

  $: fullHistory = [...$places, ...$rolls, ...$boosts].sort(
    (a, b) => b.time.getTime() - a.time.getTime()
  )
</script>

<div class="flex flex-col">
  {#each Array.from(groupByYear(fullHistory)) as [year, items]}
    {#each Array.from(groupByWeek(items)) as [week, items]}
      <section class="flex flex-1 flex-col items-center">
        <h2 class="py-4">{year} w.{week}</h2>
        <ul class="flex flex-col space-y-2 w-3/4">
          {#each items as item}
            <li data-id={item.id} class="flex flex-1 justify-between">
              <time datetime={item.time.toISOString()}>{weekday(item.time)}</time>
              <p>
                {#if item.__typename === 'Roll'}
                  {item.user.name} rolled {item.place.name}
                {:else if item.__typename === 'Boost'}
                  {item.user.name} boosted {item.place.name}
                {:else if item.__typename === 'Place' && item.user}
                  {item.user.name} added {item.name}
                {/if}
              </p>
            </li>
          {/each}
        </ul>
      </section>
    {/each}
  {/each}
</div>
