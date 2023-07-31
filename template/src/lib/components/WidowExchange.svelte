<script lang="ts">
	import type { Card } from "$lib/types";
	import CardView from "$lib/components/CardView.svelte";
	import { handSort } from "$lib/stores";
	import Dropdown from "./Dropdown.svelte";
	import { sortCards } from "$lib/utils";

    export let widow: Card[];
    export let yourHand: Card[];


	$: sortedHand = sortCards(yourHand, $handSort);

    const sortItems = [
        {action: () => $handSort = true, label: 'Ascending'},
        {action: () => $handSort = false, label: 'Descending'},
    ];

    export let toWidow: Card[];
    export let fromWidow: Card[];

    let handIndices: number[] = [];
    let widowIndices: number[] = [];
    $: toWidow = handIndices.map((i) => yourHand[i]);
    $: fromWidow = widowIndices.map((i) => widow[i]);
</script>

<div class="flex flex-row space-x-8 items-center mt-4">
	<h2 class="text-2xl">Widow</h2>
    <div>
        Select cards to take from widow
    </div>
</div>
<div class="flex flex-wrap my-8 space-y-4">
    {#each widow as card, i}
            <label class="cursor-pointer">
                <CardView {card} highlighted={widowIndices.includes(i)} />
                <input class="hidden" type="checkbox" bind:group={widowIndices} value={i} />
            </label>
    {/each}
</div>

<div class="flex flex-row space-x-8 items-center mt-4">
	<h2 class="text-2xl">Your hand</h2>
	<Dropdown title="Sort cards" items={sortItems} />
    <div>
        Select cards to give to widow
    </div>
</div>
{#key $handSort}
    <div class="flex flex-wrap my-8 space-y-4">
        {#each sortedHand as card, i}
            <label class="cursor-pointer">
                <CardView {card} highlighted={handIndices.includes(i)} />
                <input class="hidden" type="checkbox" bind:group={handIndices} value={i} />
            </label>
        {/each}
    </div>
{/key}