import { base } from "$app/paths";
import { lobbyStore } from "$lib/stores";
import type { LoadEvent } from "@sveltejs/kit";
import { get } from "svelte/store";

export function load(event: LoadEvent) {
    const subscribeToLobby = async () => {
        const lobbyInfo = get(lobbyStore);
        if (lobbyInfo === undefined) {
            return;
        }
        const sse = new EventSource(
            base + "/api/lobbies/" + lobbyInfo.id + "/subscribe",
        )
        return sse;
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

    const startGame = async () => {
        const lobbyInfo = get(lobbyStore);
        if (lobbyInfo === undefined) {
            return [false, 0];
        }
        const response = await event.fetch(
            base + "/api/games/" + lobbyInfo.id + "/bidding/start",
            {
                method: "POST",
            },
        );
        return [response.ok, response.status];
    }

    return {
        subscribeToLobby,
        swapPlayers,
        leaveLobby,
        startGame,
    };
}