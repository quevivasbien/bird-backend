<script lang="ts">
	import type { Card } from "$lib/types";
	import CardView from "./CardView.svelte";

    export let cards: Card[];
    export let leadingPlayer: number;

    let sortedCards: (Card | null)[] = [null, null, null, null];
    for (let i = 0; i < cards.length; i++) {
        const player = (leadingPlayer + i + 4) % 4;
        sortedCards[player] = cards[i];
    }
</script>

{#each sortedCards as card, i}
    <div>
        <div>Player {i + 1}</div>
        {#if card === null}
            null
        {:else}
            <CardView {card} />
        {/if}
    </div>
{/each}