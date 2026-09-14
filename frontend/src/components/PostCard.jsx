import { Link } from "react-router-dom";

export default function PostCard({ post }) {
  const date = new Date(post.createdAt);

  return (
    <article className="post-card">
      <div className="post-header">
        <Link to={`/u/${post.author.username}`} className="author-name">
          {post.author.displayName}
        </Link>
        <Link to={`/u/${post.author.username}`} className="author-handle">
          @{post.author.username}
        </Link>
        <time dateTime={post.createdAt}>
          {date.toLocaleString(undefined, {
            month: "short",
            day: "numeric",
            hour: "2-digit",
            minute: "2-digit",
          })}
        </time>
      </div>
      <p>{post.content}</p>
    </article>
  );
}
