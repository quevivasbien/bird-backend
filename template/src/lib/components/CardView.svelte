<script lang="ts">
	import type { Card } from "$lib/types";

    export let card: Card;

    const colors = ["", "Red", "Yellow", "Green", "Black"]

    let description: string;
    if (card.color === 0) {
        description = "Bird";
    } else {
        description = `${colors[card.color]} ${card.value}`;
    }
</script>

<div>
    {description}
</div>