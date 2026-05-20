import { usePlayer } from '../contexts/PlayerContext'
import { useEffect, useState } from 'react'

function formatTime(seconds) {
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

export default function BottomPlayer() {
  const { track, isPlaying, togglePlay, progress, seek, volume, setVolume } = usePlayer()
  const [duration, setDuration] = useState(0)

  useEffect(() => {
    if (!track) return
    const audio = document.querySelector('audio')
    const update = () => setDuration(audio.duration || 0)
    audio.addEventListener('loadedmetadata', update)
    return () => audio.removeEventListener('loadedmetadata', update)
  }, [track])

  if (!track) {
    return (
      <div className="border-t border-surface2 bg-surface p-4 text-sm text-muted">
        <div className="mx-auto max-w-6xl">Select a song to play from the search results.</div>
      </div>
    )
  }

  return (
    <div className="border-t border-surface2 bg-surface p-4 shadow-glow">
      <div className="mx-auto flex max-w-6xl flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <div className="text-sm uppercase tracking-[0.3em] text-muted">Now Playing</div>
          <div className="text-lg font-semibold">{track.title}</div>
          <div className="text-sm text-white/70">{track.artistName || 'Unknown artist'}</div>
        </div>

        <div className="flex items-center gap-3">
          <button
            onClick={togglePlay}
            className="rounded-full bg-white px-4 py-2 text-sm font-semibold text-black transition hover:bg-white/90"
          >
            {isPlaying ? 'Pause' : 'Play'}
          </button>
          <div className="min-w-[150px] text-right text-xs text-muted">
            {formatTime(progress)} / {formatTime(duration)}
          </div>
        </div>
      </div>

      <div className="mt-4 grid gap-3 sm:grid-cols-[1fr_auto]">
        <input
          type="range"
          min="0"
          max={duration || 100}
          value={Math.min(progress, duration || 0)}
          onChange={(event) => seek(Number(event.target.value))}
          className="w-full accent-accent"
        />
        <div className="flex items-center gap-2">
          <span className="text-xs text-muted">Volume</span>
          <input
            type="range"
            min="0"
            max="1"
            step="0.01"
            value={volume}
            onChange={(event) => setVolume(Number(event.target.value))}
            className="w-24 accent-accent"
          />
        </div>
      </div>
    </div>
  )
}
