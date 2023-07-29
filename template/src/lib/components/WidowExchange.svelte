<script lang="ts">
	import type { Card } from "$lib/types";
	import CardView from "$lib/components/CardView.svelte";

    export let widow: Card[];
    export let yourHand: Card[];

    export let toWidow: Card[];
    export let fromWidow: Card[];

    let handIndices: number[] = [];
    let widowIndices: number[] = [];
    $: toWidow = handIndices.map((i) => yourHand[i]);
    $: fromWidow = widowIndices.map((i) => widow[i]);
</script>

<h3>Widow</h3>
<div>
    {#each widow as card, i}
        <CardView {card} />
        <label>
            <input type="checkbox" bind:group={widowIndices} value={i} />
            Exchange with your hand
        </label>
    {/each}
</div>

<h3>Your hand</h3>
<div>
    {#each yourHand as card, i}
        <CardView {card} />
        <label>
            <input type="checkbox" bind:group={handIndices} value={i} />
            Exchange with widow
        </label>
    {/each}
</div>