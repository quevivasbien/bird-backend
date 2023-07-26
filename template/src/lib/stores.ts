import { writable } from "svelte/store";
import type { GameInfo, BidInfo, LobbyInfo, UserInfo } from "$lib/types";

export const userStore = writable<UserInfo | undefined>(undefined);

export const lobbyStore = writable<LobbyInfo | undefined>(undefined);

export const bidStore = writable<BidInfo | undefined>(undefined);

export const gameStore = writable<GameInfo | undefined>(undefined);