// Unit/integration tests for the Posts REST API using supertest.

import request from 'supertest';
import { app } from '../../src/api/server.js';
import { users, posts } from '../../src/api/data.js';

// Seed data snapshots — mirrors what data.js initialises with.
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

beforeEach(() => {
  // Reset both arrays to seed state.
  users.splice(0, users.length, ...SEED_USERS.map((u) => ({ ...u })));
  posts.splice(0, posts.length, ...SEED_POSTS.map((p) => ({ ...p })));
});

describe('Posts API', () => {
  it('GET /api/posts — returns 200 with an array of posts', async () => {
    const res = await request(app).get('/api/posts');
    expect(res.status).toBe(200);
    expect(Array.isArray(res.body)).toBe(true);
    expect(res.body.length).toBe(5);
    expect(res.body[0]).toMatchObject({ id: 1, userId: 1, title: 'Hello World' });
  });

  it('GET /api/posts/:id — returns 200 with a single post object', async () => {
    const res = await request(app).get('/api/posts/3');
    expect(res.status).toBe(200);
    expect(res.body).toMatchObject({ id: 3, userId: 2, title: 'Node.js Best Practices' });
  });

  it('GET /api/posts/:id — returns 404 when post does not exist', async () => {
    const res = await request(app).get('/api/posts/9999');
    expect(res.status).toBe(404);
    expect(res.body).toHaveProperty('error');
  });

  it('POST /api/posts — 201 happy path: creates post referencing a valid userId', async () => {
    const res = await request(app)
      .post('/api/posts')
      .send({ userId: 1, title: 'New Post', body: 'Post body content here.' });
    expect(res.status).toBe(201);
    expect(res.body).toMatchObject({ userId: 1, title: 'New Post', body: 'Post body content here.' });
    expect(res.body).toHaveProperty('id');
    expect(res.body).toHaveProperty('createdAt');
  });

  it('POST /api/posts — 400 validation error when userId is missing or invalid', async () => {
    // Missing userId
    const res1 = await request(app)
      .post('/api/posts')
      .send({ title: 'No User', body: 'Some body.' });
    expect(res1.status).toBe(400);
    expect(res1.body).toHaveProperty('error');

    // Non-existent userId
    const res2 = await request(app)
      .post('/api/posts')
      .send({ userId: 9999, title: 'Ghost', body: 'Body text.' });
    expect(res2.status).toBe(404);
    expect(res2.body).toHaveProperty('error');
  });

  it('PUT /api/posts/:id — 200 updates an existing post', async () => {
    const res = await request(app)
      .put('/api/posts/1')
      .send({ title: 'Updated Title' });
    expect(res.status).toBe(200);
    expect(res.body).toMatchObject({ id: 1, title: 'Updated Title' });
  });

  it('DELETE /api/posts/:id — 204 removes an existing post', async () => {
    const del = await request(app).delete('/api/posts/5');
    expect(del.status).toBe(204);

    // Confirm it is gone
    const get = await request(app).get('/api/posts/5');
    expect(get.status).toBe(404);
  });
});
