<script lang="ts">
	import Card from '$lib/components/Card.svelte';
	import { bidStore, userStore } from '$lib/stores.js';
	import { onDestroy, onMount } from 'svelte';

	export let data;

	const { subscribeToBids, submitBid } = data;

	let sse: EventSource | undefined;

	onMount(() => {
		sse = subscribeToBids();
		sse?.addEventListener('update', (e) => {
			$bidStore = JSON.parse(e.data);
		});
	});

	onDestroy(() => {
		sse?.close();
	});

	const myIndex = $bidStore?.players.indexOf($userStore?.name ?? '') ?? 0;
	const myHand = $bidStore?.hands[myIndex] ?? [];
</script>

<h1>Bidding</h1>
<h2>Your hand</h2>
{#each myHand as card}
	<Card color={card.color} value={card.value} />
{/each}
