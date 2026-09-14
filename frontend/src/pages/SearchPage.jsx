import { useState } from "react";
import { Link } from "react-router-dom";
import { api } from "../api";

function asRows(data) {
  if (Array.isArray(data)) {
    return data;
  }
  if (data && typeof data === "object") {
    if (Array.isArray(data.users)) {
      return data.users;
    }
    if (Array.isArray(data.results)) {
      return data.results;
    }
    return [data];
  }
  return [];
}

function readField(row, ...keys) {
  if (!row || typeof row !== "object") {
    return "";
  }
  for (const key of keys) {
    if (row[key] != null && row[key] !== "") {
      return String(row[key]);
    }
  }
  return "";
}

export default function SearchPage() {
  const [query, setQuery] = useState("");
  const [submitted, setSubmitted] = useState("");
  const [rows, setRows] = useState([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event) {
    event.preventDefault();
    const next = query.trim();
    setSubmitted(next);
    setError("");

    if (!next) {
      setRows([]);
      return;
    }

    setLoading(true);
    try {
      const data = await api("/api/users/search", {
        method: "POST",
        body: JSON.stringify({ q: next }),
      });
      setRows(asRows(data));
    } catch (err) {
      setRows([]);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <section>
      <header className="page-header">
        <p className="eyebrow">Find people</p>
        <h1>Search</h1>
        <p>Find accounts by username or display name.</p>
      </header>

      <form className="search-form" onSubmit={handleSubmit}>
        <input
          value={query}
          placeholder="Search users"
          onChange={(event) => setQuery(event.target.value)}
        />
        <button type="submit">Search</button>
      </form>

      {error ? <p className="error">{error}</p> : null}
      {loading ? <p className="muted">Searching...</p> : null}

      {!loading && submitted && rows.length === 0 ? (
        <p className="muted">No users matched that search.</p>
      ) : null}

      <ul className="user-list">
        {rows.map((row, index) => {
          const username = readField(row, "username");
          const displayName = readField(row, "display_name", "displayName", "username") || "Result";
          const initial = displayName.slice(0, 1).toUpperCase() || "?";
          const content = (
            <>
              <span className="avatar small">{initial}</span>
              <span className="user-row-copy">
                <strong>{displayName}</strong>
                {username ? <span className="muted">@{username}</span> : null}
                <pre className="result-json">{JSON.stringify(row, null, 2)}</pre>
              </span>
            </>
          );

          return (
            <li key={readField(row, "id") || `${username}-${index}`}>
              {username ? (
                <Link to={`/u/${username}`} className="user-row">
                  {content}
                </Link>
              ) : (
                <div className="user-row">{content}</div>
              )}
            </li>
          );
        })}
      </ul>
    </section>
  );
}
