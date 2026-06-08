# Auth Guard Design

**Date:** 2026-06-08
**Status:** Approved

## Overview

Add an authentication guard layer covering three scenarios: protected pages (redirect to `/login`), protected actions (toast with action button), and sidebar auth awareness (user info / login / logout). Auth page is the existing `/login`. No modal.

## `useRequireAuth` hook

A single `useRequireAuth()` hook reads `isAuthenticated` from the Zustand store and returns `{ requireAuth }`. Two modes:

| Mode | Usage | Behavior |
|---|---|---|
| `redirect` | Pages (useEffect) | `router.push("/login?redirect=<path>")` |
| `toast` | Inline action handlers | `toast("SIGN_IN_REQUIRED", { action: { label: "LOGIN", onClick: () => router.push("/login") } })` |

Default mode is `redirect`. Both return `false` when unauthenticated so callers can early-return.

## Protected pages

Redirect unauthenticated visitors to `/login?redirect=<path>`:

- `/profile` — my profile page
- `/settings` — settings page
- `/create` — create entry page

Each page calls `requireAuth({ redirect })` inside a `useEffect` on mount.

## Protected actions

Show a toast with "LOGIN" action button that navigates to `/login`:

- Vote thread (`ThreadDetailClient.tsx` handleVote)
- Vote comment (`ThreadDetailClient.tsx`)
- Submit reply/comment (`ThreadDetailClient.tsx`)

Guard check is inline in the event handler before `mutate()`:

```ts
const handleVote = (value: number) => {
  if (!requireAuth({ toast: "SIGN_IN_REQUIRED", description: "Authenticate to cast a signal." })) return;
  voteThread.mutate({ value });
};
```

The mutation hooks themselves remain stateless — no auth logic inside.

## Sidebar footer

Replaces static `footerItems` with auth-aware content:

**Authenticated:**
- User row: avatar/initials + username → links to `/profile`
- Log Out button → calls `useSignout().mutate()`, redirects to `/`

**Unauthenticated:**
- Log In link/button → navigates to `/login?redirect=/`

Collapsed state: tooltips on all items (existing pattern). User icon links to profile. Log In/Out uses a `LogIn`/`LogOut` icon with tooltip.

## Files changed

| File | Change |
|---|---|
| `frontend/src/hooks/use-require-auth.ts` | New — `useRequireAuth` hook |
| `frontend/src/app/(main)/profile/page.tsx` | Add auth guard useEffect |
| `frontend/src/app/(main)/settings/page.tsx` | Add auth guard useEffect |
| `frontend/src/app/(main)/create/page.tsx` | Add auth guard useEffect |
| `frontend/src/app/(main)/thread/[slug]/ThreadDetailClient.tsx` | Guard vote + comment handlers |
| `frontend/src/components/forge/Sidebar.tsx` | Auth-aware footer items |
