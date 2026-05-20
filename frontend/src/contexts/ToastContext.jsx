import { createContext, useContext, useEffect, useMemo, useState } from 'react'
import api from '../services/api'

const ToastContext = createContext(null)

export function ToastProvider({ children }) {
  const [toasts, setToasts] = useState([])

  const addToast = (message, type = 'error') => {
    const id = Date.now().toString()
    setToasts((prev) => [...prev, { id, message, type }])
    setTimeout(() => setToasts((prev) => prev.filter((toast) => toast.id !== id)), 4500)
  }

  useEffect(() => {
    api.subscribeToast(addToast)
    return () => api.subscribeToast(() => {})
  }, [])

  const value = useMemo(() => ({ toasts, addToast }), [toasts])

  return <ToastContext.Provider value={value}>{children}</ToastContext.Provider>
}

export function useToast() {
  return useContext(ToastContext)
}
