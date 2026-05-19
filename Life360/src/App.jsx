import React, { useState, useEffect } from 'react';
import UsersList from './components/UsersList.jsx';
import UserForm from './components/UserForm.jsx';
import PostsList from './components/PostsList.jsx';
import PostForm from './components/PostForm.jsx';

export default function App() {
  const [activeTab, setActiveTab] = useState('users');

  // Shared users state — PostsList and PostForm both need the users list
  const [users, setUsers] = useState([]);
  const [usersLoading, setUsersLoading] = useState(true);
  const [usersError, setUsersError] = useState(null);

  // Posts version counter — bump to trigger PostsList refresh
  const [postsVersion, setPostsVersion] = useState(0);

  const fetchUsers = () => {
    setUsersLoading(true);
    setUsersError(null);
    fetch('/api/users')
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then((data) => {
        setUsers(data);
        setUsersLoading(false);
      })
      .catch((err) => {
        setUsersError("Couldn't load users.");
        setUsersLoading(false);
      });
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleUserCreated = () => {
    fetchUsers();
  };

  const handlePostCreated = () => {
    setPostsVersion((v) => v + 1);
  };

  return (
    <div className="app">
      <header className="app-header">
        <h1>Life360</h1>
        <nav className="tab-bar" role="tablist">
          <button
            className={`tab-btn${activeTab === 'users' ? ' active' : ''}`}
            role="tab"
            aria-selected={activeTab === 'users'}
            data-testid="tab-users"
            onClick={() => setActiveTab('users')}
          >
            Users
          </button>
          <button
            className={`tab-btn${activeTab === 'posts' ? ' active' : ''}`}
            role="tab"
            aria-selected={activeTab === 'posts'}
            data-testid="tab-posts"
            onClick={() => setActiveTab('posts')}
          >
            Posts
          </button>
        </nav>
      </header>

      {activeTab === 'users' && (
        <main data-testid="users-section">
          <div className="section">
            <h2 className="section-title">Add a User</h2>
            <UserForm onUserCreated={handleUserCreated} />
          </div>
          <div className="section">
            <h2 className="section-title">All Users</h2>
            <UsersList
              users={users}
              loading={usersLoading}
              error={usersError}
            />
          </div>
        </main>
      )}

      {activeTab === 'posts' && (
        <main data-testid="posts-section">
          <div className="section">
            <h2 className="section-title">Add a Post</h2>
            <PostForm
              users={users}
              usersLoading={usersLoading}
              onPostCreated={handlePostCreated}
            />
          </div>
          <div className="section">
            <h2 className="section-title">All Posts</h2>
            <PostsList users={users} version={postsVersion} />
          </div>
        </main>
      )}
    </div>
  );
}
