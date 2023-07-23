import { base } from "$app/paths";
import { lobbyStore } from "$lib/stores";
import type { LoadEvent } from "@sveltejs/kit";
import { get } from "svelte/store";

export function load(event: LoadEvent) {
    // use this to update lobbyStore   
    const getLobbyState = async () => {
        const lobbyInfo = get(lobbyStore);
        if (lobbyInfo === undefined) {
            return;
        }
        const response = await event.fetch(
            base + "/api/lobbies/" + lobbyInfo.id,
            {
                method: "GET",
            },
        );
        console.log(response);
        if (!response.ok) {
            console.log("When fetching lobby state, got " + response.statusText);
            return;
        }
        return await response.json();
    };

    const leaveLobby = async () => {
        const lobbyInfo = get(lobbyStore);
        if (lobbyInfo === undefined) {
            return [false, 0];
        }
        const response = await event.fetch(
            base + "/api/lobbies/" + lobbyInfo.id + "/leave",
            {
                method: "POST",
            },
        );
        console.log(response);
        return [response.ok, response.status];
    };

    return {
        getLobbyState,
        leaveLobby,
    };
}