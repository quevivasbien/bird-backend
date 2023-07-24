import { base } from "$app/paths";
import { userStore } from "$lib/stores";
import type { LoadEvent } from "@sveltejs/kit";
import { get } from "svelte/store";

export const prerender = true;

export function load(event: LoadEvent) {
    const logout = async () => {
        const userInfo = get(userStore);
        if (userInfo === undefined) {
            return;
        }
        const response = await event.fetch(
            base + "/api/auth/logout",
            {
                method: "POST",
            },
        );
        if (response.ok) {
            userStore.set(undefined);
        }
        return [response.ok, response.status];
    };

    // if a jwt cookie is present, use it to get info for userStore
    const syncAuth = async () => {
        const response = await event.fetch(
            base + "/api/auth/status",
            {
                method: "GET",
            },
        );
        if (!response.ok) {
            console.log("Problem when attempting to fetch user info:", response.statusText);
            return;
        }
        const userInfo = await response.json();
        console.log(userInfo);
        // if userInfo is unexpired, use it
        if (Date.now() / 1000 < userInfo.expireTime) {
            return userInfo;
        }
    }
   


    return {
        logout,
        syncAuth,
    };
}