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
    started: boolean;
}

export const lobbyStore = writable<LobbyInfo | undefined>(undefined);

interface Card {
    color: number;
    value: number;
}

interface BidInfo {
    gameID: string;
    done: boolean;
    players: string[];
    hands: Card[][];
    widow: Card[];
    passed: boolean[];
    currentBidder: number;
    bid: number;
}

export const bidStore = writable<BidInfo | undefined>(undefined);
