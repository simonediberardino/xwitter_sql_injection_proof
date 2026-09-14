import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { api } from "../api";
import PostCard from "../components/PostCard";

export default function HomePage() {
  const [posts, setPosts] = useState([]);
  const [content, setContent] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  async function loadPosts() {
    const data = await api("/api/posts");
    setPosts(data.posts);
  }

  useEffect(() => {
    loadPosts()
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  async function handleSubmit(event) {
    event.preventDefault();
    setError("");
    try {
      const data = await api("/api/posts", {
        method: "POST",
        body: JSON.stringify({ content }),
      });
      setPosts((current) => [data.post, ...current]);
      setContent("");
    } catch (err) {
      setError(err.message);
    }
  }

  return (
    <section>
      <header className="page-header">
        <p className="eyebrow">Live feed</p>
        <h1>Home</h1>
        <p>Share a short update with everyone on Xwitter.</p>
      </header>

      <form className="composer" onSubmit={handleSubmit}>
        <textarea
          value={content}
          maxLength={280}
          rows={3}
          placeholder="What's happening?"
          onChange={(event) => setContent(event.target.value)}
        />
        <div className="composer-row">
          <span>{content.length}/280</span>
          <button type="submit" disabled={!content.trim()}>
            Post
          </button>
        </div>
      </form>

      {error ? <p className="error">{error}</p> : null}

      {loading ? (
        <p className="muted">Loading posts...</p>
      ) : posts.length === 0 ? (
        <p className="muted">
          No posts yet. Be the first, or <Link to="/search">find people</Link> to follow the conversation.
        </p>
      ) : (
        <div className="feed">
          {posts.map((post) => (
            <PostCard key={post.id} post={post} />
          ))}
        </div>
      )}
    </section>
  );
}
