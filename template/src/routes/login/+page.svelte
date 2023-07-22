<script lang="ts">
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';

	export let data;

	const { login } = data;

	let name: string = '';
	let password: string = '';

	let statusText: string = '';

	async function submitForm() {
		if (!name || !password) {
			return;
		}
		const [ok, status] = await login(name, password);
		if (ok) {
			goto(base + '/');
			return;
		}
		if (status === 401) {
			statusText = 'Invalid username and/or password';
			return;
		}
		statusText = 'Something went wrong';
	}
</script>

<form on:submit|preventDefault={submitForm}>
	<label>
		<div>Username</div>
		<input type="text" bind:value={name} />
	</label>
	<label>
		<div>Password</div>
		<input type="password" bind:value={password} />
	</label>
	<button type="submit">Login</button>
</form>
{#if statusText}
	<div>{statusText}</div>
{/if}

<div>
	New? <a href={base + '/login/register'}>Register an account</a>
</div>
