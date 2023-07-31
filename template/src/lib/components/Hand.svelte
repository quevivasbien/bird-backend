<script lang="ts">
	import type { Card } from "$lib/types";
	import CardView from "$lib/components/CardView.svelte";
	import Dropdown from "./Dropdown.svelte";

    export let cards: Card[];

	function sortCards(ascending: boolean) {
		const sign = ascending ? 1 : -1;
		cards = cards.sort((a, b) => {
			return sign * ((a.color * 100 + a.value) - (b.color * 100 + b.value));
		});
	}

	const sortItems = [
		{action: () => sortCards(true), label: 'Ascending'},
		{action: () => sortCards(false), label: 'Descending'},
	];
</script>

<div class="flex flex-row space-x-8">
	<h2 class="text-2xl">Your hand</h2>
	<Dropdown title="Sort cards" items={sortItems} />
</div>
<div class="flex flex-wrap my-8">
	{#each cards as card}
		<CardView {card} />
	{/each}
</div>
