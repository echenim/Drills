import React, { useState } from 'react';

const EMPTY = { userId: '', title: '', body: '' };

export default function PostForm({ users, usersLoading, onPostCreated }) {
  const [fields, setFields] = useState(EMPTY);
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState(null);

  const handleChange = (e) => {
    setFields((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    setFormError(null);

    const userId = Number(fields.userId);
    const title = fields.title.trim();
    const body = fields.body.trim();

    if (!userId || !title || !body) {
      setFormError('User, title, and body are all required.');
      return;
    }

    setSubmitting(true);
    fetch('/api/posts', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ userId, title, body }),
    })
      .then((res) => {
        if (!res.ok) {
          return res.json().then((payload) => {
            throw new Error(payload.error || `HTTP ${res.status}`);
          });
        }
        return res.json();
      })
      .then(() => {
        setFields(EMPTY);
        setSubmitting(false);
        onPostCreated();
      })
      .catch((err) => {
        setFormError(err.message || 'Failed to create post.');
        setSubmitting(false);
      });
  };

  return (
    <div className="form-panel" data-testid="post-form">
      <h3>New Post</h3>
      <form onSubmit={handleSubmit} noValidate>
        <div className="form-row">
          <select
            name="userId"
            value={fields.userId}
            onChange={handleChange}
            data-testid="post-user-select"
            disabled={usersLoading}
          >
            <option value="">
              {usersLoading ? 'Loading users…' : '— Select an author —'}
            </option>
            {users.map((user) => (
              <option key={user.id} value={user.id}>
                {user.name}
              </option>
            ))}
          </select>
          <input
            type="text"
            name="title"
            placeholder="Post title"
            value={fields.title}
            onChange={handleChange}
            data-testid="post-title-input"
            autoComplete="off"
          />
          <textarea
            name="body"
            placeholder="Write something…"
            value={fields.body}
            onChange={handleChange}
            data-testid="post-body-input"
          />
        </div>
        {formError && (
          <p className="form-error" data-testid="post-form-error">{formError}</p>
        )}
        <button
          type="submit"
          className="btn-submit"
          disabled={submitting || usersLoading}
          data-testid="post-submit-btn"
        >
          {submitting ? 'Posting…' : 'Create Post'}
        </button>
      </form>
    </div>
  );
}
