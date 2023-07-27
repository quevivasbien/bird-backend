<script lang="ts">
	import Hand from '$lib/components/Hand.svelte';
	import WidowExchange from '$lib/components/WidowExchange.svelte';
	import { gameStore, userStore } from '$lib/stores';
	import { onDestroy, onMount } from 'svelte';

	export let data;

	const { subscribeToGame, getWidow, startRound, getScore, playCard } = data;

	let sse: EventSource | undefined;

	onMount(() => {
		sse = subscribeToGame();
		sse?.addEventListener('update', (e) => {
			$gameStore = JSON.parse(e.data);
		});
		sse?.addEventListener('continue', (e) => {
			// todo
			console.log('got continue signal');
		});
	});

	onDestroy(() => {
		sse?.close();
	});

	const yourIndex = $gameStore?.players.indexOf($userStore?.name ?? '') ?? -1;
	const tookBid = $gameStore?.bidWinner === yourIndex;

	$: trumpSelected = $gameStore?.trump ?? 0 != 0;
	$: yourHand = $gameStore?.hand ?? [];

	let trumpSelection: number = 0;
</script>

{#if !trumpSelected}
	{#if tookBid}
		{#await getWidow()}
			loading widow...
		{:then widow}
			<WidowExchange widow={widow ?? []} {yourHand} />
		{/await}
		<label>
			<input type="radio" bind:group={trumpSelection} value="1" />
			<span>Red</span>
		</label>
		<label>
			<input type="radio" bind:group={trumpSelection} value="2" />
			<span>Yellow</span>
		</label>
		<label>
			<input type="radio" bind:group={trumpSelection} value="3" />
			<span>Green</span>
		</label>
		<label>
			<input type="radio" bind:group={trumpSelection} value="4" />
			<span>Black</span>
		</label>
	{:else}
		<div>Waiting for player {$gameStore?.bidWinner ?? -1 + 1} to choose trump color...</div>
	{/if}
{:else}
	game
{/if}

<Hand cards={yourHand} />
