<script lang="ts">
	import { lobbyStore } from '$lib/stores';
	import { onDestroy, onMount } from 'svelte';

	const UPDATE_INTVL = 1000;

	export let data;

	const { getLobbyState, leaveLobby } = data;

    let interval: number;

	onMount(async () => {
		$lobbyStore = await getLobbyState();
		interval = setInterval(async () => ($lobbyStore = await getLobbyState()), UPDATE_INTVL);
	});

	onDestroy(async () => {
        console.log("lobby ondestroy called");
        clearInterval(interval);
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
</script>

<div>
	<div>Host: {host || 'Host left'}</div>
	<div>Players</div>
	<ol>
		{#each players as player}
			<li>{player || 'Empty'}</li>
		{/each}
	</ol>
</div>
