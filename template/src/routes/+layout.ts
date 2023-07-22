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
            base + "/api/logout",
            {
                method: "POST",
            },
        );
        console.log(response);
        if (response.ok) {
            userStore.set(undefined);
        }
        return response.status;
    }

    return {
        logout,
    }
}