import { createContext, useContext, useEffect, useState } from 'react'
import api from '../services/api'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [token, setToken] = useState(() => {
    const saved = localStorage.getItem('jwt_token')
    if (saved) {
      api.setToken(saved)
    }
    return saved
  })
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const saved = localStorage.getItem('jwt_token')
    if (saved) {
      setToken(saved)
      api.setToken(saved)
      api.get('/profile/me')
        .then((response) => setUser(response.data))
        .catch(() => handleLogout())
        .finally(() => setLoading(false))
    } else {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    if (token) {
      localStorage.setItem('jwt_token', token)
      api.setToken(token)
    } else {
      localStorage.removeItem('jwt_token')
      api.clearToken()
    }
  }, [token])

  useEffect(() => {
    api.subscribeUnauthorized(() => {
      handleLogout()
    })
  }, [])

  const handleLogin = (jwt, userData) => {
    setToken(jwt)
    setUser(userData)
    api.setToken(jwt)
  }

  const handleLogout = () => {
    setToken(null)
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ token, user, loading, login: handleLogin, logout: handleLogout, setUser }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  return useContext(AuthContext)
}
