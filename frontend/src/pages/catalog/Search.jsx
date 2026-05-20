import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import api from '../../services/api'
import SongCard from '../../components/SongCard'
import LoadingSpinner from '../../components/LoadingSpinner'

export default function Search() {
  const navigate = useNavigate()
  const [query, setQuery] = useState('')
  const [results, setResults] = useState({ artists: [], albums: [], songs: [] })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!query) {
      setResults({ artists: [], albums: [], songs: [] })
      return
    }
    const timeout = setTimeout(() => {
      setLoading(true)
      api.get('/catalog/search', { params: { q: query, limit: 30 } })
        .then((response) => setResults(response.data))
        .catch((err) => setError(err.response?.data?.error || err.message || 'Search failed'))
        .finally(() => setLoading(false))
    }, 250)

    return () => clearTimeout(timeout)
  }, [query])

  return (
    <div className="space-y-8">
      <div className="rounded-3xl border border-white/10 bg-surface p-8 shadow-glow">
        <h1 className="text-3xl font-semibold text-white">Search the catalog</h1>
        <p className="mt-2 text-sm text-muted">Find artists, albums, and songs across the music library.</p>

        <div className="mt-8">
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search for artist, album, or song"
            className="w-full rounded-3xl border border-white/10 bg-[#11131b] px-5 py-4 text-white outline-none focus:border-accent"
          />
        </div>
      </div>

      {loading ? (
        <LoadingSpinner />
      ) : (
        <div className="grid gap-8">
          {error && <div className="rounded-3xl border border-red-500/30 bg-red-500/10 p-5 text-red-200">{error}</div>}

          {results.artists?.length > 0 && (
            <section className="rounded-3xl border border-white/10 bg-surface p-6">
              <h2 className="text-xl font-semibold text-white">Artists</h2>
              <div className="mt-4 grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
                {results.artists.map((artist) => (
                  <button
                    key={artist.id}
                    onClick={() => navigate(`/artist/${artist.id}`)}
                    className="rounded-3xl border border-white/10 bg-surface2 p-4 text-left text-white transition hover:border-accent"
                  >
                    <div className="text-lg font-semibold">{artist.name}</div>
                    <p className="mt-2 text-sm text-white/70">{artist.bio || 'No artist bio available.'}</p>
                  </button>
                ))}
              </div>
            </section>
          )}

          {results.albums?.length > 0 && (
            <section className="rounded-3xl border border-white/10 bg-surface p-6">
              <h2 className="text-xl font-semibold text-white">Albums</h2>
              <div className="mt-4 grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
                {results.albums.map((album) => (
                  <button
                    key={album.id}
                    onClick={() => navigate(`/album/${album.id}`)}
                    className="rounded-3xl border border-white/10 bg-surface2 p-4 text-left text-white transition hover:border-accent"
                  >
                    <div className="text-lg font-semibold">{album.title}</div>
                    <p className="mt-2 text-sm text-white/70">{album.artist_name || album.artistName || 'Unknown artist'}</p>
                  </button>
                ))}
              </div>
            </section>
          )}

          {results.songs?.length > 0 && (
            <section className="rounded-3xl border border-white/10 bg-surface p-6">
              <div className="flex items-center justify-between gap-4">
                <h2 className="text-xl font-semibold text-white">Songs</h2>
              </div>
              <div className="mt-5 grid gap-4">
                {results.songs.map((song) => (
                  <SongCard key={song.id} song={song} playlists={[]} />
                ))}
              </div>
            </section>
          )}

          {!results.artists?.length && !results.albums?.length && !results.songs?.length && query && (
            <div className="rounded-3xl border border-white/10 bg-surface p-6 text-white/70">
              No results found for “{query}”. Try another query.
            </div>
          )}
        </div>
      )}
    </div>
  )
}
