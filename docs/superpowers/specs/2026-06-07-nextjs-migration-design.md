# Next.js Migration Design

**Date:** 2026-06-07
**Project:** Community Forum (Midnight Forge)
**Status:** Approved

## Overview

Migrate the frontend from Vite + React 18 to Next.js 15 (App Router) + React 19. The Go Fiber backend remains unchanged.

## Motivation

- SSR for SEO (search engines can index forum threads, profiles)
- Faster initial page loads via server-rendered content
- File-based routing eliminates manual route config
- Ecosystem alignment (Next.js is standard for React apps)

## Approach

**Full migration with hybrid rendering.** Each page evaluated for server vs. client rendering:

- **Server components**: Page-level data fetching (SSR for initial load, SEO)
- **Client components** (`'use client'`): shadcn/ui, TanStack Query hooks, Zustand auth, interactivity (voting, comments, forms, toasts)

Chosen approach (Approach 1) — clean break from Vite, restructure entire frontend as Next.js App Router project.

## Architecture

```
frontend/
├── app/
│   ├── layout.tsx              # Root layout (html, body, Providers)
│   ├── page.tsx                # Home page (server component)
│   ├── not-found.tsx           # 404 page
│   ├── (auth)/
│   │   └── login/
│   │       └── page.tsx        # Login/signup (client component)
│   ├── (main)/
│   │   ├── layout.tsx          # App shell (TopNav + Sidebar + Footer)
│   │   ├── thread/
│   │   │   └── [slug]/
│   │   │       └── page.tsx    # Thread detail (hybrid)
│   │   ├── create/
│   │   │   └── page.tsx        # Create entry (client component)
│   │   ├── profile/
│   │   │   ├── page.tsx        # Own profile (client component)
│   │   │   └── [username]/
│   │   │       └── page.tsx    # User profile (hybrid)
│   │   └── settings/
│   │       └── page.tsx        # Settings (client component)
├── components/
│   ├── forge/                  # Layout components (TopNav, Sidebar, etc.)
│   ├── ui/                     # shadcn/ui (unchanged, 'use client')
│   └── providers.tsx           # QueryClient + Tooltip + Toaster providers
├── hooks/                      # TanStack Query hooks (minor env var changes)
├── lib/
│   ├── api.ts                  # env var changes (NEXT_PUBLIC_API_URL)
│   └── utils.ts                # unchanged
├── stores/
│   └── auth-store.ts           # Zustand (unchanged)
├── public/                     # Static assets (images from src/assets/)
├── app/globals.css             # Moved from src/index.css
├── next.config.ts
└── package.json
```

## Route Mapping

| Current (React Router) | Next.js (App Router) | Component Type |
|---|---|---|
| `/` | `app/page.tsx` | Server |
| `/login` `/signup` | `app/(auth)/login/page.tsx` | Client |
| `/thread/:slug` | `app/thread/[slug]/page.tsx` | Hybrid |
| `/create` | `app/create/page.tsx` | Client |
| `/profile` | `app/profile/page.tsx` | Client |
| `/profile/:username` | `app/profile/[username]/page.tsx` | Hybrid |
| `/settings` | `app/settings/page.tsx` | Client |
| `*` | `app/not-found.tsx` | Static |

## Data Fetching Strategy

### Server components
Pages fetch initial data directly using `fetch()` to the Go API:
- **Home**: Featured thread + trending + paginated threads (parallel fetch)
- **Thread detail**: Thread with comments by slug
- **Profile**: User data + their threads by username

Data is passed as props to client child components for hydration.

### Client components
- TanStack Query hooks remain for: mutations (create thread, vote, comment, auth), stale data refetching, cache invalidation
- Zustand auth store remains unchanged
- Server-fetched data is used as TanStack Query initial data where applicable
- Mock data fallbacks remain in hooks for API-unavailable states

## Layout Architecture

### Route groups
- **`(auth)`** — Minimal layout (no sidebar) for login/signup pages. Guest-only access.
- **`(main)`** — App shell with `TopNav`, `Sidebar`, and `Footer`. Shared across all authenticated pages.

### Providers client component
New `components/providers.tsx` wrapping:
- `QueryClientProvider` (TanStack Query)
- `TooltipProvider` (Radix)
- `Toaster` + `Sonner` (notification components)

Replaces the provider setup currently in `App.tsx`.

## Component Changes

### New / modified
- `components/providers.tsx` — Client component wrapping all providers
- `components/forge/` — TopNav and Sidebar remain but use Next.js `Link` instead of `NavLink`
- `hooks/use-threads.ts` etc. — `import.meta.env` → `process.env` for API URL

### Removed
- `src/main.tsx` — Replaced by `app/layout.tsx`
- `src/App.tsx` — Replaced by route groups + providers
- `src/components/NavLink.tsx` — Replaced by Next.js `Link` + `usePathname`
- `src/components/forge/AppLayout.tsx` — Replaced by `(main)/layout.tsx`

### Unchanged
- All `components/ui/` shadcn components
- `stores/auth-store.ts` (Zustand)
- `lib/utils.ts`
- `lib/mock-data.ts`
- TanStack Query hooks (except env var reference)
- Test files

## Asset Handling

- Images: move `src/assets/*` → `public/images/`
- Google Fonts: replace `@import` in CSS with `next/font` (JetBrains Mono + Space Grotesk)
- Favicon: `public/` directory

## Environment Variables

| Current | New |
|---|---|
| `VITE_API_URL` | `NEXT_PUBLIC_API_URL` |

Default value same: `http://localhost:8080`

## Dependency Changes

### Remove
- `vite`, `@vitejs/plugin-react-swc`
- `react-router-dom`
- `lovable-tagger`
- `postcss` (Next.js includes it)
- `autoprefixer` (Next.js includes it)

### Add
- `next` (v15)
- `@types/react` upgrade for React 19

### Keep
- All `@radix-ui/*` packages
- `@tanstack/react-query`
- `tailwindcss`, `tailwindcss-animate`, `tailwind-merge`, `clsx`, `class-variance-authority`
- `zustand`
- `zod`, `react-hook-form`, `@hookform/resolvers`
- `lucide-react`
- `sonner`
- `date-fns`
- All other existing dependencies

## Build & Dev Commands

| Action | Before | After |
|---|---|---|
| Dev server | `bun run dev` (Vite, port 5713) | `bun run dev` (Next.js, port 8080) |
| Production build | `bun run build` | `bun run build` |
| Lint | `bun run lint` | `bun run lint` |
| Test | `bun test` (Vitest) | `bun test` (Vitest, same config) |

## Testing

Vitest continues as the test runner. Tests in `src/test/` remain unchanged. Test configuration may need minor updates for path resolution.

## Migration Steps (Outline)

1. Initialize Next.js project in `frontend/` directory
2. Restructure `src/` → `app/` directory
3. Create root layout with Providers
4. Create `(auth)` and `(main)` route groups with layouts
5. Convert AppLayout → route group layout
6. Migrate each page: Home (server), Thread detail (hybrid), Profile (hybrid), Login/Create/Settings (client)
7. Replace NavLink with Next.js Link
8. Move static assets to `public/`
9. Configure `next/font` for Google Fonts
10. Update environment variable references
11. Remove Vite-specific config files
12. Configure Vitest for Next.js
13. Verify all routes render and API integration works
