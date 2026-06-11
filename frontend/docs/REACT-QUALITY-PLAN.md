# React Quality Improvement Plan

Review conducted using: react-quality-check, react-doctor, react-hooks-optimizer, code-review-and-quality, code-simplification, vercel-react-best-practices, vercel-composition-patterns.

---

## Phase 1: Correct Critical Bugs

### 1.1 Auth guard returns `true` before fetch completes
**File:** `src/hooks/use-require-auth.ts` (line 18)

Replace `!isFetched` returning `true` with `false` (or intermediate loading state) to prevent unauthenticated users from performing guarded actions before the auth query resolves.

```ts
// Before
if (!isFetched) return true;

// After
if (!isFetched) return false; // or show loading
```

### 1.2 Eliminate dual auth source of truth
**Files:** `src/components/session-loader.tsx`, `src/stores/auth-store.ts`, `src/hooks/use-auth.ts`

Pick ONE source of truth for auth state:
- **Option A (Recommended):** Remove Zustand auth store. Use TanStack Query (`useMe`) everywhere. Remove `SessionLoader` entirely. Remove `useAuthStore` imports from `Sidebar.tsx`, `CommentSection.tsx`, `settings/page.tsx`.
- **Option B:** Remove `useMe`/`useSignin`/`useSignout`. Manage all auth in Zustand with direct API calls. Remove `SessionLoader`.

Option A is preferred since the codebase already uses TanStack Query for all other data.

### 1.3 `useMobile` hook extra render
**File:** `src/hooks/use-mobile.tsx`

```ts
// Before
const [isMobile, setIsMobile] = React.useState<boolean | undefined>(undefined);

// After — synchronously initialize
const [isMobile, setIsMobile] = React.useState(false);
// Or use useSyncExternalStore for best practice
```

---

## Phase 2: Architectural Fixes

### 2.1 Extract `getInitials` to utils
**Files:** `src/app/(main)/thread/[slug]/ThreadDetailClient.tsx` (line 44), `src/app/(main)/thread/[slug]/CommentSection.tsx` (line 15-17)

Extract into `src/lib/utils.ts`:
```ts
export function getInitials(username?: string): string {
  return (username || "??").replace("@", "").slice(0, 2).toUpperCase();
}
```

Then update all callers to import from `@/lib/utils`.

### 2.2 Extract participants IIFE to `useMemo`
**File:** `src/app/(main)/thread/[slug]/ThreadDetailClient.tsx` (lines 190-226)

Pull the inline IIFE into a `useMemo`:
```ts
const participants = useMemo(() => {
  const authorMap = new Map<...>();
  // ... collection logic
  return Array.from(authorMap.values()).slice(0, 5);
}, [currentThread?.comments, currentThread?.author]);
```

### 2.3 Flatten conditional query data chains
**File:** `src/app/(main)/threads/ThreadsClient.tsx` (lines 61-71)

Pass `initialData` directly to the TanStack Query hooks instead of post-processing with `??` chains:
```ts
// In the hook definition (e.g., useThreads):
export function useThreads(page = 1, pageSize = 5, sort = "latest", initialData?: ThreadsResponse) {
  return useQuery({
    queryKey: queryKeys.threads.list({ page, pageSize, sort }),
    queryFn: () => api.get<ThreadsResponse>(`/threads?page=${page}&pageSize=${pageSize}&sort=${sort}`),
    staleTime: 60 * 1000,
    initialData,
  });
}
```

### 2.4 Normalize hook URL patterns
**Files:** All hooks in `src/hooks/`

Consistently use template literals everywhere:
```ts
// Before
api.get<ThreadDetail>("/threads/" + slug)
api.post<VoteResponse>("/threads/" + slug + "/vote", data)

// After
api.get<ThreadDetail>(`/threads/${slug}`)
api.post<VoteResponse>(`/threads/${slug}/vote`, data)
```

---

## Phase 3: Performance

### 3.1 Add `enabled` to queries for inactive tabs
**File:** `src/app/(main)/threads/ThreadsClient.tsx` (lines 36-40)

Use `enabled` to prevent query execution for inactive tabs:
```ts
const allThreads = useThreads(page, 5, sort, { enabled: activeTab !== "my-posts" });
const myThreads = useMyThreads(username ?? "", page, 5, { enabled: activeTab === "my-posts" });
```

Update hook signatures to accept `options`:
```ts
export function useThreads(page = 1, pageSize = 5, sort = "latest", options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: queryKeys.threads.list({ page, pageSize, sort }),
    queryFn: ...,
    staleTime: 60 * 1000,
    ...options,
  });
}
```

### 3.2 Add `priority` to hero images above the fold
**File:** `src/app/(main)/threads/ThreadsClient.tsx` (line 105)
**File:** `src/app/(main)/thread/[slug]/ThreadDetailClient.tsx` (line 99)

Ensure the featured hero image gets `priority` on `ThreadsClient` (it already has it on `ThreadDetailClient`).

### 3.3 Consolidate `ThreadHeader` props
**File:** `src/app/(main)/thread/[slug]/ThreadHeader.tsx` (lines 5-18)

