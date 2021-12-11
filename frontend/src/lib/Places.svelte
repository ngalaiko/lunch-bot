<script lang="ts">
  import places from './places'
  import Place from './Place.svelte'
</script>

<ul>
  {#each $places as place}
    <Place {place} />
  {/each}
</ul>
