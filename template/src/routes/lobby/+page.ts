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
        if (!response.ok) {
            console.log("When fetching lobby state, got " + response.statusText);
            return;
        }
        return await response.json();
    };

    const swapPlayers = async (i: number, j: number) => {
        const lobbyInfo = get(lobbyStore);
        if (lobbyInfo === undefined) {
            return [false, 0];
        }
        const response = await event.fetch(
            base + "/api/lobbies/" + lobbyInfo.id + "/swap",
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ i, j }),
            },
        );
        return [response.ok, response.status];
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
        return [response.ok, response.status];
    };

    return {
        getLobbyState,
        swapPlayers,
        leaveLobby,
    };
}