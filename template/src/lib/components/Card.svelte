<script lang="ts">
    export let suit: number;
    export let value: number;

    const suites = ["", "Red", "Yellow", "Green", "Black"]

    let description: string;
    if (suit === 0) {
        description = "Rook";
    } else {
        description = `${suites[suit]} ${value}`;
    }
</script>

<div>
    {description}
</div>