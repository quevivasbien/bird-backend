<script lang="ts">
	import { base } from "$app/paths";
	import { userStore } from "$lib/stores";

    export let data;

    const { logout } = data;
</script>

<svelte:head>
    <title>Bird</title>
</svelte:head>

<div>
    {#if $userStore !== undefined}
        {$userStore.name}
        <a href={`${base}/`} on:click={logout}>Log out</a>
    {:else}
        <a href={`${base}/login`}>Log in</a>
    {/if}

</div>

<slot />