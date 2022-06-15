import { writable } from "svelte/store";
import { websocket } from './protocols'
import { parseJSON as parseUserJSON, type User } from './users'

export type Room = {
    __typename: 'Room'

    id: string
    name: string
    time: Date
    user: User
    members: User[]
}

const isSameRoom = (room: Room, otherRoom: Room) => room.id === otherRoom.id

const store = writable<Room[]>([]);

const parseJSON = (data: any): Room =>  ({
    __typename: 'Room',
    id: data.id,
    name: data.name,
    time: new Date(data.time),
    user: parseUserJSON(data.user),
    members: data.members.map(parseUserJSON),
})

const storeResponse = (response: any) => 
  response.rooms && response.rooms.map(parseJSON).forEach((room : Room)  =>  
    store.update(rooms => 
      rooms
        .filter(r => !isSameRoom(r, room))
        .concat(room)
        .sort((a, b) => a.time.getTime() - b.time.getTime())
    )
)

let listPromise: Promise<any> | undefined
const list = async (): Promise<void> => {
    await websocket.open()
    if (listPromise) return await listPromise
    listPromise = websocket.request({method: 'rooms/list'})
    const response = await listPromise
    if (response.error) throw new Error(response.error)
    storeResponse(response)
}

const create = async (name: string): Promise<void> => {
    await websocket.open()
    const response = await websocket.request({method: 'rooms/create', params: {name}})
    if (response.error) throw new Error(response.error)
    storeResponse(response)
}

websocket.onMessage(storeResponse)

export default {
    create,
    list,
    subscribe: store.subscribe
}
