import { writable } from "svelte/store";

interface UserInfo {
    name: string;
}

export const userStore = writable<UserInfo | undefined>(undefined);