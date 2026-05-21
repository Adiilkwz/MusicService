import { useState, useEffect } from 'react';

export default function Admin({ token }) {
  const [users, setUsers] = useState([]);
  const [error, setError] = useState('');

  const fetchUsers = async () => {
    const res = await fetch('/api/admin/users/?limit=50', {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    const data = await res.json();
    if (res.ok) {
      setUsers(data || []);
    } else {
      setError(data.error || 'Failed to load users');
    }
  };

  useEffect(() => { fetchUsers(); }, []);

  const handleRoleChange = async (userId, newRole) => {
    const res = await fetch(`/api/admin/users/${userId}/role`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ new_role: newRole })
    });
    
    if (res.ok) {
      fetchUsers();
    } else {
      const data = await res.json();
      alert(data.error);
    }
  };

  return (
    <div style={{ padding: '20px', maxWidth: '800px', margin: '0 auto' }}>
      <h2>🛡️ Admin Control Panel</h2>
      {error ? (
        <p style={{ color: 'red' }}>{error} (Are you logged in as an Admin?)</p>
      ) : (
        <table style={{ width: '100%', textAlign: 'left', borderCollapse: 'collapse' }}>
          <thead>
            <tr style={{ borderBottom: '2px solid #ccc' }}>
              <th>ID</th>
              <th>Email</th>
              <th>Role</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {users.map(user => (
              <tr key={user.id} style={{ borderBottom: '1px solid #eee' }}>
                <td style={{ padding: '10px 0' }}>{user.id}</td>
                <td>{user.email}</td>
                <td>{user.role}</td>
                <td>
                  <select 
                    value={user.role} 
                    onChange={(e) => handleRoleChange(user.id, e.target.value)}
                    style={{ padding: '5px' }}
                  >
                    <option value="user">User</option>
                    <option value="artist">Artist</option>
                    <option value="admin">Admin</option>
                  </select>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}