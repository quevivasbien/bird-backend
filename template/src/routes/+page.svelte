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
{#if $userStore !== undefined}
	<div class="mb-8">
		<h2 class="text-3xl">Create new game</h2>
		<form class="flex space-x-4" on:submit={attemptCreateLobby}>
			<label class="flex flex-col">
				<div class="flex">Game name</div>
				<input class="flex" type="text" bind:value={newGameID} />
			</label>
			<button class="h-12 w-24 self-end" type="submit">Create</button>
		</form>
		{#if createGameStatusText}
			<div class="text-red-800 m-1">{createGameStatusText}</div>
		{/if}
	</div>
	<div class="mb-8">
		<h2 class="text-3xl">Join game</h2>
		<form class="flex space-x-4" on:submit={attemptJoinLobby}>
			<label class="flex flex-col">
				<div class="flex">Game name</div>
				<input class="flex" type="text" bind:value={joinGameID} />
			</label>
			<button class="h-12 w-24 self-end" type="submit">Join</button>
		</form>
		{#if joinGameStatusText}
			<div class="text-red-800 m-1">{joinGameStatusText}</div>
		{/if}
	</div>
{:else}
	<div class="text-center">
		<a class="text-purple-900 hover:underline" href={base + '/login'}>Log in</a> to get started
	</div>
{/if}
