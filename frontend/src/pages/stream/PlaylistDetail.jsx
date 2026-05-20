import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import api from '../../services/api'
import SongCard from '../../components/SongCard'
import LoadingSpinner from '../../components/LoadingSpinner'
import { useToast } from '../../contexts/ToastContext'

export default function PlaylistDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { addToast } = useToast()
  const [playlist, setPlaylist] = useState(null)
  const [songs, setSongs] = useState([])
  const [playlists, setPlaylists] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const loadPlaylist = async () => {
    setLoading(true)
    try {
      const [resp, playlistResp] = await Promise.all([
        api.get('/playlists'),
        api.get(`/playlists/${id}`),
      ])
      setPlaylists(resp.data.playlists || resp.data)
      setPlaylist(playlistResp.data)
      setSongs(playlistResp.data.songs || playlistResp.data.song_ids || [])
    } catch (err) {
      setError(err.response?.data?.error || err.message || 'Cannot load playlist')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadPlaylist()
  }, [id])

  const handleDeleteSong = async (songId) => {
    try {
      await api.delete(`/playlists/${id}/songs/${songId}`)
      addToast('Song removed from playlist', 'success')
      loadPlaylist()
    } catch (err) {
      // handled globally
    }
  }

  if (loading) return <LoadingSpinner />

  if (error) {
    return (
      <div className="rounded-3xl border border-red-500/30 bg-red-500/10 p-6 text-red-200">{error}</div>
    )
  }

  if (!playlist) {
    return (
      <div className="rounded-3xl border border-white/10 bg-surface p-6 text-white/70">Playlist not found.</div>
    )
  }

  return (
    <div className="space-y-8">
      <div className="rounded-3xl border border-white/10 bg-surface p-8 shadow-glow">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-3xl font-semibold text-white">{playlist.title}</h1>
            <p className="mt-2 text-sm text-muted">{playlist.songs?.length ? `${playlist.songs.length} songs` : 'A quick listening list'}</p>
          </div>
          <button
            onClick={() => navigate('/playlists')}
            className="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm text-white transition hover:bg-white/10"
          >
            Back to playlists
          </button>
        </div>
      </div>

      <section className="rounded-3xl border border-white/10 bg-surface p-8 shadow-glow">
        <h2 className="text-2xl font-semibold text-white">Tracks</h2>
        <div className="mt-5 grid gap-4">
          {songs.length > 0 ? (
            songs.map((song) => (
              <div key={song.id || song.song_id} className="space-y-3 rounded-3xl border border-white/10 bg-surface2 p-5">
                <SongCard song={{ ...song, id: song.id || song.song_id }} playlists={playlists} />
                <button
                  onClick={() => handleDeleteSong(song.id || song.song_id)}
                  className="rounded-full bg-red-500 px-4 py-2 text-sm font-semibold text-white transition hover:bg-red-400"
                >
                  Remove from playlist
                </button>
              </div>
            ))
          ) : (
            <div className="text-white/70">No songs in this playlist yet.</div>
          )}
        </div>
      </section>
    </div>
  )
}
