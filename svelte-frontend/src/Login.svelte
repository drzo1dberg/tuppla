<script>
    import {authToken} from './stores';
    let username = '';
    let password = '';

    async function login(){
        const res = await fetch('http://localhost:8080/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password}),
        });
        const data = await res.json();
        if (res.ok) {
            authToken.set(data.token);
            //SECURITY RISK need to scale up to httpOnlyCookies but ok for now
            localStorage.setItem('token', data.token);
            alert('successfully login!');
            //redirect to timeline
        } else {
            alert('Error: ' + data.error);
        }
    }
</script>
<h1>Login</h1>
<form on:submit|preventDefault={login}>
  <input type="text" bind:value={username} placeholder="username" required />
  <input type="password" bind:value={password} placeholder="password" required />
  <button type="submit">Login</button>
</form>