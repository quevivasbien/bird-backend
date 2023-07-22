import { writable } from "svelte/store";

interface UserInfo {
    name: string;
    admin: boolean;
    expireTime: number;
}

export const userStore = writable<UserInfo | undefined>(undefined);

interface LobbyInfo {
    id: string;
    host: string;
    players: string[];
}

export const lobbyStore = writable<LobbyInfo | undefined>(undefined);
