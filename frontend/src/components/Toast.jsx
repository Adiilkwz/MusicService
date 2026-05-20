import { useToast } from '../contexts/ToastContext'

export default function Toast() {
  const { toasts } = useToast()

  return (
    <div className="pointer-events-none fixed inset-x-0 top-4 z-50 flex flex-col items-center gap-3 px-4">
      {toasts.map((toast) => (
        <div key={toast.id} className="pointer-events-auto rounded-3xl border border-white/10 bg-surface px-5 py-3 text-sm text-white shadow-xl shadow-black/20">
          {toast.message}
        </div>
      ))}
    </div>
  )
}
