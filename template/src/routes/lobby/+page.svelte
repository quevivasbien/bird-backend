<script lang="ts">
	import { lobbyStore, userStore } from '$lib/stores';
	import { onDestroy, onMount } from 'svelte';

	export let data;

	const { subscribeToLobby, swapPlayers, leaveLobby, startGame } = data;

	let sse: EventSource | undefined;

	onMount(async () => {
		sse = await subscribeToLobby();
		sse?.addEventListener("update", (e) => {
            $lobbyStore = JSON.parse(e.data);
        });
	});

	onDestroy(async () => {
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

	function movePlayerUp(i: number) {
		const newPos = (i - 1 + 4) % 4;
		swap(i, newPos);
	}

	function movePlayerDown(i: number) {
		const newPos = (i + 1 + 4) % 4;
		swap(i, newPos);
	}

	$: readyToStart = $lobbyStore?.players.reduce((acc, x) => acc && x !== '', true) ?? false;

	function attemptStartGame() {
		// todo
	}
</script>

<div>
	<div>Players</div>
	{#each players as player, i}
		{#if i === 0}
			<h2>Team 1</h2>
		{:else if i === 2}
			<h2>Team 2</h2>
		{/if}
		<div>
			<div>
				{i}. {player || 'Empty'}{#if player === host} (host){/if}
			</div>
			{#if amHost}
				<button on:click={() => movePlayerUp(i)}>Up</button>
				<button on:click={() => movePlayerDown(i)}>Down</button>
			{/if}
		</div>
	{/each}
</div>

{#if amHost}
	<button on:click={attemptStartGame} disabled={!readyToStart}>Start game</button>
{/if}
