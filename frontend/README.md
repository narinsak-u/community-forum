# The Lands Between — Frontend

Next.js 15 + React 19 frontend for The Lands Between community forum.

## Tech Stack

- **Next.js 15** (App Router)
- **React 19** + TypeScript
- **Tailwind CSS** + shadcn/ui
- **TanStack Query** — Server state management
- **react-hook-form** + zod — Form validation
- **WebSocket** — Real-time chat (via `useChat` hook)

## Project Structure

```
frontend/
├── src/
│   ├── app/                     # Next.js App Router pages
│   │   ├── (main)/              # Main layout (TopNav + Sidebar + Footer)
│   │   │   ├── threads/         # Thread listing
│   │   │   ├── thread/[slug]/   # Thread detail + edit
│   │   │   ├── create/          # New thread
│   │   │   ├── discussions/     # Real-time chat
│   │   │   ├── network/         # User directory
│   │   │   ├── profile/         # User profiles
│   │   │   └── settings/        # Account settings
│   │   ├── (auth)/              # Auth pages (no TopNav/Sidebar)
│   │   │   └── login/
│   │   └── layout.tsx           # Root layout (fonts, providers)
│   ├── components/
│   │   ├── forge/               # App-specific (TopNav, Sidebar, EntryEditor)
│   │   └── ui/                  # shadcn/ui components
│   ├── hooks/                   # Custom hooks (useChat, useAuth, useThread, etc.)
│   ├── stores/                  # Zustand stores (auth)
│   └── lib/                     # Utilities (api, mock-data, server-fetch)
├── public/images/               # Static assets
└── package.json
```

## Commands

```bash
bun install           # Install dependencies
bun run dev           # Dev server (port 3000 or next available)
bun run build         # Production build
bun test              # Run Vitest tests
bun run lint          # ESLint
```

## Key Features

- **Landing page** (`/`) — Hero with video background
- **Thread listing** (`/threads`) — Latest, popular, and user's posts with draft states
- **Thread detail** (`/thread/[slug]`) — Full thread with comments, votes, author actions
- **Thread creation** (`/create`) — Markdown editor with draft/publish workflow
- **Thread editing** (`/thread/[slug]/edit`) — Edit existing threads with tag management
- **Discussions** (`/discussions`) — Real-time group chat via WebSocket with online presence
- **Network** (`/network`) — Browse all registered users
- **User profiles** (`/profile/[username]`) — Bio, tech stacks, contribution history
- **Settings** (`/settings`) — Account security, signal preferences, theme

## Page Layouts

All app pages use the `(main)` layout which includes:

- **TopNav** — Sticky header with logo, nav links, search, notifications, theme toggle, user menu
- **Sidebar** — Collapsible nav with Dashboard, Discussions, Categories, Documentation, New Entry
- **Footer** — Manifesto, Privacy, Security links + copyright

The landing page (`/`) and login page (`/login`) are outside the `(main)` layout — no nav or sidebar.

## API Integration

- **REST:** All HTTP calls go through `src/lib/api.ts` (credentials: include for cookies)
- **WebSocket:** Chat uses `src/hooks/use-chat.ts` connecting to `ws://host/ws/chat`
- **Server-side fetching:** `src/lib/server-fetch.ts` for initial page data
- **Auth:** JWT stored in `forge_token` HTTP-only cookie; Zustand store manages client state

## Theme

Dark theme with primary accent color. Theme preference persisted via `next-themes`. Configurable accent intensity in Settings.
