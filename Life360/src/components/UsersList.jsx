import React from 'react';

export default function UsersList({ users, loading, error }) {
  if (loading) {
    return <p className="status-loading" data-testid="users-loading">Loading users…</p>;
  }

  if (error) {
    return <p className="status-error" data-testid="users-error">{error}</p>;
  }

  if (users.length === 0) {
    return <p className="status-loading" data-testid="users-empty">No users yet.</p>;
  }

  return (
    <div className="cards-grid" data-testid="users-list">
      {users.map((user) => (
        <div className="card" key={user.id} data-testid={`user-card-${user.id}`}>
          <div className="card-name">{user.name}</div>
          <div className="card-email">{user.email}</div>
          <div className="card-meta">
            ID: {user.id} &middot; Joined {new Date(user.createdAt).toLocaleDateString()}
          </div>
        </div>
      ))}
    </div>
  );
}
