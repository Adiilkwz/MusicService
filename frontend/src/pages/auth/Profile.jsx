import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import api from '../../services/api'
import { useAuth } from '../../contexts/AuthContext'
import { useToast } from '../../contexts/ToastContext'
import LoadingSpinner from '../../components/LoadingSpinner'

function GuestAuthPrompt() {
  return (
    <div className="mx-auto flex min-h-[calc(100vh-80px)] max-w-3xl items-center justify-center px-4 py-12">
      <div className="w-full rounded-3xl border border-white/10 bg-surface p-10 text-center shadow-xl shadow-black/20">
        <h1 className="text-3xl font-semibold text-white">Профиль</h1>
        <p className="mt-3 text-sm text-muted">Войдите в аккаунт, чтобы просматривать и редактировать профиль.</p>
        <div className="mt-8 flex flex-col gap-3 sm:flex-row sm:justify-center">
          <Link
            to="/auth/login"
            className="rounded-full bg-accent px-8 py-3 text-sm font-semibold text-black transition hover:bg-accent/90"
          >
            Войти
          </Link>
          <Link
            to="/auth/register"
            className="rounded-full border border-white/10 px-8 py-3 text-sm font-semibold text-white transition hover:border-accent"
          >
            Регистрация
          </Link>
        </div>
      </div>
    </div>
  )
}

export default function Profile() {
  const { user, token, logout, setUser, loading: authLoading } = useAuth()
  const { addToast } = useToast()
  const [profile, setProfile] = useState(user)
  const [isEditing, setIsEditing] = useState(false)
  const [displayName, setDisplayName] = useState(user?.display_name || user?.name || '')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (authLoading) return

    if (!token) {
      setProfile(null)
      setLoading(false)
      return
    }

    if (user) {
      setProfile(user)
      setDisplayName(user.display_name || user.name || '')
      setLoading(false)
      return
    }

    if (!profile) {
      setLoading(true)
      api
        .get('/profile/me')
        .then((response) => {
          setProfile(response.data)
          setDisplayName(response.data.display_name || response.data.name || '')
          setUser(response.data)
        })
        .catch(() => addToast('Не удалось загрузить профиль', 'error'))
        .finally(() => setLoading(false))
    } else {
      setLoading(false)
    }
  }, [profile, token, user, authLoading, addToast, setUser])

  const handleSave = async () => {
    try {
      await api.put('/profile', { display_name: displayName })
      setProfile((prev) => ({ ...prev, display_name: displayName }))
      setUser((prev) => ({ ...prev, display_name: displayName }))
      setIsEditing(false)
      addToast('Profile saved', 'success')
    } catch {
      // handled globally
    }
  }

  const handleDelete = async () => {
    if (!window.confirm('Delete account? This cannot be undone.')) return

    try {
      await api.delete('/profile')
      logout()
      addToast('Account deleted', 'success')
    } catch {
      // handled globally
    }
  }

  const handleLogout = () => {
    logout()
    setProfile(null)
    setIsEditing(false)
  }

  if (authLoading || (token && loading)) {
    return (
      <div className="mx-auto flex min-h-[calc(100vh-80px)] max-w-4xl items-center justify-center px-4 py-12">
        <LoadingSpinner />
      </div>
    )
  }

  if (!token) {
    return <GuestAuthPrompt />
  }

  if (!profile) {
    return (
      <div className="mx-auto flex min-h-[calc(100vh-80px)] max-w-3xl items-center justify-center px-4 py-12">
        <div className="w-full rounded-3xl border border-white/10 bg-surface p-10 text-center shadow-xl shadow-black/20">
          <p className="text-sm text-white/80">Сессия недействительна или профиль недоступен.</p>
          <div className="mt-8 flex flex-col gap-3 sm:flex-row sm:justify-center">
            <Link
              to="/auth/login"
              className="rounded-full bg-accent px-8 py-3 text-sm font-semibold text-black transition hover:bg-accent/90"
            >
              Войти
            </Link>
            <Link
              to="/auth/register"
              className="rounded-full border border-white/10 px-8 py-3 text-sm font-semibold text-white transition hover:border-accent"
            >
              Регистрация
            </Link>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-4xl space-y-8 px-4 py-6 sm:px-8">
      <div className="rounded-3xl border border-white/10 bg-surface p-8 shadow-xl shadow-black/10">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-3xl font-semibold text-white">Ваш профиль</h1>
            <p className="mt-2 text-sm text-muted">Управление настройками аккаунта.</p>
          </div>
          <div className="flex flex-wrap gap-3">
            <button
              type="button"
              onClick={handleLogout}
              className="rounded-full border border-white/10 px-5 py-3 text-sm font-semibold text-white transition hover:border-accent"
            >
              Выйти
            </button>
            <button
              type="button"
              onClick={handleDelete}
              className="rounded-full bg-red-500 px-5 py-3 text-sm font-semibold text-white transition hover:bg-red-400"
            >
              Удалить аккаунт
            </button>
          </div>
        </div>

        <div className="mt-8 grid gap-6 sm:grid-cols-[1fr_1fr]">
          <div className="space-y-4">
            <div>
              <div className="text-sm uppercase tracking-[0.35em] text-muted">Имя</div>
              {isEditing ? (
                <input
                  value={displayName}
                  onChange={(e) => setDisplayName(e.target.value)}
                  className="mt-2 w-full rounded-3xl border border-white/10 bg-[#11131b] px-4 py-3 text-white outline-none transition focus:border-accent"
                />
              ) : (
                <p className="mt-2 text-lg text-white">{profile?.display_name || profile?.name || 'Имя не указано'}</p>
              )}
            </div>

            <div>
              <div className="text-sm uppercase tracking-[0.35em] text-muted">Email</div>
              <p className="mt-2 text-lg text-white">{profile?.email}</p>
            </div>

            <div>
              <div className="text-sm uppercase tracking-[0.35em] text-muted">Дата регистрации</div>
              <p className="mt-2 text-lg text-white">
                {new Date(profile?.created_at || profile?.createdAt || Date.now()).toLocaleDateString()}
              </p>
            </div>
          </div>

          <div className="space-y-4 rounded-3xl border border-white/10 bg-surface2 p-6">
            <div className="flex items-center justify-between gap-4">
              <div>
                <p className="text-sm uppercase tracking-[0.35em] text-muted">ID аккаунта</p>
                <p className="mt-2 text-base text-white/80">{profile?.id || profile?.user_id || profile?.userId}</p>
              </div>
              <button
                type="button"
                onClick={() => setIsEditing((prev) => !prev)}
                className="rounded-full border border-white/10 px-4 py-2 text-sm text-white transition hover:border-accent"
              >
                {isEditing ? 'Отмена' : 'Изменить'}
              </button>
            </div>

            {isEditing && (
              <button
                type="button"
                onClick={handleSave}
                className="w-full rounded-full bg-accent px-5 py-3 text-sm font-semibold text-black transition hover:bg-accent/90"
              >
                Сохранить
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
