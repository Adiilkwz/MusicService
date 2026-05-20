import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import api from '../../services/api'
import { useAuth } from '../../contexts/AuthContext'
import { useToast } from '../../contexts/ToastContext'

export default function Register() {
  const navigate = useNavigate()
  const { token } = useAuth()
  const { addToast } = useToast()
  const [displayName, setDisplayName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)

  const handleRegister = async (e) => {
    e.preventDefault()
    setLoading(true)

    try {
      await api.post('/auth/register', { display_name: displayName, email, password })
      addToast('Account created. Log in to continue.', 'success')
      navigate('/auth/login')
    } catch (error) {
      // toast shown by interceptor
    } finally {
      setLoading(false)
    }
  }

  if (token) {
    navigate('/search')
    return null
  }

  return (
    <div className="mx-auto flex min-h-[calc(100vh-80px)] max-w-3xl items-center justify-center px-4 py-12">
      <div className="w-full rounded-3xl border border-white/10 bg-surface p-10 shadow-xl shadow-black/20">
        <h1 className="text-3xl font-semibold text-white">Create your account</h1>
        <p className="mt-3 text-sm text-muted">Register and start building playlists.</p>
        <form onSubmit={handleRegister} className="mt-8 space-y-5">
          <label className="block text-sm text-white/80">
            Display Name
            <input
              value={displayName}
              onChange={(e) => setDisplayName(e.target.value)}
              required
              className="mt-2 w-full rounded-3xl border border-white/10 bg-[#11131b] px-4 py-3 text-white outline-none transition focus:border-accent"
              placeholder="Your name"
            />
          </label>
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
              placeholder="Minimum 6 characters"
            />
          </label>
          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-full bg-accent px-5 py-3 text-sm font-semibold text-black transition hover:bg-accent/90 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {loading ? 'Creating account...' : 'Create account'}
          </button>
        </form>
        <p className="mt-6 text-center text-sm text-white/70">
          Already have an account?{' '}
          <Link to="/auth/login" className="text-accent hover:underline">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  )
}
