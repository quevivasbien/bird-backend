<script lang="ts">
	import type { Card } from "$lib/types";
    import BirdLogo from "$lib/assets/bird-logo.svg";

    export let card: Card;
    export let highlighted: boolean = false;

    const colors = ["", "#7f1d1d", "#fbbf24", "#14532d", "#000000"];
    const colorNames = ["", "R", "Y", "G", "B"];
    const color = colors[card.color];

    const angle = Math.floor(Math.random() * 20 - 10);
    const rotation = `transform: rotate(${angle}deg)`;
</script>

<div class={`relative sm:w-48 sm:h-64 w-32 h-48 border-2 drop-shadow rounded p-2 m-1 bg-white font-bold ${highlighted ? "border-violet-400" : ""}`} style={`${rotation}; color: ${color}`}>
    {#if card.color !== 0}
        <div class="absolute top-4 left-4">
            {colorNames[card.color]}{card.value}
        </div>
        <div class="absolute bottom-4 right-4 rotate-180">
            {colorNames[card.color]}{card.value}
        </div>
    {/if}
    <div class="flex flex-col w-full h-full justify-center items-center">
        {#if card.color === 0}
            <div class="font-cursive text-3xl text-violet-900">Bird!</div>
            <img src={BirdLogo} alt="Bird" />
        {:else}
            <div class="font-cursive font-bold text-6xl" style={`color: ${color}`}>{card.value}</div>
        {/if}
    </div>
</div>
