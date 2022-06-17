<script lang="ts" context="module">
  import ws from '../../api/protocols/websocket'
  ws.open()
</script>

<script lang="ts">
  import Auth from './Auth.svelte'
  import Header from './Header.svelte'
</script>

<Auth>
  <div class="m-auto max-w-screen-sm">
    <Header>
      <slot />
    </Header>
  </div>
</Auth>
