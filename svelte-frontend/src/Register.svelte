<script>
    let username = '';
    let password = '';

    async function register () {
        const res = await fetch('http://localhost:8080/api/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json'},
            body: JSON.stringify({ username, password}),
        });
        const data = await res.json();
        if (res.ok) {
            alert('successfully registered!');
            // Redirect to Homepage
        } else {
            alert('Error: ' + data.error);
        }
    }
</script>

<h1>register</h1>
<form on:submit|preventDefault={register}>
    <input type="text" bind:value={username} placeholder="username" required />
    <input type="password" bind:value={password} placeholder="password" required />
    <button type="submit">submit</button>
</form>