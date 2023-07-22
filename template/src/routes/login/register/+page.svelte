<script lang="ts">
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';

	export let data;

	const { register } = data;

	let name: string = '';
	let password: string = '';
	let passwordVerify: string = '';

	let statusText: string = '';

	async function submitForm() {
		if (!name || !password) {
			return;
		}
		if (password !== passwordVerify) {
			statusText = 'Passwords do not match';
			return;
		}
		const [ok, status] = await register(name, password);
		if (ok) {
			goto(base + '/login');
			return;
		}
		if (status === 409) {
			statusText = 'Username is already taken';
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
	<label>
		<div>Verify password</div>
		<input type="password" bind:value={passwordVerify} />
	</label>
	<button type="submit">Register</button>
</form>
{#if statusText}
	<div>{statusText}</div>
{/if}
