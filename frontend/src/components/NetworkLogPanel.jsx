import { useEffect, useState } from "react";
import { clearNetworkLogs, getNetworkLogs, subscribeNetworkLogs } from "../api";

function formatBody(body) {
  if (body == null || body === "") {
    return "(empty)";
  }
  if (typeof body === "string") {
    return body;
  }
  return JSON.stringify(body, null, 2);
}

export default function NetworkLogPanel() {
  const [entries, setEntries] = useState(getNetworkLogs);
  const [open, setOpen] = useState(true);

  useEffect(() => subscribeNetworkLogs(setEntries), []);

  return (
    <aside className={`network-log ${open ? "open" : "collapsed"}`}>
      <header className="network-log-header">
        <button type="button" className="ghost" onClick={() => setOpen((value) => !value)}>
          {open ? "Hide" : "Show"} network log
        </button>
        <span className="muted">{entries.length} responses</span>
        <button type="button" className="ghost" onClick={clearNetworkLogs} disabled={!entries.length}>
          Clear
        </button>
      </header>

      {open ? (
        <div className="network-log-list">
          {entries.length === 0 ? (
            <p className="muted">Responses from every API call will show up here.</p>
          ) : (
            entries.map((entry) => (
              <article key={entry.id} className={`network-log-entry ${entry.ok ? "ok" : "fail"}`}>
                <div className="network-log-meta">
                  <strong>
                    {entry.method} {entry.path}
                  </strong>
                  <span>{entry.status}</span>
                  <time>{entry.at}</time>
                </div>
                <pre>{formatBody(entry.body)}</pre>
              </article>
            ))
          )}
        </div>
      ) : null}
    </aside>
  );
}
