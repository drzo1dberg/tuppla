<script>
    import { Router, Link, Route } from 'svelte-routing';
    import Home from "./Home.svelte";
    import Register from "./Register.svelte";
    import Login from './Login.svelte';
    import Timeline from './Timeline.svelte';    
</script>

<Router>
    <nav>
        <Link to= "/">Home</Link>
        <Link to= "/register">Registrieren</Link>
        <Link to= "/login">Login</Link>
        <Link to= "/timeline">Timeline</Link>
    </nav>

    <Route path="/" component={Home} />
    <Route path="/register" component={Register} />
    <Route path="/login" component={Login} />
    <Route path="/timeline" component={Timeline} />
</Router>