<script lang="ts">
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';
	import Hand from '$lib/components/Hand.svelte';
	import { bidStore, userStore } from '$lib/stores.js';
	import type { BidInfo } from '$lib/types.js';
	import { onDestroy, onMount } from 'svelte';

	export let data;

	const { subscribeToBids, submitBid, receiveGameState } = data;

	let sse: EventSource | undefined;

	onMount(() => {
		sse = subscribeToBids();
		if (sse === undefined) {
			// no valid bidstate; navigate home
			goto(`${base}/`);
			return;
		}
		sse.addEventListener('update', (e) => {
			$bidStore = JSON.parse(e.data);
		});
		sse.addEventListener('continue', (e) => {
			console.log('winner is ', $bidStore?.currentBidder);
			biddingDone = true;
			receiveGameState().then(([ok, status]) => {
				if (ok) {
					setTimeout(() => goto(`${base}/game`), 2000);
				} else {
					console.log('Problem getting game info, status = ' + status);
				}
			});
		});
	});

	onDestroy(() => {
		sse?.close();
	});

	let biddingDone = false;

	const yourIndex = $bidStore?.players.indexOf($userStore?.name ?? '') ?? -1;
	const yourHand = $bidStore?.hand ?? [];

	$: currentBid = $bidStore?.bid ?? 0;
	$: currentBidder = $bidStore?.currentBidder ?? -1;
	$: console.log('currentBidder:', currentBidder);

	$: bidLeader = getBidLeader($bidStore);

	function getBidLeader(bidState: BidInfo | undefined) {
		if (bidState === undefined) {
			return -1;
		}
		let leader = (bidState.currentBidder - 1 + 4) % 4;
		while (bidState.passed[leader]) {
			leader = (leader - 1 + 4) % 4;
		}
		return leader;
	}

	let yourBid: number;
	$: updateYourBid(currentBid);

	function updateYourBid(b: number) {
		yourBid = Math.max(100, b + 5);
	}

	function attemptSubmitBid(b?: number) {
		if (b === undefined) {
			b = yourBid;
		}
		submitBid(b).then(([ok, status]) => {
			if (!ok) {
				console.log('Problem submitting bid; status = ', status);
			}
		});
	}

	function pass() {
		attemptSubmitBid(0);
	}
</script>

<h1>Bidding</h1>

<div>
	{#if !biddingDone}
		<div>
			{#if currentBidder === yourIndex}Your{:else}Player {currentBidder + 1}'s{/if} turn to bid
		</div>
		{#if currentBid > 0}
			<div>Current bid: {currentBid} (Player {bidLeader + 1})</div>
		{/if}
		{#if currentBidder === yourIndex}
			<form on:submit|preventDefault={() => attemptSubmitBid()}>
				<input type="text" bind:value={yourBid} />
				<button type="button" on:click={() => (yourBid -= 5)} disabled={yourBid <= currentBid + 5}
					>Down</button
				>
				<button type="button" on:click={() => (yourBid += 5)}>Up</button>
				<button type="submit">Submit</button>
				<button type="button" on:click={pass}>Pass</button>
			</form>
		{/if}
	{:else}
		<div>
			Player {currentBidder + 1} won the bid for {currentBid}!
		</div>
	{/if}
</div>

<Hand cards={yourHand} />
