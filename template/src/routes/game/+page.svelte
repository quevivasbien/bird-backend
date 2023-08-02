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
	import { fade } from 'svelte/transition';

	export let data;

	const { subscribeToGame, getWidow, startRound, getScore, playCard, finishPlay } = data;

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
			console.log('data update:', data);
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

	let players: string[] = [];
	let currentPlayer: number = -1;
	let lastWinner: number = -1;
	let table: Card[] = [];
	let done: boolean = false;
	$: if ($gameStore !== undefined) {
		({ players, currentPlayer, lastWinner, table, done } = $gameStore);
	}
	$: players = players.map((p) => p === '' ? 'AI' : p);

	let toWidow: Card[] = [];
	let fromWidow: Card[] = [];
	let startGameStatus = '';
	async function submitCreateGame() {
		if (toWidow.length !== fromWidow.length) {
			startGameStatus =
				'You must take the same number of cards out of your hand as you take out of the widow.';
			return;
		}
		if (trumpSelection === 0) {
			startGameStatus = 'You need to choose the trump color.';
			return;
		}
		const [ok, status] = await startRound(trumpSelection, toWidow, fromWidow);
		if (!ok) {
			console.log('When trying to start round, got status = ' + status);
		}
	}

	let selectedCard: Card | null;
	let cardSelectStatus = '';
	async function submitSelectCard() {
		if (selectedCard === undefined || selectedCard === null) {
			return;
		}
		if (!selectedCardIsOk()) {
			cardSelectStatus = 'Selected card is not a legal play';
			return;
		}
		cardSelectStatus = '';
		const [ok, status] = await playCard(selectedCard);
		if (!ok) {
			console.log('When trying to play card, got status = ' + status);
		}
	}

	function selectedCardIsOk() {
		if (selectedCard === undefined || selectedCard === null) {
			return false;
		}
		if (table.length === 0) {
			return true;
		}
		const leadingColor = table[0].color;
		const haveLeadingColor = yourHand.reduce((acc, x) => acc || x.color === leadingColor, false);
		if (haveLeadingColor) {
			console.log("haveLeadingColor");
			return selectedCard.color === leadingColor || selectedCard.color === 0 && leadingColor === $gameStore?.trump;
		}
		const haveTrump = yourHand.reduce((acc, x) => acc || x.color === $gameStore?.trump || x.color === 0, false);
		if (haveTrump) {
			console.log("haveTrump");
			return selectedCard.color === $gameStore?.trump || selectedCard.color === 0;
		}
		return true;
	}

	async function attemptFinishPlay() {
		const [ok, status] = await finishPlay();
		if (!ok) {
			console.log("Problem finishing round; got status = " + status);
		}
	}

	const SCORE_TIMEOUT = 1000;
	async function initGetScores() {
		const scores = await getScore();
		setTimeout(() => getFinalScores(scores as [number, number]), SCORE_TIMEOUT);
		return scores;
	}

	let finalScoreUpdateText = '';
	let finalScores: [number, number] | undefined = undefined;
	function getFinalScores(scores: [number, number]) {
		if ($gameStore === undefined) {
			console.log("Game state is undefined when calculating final scores");
			return;
		}	
		finalScores = [...scores];
		const teamTookBid = $gameStore.bidWinner % 2;
		if (scores[teamTookBid] < $gameStore.bid) {
			finalScoreUpdateText = `Team ${teamTookBid + 1} failed to make the bid and loses ${$gameStore.bid} points!`;
			finalScores[teamTookBid] = -$gameStore.bid;
		}
		else {
			finalScoreUpdateText = `Team ${teamTookBid + 1} made the bid!`;
		}
	}
</script>

{#if !trumpSelected}
	<!-- displayed while bid winner exchanges cards with widow and chooses trump -->
	{#if tookBid}
		<form on:submit={submitCreateGame}>
			<h1 class="text-3xl">Choose cards to exchange with widow</h1>
			{#await getWidow()}
				loading widow...
			{:then widow}
				<WidowExchange widow={widow ?? []} {yourHand} bind:toWidow bind:fromWidow />
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
		<div class="text-xl my-4">Waiting for player {($gameStore?.bidWinner ?? -1) + 1} to choose trump color...</div>
		<Hand cards={yourHand} />
	{/if}
{:else if !done}
	<!-- displayed after player chooses trump; this is the actual game -->
	<div class="text-2xl">{trumpColor} is trump.</div>
	{#key [table, lastWinner]}
		<Table cards={table} {players} leadingPlayer={lastWinner} />
	{/key}
	{#if table.length === 4}
		<!-- displayed at end of each play -->
		{#if yourIndex === 0}
			<button
				class="p-2 my-4 drop-shadow-lg rounded text-white bg-violet-800"
				on:click={attemptFinishPlay}>Finish play</button
			>
		{:else}
			<div class="text-3xl my-4">Waiting for next play</div>
		{/if}
		<Hand cards={yourHand} />
	{:else}
		<!-- card select stage -->
		{#if lastWinner !== undefined}
			<div class="text-2xl my-2">Player {lastWinner + 1} ({players[lastWinner]}) won last play.</div>
		{/if}
		{#if currentPlayer === yourIndex}
			<div class="text-3xl my-4">Your turn</div>
			<form on:submit={submitSelectCard}>
				<CardSelect cards={yourHand} bind:selection={selectedCard} />
				<button class="my-4" type="submit" disabled={selectedCard === null}>Play card</button>
				{#if cardSelectStatus}
					<div class="text-red-800">{cardSelectStatus}</div>
				{/if}
			</form>
		{:else}
			<div class="text-3xl my-4">Player {currentPlayer + 1} ({players[currentPlayer]})'s turn</div>
			<Hand cards={yourHand} />
		{/if}
	{/if}

	<div class="m-8">
		<h2 class="text-2xl">Discard piles</h2>
		<div>Team 1 ({players[0]} and {players[2]}): {$gameStore?.discardSize[0] ?? 0} cards</div>
		<div>Team 2 ({players[1]} and {players[3]}): {$gameStore?.discardSize[1] ?? 0} cards</div>
	</div>
{:else}
	<!-- displayed when round is complete -->
	<h1 class="3xl">Round complete!</h1>
	{#await initGetScores()}
		<div>Calculating scores...</div>
	{:then scores}
		<div class="my-8" transition:fade>
			<div class="text-lg">Preliminary scores:</div>
			<div>Team 1: {scores[0]}</div>
			<div>Team 2: {scores[1]}</div>
		</div>
		<div class="text-lg">Final scores:</div>
		{#if finalScores !== undefined}
			<div class="my-8" transition:fade>
				<div class="italic">{finalScoreUpdateText}</div>
				<div>Team 1: {finalScores[0]}</div>
				<div>Team 2: {finalScores[1]}</div>
			</div>
		{/if}
	{/await}
	<a class="my-8" href={`${base}/`}>Back to home</a>
{/if}
