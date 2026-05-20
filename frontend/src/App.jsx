import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './contexts/AuthContext'
import { PlayerProvider } from './contexts/PlayerContext'
import { ToastProvider } from './contexts/ToastContext'
import Layout from './components/Layout'
import Login from './pages/auth/Login'
import Register from './pages/auth/Register'
import Profile from './pages/auth/Profile'
import Search from './pages/catalog/Search'
import ArtistView from './pages/catalog/ArtistView'
import AlbumView from './pages/catalog/AlbumView'
import Playlists from './pages/stream/Playlists'
import PlaylistDetail from './pages/stream/PlaylistDetail'
import LikedSongs from './pages/stream/LikedSongs'
import NotFound from './pages/NotFound'

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <PlayerProvider>
          <ToastProvider>
            <Routes>
              <Route path="/auth/login" element={<Login />} />
              <Route path="/auth/register" element={<Register />} />
              <Route path="/" element={<Layout />}>
                <Route index element={<Navigate to="/search" replace />} />
                <Route path="search" element={<Search />} />
                <Route path="artist/:id" element={<ArtistView />} />
                <Route path="album/:id" element={<AlbumView />} />
                <Route path="playlists" element={<Playlists />} />
                <Route path="playlist/:id" element={<PlaylistDetail />} />
                <Route path="likes" element={<LikedSongs />} />
                <Route path="profile" element={<Profile />} />
                <Route path="*" element={<NotFound />} />
              </Route>
            </Routes>
          </ToastProvider>
        </PlayerProvider>
      </AuthProvider>
    </BrowserRouter>
  )
}

export default App;
