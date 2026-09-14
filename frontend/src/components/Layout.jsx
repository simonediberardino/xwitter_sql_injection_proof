import { NavLink, Outlet } from "react-router-dom";
import { useAuth } from "../auth";

function IconHome() {
  return (
    <svg className="nav-icon" viewBox="0 0 24 24" aria-hidden="true">
      <path d="M4 11.5 12 4l8 7.5V20a1 1 0 0 1-1 1h-5v-6H10v6H5a1 1 0 0 1-1-1z" />
    </svg>
  );
}

function IconSearch() {
  return (
    <svg className="nav-icon" viewBox="0 0 24 24" aria-hidden="true">
      <circle cx="10.5" cy="10.5" r="6.5" />
      <path d="m15.5 15.5 5 5" />
    </svg>
  );
}

function IconProfile() {
  return (
    <svg className="nav-icon" viewBox="0 0 24 24" aria-hidden="true">
      <circle cx="12" cy="8" r="3.5" />
      <path d="M5 19.5c1.4-3.2 4-5 7-5s5.6 1.8 7 5" />
    </svg>
  );
}

function IconLogout() {
  return (
    <svg className="nav-icon" viewBox="0 0 24 24" aria-hidden="true">
      <path d="M10 5H6v14h4" />
      <path d="M10 12h9" />
      <path d="m16 8 4 4-4 4" />
    </svg>
  );
}

export default function Layout() {
  const { user, logout } = useAuth();

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <span className="brand-mark">X</span>
          <span>
            Xwitter
            <span className="brand-sub">Public feed</span>
          </span>
        </div>
        <nav className="nav">
          <NavLink to="/" end>
            <IconHome />
            Home
          </NavLink>
          <NavLink to="/search">
            <IconSearch />
            Search
          </NavLink>
          <NavLink to="/profile">
            <IconProfile />
            Profile
          </NavLink>
        </nav>
        <div className="sidebar-footer">
          <div className="sidebar-account">
            <IconProfile />
            <div>
              <strong>{user.displayName}</strong>
              <span>@{user.username}</span>
            </div>
          </div>
          <button type="button" className="sidebar-logout" onClick={logout}>
            <IconLogout />
            Log out
          </button>
        </div>
      </aside>
      <main className="main">
        <Outlet />
      </main>
    </div>
  );
}
