import { Outlet } from 'react-router-dom'
import Sidebar from './Sidebar'
import BottomPlayer from './BottomPlayer'
import Toast from './Toast'

export default function Layout() {
  return (
    <div className="min-h-screen bg-bg text-white">
      <div className="flex min-h-screen">
        <Sidebar />
        <div className="flex-1 flex flex-col">
          <div className="flex-1 p-5 sm:p-8">
            <Outlet />
          </div>
          <BottomPlayer />
        </div>
      </div>
      <Toast />
    </div>
  )
}
