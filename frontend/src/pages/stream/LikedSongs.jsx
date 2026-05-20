import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import api from '../../services/api'
import { useAuth } from '../../contexts/AuthContext'
import SongCard from '../../components/SongCard'
import LoadingSpinner from '../../components/LoadingSpinner'

export default function LikedSongs() {
  const { token, loading: authLoading } = useAuth()
  const [songs, setSongs] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    if (authLoading) return
    if (!token) {
      setLoading(false)
      return
    }
    api.get('/likes')
      .then((response) => setSongs(response.data.songs || response.data))
      .catch((err) => setError(err.response?.data?.error || err.message || 'Unable to load liked songs'))
      .finally(() => setLoading(false))
  }, [token, authLoading])

  if (loading || authLoading) return <LoadingSpinner />

  if (!token) {
    return (
      <div className="rounded-3xl border border-white/10 bg-surface p-8 text-white/80">
        <p>Sign in to see your liked songs.</p>
        <Link to="/auth/login" className="mt-4 inline-block text-accent hover:underline">
          Sign in
        </Link>
      </div>
    )
  }

  return (
    <div className="space-y-8">
      <div className="rounded-3xl border border-white/10 bg-surface p-8 shadow-glow">
        <h1 className="text-3xl font-semibold text-white">Liked songs</h1>
        <p className="mt-2 text-sm text-muted">Your favorite tracks are saved here.</p>
      </div>

      {error && <div className="rounded-3xl border border-red-500/30 bg-red-500/10 p-6 text-red-200">{error}</div>}

      <div className="grid gap-4">
        {songs.length === 0 ? (
          <div className="rounded-3xl border border-white/10 bg-surface p-8 text-white/70">You haven't liked any songs yet.</div>
        ) : (
          songs.map((song) => <SongCard key={song.id} song={song} playlists={[]} />)
        )}
      </div>
    </div>
  )
}
