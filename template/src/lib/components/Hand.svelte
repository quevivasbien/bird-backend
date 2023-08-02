<script lang="ts">
	import type { Card } from "$lib/types";
	import CardView from "$lib/components/CardView.svelte";
	import Dropdown from "./Dropdown.svelte";
	import { handSort } from "$lib/stores";
	import { sortCards } from "$lib/utils";

    export let cards: Card[];

	$: sortedCards = sortCards(cards, $handSort);

	const sortItems = [
		{action: () => $handSort = true, label: 'Ascending'},
		{action: () => $handSort = false, label: 'Descending'},
	];
</script>

<div class="flex flex-row space-x-8">
	<h2 class="text-2xl">Your hand</h2>
	<Dropdown title="Sort cards" items={sortItems} />
</div>
{#key $handSort}
	<div class="flex flex-wrap my-8">
		{#each sortedCards as card}
			<CardView {card} />
		{/each}
	</div>
{/key}
