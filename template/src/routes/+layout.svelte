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
</svelte:head>

<div>
	<a href={base + '/'}>Home</a>
	{#if $userStore !== undefined}
		{$userStore.name}
		<a href={base + '/'} on:click={logout}>Log out</a>
	{:else}
		<a href={base + '/login'}>Log in</a>
	{/if}
</div>

<slot />
