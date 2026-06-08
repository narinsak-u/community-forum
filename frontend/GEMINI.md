# Midnight Forge - Frontend Context

This directory contains the Next.js frontend for the Community Forum application, styled with a terminal-inspired aesthetic called "Midnight Forge".

## Project Overview

- **Framework**: **Next.js 15** (App Router) with **React 19**.
- **Styling**: **Tailwind CSS** with **shadcn/ui** (Radix UI primitives).
- **Icons**: **Lucide React**.
- **State Management**:
  - **Server State**: **TanStack Query (v5)** for data fetching and caching.
  - **Client State**: **Zustand** for lightweight global state (e.g., Auth).
- **Forms**: **React Hook Form** with **Zod** for validation.
- **Testing**: **Vitest** with **React Testing Library**.
- **Runtime/Package Manager**: **Bun**.

## Directory Structure

- `src/app/`: Next.js App Router routes and layouts.
  - `(auth)/`: Authentication-related pages (Login).
  - `(main)/`: Main application content (Home, Thread Detail, Profile, Create).
- `src/components/`:
  - `ui/`: Base UI components from **shadcn/ui**.
  - `forge/`: Custom theme-specific components (Sidebar, TopNav, etc.).
  - `providers.tsx`: Global context providers (QueryClient, Theme, etc.).
- `src/hooks/`: Custom hooks for API interactions (e.g., `use-threads.ts`, `use-auth.ts`).
- `src/lib/`: Utilities and API client configuration (`api.ts`, `server-fetch.ts`).
- `src/stores/`: Zustand store definitions.
- `src/test/`: Vitest configuration and test setup.

## Building and Running

### Prerequisites
- [Bun](https://bun.sh/)

### Commands
- **Install Dependencies**: `bun install`
- **Development Server**: `bun run dev` (Runs on `http://localhost:3000` by default).
- **Production Build**: `bun run build`
- **Start Production Server**: `bun run start`
- **Run Tests**: `bun run test`
- **Linting**: `bun run lint`

## Development Conventions

- **Component Patterns**:
  - Use functional components with TypeScript.
  - Keep logic in hooks (`src/hooks/`) and UI in components.
  - Use `use client` directive only when necessary (interactive components).
- **Data Fetching**:
  - Use **TanStack Query** hooks for most server-side data interactions.
  - Use `@/lib/api` for client-side fetches.
  - Use `@/lib/server-fetch` for Server Component data fetching.
- **Styling**:
  - Strictly follow the terminal aesthetic using Tailwind classes.
  - Fonts used: `JetBrains Mono` (Mono) and `Space Grotesk` (Display).
- **Paths**: Use the `@/` alias to reference the `src/` directory.
- **Validation**: Use **Zod** schemas for all form inputs and API response validation where critical.

## Key Configurations
- `next.config.ts`: Next.js configuration.
- `tailwind.config.ts`: Custom Tailwind theme including "Midnight Forge" colors and fonts.
- `components.json`: shadcn/ui configuration.
- `vitest.config.ts`: Vitest environment setup.
