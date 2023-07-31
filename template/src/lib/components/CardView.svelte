<script lang="ts">
	import type { Card } from "$lib/types";
    import BirdLogo from "$lib/assets/bird-logo.svg";
	import { slide } from "svelte/transition";

    export let card: Card;

    const colors = ["", "red", "yellow", "green", "black"];
    const color = colors[card.color];

    const angle = Math.floor(Math.random() * 30 - 15);
    const rotation = `transform: rotate(${angle}deg)`;
</script>

<div class="flex flex-col sm:w-48 sm:h-64 w-32 h-48 justify-center items-center border drop-shadow rounded p-2 m-1 bg-white" style={rotation}>
    {#if card.color === 0}
        <div class="font-cursive text-3xl text-violet-900">Bird!</div>
        <img src={BirdLogo} alt="Bird" />
    {:else}
        <div class={`font-cursive font-bold text-6xl text-${color}-800`}>{card.value}</div>
    {/if}
</div>