// Integration tests — exercises the full stack (Express API + React UI) together.
// fetch is bridged to supertest so no real TCP port is needed.

import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import supertest from 'supertest';

import App from '../../src/App.jsx';
import { app } from '../../src/api/server.js';
import { users, posts } from '../../src/api/data.js';

// ---------------------------------------------------------------------------
// Seed snapshots
// ---------------------------------------------------------------------------
const SEED_USERS = [
  { id: 1, name: 'Alice Johnson', email: 'alice@example.com', createdAt: '2024-01-10T08:00:00.000Z' },
  { id: 2, name: 'Bob Smith',     email: 'bob@example.com',   createdAt: '2024-01-11T09:00:00.000Z' },
  { id: 3, name: 'Carol White',   email: 'carol@example.com', createdAt: '2024-01-12T10:00:00.000Z' },
];

const SEED_POSTS = [
  { id: 1, userId: 1, title: 'Hello World',            body: 'My first post on this platform!',            createdAt: '2024-01-15T08:30:00.000Z' },
  { id: 2, userId: 1, title: 'React Tips',             body: 'Here are some tips for building React apps.', createdAt: '2024-01-16T09:00:00.000Z' },
  { id: 3, userId: 2, title: 'Node.js Best Practices', body: 'Things I learned building REST APIs with Node.js.', createdAt: '2024-01-17T10:00:00.000Z' },
  { id: 4, userId: 3, title: 'CSS Grid vs Flexbox',    body: 'A comparison of modern CSS layout techniques.',  createdAt: '2024-01-18T11:00:00.000Z' },
  { id: 5, userId: 2, title: 'Testing with Vitest',    body: 'How to write fast unit tests using Vitest.',     createdAt: '2024-01-19T12:00:00.000Z' },
];

// ---------------------------------------------------------------------------
// Supertest → fetch bridge
// Converts relative `/api/*` paths to real supertest calls and returns a
// Response-compatible object so React components work unchanged.
// ---------------------------------------------------------------------------
async function supertestFetch(url, opts = {}) {
  const path = url.startsWith('/') ? url : `/${new URL(url).pathname}`;
  const method = (opts.method || 'GET').toLowerCase();
  let req = supertest(app)[method](path);

  if (opts.headers) {
    Object.entries(opts.headers).forEach(([k, v]) => {
      req = req.set(k, v);
    });
  }

  if (opts.body) {
    req = req.send(JSON.parse(opts.body));
  }

  const stRes = await req;

  const body = stRes.body;
  const text = typeof body === 'string' ? body : JSON.stringify(body);

  return {
    ok: stRes.status >= 200 && stRes.status < 300,
    status: stRes.status,
    json: () => Promise.resolve(body),
    text: () => Promise.resolve(text),
  };
}

// ---------------------------------------------------------------------------
// Lifecycle hooks
// ---------------------------------------------------------------------------
beforeEach(() => {
  // Reset in-memory data to seed state.
  users.splice(0, users.length, ...SEED_USERS.map((u) => ({ ...u })));
  posts.splice(0, posts.length, ...SEED_POSTS.map((p) => ({ ...p })));

  // Route all fetch calls through supertest.
  global.fetch = vi.fn((url, opts) => supertestFetch(url, opts));
});

afterEach(() => {
  vi.restoreAllMocks();
});

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
describe('Full-stack integration flows', () => {
  it('full user flow — load existing users, create a new user, see it appear in the list', async () => {
    const user = userEvent.setup();

    render(<App />);

    // Initial user list loads from seed data.
    await waitFor(() => {
      expect(screen.getByTestId('users-list')).toBeInTheDocument();
    });

    expect(screen.getByTestId('user-card-1')).toBeInTheDocument();
    expect(screen.getByTestId('user-card-2')).toBeInTheDocument();
    expect(screen.getByTestId('user-card-3')).toBeInTheDocument();

    // Create a new user via the form.
    await user.type(screen.getByTestId('user-name-input'), 'Eve Newuser');
    await user.type(screen.getByTestId('user-email-input'), 'eve@test.com');
    await user.click(screen.getByTestId('user-submit-btn'));

    // After creation + refresh, new card must appear (id 4 from seed nextUserId).
    await waitFor(() => {
      expect(screen.getByTestId('user-card-4')).toBeInTheDocument();
    });

    expect(screen.getByText('Eve Newuser')).toBeInTheDocument();
  });

  it('full post flow — load existing posts, create a post linked to a user, see it appear', async () => {
    const user = userEvent.setup();

    render(<App />);

    // Wait for users to load (needed for the Posts tab to function).
    await waitFor(() => {
      expect(screen.getByTestId('users-list')).toBeInTheDocument();
    });

    // Switch to Posts tab.
    await user.click(screen.getByTestId('tab-posts'));

    // Existing posts should render.
    await waitFor(() => {
      expect(screen.getByTestId('posts-list')).toBeInTheDocument();
    });

    expect(screen.getByTestId('post-card-1')).toBeInTheDocument();

    // Create a new post.
    await user.selectOptions(screen.getByTestId('post-user-select'), '1');
    await user.type(screen.getByTestId('post-title-input'), 'Integration Post');
    await user.type(screen.getByTestId('post-body-input'), 'Created in integration test.');
    await user.click(screen.getByTestId('post-submit-btn'));

    // New post card should appear (id 6 from seed nextPostId).
    await waitFor(() => {
      expect(screen.getByTestId('post-card-6')).toBeInTheDocument();
    });

    expect(screen.getByText('Integration Post')).toBeInTheDocument();
  });

  it('API + UI happy path — app boots, fetches data from the live API, and renders without errors', async () => {
    render(<App />);

    // Users section renders by default.
    await waitFor(() => {
      expect(screen.getByTestId('users-section')).toBeInTheDocument();
    });

    // Users list (not error, not loading) should be present.
    await waitFor(() => {
      expect(screen.getByTestId('users-list')).toBeInTheDocument();
    });

    expect(screen.queryByTestId('users-error')).not.toBeInTheDocument();
    expect(screen.queryByTestId('users-loading')).not.toBeInTheDocument();

    // Health endpoint also works.
    const healthRes = await supertestFetch('/api/health');
    expect(healthRes.ok).toBe(true);
    const healthBody = await healthRes.json();
    expect(healthBody).toMatchObject({ status: 'ok' });
  });
});
