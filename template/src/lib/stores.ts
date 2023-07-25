import { writable } from "svelte/store";
import type { BidInfo, LobbyInfo, UserInfo } from "$lib/types";

export const userStore = writable<UserInfo | undefined>(undefined);

export const lobbyStore = writable<LobbyInfo | undefined>(undefined);

export const bidStore = writable<BidInfo | undefined>(undefined);