Replace 17 individual props with a thread summary object:
```ts
interface ThreadHeaderProps {
  thread: {
    title?: string;
    tagName?: string;
    author?: { username?: string; id?: number };
    createdAt?: string;
    viewCount: number;
    repliesCount: number;
    upvotes: number;
    downvotes: number;
  };
  slug?: string;
  currentUserId?: number;
}
```

---

## Phase 4: Code Simplification

### 4.1 Remove `useMemo` from trivial transformations
**File:** `src/components/forge/EntryEditor.tsx` (line 38)

```ts
// Before (unnecessary memoization)
const previewTitle = useMemo(
  () => (title || "NEW_ENTRY").replace(/\s+/g, "_").toUpperCase(),
  [title],
);

// After
const previewTitle = (title || "NEW_ENTRY").replace(/\s+/g, "_").toUpperCase();
```

### 4.2 Extract `formatCount` outside component
**File:** `src/app/(main)/thread/[slug]/ThreadHeader.tsx` (line 39)

Move `formatCount` to module-level or utils to avoid recreating on every render.

### 4.3 Split/full toggle button dedup
**File:** `src/components/forge/EntryEditor.tsx` (lines 92-114)

Map over an array for the two view mode buttons to eliminate duplicate JSX:
```ts
{["split", "full"].map((mode) => (
  <button
    key={mode}
    onClick={() => onViewModeChange(mode as "split" | "full")}
    className={`px-3 py-1.5 transition-colors ${viewMode === mode ? "bg-primary text-primary-foreground font-bold" : "text-muted-foreground hover:bg-secondary"}`}
  >
    {mode}
  </button>
))}
```

### 4.4 Remove `console.log` in production
**File:** `src/hooks/use-chat.ts` (line 55)

```ts
// Remove
console.log("chat: WebSocket connected to", WS_URL);
```

### 4.5 Remove commented-out code
**File:** `src/components/forge/Sidebar.tsx` (lines 69-73)

Remove the commented `<span>NAV</span>` block.

---

## Phase 5: Composition Patterns

### 5.1 Replace boolean prop with slot/children
**File:** `src/components/forge/Sidebar.tsx`

Instead of `showNewEntry` boolean, accept a `newEntrySlot` ReactNode:
```ts
interface SidebarProps {
  newEntrySlot?: React.ReactNode;
}
```

Caller decides whether to render the button, not the sidebar.

### 5.2 `ThreadDetailClient` — extract participants panel
**File:** `src/app/(main)/thread/[slug]/ThreadDetailClient.tsx`

Extract the participants panel into its own component (`ParticipantsPanel`) for clarity and reusability.

---

## Phase 6: Error Handling & Edge Cases

### 6.1 Add error state to thread detail
**File:** `src/app/(main)/thread/[slug]/ThreadDetailClient.tsx`

Handle the `isError` case from `useThread`:
```ts
const { data: thread, isLoading, isError } = useThread(slug);

if (isError) {
  return <ErrorState message="Failed to load thread" />;
}
```

### 6.2 Add client-side validation to signup
**File:** `src/hooks/use-auth.ts` (and/or the signup form)

Validate `username` length, `email` format, `password` strength before sending to API.

### 6.3 Add sizes prop to Image components
**File:** `src/app/(main)/thread/[slug]/CommentSection.tsx` — add `sizes` prop where `fill` is used for images.

Check all `<Image>` components with `fill` that are missing `sizes`.

---

## Phase 7: Testing

### 7.1 Add test coverage for critical paths

| Priority | Area | What to Test |
|----------|------|-------------|
| P0 | Auth | `useSignin`, `useSignout`, `useRequireAuth` redirect/toast behavior |
| P0 | Threads | `useThreads` query construction, `useCreateThread` mutation |
| P1 | Voting | `useVoteThread`, `useVoteComment` mutation + cache invalidation |
| P1 | Comments | `useCreateComment`, `useDeleteComment` with `parentId` |
| P1 | Auth store | `setUser`, `logout`, `isAuthenticated` computed value |
| P2 | Utils | `getInitials`, `timeAgo`, `formatCount`, `cn` |

### 7.2 Add integration tests for key flows
- User signup → redirect to login
- User login → session persists → can create thread
- Thread detail → load comments → vote → cache updates

---

## Execution Order

```
Phase 1: Critical Bugs (fix now, highest risk)
    └─ 1.1, 1.2, 1.3

Phase 2: Architecture (fix next, medium effort, high payoff)
    └─ 2.1, 2.2, 2.3, 2.4

Phase 3: Performance (quick wins)
    └─ 3.1, 3.2, 3.3

Phase 4: Simplification (cleanup)
    └─ 4.1, 4.2, 4.3, 4.4, 4.5

Phase 5: Composition (design improvements)
    └─ 5.1, 5.2

Phase 6: Error handling (defensive)
    └─ 6.1, 6.2, 6.3

Phase 7: Testing (ongoing)
    └─ 7.1, 7.2
```

Each phase should be committed separately to keep diffs focused and reviewable.
