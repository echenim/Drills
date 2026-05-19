import React, { useState, useEffect } from 'react';

export default function PostsList({ users, version }) {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    setLoading(true);
    setError(null);
    fetch('/api/posts')
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then((data) => {
        setPosts(data);
        setLoading(false);
      })
      .catch(() => {
        setError("Couldn't load posts.");
        setLoading(false);
      });
  }, [version]);

  const getAuthorName = (userId) => {
    if (!users || users.length === 0) return `User #${userId}`;
    const found = users.find((u) => u.id === userId);
    return found ? found.name : `User #${userId}`;
  };

  if (loading) {
    return <p className="status-loading" data-testid="posts-loading">Loading posts…</p>;
  }

  if (error) {
    return <p className="status-error" data-testid="posts-error">{error}</p>;
  }

  if (posts.length === 0) {
    return <p className="status-loading" data-testid="posts-empty">No posts yet.</p>;
  }

  return (
    <div className="cards-grid" data-testid="posts-list">
      {posts.map((post) => (
        <div className="card" key={post.id} data-testid={`post-card-${post.id}`}>
          <div className="card-title">{post.title}</div>
          <div className="card-body">{post.body}</div>
          <div className="card-author" data-testid={`post-author-${post.id}`}>
            By {getAuthorName(post.userId)}
          </div>
          <div className="card-meta">
            Post #{post.id} &middot; {new Date(post.createdAt).toLocaleDateString()}
          </div>
        </div>
      ))}
    </div>
  );
}
