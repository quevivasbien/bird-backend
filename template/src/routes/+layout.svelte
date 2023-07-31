<script lang="ts">
	import '../app.css'
	import { base } from '$app/paths';
	import { userStore } from '$lib/stores';
	import { onMount } from 'svelte';

	export let data;

	const { logout, syncAuth } = data;

	onMount(() => {
		if ($userStore === undefined) {
			syncAuth().then((v) => $userStore = v);
		}
	});
</script>

<svelte:head>
	<title>Bird</title>
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Gloria+Hallelujah&family=Rubik:ital,wght@0,400;0,700;1,400;1,700&display=swap" rel="stylesheet">
</svelte:head>

<div class="flex flex-row sm:mx-auto mb-4 max-w-3xl px-8 py-4 border-x border-b rounded sticky top-0 z-10 bg-white">
	<div class="flex flex-grow justify-start items-center space-x-4">
		<a class="text-4xl font-cursive" style="text-decoration: none !important;" href={base + '/'}>Bird</a>
		<a class="font-bold" href={base + '/'}>Home</a>
	</div>
	<div class="flex justify-end items-center space-x-4">
		{#if $userStore !== undefined}
			<div>{$userStore.name}</div>
			<a class="font-bold" href={base + '/'} on:click={logout}>Log out</a>
		{:else}
			<a class="font-bold" href={base + '/login'}>Log in</a>
		{/if}
	</div>
</div>
<div class="sm:mx-auto mx-8 max-w-2xl">
	<slot />
</div>

<style lang="postcss">
	:global(a) {
		@apply text-violet-900;
	}
	:global(a:hover) {
		@apply underline;
	}

	:global(input[type=text]) {
		@apply p-2;
		@apply border;
		@apply drop-shadow-md;
		@apply rounded;
	}
	:global(input[type=password]) {
		@apply p-2;
		@apply border;
		@apply drop-shadow-md;
		@apply rounded;
	}

	:global(button[type=submit]) {
		@apply p-2;
		@apply drop-shadow-lg;
		@apply rounded;
		@apply text-white;
		@apply bg-violet-800;
	}
	:global(button[type=submit]:hover) {
		@apply bg-violet-900;
	}
	:global(button[type=submit]:disabled) {
		@apply bg-gray-400;
	}
</style>