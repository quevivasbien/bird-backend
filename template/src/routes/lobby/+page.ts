import { base } from "$app/paths";
import { bidStore, lobbyStore } from "$lib/stores";
import type { LoadEvent } from "@sveltejs/kit";
import { get } from "svelte/store";

export function load(event: LoadEvent) {
    const subscribeToLobby = () => {
        const lobbyInfo = get(lobbyStore);
        if (lobbyInfo === undefined) {
            return;
        }
        const sse = new EventSource(
           `${base}/api/lobbies/${lobbyInfo.id}/subscribe`,
        );
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

    const startBidding = async () => {
        const lobbyInfo = get(lobbyStore);
        if (lobbyInfo === undefined) {
            return [false, 0];
        }
        const response = await event.fetch(
            `${base}/api/bidding/${lobbyInfo.id}`,
            {
                method: "PUT",
            },
        );
        return [response.ok, response.status];
    }

    const receiveBidState = async () => {
        const lobbyInfo = get(lobbyStore);
        if (lobbyInfo === undefined) {
            console.log("Lobby info undefined when trying to get bid state!")
            return [false, 0];
        }
        const response = await event.fetch(
            `${base}/api/bidding/${lobbyInfo.id}`,
            {
                method: "GET",
            },
        );
        if (response.ok) {
            const bidState = await response.json();
            bidStore.set(bidState);
        }
        return [response.ok, response.status];
    };

    return {
        subscribeToLobby,
        swapPlayers,
        leaveLobby,
        startBidding,
        receiveBidState,
    };
}