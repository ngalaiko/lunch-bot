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
    isValid ? '' : 'text-red-900 placeholder-red-300'
  } m-0 p-0 border-0 focus:shadow-transparent block w-full sm:text-sm rounded-md`}
  {placeholder}
  on:input={validate}
  on:paste={validate}
  on:change={validate}
/>
