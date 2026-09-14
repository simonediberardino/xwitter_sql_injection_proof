const TOKEN_KEY = "xwitter_token";

const MAX_LOGS = 80;
let networkLogs = [];
const networkLogListeners = new Set();

export function getToken() {
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token) {
  if (token) {
    localStorage.setItem(TOKEN_KEY, token);
  } else {
    localStorage.removeItem(TOKEN_KEY);
  }
}

export function getNetworkLogs() {
  return networkLogs;
}

export function subscribeNetworkLogs(listener) {
  networkLogListeners.add(listener);
  return () => networkLogListeners.delete(listener);
}

export function clearNetworkLogs() {
  networkLogs = [];
  notifyNetworkLogs();
}

function notifyNetworkLogs() {
  for (const listener of networkLogListeners) {
    listener(networkLogs);
  }
}

function pushNetworkLog(entry) {
  networkLogs = [entry, ...networkLogs].slice(0, MAX_LOGS);
  notifyNetworkLogs();
}

function parseBody(text) {
  if (!text) {
    return null;
  }

  try {
    return JSON.parse(text);
  } catch {
    return text;
  }
}

export async function api(path, options = {}) {
  const headers = { ...(options.headers || {}) };
  const token = getToken();

  if (options.body && !(options.body instanceof FormData)) {
    headers["Content-Type"] = "application/json";
  }
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const method = (options.method || "GET").toUpperCase();
  const response = await fetch(path, { ...options, headers });
  const text = await response.text();
  const data = parseBody(text);

  const logEntry = {
    id: `${Date.now()}-${Math.random().toString(16).slice(2)}`,
    at: new Date().toISOString(),
    method,
    path,
    status: response.status,
    ok: response.ok,
    body: data,
  };

  pushNetworkLog(logEntry);
  console.log(`[network] ${method} ${path} ${response.status}`, data);

  if (!response.ok) {
    const message =
      data && typeof data === "object" && data.error
        ? data.error
        : "Request failed";
    throw new Error(message);
  }

  return data;
}
