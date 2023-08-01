<script lang="ts">
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';
	import Dropdown from '$lib/components/Dropdown.svelte';
	import { lobbyStore, userStore } from '$lib/stores';
	import { onDestroy, onMount } from 'svelte';

	export let data;

	const { subscribeToLobby, swapPlayers, leaveLobby, startBidding, receiveBidState } = data;

	let sse: EventSource | undefined;

	onMount(() => {
		sse = subscribeToLobby();
        if (sse === undefined) {
            // no valid lobby info; navigate home
            goto(`${base}/`);
            return;
        }
		sse.addEventListener("update", (e) => {
            $lobbyStore = JSON.parse(e.data);
        });
		sse.addEventListener("continue", (e) => {
			receiveBidState().then(([ok, status]) => {
				if (ok) {
					goto(`${base}/bidding`);
				}
				else {
					console.log("Something went wrong when trying to fetch bid state. Status = " + status);
					goto(`${base}/`);
				}
			});
		});
		sse?.addEventListener("message", (e) => console.log("message:", e.data));
	});

	onDestroy(async () => {
        sse?.close();
		const [ok, status] = await leaveLobby();
		if (!ok) {
			console.log('When attempting to leave lobby, got status', status);
		}
	});

	let host: string = '';
	let players: string[] = [];
	$: if ($lobbyStore !== undefined) {
		({ host, players } = $lobbyStore);
	}

	$: amHost = $userStore?.name === host;

	async function swap(i: number, j: number) {
		const [ok, status] = await swapPlayers(i, j);
		if (!ok) {
			console.log('When attempting to swap players, got status', status);
		}
	}

	function switchPlayers(i: number, j: number) {
		if (j === i) {
			return;
		}
		swap(i, j);
	}

	function itemsForPlayer(i: number) {
		return [0, 1, 2, 3]
			.filter((j) => j !== i)
			.map((j) => { return {
				'action': () => switchPlayers(i, j),
				'label': `Player ${j + 1}`,
			}});
	}

	async function attemptStartBidding() {
		const [ok, status] = await startBidding();
        if (!ok) {
			console.log('When attempting to start bidding, got status', status);
		}
	}
</script>

<h1 class="text-3xl">Lobby for game <span class="italic">{$lobbyStore?.id ?? ''}</span></h1>
<div class="max-w-sm m-8">
    {#each players as player, i}
        {#if i === 0}
            <h3 class="text-xl font-bold">Team 1</h3>
        {:else if i === 2}
            <h3 class="text-xl font-bold">Team 2</h3>
        {/if}
        <div class="flex flex-row ml-4 my-4 items-center space-x-8">
            <div class="flex flex-grow justify-start">
                {i + 1}. {player || 'Empty (AI)'}{#if player === host}&ThickSpace;(host){/if}
            </div>
            {#if amHost}
                <div class="flex justify-end">
                    <Dropdown title="Swap position" items={itemsForPlayer(i)} />
                </div>
            {/if}
        </div>
    {/each}
    {#if amHost}
        <div class="pt-4 border-t" />
        <button class="p-2 drop-shadow-lg rounded text-white bg-violet-800 hover:bg-violet-900 disabled:bg-gray-400" on:click={attemptStartBidding}>Start game</button>
    {/if}
</div>
