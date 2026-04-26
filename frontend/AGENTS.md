# Agent Guidelines for Community Forum

## Build, Lint, and Test Commands

### Development
```bash
bun run dev          # Start dev server on port 8080
bun run preview      # Preview production build
```

### Build
```bash
bun run build        # Production build to dist/
bun run build:dev   # Development build
```

### Linting
```bash
bun run lint         # Run ESLint on all files
```

### Testing (Vitest)
```bash
bun test             # Run all tests once
bun run test:watch   # Run tests in watch mode

# Run a single test file
bunx vitest run src/test/example.test.ts

# Run tests matching a pattern
bunx vitest run --grep "example"

# Run tests in a specific directory
bunx vitest run src/test/
```

## Code Style Guidelines

### General
- This is a React 18 + TypeScript + Vite project with Tailwind CSS and shadcn/ui components
- Use the `@/` path alias for imports (configured in tsconfig.json and vite.config.ts)
- ESLint is configured with TypeScript ESLint, react-hooks, and react-refresh plugins
- TypeScript strict mode is disabled - implicit `any` is allowed

### Imports
- Use absolute imports with `@/` prefix: `import Button from "@/components/ui/button"`
- Group imports: external libs → internal components → local utilities
- Use `type` keyword for type-only imports when possible

### Components
- Use functional components with arrow functions or `function` declarations
- Use React.forwardRef for components that need ref forwarding
- Use CVA (class-variance-authority) for component variants (see button.tsx)
- Export components as named exports with default for page components

### Naming Conventions
- **Components**: PascalCase (e.g., `Button`, `AppLayout`)
- **Hooks**: camelCase with `use` prefix (e.g., `useToast`, `useMobile`)
- **Utilities**: camelCase (e.g., `cn`, `formatDate`)
- **Constants**: UPPER_SNAKE_CASE for true constants, PascalCase for object configs
- **Files**: PascalCase for components, camelCase for utilities/hooks

### Types
- Use TypeScript interfaces for object shapes
- Use type for unions, tuples, and type aliases
- Use `React.FC<Props>` or inline prop types for component typing
- Use `VariantProps<T>` from CVA for variant types

### CSS / Tailwind
- Use Tailwind utility classes primarily
- Use `cn()` from `@/lib/utils` to merge Tailwind classes
- Custom colors are defined in tailwind.config.ts (follow existing patterns)
- Use `panel`, `terminal-label`, `heading-display` custom classes where appropriate

### Error Handling
- Use toast notifications via `toast()` from `@/hooks/use-toast`
- Wrap async operations in try/catch blocks
- Provide user feedback for failures

### React Patterns
- Use TanStack Query (@tanstack/react-query) for data fetching
- Use React Router v6 for routing
- Use controlled forms with react-hook-form + zod validation
- Follow react-hooks rules (exhaustive deps)

### Project Structure
```
src/
├── components/
│   ├── ui/        # shadcn/ui components
│   ├── forge/     # Custom app components
│   └── NavLink.tsx
├── hooks/         # Custom React hooks
├── lib/           # Utilities (utils.ts)
├── pages/         # Route pages (Index, Login, etc.)
├── test/          # Test files and setup
├── assets/        # Static assets
└── App.tsx        # Root app with routing
```

### Testing
- Tests go in `src/test/` directory
- Use Vitest with @testing-library/react
- Test files should be named `*.test.ts` or `*.test.tsx`
- Setup file is `src/test/setup.ts` (includes jest-dom matchers)

### What NOT to Do
- Don't use default exports for components (use named exports)
- Don't use CSS files - use Tailwind classes or inline styles only when needed
- Don't create new component folders at root level - add to existing directories
- Don't use `any` unless necessary (TypeScript is lenient here)