<script lang="ts">
	import type { Card } from "$lib/types";
	import CardView from "$lib/components/CardView.svelte";

    export let cards: Card[];
    export let selection: Card;
</script>

{#each cards as card, i}
    <label>
        <input type="radio" bind:group={selection} value={card} />
        <CardView {card} />
    </label>
{/each}