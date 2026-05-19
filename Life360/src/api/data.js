// In-memory data store with seed data

export const users = [
  { id: 1, name: "Alice Johnson", email: "alice@example.com", createdAt: "2024-01-10T08:00:00.000Z" },
  { id: 2, name: "Bob Smith", email: "bob@example.com", createdAt: "2024-01-11T09:00:00.000Z" },
  { id: 3, name: "Carol White", email: "carol@example.com", createdAt: "2024-01-12T10:00:00.000Z" },
];

export const posts = [
  { id: 1, userId: 1, title: "Hello World", body: "My first post on this platform!", createdAt: "2024-01-15T08:30:00.000Z" },
  { id: 2, userId: 1, title: "React Tips", body: "Here are some tips for building React apps.", createdAt: "2024-01-16T09:00:00.000Z" },
  { id: 3, userId: 2, title: "Node.js Best Practices", body: "Things I learned building REST APIs with Node.js.", createdAt: "2024-01-17T10:00:00.000Z" },
  { id: 4, userId: 3, title: "CSS Grid vs Flexbox", body: "A comparison of modern CSS layout techniques.", createdAt: "2024-01-18T11:00:00.000Z" },
  { id: 5, userId: 2, title: "Testing with Vitest", body: "How to write fast unit tests using Vitest.", createdAt: "2024-01-19T12:00:00.000Z" },
];

let nextUserId = 4;
let nextPostId = 6;

export function getNextUserId() {
  return nextUserId++;
}

export function getNextPostId() {
  return nextPostId++;
}
