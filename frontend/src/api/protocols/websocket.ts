const websocketUri = 'wss://localhost:8000/ws'

import WebSocketAsPromised from 'websocket-as-promised'
import ReconnectingWebSocket from 'reconnecting-websocket'

const wsp = new WebSocketAsPromised(websocketUri, {
  // replace websocket implementation
  createWebSocket: url => new ReconnectingWebSocket(url) as WebSocket,
  // use json encoding
  packMessage: data => JSON.stringify(data),
  unpackMessage: data => JSON.parse(data as string),
  // attach requestId to message as `id` field
  attachRequestId: (data, requestId) => Object.assign({ id: requestId }, data),
  // read requestId from message `id` field
  extractRequestId: data => data && data.id
})

export default {
  open: (): Promise<Event> => wsp.open(),
  request: (request: any): Promise<any> => wsp.sendRequest(request)
}
