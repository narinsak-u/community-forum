export const queryKeys = {
  auth: {
    me: ["me"] as const,
  },
  threads: {
    all: ["threads"] as const,
    trending: ["threads", "trending"] as const,
    featured: ["threads", "featured"] as const,
    list: (params: { page: number; pageSize: number; sort: string }) =>
      ["threads", params] as const,
    detail: (slug: string) => ["threads", slug] as const,
    comments: (slug: string) => ["threads", slug, "comments"] as const,
  },
  users: {
    profile: (username: string) => ["users", username] as const,
    threads: (username: string) => ["users", username, "threads"] as const,
  },
  votes: {
    thread: (slug: string) => ["votes", "thread", slug] as const,
  },
} as const;
