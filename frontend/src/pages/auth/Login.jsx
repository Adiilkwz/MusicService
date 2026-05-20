import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import api from '../../services/api'
import { useAuth } from '../../contexts/AuthContext'
import { useToast } from '../../contexts/ToastContext'

export default function Login() {
  const navigate = useNavigate()
  const { login, token } = useAuth()
  const { addToast } = useToast()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (token) {
      navigate('/search', { replace: true })
    }
  }, [token, navigate])

  const handleLogin = async (e) => {
    e.preventDefault()
    setLoading(true)

    try {
      const response = await api.post('/auth/login', { email, password })
      login(response.data.token, response.data.user)
      addToast('Welcome back!', 'success')
      navigate('/search')
    } catch (error) {
      // toast shown by API interceptor
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="mx-auto flex min-h-[calc(100vh-80px)] max-w-3xl items-center justify-center px-4 py-12">
      <div className="w-full rounded-3xl border border-white/10 bg-surface p-10 shadow-xl shadow-black/20">
        <h1 className="text-3xl font-semibold text-white">Sign in</h1>
        <p className="mt-3 text-sm text-muted">Use your account to continue listening.</p>
        <form onSubmit={handleLogin} className="mt-8 space-y-5">
          <label className="block text-sm text-white/80">
            Email
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className="mt-2 w-full rounded-3xl border border-white/10 bg-[#11131b] px-4 py-3 text-white outline-none transition focus:border-accent"
              placeholder="you@example.com"
            />
          </label>
          <label className="block text-sm text-white/80">
            Password
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className="mt-2 w-full rounded-3xl border border-white/10 bg-[#11131b] px-4 py-3 text-white outline-none transition focus:border-accent"
              placeholder="••••••••"
            />
          </label>
          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-full bg-accent px-5 py-3 text-sm font-semibold text-black transition hover:bg-accent/90 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {loading ? 'Signing in...' : 'Sign in'}
          </button>
        </form>
        <p className="mt-6 text-center text-sm text-white/70">
          Don’t have an account?{' '}
          <Link to="/auth/register" className="text-accent hover:underline">
            Create one
          </Link>
        </p>
      </div>
    </div>
  )
}
