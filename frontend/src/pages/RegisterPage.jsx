import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../auth";

export default function RegisterPage() {
  const { register } = useAuth();
  const navigate = useNavigate();
  const [form, setForm] = useState({
    username: "",
    email: "",
    displayName: "",
    bio: "",
    password: "",
  });
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  function update(field, value) {
    setForm((current) => ({ ...current, [field]: value }));
  }

  async function handleSubmit(event) {
    event.preventDefault();
    setError("");
    setSubmitting(true);
    try {
      await register(form);
      navigate("/", { replace: true });
    } catch (err) {
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="auth-page">
      <form className="auth-card" onSubmit={handleSubmit}>
        <div className="brand">
          <span className="brand-mark">X</span>
          <span>
            Xwitter
            <span className="brand-sub">Public feed</span>
          </span>
        </div>
        <h1>Create account</h1>
        <p className="muted">Join Xwitter and start posting short updates.</p>

        <label>
          Username
          <input
            value={form.username}
            autoComplete="username"
            onChange={(event) => update("username", event.target.value)}
            required
          />
        </label>
        <label>
          Email
          <input
            type="email"
            value={form.email}
            autoComplete="email"
            onChange={(event) => update("email", event.target.value)}
            required
          />
        </label>
        <label>
          Display name
          <input
            value={form.displayName}
            autoComplete="name"
            onChange={(event) => update("displayName", event.target.value)}
            required
          />
        </label>
        <label>
          Bio (optional)
          <textarea
            value={form.bio}
            maxLength={160}
            rows={3}
            onChange={(event) => update("bio", event.target.value)}
          />
        </label>
        <label>
          Password
          <input
            type="password"
            value={form.password}
            autoComplete="new-password"
            onChange={(event) => update("password", event.target.value)}
            required
            minLength={8}
          />
        </label>

        {error ? <p className="error">{error}</p> : null}

        <button type="submit" disabled={submitting}>
          {submitting ? "Creating account..." : "Sign up"}
        </button>
        <p className="muted">
          Already have an account? <Link to="/login">Sign in</Link>
        </p>
      </form>
    </div>
  );
}
