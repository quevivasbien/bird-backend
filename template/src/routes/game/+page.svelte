<script lang="ts">
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';
	import CardSelect from '$lib/components/CardSelect.svelte';
	import Hand from '$lib/components/Hand.svelte';
	import Table from '$lib/components/Table.svelte';
	import WidowExchange from '$lib/components/WidowExchange.svelte';
	import { gameStore, userStore } from '$lib/stores';
	import type { Card, GameInfo } from '$lib/types';
	import { onDestroy, onMount } from 'svelte';

	export let data;

	const { subscribeToGame, getWidow, startRound, getScore, playCard } = data;

	let sse: EventSource | undefined;

	onMount(() => {
		sse = subscribeToGame();
		if (sse === undefined) {
			// no valid game info; navigate home
			goto(`${base}/`);
			return;
		}
		sse.addEventListener('update', (e) => {
			const data = JSON.parse(e.data);
			data.table = data.table ?? [];
			console.log("data update:", data);
			$gameStore = data;
		});
		sse.addEventListener('continue', (e) => {
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

	$: currentPlayer = $gameStore?.currentPlayer ?? -1;

	let trumpColor = '';
	$: if (trumpSelected) {
		trumpColor = getTrumpColor($gameStore);
	}
	function getTrumpColor(gameInfo: GameInfo | undefined) {
		if (gameInfo === undefined) {
			return '';
		}
		const trump = gameInfo.trump;
		if (trump === 1) {
			return 'Red';
		}
		if (trump === 2) {
			return 'Yellow';
		}
		if (trump === 3) {
			return 'Green';
		}
		if (trump === 4) {
			return 'Black';
		}
		return '';
	}

	let toWidow: Card[] = [];
	let fromWidow: Card[]  = [];
	let startGameStatus = "";
	async function submitCreateGame() {
		if (toWidow.length !== fromWidow.length) {
			startGameStatus = "You must take the same number of cards out of your hand as you take out of the widow.";
			return;
		}
		if (trumpSelection === 0) {
			startGameStatus = "You need to choose the trump color.";
			return;
		}
		const [ok, status] = await startRound(trumpSelection, toWidow, fromWidow);
		if (!ok) {
			console.log("When trying to start round, got status = " + status);
		}
	}

	let selectedCard: Card;
	async function submitSelectCard() {
		if (selectedCard === undefined) {
			return;
		}
		const [ok, status] = await playCard(selectedCard);
		if (!ok) {
			console.log("When trying to play card, got status = " + status);
		}
	}
</script>

{#if !trumpSelected}
	{#if tookBid}
    <form on:submit|preventDefault={submitCreateGame}>
            <h1 class="text-3xl">Choose cards to exchange with widow</h1>
			{#await getWidow()}
				loading widow...
			{:then widow}
				<WidowExchange widow={widow ?? []} {yourHand} bind:toWidow={toWidow} bind:fromWidow={fromWidow} />
			{/await}
            <h1 class="text-3xl">Choose trump color</h1>
            <div class="flex flex-row space-x-4 items-center m-4">
                <label>
                    <input type="radio" bind:group={trumpSelection} value={1} />
                    <span>Red</span>
                </label>
                <label>
                    <input type="radio" bind:group={trumpSelection} value={2} />
                    <span>Yellow</span>
                </label>
                <label>
                    <input type="radio" bind:group={trumpSelection} value={3} />
                    <span>Green</span>
                </label>
                <label>
                    <input type="radio" bind:group={trumpSelection} value={4} />
                    <span>Black</span>
                </label>
                <div class="flex-grow text-center">
                    <button type="submit" disabled={trumpSelection == 0}>Submit</button>
                </div>
            </div>
		</form>
		{#if startGameStatus}
			<div>{startGameStatus}</div>
		{/if}
	{:else}
		<div>Waiting for player {($gameStore?.bidWinner ?? -1) + 1} to choose trump color...</div>
		<Hand cards={yourHand} />
	{/if}
{:else}
	<div class="text-2xl">{trumpColor} is trump.</div>
	<Table />
	{#if currentPlayer === yourIndex}
		<div class="text-3xl my-4">Your turn</div>
		<form on:submit|preventDefault={submitSelectCard}>
			<CardSelect cards={yourHand} bind:selection={selectedCard} />
			<button type="submit">Play card</button>
		</form>
	{:else}
		<div class="text-3xl my-4">Player {currentPlayer + 1}'s turn</div>
		<Hand cards={yourHand} />
	{/if}
{/if}
