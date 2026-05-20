import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import api from '../../services/api'
import SongCard from '../../components/SongCard'
import LoadingSpinner from '../../components/LoadingSpinner'

export default function AlbumView() {
  const { id } = useParams()
  const [album, setAlbum] = useState(null)
  const [songs, setSongs] = useState([])
  const [playlists, setPlaylists] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    setLoading(true)
    Promise.all([
      api.get(`/catalog/albums/${id}`),
      api.get(`/catalog/albums/${id}/songs`),
      api.get('/playlists'),
    ])
      .then(([albumResp, songsResp, playlistsResp]) => {
        setAlbum(albumResp.data)
        setSongs(songsResp.data.songs || songsResp.data)
        setPlaylists(playlistsResp.data.playlists || playlistsResp.data)
      })
      .catch((err) => setError(err.response?.data?.error || err.message || 'Unable to load album'))
      .finally(() => setLoading(false))
  }, [id])

  if (loading) return <LoadingSpinner />
  if (error) return <div className="rounded-3xl border border-red-500/30 bg-red-500/10 p-6 text-red-200">{error}</div>
  if (!album) return <div className="rounded-3xl border border-white/10 bg-surface p-6">Album not found.</div>

  return (
    <div className="space-y-8">
      <div className="rounded-3xl border border-white/10 bg-surface p-8 shadow-glow">
        <div className="space-y-4">
          <div className="text-sm uppercase tracking-[0.35em] text-muted">Album</div>
          <h1 className="text-3xl font-semibold text-white">{album.title}</h1>
          <div className="text-sm text-white/70">{album.artist_name || album.artistName || 'Unknown artist'}</div>
          <div className="text-sm text-muted">Released {album.release_year || 'unknown'}</div>
        </div>
      </div>

      <section className="rounded-3xl border border-white/10 bg-surface p-8 shadow-glow">
        <h2 className="text-2xl font-semibold text-white">Tracks</h2>
        <div className="mt-5 grid gap-4">
          {songs.length > 0 ? (
            songs.map((song) => <SongCard key={song.id} song={song} playlists={playlists} />)
          ) : (
            <div className="text-white/70">No tracks found for this album.</div>
          )}
        </div>
      </section>
    </div>
  )
}
