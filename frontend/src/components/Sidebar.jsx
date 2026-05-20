import { NavLink, useLocation } from 'react-router-dom'
import { useEffect, useState } from 'react'
import api from '../services/api'
import { useAuth } from '../contexts/AuthContext'

const navItems = [
  { label: 'Search', to: '/search', icon: '🔍' },
  { label: 'Favorites', to: '/likes', icon: '❤️' },
  { label: 'Playlists', to: '/playlists', icon: '🎧' },
  { label: 'Profile', to: '/profile', icon: '👤' },
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
    <aside className="w-72 min-h-screen border-r border-surface/60 bg-surface p-6 hidden sm:flex sm:flex-col">
      <div className="mb-10">
        <div className="text-xs uppercase tracking-[0.35em] text-muted">MusicHub</div>
        <h1 className="mt-4 text-3xl font-semibold tracking-tight text-white">Stream Studio</h1>
      </div>

      <nav className="space-y-2">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
              `flex items-center gap-3 rounded-3xl px-4 py-3 text-sm font-medium transition ${
                isActive ? 'bg-accent text-black shadow-lg' : 'text-white/80 hover:bg-white/5'
              }`
            }
          >
            <span className="text-base">{item.icon}</span>
            <span>{item.label}</span>
          </NavLink>
        ))}
      </nav>

      <div className="mt-auto rounded-3xl border border-white/10 bg-[#0f111d] p-4 text-sm text-white/70">
        {user ? `Logged in as ${user.name || user.display_name || 'User'}` : 'Login to unlock playlists and save favorites'}
      </div>
    </aside>
  )
}
