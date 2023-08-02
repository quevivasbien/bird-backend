import { get } from "svelte/store";
import type { Card } from "./types";
import { gameStore } from "./stores";

function sortCardValue(card: Card, sign: number, trump: number) {
    let { color, value } = card;
    if (color === 0 || (trump !== 0 && color === trump)) {
        color = 10;
    }
    if (value === 1) {
        value = 15;
    }
    return sign * (color * 100 + value);
}

export function sortCards(cards: Card[], ascending: boolean | null) {
    if (ascending === null) {
        return [...cards];
    } 
    const sign = ascending ? 1 : -1;
    const trump = get(gameStore)?.trump ?? 0;
    return [...cards].sort((a, b) => {
        const v = sortCardValue(a, sign, trump) - sortCardValue(b, sign, trump);
        return v;
    });
}
