import { base } from "$app/paths";
import { lobbyStore } from "$lib/stores";
import type { LoadEvent } from "@sveltejs/kit";

export function load(event: LoadEvent) {
    const createGame = async () => {
        const response = await event.fetch(
            base + "/api/games/create",
            {
                method: "POST",
            }
        );
        console.log(response);
        if (response.ok) {
            const lobbyInfo = await response.json();
            lobbyStore.set(lobbyInfo);
        }
        return [response.ok, response.status]; 
    };
    return {
        createGame,
    }
}