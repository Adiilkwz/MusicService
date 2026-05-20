import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: { 'Content-Type': 'application/json' },
  withCredentials: true,
})

let currentToken = null
let onUnauthorized = () => {}
let onLoading = () => {}
let onToast = () => {}

api.interceptors.request.use(
  (config) => {
    onLoading(true)
    if (currentToken) {
      config.headers.Authorization = `Bearer ${currentToken}`
    }
    return config
  },
  (error) => {
    onLoading(false)
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  (response) => {
    onLoading(false)
    return response
  },
  (error) => {
    onLoading(false)
    const status = error.response?.status
    if (status === 401) {
      onUnauthorized()
    }
    const message = error.response?.data?.message || error.message || 'Request failed'
    onToast(message, 'error')
    return Promise.reject(error)
  }
)

const setToken = (token) => {
  currentToken = token
}

const clearToken = () => {
  currentToken = null
}

const subscribeLoading = (handler) => {
  onLoading = handler
}

const subscribeUnauthorized = (handler) => {
  onUnauthorized = handler
}

const subscribeToast = (handler) => {
  onToast = handler
}

export default {
  get: api.get,
  post: api.post,
  put: api.put,
  delete: api.delete,
  setToken,
  clearToken,
  subscribeLoading,
  subscribeUnauthorized,
  subscribeToast,
}
