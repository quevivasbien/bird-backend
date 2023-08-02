<script lang="ts">
	import type { Card } from '$lib/types';
	import CardView from '$lib/components/CardView.svelte';
	import { onMount } from 'svelte';

	export let cards: Card[];
	export let players: string[];
	export let leadingPlayer: number;

	let sortedCards: (Card | null)[] = [null, null, null, null];
	function sortCards(cards: Card[]) {
		for (let i = 0; i < cards.length; i++) {
			const player = (leadingPlayer + i + 4) % 4;
			sortedCards[player] = cards[i];
		}
		return sortedCards;
	}

	onMount(() => {
		sortedCards = sortCards(cards);
	});
</script>

<div class="flex flex-row space-x-8 items-center mt-4">
	<h2 class="text-2xl">Cards on table</h2>
</div>
<!-- TODO: change table layout to grid -->
<div class="flex flex-wrap">
	{#each sortedCards as card, i}
		<div class="flex flex-col space-y-2 m-4">
			<div class="text-xl text-center">Player {i + 1} ({players[i]})</div>
			{#if i === leadingPlayer}
				<div class="mb-2 text-center">Leading player</div>
			{/if}
			{#if card === null}
				<div class="text-center">No card yet</div>
			{:else}
				<CardView {card} />
			{/if}
		</div>
	{/each}
</div>
