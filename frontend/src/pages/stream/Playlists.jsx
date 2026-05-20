import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import api from '../../services/api'
import LoadingSpinner from '../../components/LoadingSpinner'

export default function Playlists() {
  const [playlists, setPlaylists] = useState([])
  const [title, setTitle] = useState('')
  const [loading, setLoading] = useState(true)
  const [creating, setCreating] = useState(false)
  const [error, setError] = useState('')

  const loadPlaylists = async () => {
    setLoading(true)
    try {
      const response = await api.get('/playlists')
      setPlaylists(response.data.playlists || response.data)
    } catch (err) {
      setError(err.response?.data?.error || err.message || 'Unable to load playlists')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { loadPlaylists() }, [])

  const handleCreate = async (e) => {
    e.preventDefault()
    if (!title.trim()) return
    setCreating(true)
    try {
      await api.post('/playlists', { title })
      setTitle('')
      loadPlaylists()
    } catch (err) {
      setError(err.response?.data?.error || err.message || 'Unable to create playlist')
    } finally {
      setCreating(false)
    }
  }

  return (
    <div className="space-y-8">
      <div className="rounded-3xl border border-white/10 bg-surface p-8 shadow-glow">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-3xl font-semibold text-white">Your Playlists</h1>
            <p className="mt-2 text-sm text-muted">Create and manage playlists for every listening mood.</p>
          </div>
          <form onSubmit={handleCreate} className="flex flex-col gap-3 sm:flex-row sm:items-center">
            <input
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="New playlist title"
              className="rounded-3xl border border-white/10 bg-[#11131b] px-4 py-3 text-white outline-none focus:border-accent"
            />
            <button
              type="submit"
              disabled={creating}
              className="rounded-full bg-accent px-5 py-3 text-sm font-semibold text-black transition hover:bg-accent/90 disabled:opacity-60"
            >
              {creating ? 'Creating…' : 'Create'}
            </button>
          </form>
        </div>
      </div>

      {loading ? (
        <LoadingSpinner />
      ) : (
        <div className="grid gap-4">
          {error && <div className="rounded-3xl border border-red-500/30 bg-red-500/10 p-5 text-red-200">{error}</div>}
          {playlists.length === 0 ? (
            <div className="rounded-3xl border border-white/10 bg-surface p-8 text-white/70">No playlists yet. Create one to get started.</div>
          ) : (
            playlists.map((playlist) => (
              <Link
                key={playlist.id}
                to={`/playlist/${playlist.id}`}
                className="rounded-3xl border border-white/10 bg-surface2 p-6 text-white transition hover:border-accent"
              >
                <div className="text-lg font-semibold">{playlist.title}</div>
                <p className="mt-2 text-sm text-muted">{playlist.song_count ? `${playlist.song_count} songs` : 'Tap to open playlist'}</p>
              </Link>
            ))
          )}
        </div>
      )}
    </div>
  )
}
