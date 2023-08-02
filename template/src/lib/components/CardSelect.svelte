<script lang="ts">
	import type { Card } from "$lib/types";
	import CardView from "$lib/components/CardView.svelte";
	import Dropdown from "$lib/components/Dropdown.svelte";
	import { handSort } from "$lib/stores";
	import { sortCards } from "$lib/utils";

    export let cards: Card[];
    export let selection: Card;

    $: sortedCards = sortCards(cards, $handSort);

    const sortItems = [
        {action: () => $handSort = true, label: 'Ascending'},
        {action: () => $handSort = false, label: 'Descending'},
    ];

    let selectionIndex: number = 0;
    $: selection = sortedCards[selectionIndex];
</script>

<div class="flex flex-row space-x-8 items-center mt-4">
	<h2 class="text-2xl">Your hand</h2>
	<Dropdown title="Sort cards" items={sortItems} />
    <div>Select card to play</div>
</div>
{#key $handSort}
    <div class="flex flex-wrap my-8 space-y-4">
        {#each sortedCards as card, i}
            <label class="flex flex-row space-x-2 cursor-pointer">
                <input class="hidden" type="radio" bind:group={selectionIndex} value={i} />
                <CardView {card} highlighted={selectionIndex === i}/>
            </label>
        {/each}
    </div>
{/key}