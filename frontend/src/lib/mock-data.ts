export interface Author {
  id: number;
  username: string;
  avatar: string;
}

export interface Tag {
  id: number;
  name: string;
  color: string;
}

export interface ThreadItem {
  id: number;
  title: string;
  slug: string;
  content: string;
  status: string;
  view_count: number;
  upvotes: number;
  downvotes: number;
  replies_count: number;
  author: Author;
  tags: Tag[];
  created_at: string;
}

export interface CommentItem {
  id: number;
  content: string;
  upvotes: number;
  downvotes: number;
  author: Author;
  replies: CommentItem[];
  created_at: string;
}

export interface ThreadDetail extends ThreadItem {
  comments: CommentItem[];
}

export interface Pagination {
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
}

export interface ThreadsResponse {
  threads: ThreadItem[];
  pagination: Pagination;
}

export interface UserProfile {
  id: number;
  username: string;
  avatar: string;
  bio: string;
  role: string;
  stacks: string[];
  created_at: string;
}

// ---------- Mock Data ----------

export const MOCK_TRENDING: ThreadItem[] = [
  {
    id: 1,
    title: "State management in Rust-based UI",
    slug: "state-management-rust-ui",
    content: "Discussion on state management approaches in Rust-based UI frameworks.",
    status: "published",
    view_count: 892,
    upvotes: 89,
    downvotes: 3,
    replies_count: 42,
    author: { id: 2, username: "@system_admin", avatar: "" },
    tags: [{ id: 1, name: "ANNOUNCEMENT", color: "#6366f1" }],
    created_at: "2026-05-06T10:00:00Z",
  },
  {
    id: 2,
    title: "Legacy code archive is now live for v1 members",
    slug: "legacy-code-archive-live",
    content: "The legacy code archive has been published for v1 members access.",
    status: "published",
    view_count: 654,
    upvotes: 72,
    downvotes: 1,
    replies_count: 28,
    author: { id: 3, username: "@void_strider", avatar: "" },
    tags: [{ id: 2, name: "ANNOUNCEMENT", color: "#6366f1" }],
    created_at: "2026-05-05T14:00:00Z",
  },
  {
    id: 3,
    title: "Shader optimization techniques for WebGL",
    slug: "shader-optimization-webgl",
    content: "Exploring optimization techniques for WebGL shader performance.",
    status: "published",
    view_count: 421,
    upvotes: 56,
    downvotes: 2,
    replies_count: 15,
    author: { id: 4, username: "@codex_null", avatar: "" },
    tags: [{ id: 3, name: "DESIGN DOCS", color: "#22c55e" }],
    created_at: "2026-05-04T08:00:00Z",
  },
];

