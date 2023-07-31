<script lang="ts">
	import { gameStore } from "$lib/stores";
	import type { Card, GameInfo } from "$lib/types";
	import CardView from "$lib/components/CardView.svelte";
	import { onMount } from "svelte";

    $: cards = $gameStore?.table ?? [];

    let leadingPlayer: number = -1;
	$: leadingPlayer = getLeadingPlayer($gameStore);
	function getLeadingPlayer(gameInfo: GameInfo | undefined) {
		if (gameInfo === undefined) {
			return -1;
		}
		if (gameInfo.table.length === 0) {
			return gameInfo.currentPlayer;
		}
		return leadingPlayer;
	}

    let sortedCards: (Card | null)[] = [null, null, null, null];
    function sortCards(cards: Card[]) {
        for (let i = 0; i < cards.length; i++) {
            const player = (leadingPlayer + i + 4) % 4;
            sortedCards[player] = cards[i];
        }
        return sortedCards;
    }

    $: sortedCards = sortCards(cards);
</script>

<div class="flex flex-row space-x-8 items-center mt-4">
	<h2 class="text-2xl">Cards on table</h2>
</div>
<div class="flex flex-wrap">
    {#each sortedCards as card, i}
        <div class="flex flex-col space-y-2 m-4">
            <div class="text-xl text-center">Player {i + 1}</div>
            {#if card === null}
                No card yet
            {:else}
                <CardView {card} />
            {/if}
        </div>
    {/each}
</div>