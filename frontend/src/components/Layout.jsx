import { Outlet } from 'react-router-dom'
import Sidebar from './Sidebar'
import BottomPlayer from './BottomPlayer'
import Toast from './Toast'

export default function Layout() {
  return (
    <div className="min-h-screen bg-bg text-white">
      <div className="flex min-h-screen pb-28">
        <Sidebar />
        <main className="flex-1 flex flex-col">
          <div className="flex-1 p-5 sm:p-8">
            <Outlet />
          </div>
        </main>
      </div>

      <div className="fixed inset-x-0 bottom-0 z-40 border-t border-surface2 bg-surface backdrop-blur-sm">
        <BottomPlayer />
      </div>

      <Toast />
    </div>
  )
}
