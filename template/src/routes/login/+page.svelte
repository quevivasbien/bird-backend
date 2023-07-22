<script lang="ts">
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';
    import { userStore } from '$lib/stores'

    export let data;

    const { login } = data;

    let name: string = "";
    let password: string = "";

    let statusText: string = "";

    async function submitForm() {
        if (!name || !password) {
            return;
        }
        const status = await login(name, password);
        if (status === 200) {
            $userStore = { name };
            goto(`${base}/`);
            return;
        }
        if (status === 401) {
            statusText = "Invalid username and/or password";
            return;
        }
        statusText = "Something went wrong"
    }
</script>

<form on:submit|preventDefault={submitForm}>
    <input type="text" bind:value={name}>
    <input type="password" bind:value={password}>
    <button type="submit">Login</button>
</form>
{#if statusText}
    <div>{statusText}</div>
{/if}

