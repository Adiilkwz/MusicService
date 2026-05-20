import { Link } from 'react-router-dom'

export default function NotFound() {
  return (
    <div className="mx-auto flex min-h-[calc(100vh-80px)] max-w-4xl items-center justify-center px-4 py-12">
      <div className="rounded-3xl border border-white/10 bg-surface p-10 text-center shadow-xl shadow-black/20">
        <h1 className="text-4xl font-semibold text-white">Page not found</h1>
        <p className="mt-4 text-sm text-muted">We couldn't find that page. Head back to the search page.</p>
        <Link
          to="/search"
          className="mt-8 inline-flex rounded-full bg-accent px-6 py-3 text-sm font-semibold text-black transition hover:bg-accent/90"
        >
          Go to search
        </Link>
      </div>
    </div>
  )
}
