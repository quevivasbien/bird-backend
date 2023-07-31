import type { Card } from "./types";

export function sortCards(cards: Card[], ascending: boolean | null) {
    if (ascending === null) {
        return cards;
    } 
    const sign = ascending ? 1 : -1;
    return cards.sort((a, b) => {
        const v = (a.color * 100 + sign * a.value) - (b.color * 100 + sign * b.value);
        return v;
    });
}