<script lang="ts">
	import { gameStore } from "$lib/stores";
	import type { Card, GameInfo } from "$lib/types";
	import CardView from "$lib/components/CardView.svelte";

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