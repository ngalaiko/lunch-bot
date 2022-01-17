<script lang="ts">
  import { Link } from 'svelte-routing'
  import { Icon } from '../../assets'
  import { users } from '../../api'

  const routes = [
    {
      title: 'Roll',
      path: '/',
      active: document.location.pathname === '/'
    },
    {
      title: 'List',
      path: '/list',
      active: document.location.pathname === '/list'
    }
  ]
  const onLogoutClick = users.logout
</script>

<nav class="flex flex-1 m-auto p-5 max-w-screen-sm">
  <Link to="/">
    <Icon class="w-10 h-10" />
  </Link>

  <ul class="flex flex-1 space-x-5 justify-center">
    {#each routes as route}
      <li
        class="{route.active
          ? 'border-yellow-500'
          : 'border-transparent hover:border-yellow-300 hover:text-gray-700'} inline-flex items-center border-b-2"
      >
        <Link to={route.path}>{route.title}</Link>
      </li>
    {/each}
  </ul>

  <button
    class="border-transparent hover:border-yellow-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
    on:click|preventDefault={onLogoutClick}>Logout</button
  >
</nav>

<slot />
