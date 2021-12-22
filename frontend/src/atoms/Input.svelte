<script lang="ts">
  export let name: string
  export let placeholder: string
  export let value: string
  export let valid: (value: string) => boolean = () => true

  let isValid = true
  const validate = () => {
    isValid = valid(value)
  }
</script>

<input
  type="text"
  {name}
  bind:value
  id={name}
  class={`${
    isValid
      ? ''
      : 'border-red-300 text-red-900 placeholder-red-300 focus:ring-red-500 focus:border-red-500'
  } shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md`}
  {placeholder}
  on:input={validate}
  on:paste={validate}
  on:change={validate}
/>
