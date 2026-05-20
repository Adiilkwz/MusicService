import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import api from '../../services/api'
import LoadingSpinner from '../../components/LoadingSpinner'

export default function ArtistView() {
  const { id } = useParams()
  const [artist, setArtist] = useState(null)
  const [albums, setAlbums] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    setLoading(true)
    Promise.all([
      api.get(`/catalog/artists/${id}`),
      api.get(`/catalog/artists/${id}/albums`),
    ])
      .then(([artistResp, albumsResp]) => {
        setArtist(artistResp.data)
        setAlbums(albumsResp.data.albums || albumsResp.data)
      })
      .catch((err) => setError(err.response?.data?.error || err.message || 'Unable to load artist'))
      .finally(() => setLoading(false))
  }, [id])

  if (loading) return <LoadingSpinner />
  if (error) return <div className="rounded-3xl border border-red-500/30 bg-red-500/10 p-6 text-red-200">{error}</div>
  if (!artist) return <div className="rounded-3xl border border-white/10 bg-surface p-6">Artist not found.</div>

  return (
    <div className="space-y-8">
      <div className="rounded-3xl border border-white/10 bg-surface p-8 shadow-glow">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-3xl font-semibold text-white">{artist.name}</h1>
            <p className="mt-2 text-sm text-muted">{artist.bio || 'No biography available.'}</p>
          </div>
          <div className="rounded-3xl bg-[#11131b] px-5 py-3 text-sm text-white/70">
            Artist ID {artist.id}
          </div>
        </div>
      </div>

      <section className="rounded-3xl border border-white/10 bg-surface p-8 shadow-glow">
        <div className="flex items-center justify-between gap-4">
          <div>
            <h2 className="text-2xl font-semibold text-white">Albums</h2>
            <p className="mt-2 text-sm text-muted">Browse the artist’s albums and open them to play tracks.</p>
          </div>
          <Link
            to="/search"
            className="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm text-white transition hover:bg-white/10"
          >
            Back to search
          </Link>
        </div>

        <div className="mt-6 grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
          {albums.map((album) => (
            <Link
              key={album.id}
              to={`/album/${album.id}`}
              className="rounded-3xl border border-white/10 bg-surface2 p-5 text-white transition hover:border-accent"
            >
              <div className="text-lg font-semibold">{album.title}</div>
              <p className="mt-2 text-sm text-white/70">{album.release_year ? `Released ${album.release_year}` : 'Release date unknown'}</p>
            </Link>
          ))}
          {albums.length === 0 && <div className="text-muted">No albums found for this artist.</div>}
        </div>
      </section>
    </div>
  )
}
