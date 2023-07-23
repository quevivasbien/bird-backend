<script lang="ts">
	import { lobbyStore } from '$lib/stores';
	import { onMount } from 'svelte';

    const UPDATE_INTVL = 500;

    export let data;

    const { getLobbyState } = data;

    let t: string | undefined = "";

    onMount(async () => {
        t = await getLobbyState();
        // setInterval(
        //     async () => t = await getLobbyState(),
        //     UPDATE_INTVL,
        // );
    });
</script>

{t ?? ""}