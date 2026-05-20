import { NavLink, useLocation } from 'react-router-dom'
import { useEffect, useState } from 'react'
import api from '../services/api'
import { useAuth } from '../contexts/AuthContext'

const navItems = [
  { label: 'Search', to: '/search' },
  { label: 'Favorites', to: '/likes' },
  { label: 'Playlists', to: '/playlists' },
  { label: 'Profile', to: '/profile' },
]

export default function Sidebar() {
  const { user } = useAuth()
  const location = useLocation()
  const [playlists, setPlaylists] = useState([])

  useEffect(() => {
    api.get('/playlists')
      .then((response) => setPlaylists(response.data || []))
      .catch(() => setPlaylists([]))
  }, [location.pathname])

  return (
    <aside className="w-72 min-h-screen border-r border-surface/60 bg-surface p-5 hidden sm:block">
      <div className="mb-8">
        <div className="text-sm uppercase tracking-[0.3em] text-muted">MusicHub</div>
        <h1 className="mt-4 text-2xl font-semibold">Stream Studio</h1>
      </div>

      <nav className="space-y-1">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
              `block rounded-2xl px-4 py-3 text-sm font-medium transition ${
                isActive ? 'bg-accent text-black' : 'text-white/80 hover:bg-white/5'
              }`
            }
          >
            {item.label}
          </NavLink>
        ))}
      </nav>

      <div className="mt-8 pb-10 border-b border-surface2" />

      <div>
        <div className="text-xs uppercase tracking-[0.3em] text-muted mb-3">Your playlists</div>
        <div className="space-y-2">
          {playlists.length === 0 ? (
            <div className="text-sm text-muted">No playlists yet</div>
          ) : (
            playlists.map((playlist) => (
              <NavLink
                key={playlist.id}
                to={`/playlist/${playlist.id}`}
                className="block rounded-2xl px-4 py-3 text-sm text-white/80 hover:bg-white/5"
              >
                {playlist.title}
              </NavLink>
            ))
          )}
        </div>
      </div>

      <div className="mt-auto pt-8 text-sm text-muted">
        {user ? `Logged as ${user.name}` : 'Login to unlock playlists'}
      </div>
    </aside>
  )
}
