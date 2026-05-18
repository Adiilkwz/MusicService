import { useState, useEffect } from 'react';

export default function Dashboard({ token }) {
  const [likedSongs, setLikedSongs] = useState([]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const fetchLikedSongs = async () => {
    setLoading(true);
    setError('');
    
    try {
      const response = await fetch('/api/likes', {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Failed to fetch liked songs');
      }

      setLikedSongs(data.songs || []);
      
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ padding: '20px', fontFamily: 'sans-serif', maxWidth: '800px', margin: '0 auto' }}>
      <h2>My Music Dashboard</h2>
      
      <button 
        onClick={fetchLikedSongs}
        style={{ padding: '10px 20px', background: '#1DB954', color: 'white', border: 'none', borderRadius: '20px', cursor: 'pointer', marginBottom: '20px' }}
      >
        {loading ? 'Loading...' : 'Fetch My Liked Songs'}
      </button>

      {error && <p style={{ color: 'red' }}>{error}</p>}

      <div style={{ display: 'grid', gap: '10px' }}>
        {likedSongs.length === 0 && !loading && !error && (
          <p style={{ color: '#666' }}>Click the button to load your liked songs, or you don't have any likes yet!</p>
        )}
        
        {likedSongs.map((song, index) => (
          <div key={index} style={{ padding: '15px', background: '#282c34', color: 'white', borderRadius: '8px', display: 'flex', justifyContent: 'space-between' }}>
            <div>
              <strong style={{ fontSize: '18px' }}>{song.title}</strong>
              <p style={{ margin: '5px 0 0 0', color: '#aaa', fontSize: '14px' }}>Genre: {song.genre}</p>
            </div>
            <div style={{ color: '#1DB954' }}>▶ Play</div>
          </div>
        ))}
      </div>
    </div>
  );
}