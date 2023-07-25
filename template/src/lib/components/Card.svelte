<script lang="ts">
    export let color: number;
    export let value: number;

    const colors = ["", "Red", "Yellow", "Green", "Black"]

    let description: string;
    if (color === 0) {
        description = "Bird";
    } else {
        description = `${colors[color]} ${value}`;
    }
</script>

<div>
    {description}
</div>