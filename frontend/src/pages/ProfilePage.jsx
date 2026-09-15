import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { api } from "../api";
import { useAuth } from "../auth";
import PostCard from "../components/PostCard";

export default function ProfilePage() {
  const { username } = useParams();
  const { user: currentUser, updateProfile } = useAuth();
  const profileUsername = username || currentUser.username;
  const isOwnProfile = profileUsername.toLowerCase() === currentUser.username.toLowerCase();

  const [profile, setProfile] = useState(null);
  const [posts, setPosts] = useState([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const [form, setForm] = useState({
    username: "",
    displayName: "",
    bio: "",
  });
  const [formError, setFormError] = useState("");
  const [formSuccess, setFormSuccess] = useState("");
  const [saving, setSaving] = useState(false);

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
        setForm({
          username: profileData.user.username,
          displayName: profileData.user.displayName,
          bio: profileData.user.bio || "",
        });
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

  async function handleSave(event) {
    event.preventDefault();
    setFormError("");
    setFormSuccess("");
    setSaving(true);
    try {
      const updated = await updateProfile({
        username: form.username.trim(),
        displayName: form.displayName.trim(),
        bio: form.bio.trim(),
      });
      setProfile((current) => ({
        ...updated,
        postCount: current?.postCount || updated.postCount,
      }));
      setFormSuccess("Profile updated.");
    } catch (err) {
      setFormError(err.message);
    } finally {
      setSaving(false);
    }
  }

  return (
    <section>
      <header className="page-header">
        <p className="eyebrow">{isOwnProfile ? "Your account" : "Public profile"}</p>
        <h1>{isOwnProfile ? "Your profile" : "Profile"}</h1>
        <p>
          {isOwnProfile
            ? "This is how other people see your account. You can change your username, display name, and bio here."
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

          {isOwnProfile ? (
            <form className="settings-form" onSubmit={handleSave}>
              <h3 className="section-title">Edit profile</h3>
              <p className="muted">No password is required to change these fields.</p>
              <label>
                Username
                <input
                  value={form.username}
                  autoComplete="username"
                  onChange={(event) =>
                    setForm((current) => ({ ...current, username: event.target.value }))
                  }
                  required
                />
              </label>
              <label>
                Display name
                <input
                  value={form.displayName}
                  autoComplete="name"
                  onChange={(event) =>
                    setForm((current) => ({ ...current, displayName: event.target.value }))
                  }
                  required
                />
              </label>
              <label>
                Bio
                <textarea
                  value={form.bio}
                  maxLength={160}
                  rows={3}
                  onChange={(event) =>
                    setForm((current) => ({ ...current, bio: event.target.value }))
                  }
                />
              </label>
              {formError ? <p className="error">{formError}</p> : null}
              {formSuccess ? <p className="success">{formSuccess}</p> : null}
              <button type="submit" disabled={saving}>
                {saving ? "Saving..." : "Save profile"}
              </button>
            </form>
          ) : null}

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
