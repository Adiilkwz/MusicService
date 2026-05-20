export default function LoadingSpinner() {
  return (
    <div className="flex min-h-[250px] items-center justify-center">
      <div className="h-12 w-12 animate-spin rounded-full border-4 border-accent border-t-transparent" />
    </div>
  )
}
