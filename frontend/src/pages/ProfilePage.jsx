import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { api } from "../api";
import { useAuth } from "../auth";
import PostCard from "../components/PostCard";

export default function ProfilePage() {
  const { username } = useParams();
  const { user: currentUser } = useAuth();
  const profileUsername = username || currentUser.username;
  const isOwnProfile = profileUsername.toLowerCase() === currentUser.username.toLowerCase();

  const [profile, setProfile] = useState(null);
  const [posts, setPosts] = useState([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    setError("");

    Promise.all([
      api(`/api/users/${encodeURIComponent(profileUsername)}`),
      api(`/api/users/${encodeURIComponent(profileUsername)}/posts`),
    ])
      .then(([profileData, postsData]) => {
        if (cancelled) return;
        setProfile(profileData.user);
        setPosts(postsData.posts);
      })
      .catch((err) => {
        if (!cancelled) setError(err.message);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });

    return () => {
      cancelled = true;
    };
  }, [profileUsername]);

  return (
    <section>
      <header className="page-header">
        <p className="eyebrow">{isOwnProfile ? "Your account" : "Public profile"}</p>
        <h1>{isOwnProfile ? "Your profile" : "Profile"}</h1>
        <p>
          {isOwnProfile
            ? "This is how other people see your account."
            : "Public posts from this account."}
        </p>
      </header>

      {loading ? <p className="muted">Loading profile...</p> : null}
      {error ? <p className="error">{error}</p> : null}

      {profile ? (
        <>
          <div className="profile-card">
            <div className="avatar">{profile.displayName.slice(0, 1).toUpperCase()}</div>
            <div>
              <h2>{profile.displayName}</h2>
              <p className="muted">@{profile.username}</p>
              {profile.bio ? <p className="bio">{profile.bio}</p> : null}
              <p className="muted">
                {profile.postCount} {profile.postCount === "1" ? "post" : "posts"} · joined{" "}
                {new Date(profile.createdAt).toLocaleDateString()}
              </p>
            </div>
          </div>

          <h3 className="section-title">Posts</h3>
          {posts.length === 0 ? (
            <p className="muted">
              {isOwnProfile ? (
                <>
                  You have not posted yet. <Link to="/">Write something on Home</Link>.
                </>
              ) : (
                "This user has not posted yet."
              )}
            </p>
          ) : (
            <div className="feed">
              {posts.map((post) => (
                <PostCard key={post.id} post={post} />
              ))}
            </div>
          )}
        </>
      ) : null}
    </section>
  );
}
