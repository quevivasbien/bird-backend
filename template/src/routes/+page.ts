import { base } from "$app/paths";
import { lobbyStore } from "$lib/stores";
import type { LoadEvent } from "@sveltejs/kit";

export function load(event: LoadEvent) {
    const createLobby = async (id: string) => {
        const response = await event.fetch(
            base + "/api/lobbies/" + id,
            {
                method: "PUT",
            }
        );
        if (response.ok) {
            const lobbyInfo = await response.json();
            lobbyStore.set(lobbyInfo);
        }
        return [response.ok, response.status]; 
    };

    const joinLobby = async (id: string) => {
        const response = await event.fetch(
            base + "/api/lobbies/" + id + "/join",
            {
                method: "POST",
            }
        );
        if (response.ok) {
            const lobbyInfo = await response.json();
            lobbyStore.set(lobbyInfo);
        }
        return [response.ok, response.status];
    };

    return {
        createLobby,
        joinLobby,
    }
}