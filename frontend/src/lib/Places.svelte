<script lang="ts">
  import places from '../stores/places'
  import Place from './Place.svelte'
</script>

<ul>
  {#each $places as place}
    <Place {place} />
  {/each}
</ul>
