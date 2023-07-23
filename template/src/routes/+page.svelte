<script lang="ts">
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';
	import { userStore } from '$lib/stores';

	export let data;

	const { createLobby, joinLobby } = data;

	let newGameID: string = '';
	let createGameStatusText: string = '';

	function attemptCreateLobby() {
		if (!newGameID) {
			return;
		}
		createLobby(newGameID).then(([ok, status]) => {
			if (ok) {
				goto(base + '/lobby');
			} else if (status === 409) {
				createGameStatusText = 'Game name is already taken';
			} else if (status === 401) {
				$userStore = undefined;
				goto(base + '/login');
			} else {
				createGameStatusText = 'Something went wrong: status ' + status;
			}
		});
	}

	let joinGameID: string = '';
	let joinGameStatusText: string = '';

	function attemptJoinLobby() {
		if (!joinGameID) {
			return;
		}
		joinLobby(joinGameID).then(([ok, status]) => {
			if (ok) {
				goto(base + '/lobby');
			} else if (status === 404) {
				joinGameStatusText = 'No lobby found with name ' + joinGameID;
			} else if (status === 409) {
				createGameStatusText = 'Lobby is full';
			} else if (status === 401) {
				$userStore = undefined;
				goto(base + '/login');
			} else {
				joinGameStatusText = 'Something went wrong: status ' + status;
			}
		});
	}
</script>

<h1>Bird</h1>

{#if $userStore !== undefined}
	<div>
		<h2>Create new game</h2>
		<form on:submit|preventDefault={attemptCreateLobby}>
			<label>
				<span>Game name</span>
				<input type="text" bind:value={newGameID} />
			</label>
			<button type="submit">Create</button>
		</form>
		{#if createGameStatusText}
			<div>{createGameStatusText}</div>
		{/if}
	</div>
	<h2>Join game</h2>
	<form on:submit|preventDefault={attemptJoinLobby}>
		<label>
			<span>Game name</span>
			<input type="text" bind:value={joinGameID} />
		</label>
		<button type="submit">Join</button>
	</form>
	{#if joinGameStatusText}
		<div>{joinGameStatusText}</div>
	{/if}
{:else}
	<div>
		<a href={base + '/login'}>Log in</a> to get started
	</div>
{/if}
