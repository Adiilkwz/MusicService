import { createContext, useContext, useEffect, useMemo, useRef, useState } from 'react'

const PlayerContext = createContext(null)

export function PlayerProvider({ children }) {
  const [track, setTrack] = useState(null)
  const [isPlaying, setIsPlaying] = useState(false)
  const [volume, setVolume] = useState(0.7)
  const [progress, setProgress] = useState(0)
  const audioRef = useRef(null)

  useEffect(() => {
    const audio = audioRef.current
    if (!audio) return

    const handleTimeUpdate = () => setProgress(audio.currentTime)
    const handleEnded = () => setIsPlaying(false)
    audio.addEventListener('timeupdate', handleTimeUpdate)
    audio.addEventListener('ended', handleEnded)

    return () => {
      audio.removeEventListener('timeupdate', handleTimeUpdate)
      audio.removeEventListener('ended', handleEnded)
    }
  }, [])

  useEffect(() => {
    const audio = audioRef.current
    if (!audio) return
    audio.volume = volume
  }, [volume])

  useEffect(() => {
    const audio = audioRef.current
    if (!audio) return
    if (track?.streamUrl) {
      audio.src = track.streamUrl
      audio.load()
    }
    if (isPlaying && track) {
      audio.play().catch(() => setIsPlaying(false))
    }
  }, [track, isPlaying])

  const play = (song) => {
    setTrack(song)
    setIsPlaying(true)
  }

  const pause = () => setIsPlaying(false)

  const togglePlay = () => setIsPlaying((prev) => !prev)

  const seek = (seconds) => {
    const audio = audioRef.current
    if (!audio) return
    audio.currentTime = seconds
    setProgress(seconds)
  }

  const state = useMemo(
    () => ({ track, isPlaying, volume, progress, audioRef, play, pause, togglePlay, seek, setVolume }),
    [track, isPlaying, volume, progress]
  )

  return <PlayerContext.Provider value={state}>{children}<audio ref={audioRef} hidden /></PlayerContext.Provider>
}

export function usePlayer() {
  return useContext(PlayerContext)
}
