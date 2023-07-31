<script lang="ts">
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';
	import Dropdown from '$lib/components/Dropdown.svelte';
	import Hand from '$lib/components/Hand.svelte';
	import { bidStore, lobbyStore, userStore } from '$lib/stores';
	import type { BidInfo } from '$lib/types';
	import { onDestroy, onMount } from 'svelte';

    $userStore = {
        name: "admin",
        admin: true,
        expireTime: 0,
    };
    $lobbyStore = {
        id: "test",
        host: $userStore?.name ?? "",
        players: [$userStore?.name ?? "", "dummy1", "dummy2", "dummy3"],
        started: false,
    };
    $bidStore = {
        id: "admin",
        done: false,
        players: ["admin", "dummy1", "dummy2", "dummy3"],
        hand: [{color: 0, value: 0}, {color: 1, value: 1}, {color: 2, value: 2}],
        passed: [false, false, false, false],
        currentBidder: 0,
        bid: 0,
    };

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
		
	}

	function pass() {
		attemptSubmitBid(0);
	}
</script>

<h1 class="text-3xl">Bidding</h1>

<div class="mt-8 mb-16">
	{#if !biddingDone}
		<div class="text-lg">
			{#if currentBidder === yourIndex}Your{:else}Player {currentBidder + 1}'s{/if} turn to bid
		</div>
		{#if currentBid > 0}
			<div>Current bid: {currentBid} (Player {bidLeader + 1})</div>
		{/if}
		{#if currentBidder === yourIndex}
			<div class="flex flex-col space-y-4">
                <div class="flex flex-row space-x-2 items-center">
                    <button class="p-1 border rounded w-10" type="button" on:click={() => (yourBid -= 5)} disabled={yourBid <= currentBid + 5}>&#8595;</button>
                    <div class="text-2xl w-16 text-center">{yourBid}</div>
                    <button class="p-1 border rounded w-10" type="button" on:click={() => (yourBid += 5)}>&#8593;</button>
                </div>
                <div class="flex-flex-row">
                    <button class="p-2 drop-shadow-lg rounded text-white bg-violet-800 hover:bg-violet-900 disabled:bg-gray-400" on:click={() => attemptSubmitBid(yourBid)}>Submit bid</button>
                    <button class="p-2 drop-shadow-lg rounded text-white bg-violet-800 hover:bg-violet-900 disabled:bg-gray-400" on:click={pass}>Pass</button>
                </div>
			</div>
		{/if}
	{:else}
		<div>
			Player {currentBidder + 1} won the bid for {currentBid}!
		</div>
	{/if}
</div>

<Hand cards={yourHand} />
