<script>
    import { authToken } from './stores';
    let posts = [];
    let content = '';
  
    $: token = $authToken;
  
    async function loadPosts() {
      const res = await fetch('http://localhost:8080/api/posts', {
        headers: {
          'Authorization': token,
        },
      });
      const data = await res.json();
      if (res.ok) {
        posts = data.posts;
      } else {
        alert('Error loading posts');
      }
    }
  
    async function createPost() {
      const res = await fetch('http://localhost:8080/api/posts', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': token,
        },
        body: JSON.stringify({ content }),
      });
      const data = await res.json();
      if (res.ok) {
        content = '';
        await loadPosts();
      } else {
        alert('error creating post');
      }
    }
  
    // posts load when component finished loading
    loadPosts();
  </script>
  
  <h1>Timeline</h1>
  
  <form on:submit|preventDefault={createPost}>
    <textarea bind:value={content} placeholder="Got something to say?" required></textarea>
    <button type="submit">Posten</button>
  </form>
  
  <ul>
    {#each posts as post}
      <li>
        <strong>{post.User.Username}</strong> ({new Date(post.CreatedAt).toLocaleString()}):
        <p>{post.Content}</p>
      </li>
    {/each}
  </ul>
  