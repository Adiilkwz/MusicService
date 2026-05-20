import { useAuth } from '../contexts/AuthContext'
import { usePlayer } from '../contexts/PlayerContext'
import api from '../services/api'
import { useToast } from '../contexts/ToastContext'
import { useState } from 'react'

export default function SongCard({ song, onLikeUpdate, playlists, onAddToPlaylist }) {
  const { isPlaying, play } = usePlayer()
  const { user } = useAuth()
  const { addToast } = useToast()
  const [liked, setLiked] = useState(song.liked || false)
  const [saving, setSaving] = useState(false)

  const handlePlay = async () => {
    play({
      ...song,
      streamUrl: song.streamUrl || `https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3`,
    })
    try {
      await api.post('/stream/recordplay', { song_id: song.id })
    } catch (error) {
      addToast('Unable to register play event', 'error')
    }
  }

  const toggleLike = async () => {
    if (!user) {
      addToast('Login to like songs', 'error')
      return
    }

    setSaving(true)
    try {
      if (liked) {
        await api.delete(`/likes/${song.id}`)
      } else {
        await api.post('/likes', { song_id: song.id })
      }
      setLiked(!liked)
      onLikeUpdate?.()
    } catch (error) {
      addToast('Failed to update like', 'error')
    } finally {
      setSaving(false)
    }
  }

  const handleAddToPlaylist = async (playlistId) => {
    if (!playlistId) return
    if (!user) {
      addToast('Login to manage playlists', 'error')
      return
    }
    try {
      await api.post(`/playlists/${playlistId}/songs`, { song_id: song.id })
      addToast('Added to playlist', 'success')
      onAddToPlaylist?.()
    } catch (error) {
      addToast('Failed to add to playlist', 'error')
    }
  }

  return (
    <div className="rounded-3xl border border-white/10 bg-surface2 p-4 shadow-lg shadow-black/10">
      <div className="flex items-start justify-between gap-3">
        <div>
          <div className="text-base font-semibold text-white">{song.title}</div>
          <div className="mt-1 text-sm text-muted">{song.genre || 'Unknown genre'}</div>
          <div className="mt-2 text-xs text-white/60">{Math.floor((song.duration || 0) / 60)}:{((song.duration || 0) % 60).toString().padStart(2, '0')} min</div>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={handlePlay}
            className="rounded-full bg-accent px-4 py-2 text-sm font-semibold text-black transition hover:bg-accent/90"
          >
            Play
          </button>
          <button
            onClick={toggleLike}
            disabled={saving}
            className={`rounded-full px-3 py-2 text-sm transition ${liked ? 'bg-red-500 text-white' : 'bg-white/5 text-white/80 hover:bg-white/10'}`}
          >
            {liked ? '♥' : '♡'}
          </button>
        </div>
      </div>
      <div className="mt-3 flex items-center justify-between text-sm text-white/60">
        <span>Plays: {song.plays_count ?? 0}</span>
        {playlists?.length > 0 && (
          <select onChange={(e) => handleAddToPlaylist(e.target.value)} className="rounded-2xl border border-white/10 bg-surface px-3 py-2 text-sm text-white outline-none">
            <option value="">Add to playlist</option>
            {playlists.map((playlist) => (
              <option key={playlist.id} value={playlist.id}>{playlist.title}</option>
            ))}
          </select>
        )}
      </div>
    </div>
  )
}
