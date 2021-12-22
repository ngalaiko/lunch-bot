<script lang="ts">
  type Size = 'sm' | 'md' | 'lg'
  export let size: Size = 'md'

  const padding = {
    sm: 'px-2.5 py-1.5',
    md: 'px-4 py-2',
    lg: 'px-6 py-3'
  }

  const text = {
    sm: 'text-xs',
    md: 'text-base',
    lg: 'text-lg'
  }

  const css = `${padding[size]} ${text[size]}`
</script>

<button
  type="button"
  class={`${css} inline-flex items-center border border-gray-300 shadow-sm font-medium rounded text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500`}
  on:click
>
  <slot />
</button>
