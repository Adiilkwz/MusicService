export default {
  content: ['./index.html', './src/**/*.{js,jsx}'],
  theme: {
    extend: {
      colors: {
        bg: '#0e0f14',
        surface: '#141720',
        surface2: '#1f2230',
        muted: '#8b94a8',
        accent: '#1db954',
        accentSoft: '#1db95433',
      },
      boxShadow: {
        glow: '0 0 50px rgba(29, 185, 84, 0.18)',
      },
    },
  },
  plugins: [],
}
