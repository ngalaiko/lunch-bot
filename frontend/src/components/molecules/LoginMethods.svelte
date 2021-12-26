<script lang="ts">
  import { SlackIcon } from '../atoms'
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  const handleClickOnSlack = () => {
    dispatch('slack')
  }
</script>

<div class="flex flex-col items-center">
  <h1 class="text-xl">Login</h1>
  <ul class="mt-2">
    <li>
      <button
        on:click|preventDefault={handleClickOnSlack}
        class="inline-flex items-center px-3 py-2 border border-transparent shadow-sm text-sm leading-4 font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        ><SlackIcon class="-ml-0.5 mr-2 h-4 w-4" />Slack</button
      >
    </li>
  </ul>
</div>
