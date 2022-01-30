<script lang="ts">
  import { Loading } from '../molecules'
  let loading: Promise<void>
  let className: string = ''
  export { loading, className as class }
</script>

{#await loading}
  <Loading class={className} />
{:then}
  <slot />
{/await}
