<script lang="ts">
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';
	import { userStore } from '$lib/stores';

	export let data;

	const { createGame } = data;

	let createGameStatusText: string = '';

	function createAndJoinGame() {
		createGame().then(([ok, status]) => {
			if (ok) {
				goto(base + '/lobby');
			} else {
				createGameStatusText = 'Something went wrong: status ' + status;
			}
		});
	}
</script>

<h1>Bird</h1>

{#if $userStore !== undefined}
	<div>
		<button on:click={createAndJoinGame}>Create new game</button>
		{#if createGameStatusText}
			<div>{createGameStatusText}</div>
		{/if}
	</div>
	<h2>Join game</h2>
{:else}
	<div>
		<a href={base + '/login'}>Log in</a> to get started
	</div>
{/if}
