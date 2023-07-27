<script lang="ts">
	import type { Card } from "$lib/types";
	import CardView from "$lib/components/CardView.svelte";

    export let cards: Card[];
</script>


<h2>Your hand</h2>
{#each cards as card}
	<CardView {card} />
{/each}