export const MOCK_THREADS: ThreadItem[] = [
  {
    id: 10,
    title: "Best practices for CSS containment in complex dashboards?",
    slug: "css-containment-dashboards",
    content: "Looking for best practices on CSS containment for dashboards.",
    status: "published",
    view_count: 245,
    upvotes: 24,
    downvotes: 1,
    replies_count: 12,
    author: { id: 3, username: "@void_strider", avatar: "" },
    tags: [{ id: 11, name: "DESIGN", color: "#f59e0b" }],
    created_at: "2026-05-06T04:00:00Z",
  },
  {
    id: 11,
    title: "Server Maintenance: Scheduled migration to the new edge nodes",
    slug: "server-maintenance-migration",
    content: "Scheduled server maintenance for migrating to new edge nodes.",
    status: "published",
    view_count: 567,
    upvotes: 56,
    downvotes: 0,
    replies_count: 8,
    author: { id: 2, username: "@system_admin", avatar: "" },
    tags: [{ id: 1, name: "ANNOUNCEMENTS", color: "#6366f1" }],
    created_at: "2026-05-05T12:00:00Z",
  },
  {
    id: 12,
    title: "Optimizing garbage collection for high-throughput node.js streams",
    slug: "optimizing-gc-nodejs-streams",
    content: "Tips for optimizing GC in high-throughput Node.js streams.",
    status: "published",
    view_count: 189,
    upvotes: 8,
    downvotes: 2,
    replies_count: 45,
    author: { id: 4, username: "@codex_null", avatar: "" },
    tags: [{ id: 4, name: "TECHNICAL", color: "#ef4444" }],
    created_at: "2026-05-05T08:00:00Z",
  },
  {
    id: 13,
    title: "Has anyone successfully implemented eye-tracking navigation in VR?",
    slug: "eye-tracking-navigation-vr",
    content: "I'm trying to implement eye-tracking for VR navigation. Any tips?",
    status: "published",
    view_count: 2340,
    upvotes: 112,
    downvotes: 8,
    replies_count: 301,
    author: { id: 5, username: "@neural_link", avatar: "" },
    tags: [{ id: 5, name: "THE VOID", color: "#8b5cf6" }],
    created_at: "2026-05-04T10:00:00Z",
  },
  {
    id: 14,
    title: "Edge-cached event streams: predictable hydration patterns explored",
    slug: "edge-cached-event-streams",
    content: "Exploring predictable hydration patterns with edge-cached event streams.",
    status: "published",
    view_count: 312,
    upvotes: 38,
    downvotes: 1,
    replies_count: 22,
    author: { id: 6, username: "@quantum_byte", avatar: "" },
    tags: [{ id: 4, name: "TECHNICAL", color: "#ef4444" }],
    created_at: "2026-05-03T16:00:00Z",
  },
];

export const MOCK_FEATURED: ThreadDetail = {
  id: 20,
  title: "Architectural Shift: Transitioning to Reactive Components in v2.4",
  slug: "architectural-shift",
  content: `The transition to version 2.4 marks a fundamental departure from our previous imperative component lifecycle. We are moving towards a fully reactive, state-driven architecture that prioritizes deterministic rendering over performance-heavy mutations.

This shift necessitates a re-evaluation of how we handle global state injections. In the new reactive model, components no longer "fetch" their dependencies; they are "hydrated" by the system's core event bus through a series of prioritized streams.`,
  status: "published",
  view_count: 52400,
  upvotes: 1482,
  downvotes: 18,
  replies_count: 142,
  author: { id: 1, username: "@alpha_lead", avatar: "" },
  tags: [{ id: 6, name: "System Core", color: "#6366f1" }],
  created_at: "2026-05-06T02:00:00Z",
  comments: [
    {
      id: 101,
      content:
        "This transition is overdue. The imperative bloat in v2.3 was becoming unmanageable for large-scale telemetry dashboards. My only concern is the overhead on the primary event bus—have we stress-tested the reactive hydration with >10k concurrent nodes?",
      upvotes: 82,
      downvotes: 3,
      author: { id: 7, username: "@void_iterator", avatar: "" },
      created_at: "2026-05-06T01:00:00Z",
      replies: [
        {
          id: 102,
          content:
            "@void_iterator Yes, internal tests on the 'KRYPTON' cluster showed a 14% memory reduction compared to v2.3 under high load. The event bus is now sharded by default.",
          upvotes: 45,
          downvotes: 1,
          author: { id: 1, username: "@alpha_lead", avatar: "" },
          replies: [],
          created_at: "2026-05-06T00:45:00Z",
        },
      ],
    },
    {
      id: 103,
      content:
        "Will there be an automated migration script for legacy v2.2 legacy_hooks? Or is it a manual rewrite?",
      upvotes: 12,
      downvotes: 0,
      author: { id: 8, username: "@null_pointer", avatar: "" },
      replies: [],
      created_at: "2026-05-06T00:35:00Z",
    },
  ],
};

export const MOCK_USER: UserProfile = {
  id: 1,
  username: "NODE_8829",
  avatar: "",
  bio: "Specializing in decentralized infrastructure and encrypted communication protocols. Former lead dev for Project Wraith. Currently optimizing neural-link throughput for the ARCHITECT_FORUM core.",
  role: "SENIOR ARCHITECT",
  stacks: ["Rust", "Solidity", "Go", "Post-Quantum Cryptography"],
  created_at: "2026-01-01T00:00:00Z",
};
