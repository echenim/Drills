import React, { useState } from 'react';

const EMPTY = { name: '', email: '' };

export default function UserForm({ onUserCreated }) {
  const [fields, setFields] = useState(EMPTY);
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState(null);

  const handleChange = (e) => {
    setFields((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    setFormError(null);

    const name = fields.name.trim();
    const email = fields.email.trim();

    if (!name || !email) {
      setFormError('Both name and email are required.');
      return;
    }

    setSubmitting(true);
    fetch('/api/users', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, email }),
    })
      .then((res) => {
        if (!res.ok) {
          return res.json().then((body) => {
            throw new Error(body.error || `HTTP ${res.status}`);
          });
        }
        return res.json();
      })
      .then(() => {
        setFields(EMPTY);
        setSubmitting(false);
        onUserCreated();
      })
      .catch((err) => {
        setFormError(err.message || 'Failed to create user.');
        setSubmitting(false);
      });
  };

  return (
    <div className="form-panel" data-testid="user-form">
      <h3>New User</h3>
      <form onSubmit={handleSubmit} noValidate>
        <div className="form-row">
          <input
            type="text"
            name="name"
            placeholder="Full name"
            value={fields.name}
            onChange={handleChange}
            data-testid="user-name-input"
            autoComplete="off"
          />
          <input
            type="email"
            name="email"
            placeholder="Email address"
            value={fields.email}
            onChange={handleChange}
            data-testid="user-email-input"
            autoComplete="off"
          />
        </div>
        {formError && (
          <p className="form-error" data-testid="user-form-error">{formError}</p>
        )}
        <button
          type="submit"
          className="btn-submit"
          disabled={submitting}
          data-testid="user-submit-btn"
        >
          {submitting ? 'Creating…' : 'Create User'}
        </button>
      </form>
    </div>
  );
}
