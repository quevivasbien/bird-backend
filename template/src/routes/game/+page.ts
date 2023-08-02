import { base } from "$app/paths";
import { gameStore } from "$lib/stores";
import type { Card } from "$lib/types";
import type { LoadEvent } from "@sveltejs/kit";
import { get } from "svelte/store";

export function load(event: LoadEvent) {
    const subscribeToGame = () => {
        const gameInfo = get(gameStore);
        if (gameInfo === undefined) {
            return;
        }
        const sse = new EventSource(
            `${base}/api/games/${gameInfo.id}/subscribe`,
        );
        return sse;
    };

    const getWidow = async () => {
        const gameInfo = get(gameStore);
        if (gameInfo === undefined) {
            return;
        }
        const response = await fetch(
            `${base}/api/games/${gameInfo.id}/widow`,
            {
                method: "GET",
            },
        );
        if (!response.ok) {
            console.log("Problem getting widow; status =" + response.status);
            return;
        }
        const data = await response.json();
        return data;
    };

    const startRound = async (trump: number, toWidow: Card[], fromWidow: Card[]) => {
        const gameInfo = get(gameStore);
        if (gameInfo === undefined) {
            return [false, 0];
        }
        const response = await event.fetch(
            `${base}/api/games/${gameInfo.id}/start`,
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ trump, toWidow, fromWidow }),
            }
        );
        return [response.ok, response.status];
    };

    const getScore = async () => {
        const gameInfo = get(gameStore);
        if (gameInfo === undefined) {
            return [];
        }
        const response = await fetch(
            `${base}/api/games/${gameInfo.id}/score`,
            {
                method: "GET",
            },
        );
        if (!response.ok) {
            console.log("Problem getting end-of-game score; status = " + response.status);
            return [];
        }
        const data = await response.json();
        return [data.score0, data.score1];
    };

    const playCard = async (card: Card) => {
        const gameInfo = get(gameStore);
        if (gameInfo === undefined) {
            return [false, 0];
        }
        const response = await fetch(
            `${base}/api/games/${gameInfo.id}/play`,
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(card),
            },
        );
        return [response.ok, response.status];
    };

    const finishPlay = async () => {
        const gameInfo = get(gameStore);
        if (gameInfo === undefined) {
            return [false, 0];
        }
        const response = await fetch(
            `${base}/api/games/${gameInfo.id}/finish`,
            {
                method: "POST",
            },
        );
        return [response.ok, response.status];
    };

    return {
        subscribeToGame,
        getWidow,
        startRound,
        getScore,
        playCard,
        finishPlay,
    };
}
