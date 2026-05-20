import { useEffect, useState } from 'react'
import api from '../../services/api'
import { useAuth } from '../../contexts/AuthContext'
import { useToast } from '../../contexts/ToastContext'

export default function Profile() {
  const { user, logout, setUser } = useAuth()
  const { addToast } = useToast()
  const [profile, setProfile] = useState(user)
  const [isEditing, setIsEditing] = useState(false)
  const [displayName, setDisplayName] = useState(user?.display_name || user?.name || '')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!profile) {
      api.get('/profile/me')
        .then((response) => {
          setProfile(response.data)
          setDisplayName(response.data.display_name || response.data.name || '')
          setUser(response.data)
        })
        .catch(() => addToast('Unable to load profile', 'error'))
        .finally(() => setLoading(false))
    } else {
      setLoading(false)
    }
  }, [profile, addToast, setUser])

  const handleSave = async () => {
    try {
      await api.put('/profile', { display_name: displayName })
      setProfile((prev) => ({ ...prev, display_name: displayName }))
      setUser((prev) => ({ ...prev, display_name: displayName }))
      setIsEditing(false)
      addToast('Profile saved', 'success')
    } catch (error) {
      // handled globally
    }
  }

  const handleDelete = async () => {
    if (!window.confirm('Delete account? This cannot be undone.')) return

    try {
      await api.delete('/profile')
      logout()
      addToast('Account deleted', 'success')
    } catch (error) {
      // handled globally
    }
  }

  if (loading) {
    return (
      <div className="mx-auto flex min-h-[calc(100vh-80px)] max-w-4xl items-center justify-center px-4 py-12">
        <div className="rounded-3xl border border-white/10 bg-surface p-8">Loading profile…</div>
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-4xl space-y-8 px-4 py-6 sm:px-8">
      <div className="rounded-3xl border border-white/10 bg-surface p-8 shadow-xl shadow-black/10">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-3xl font-semibold text-white">Your profile</h1>
            <p className="mt-2 text-sm text-muted">Manage your personal account settings.</p>
          </div>
          <button
            onClick={handleDelete}
            className="rounded-full bg-red-500 px-5 py-3 text-sm font-semibold text-white transition hover:bg-red-400"
          >
            Delete account
          </button>
        </div>

        <div className="mt-8 grid gap-6 sm:grid-cols-[1fr_1fr]">
          <div className="space-y-4">
            <div>
              <div className="text-sm uppercase tracking-[0.35em] text-muted">Name</div>
              {isEditing ? (
                <input
                  value={displayName}
                  onChange={(e) => setDisplayName(e.target.value)}
                  className="mt-2 w-full rounded-3xl border border-white/10 bg-[#11131b] px-4 py-3 text-white outline-none"
                />
              ) : (
                <p className="mt-2 text-lg text-white">{profile.display_name || profile.name || 'No name set'}</p>
              )}
            </div>

            <div>
              <div className="text-sm uppercase tracking-[0.35em] text-muted">Email</div>
              <p className="mt-2 text-lg text-white">{profile.email}</p>
            </div>

            <div>
              <div className="text-sm uppercase tracking-[0.35em] text-muted">Member since</div>
              <p className="mt-2 text-lg text-white">{new Date(profile.created_at || profile.createdAt || Date.now()).toLocaleDateString()}</p>
            </div>
          </div>

          <div className="space-y-4 rounded-3xl border border-white/10 bg-surface2 p-6">
            <div className="flex items-center justify-between gap-4">
              <div>
                <p className="text-sm uppercase tracking-[0.35em] text-muted">Account ID</p>
                <p className="mt-2 text-base text-white/80">{profile.id || profile.user_id || profile.userId}</p>
              </div>
              <button
                onClick={() => setIsEditing((prev) => !prev)}
                className="rounded-full border border-white/10 px-4 py-2 text-sm text-white transition hover:border-accent"
              >
                {isEditing ? 'Cancel' : 'Edit'}
              </button>
            </div>

            {isEditing && (
              <button
                onClick={handleSave}
                className="w-full rounded-full bg-accent px-5 py-3 text-sm font-semibold text-black transition hover:bg-accent/90"
              >
                Save changes
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
