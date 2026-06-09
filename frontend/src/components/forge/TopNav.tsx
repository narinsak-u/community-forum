"use client";

import { useEffect, useState } from "react";
import { ActiveLink } from "@/components/ActiveLink";
import { Bell, Settings, Search, User, Sun, Moon } from "lucide-react";
import { Input } from "@/components/ui/input";
import { useTheme } from "@/components/theme-provider";
import { useMe } from "@/hooks/use-auth";

const navItems = [
  { to: "/threads", label: "Threads" },
  { to: "/network", label: "Network" },
];

export const TopNav = () => {
  const [mounted, setMounted] = useState(false);
  const [searchOpen, setSearchOpen] = useState(false);
  const { theme, setTheme } = useTheme();
  const { data: currentUser } = useMe();

  useEffect(() => setMounted(true), []);

  return (
    <header className="sticky top-0 z-40 border-b border-border/60 bg-background/85 backdrop-blur-xl">
      <div className="flex h-16 items-center gap-6 px-3">
        {/* logo */}
        <ActiveLink href="/" className="flex items-center gap-2 group">
          <div className="w-56">
            <div className="flex items-center mx-auto gap-3 py-2 rounded-sm ">
              <div className="h-10 w-10 bg-gradient-signal grid place-items-center rounded-sm">
                <div className="h-4 w-4 bg-background/30 rotate-45" />
              </div>
              <div>
                <div className="text-sm font-bold text-foreground uppercase">
                  Terminal
                </div>
                <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                  Technical Forum
                </div>
              </div>
            </div>
          </div>
        </ActiveLink>

        <nav className="hidden md:flex items-center gap-1">
          {navItems.map((item) => (
            <ActiveLink
              key={item.to}
              href={item.to}
              end={item.to === "/"}
              className="px-3 py-1.5 text-xs uppercase tracking-[0.18em] text-muted-foreground hover:text-foreground transition-colors rounded-sm"
              activeClassName="!text-primary border-b-2 border-primary"
            >
              {item.label}
            </ActiveLink>
          ))}
        </nav>

        <div className="ml-auto relative flex items-center">
          <button
            type="button"
            onClick={() => setSearchOpen((o) => !o)}
            className="h-9 w-9 grid place-items-center text-muted-foreground hover:text-primary transition-colors shrink-0"
          >
            <Search className="h-4 w-4" />
          </button>
          <div
            className="overflow-hidden transition-all duration-300 ease-in-out"
            style={{ width: searchOpen ? 280 : 0, opacity: searchOpen ? 1 : 0 }}
          >
            <Input
              placeholder="TERMINAL_SEARCH..."
              className="h-9 bg-secondary/60 border-border/80 text-xs uppercase tracking-wider placeholder:text-muted-foreground/60 font-mono focus-visible:ring-primary/40 ml-2"
              autoFocus={searchOpen}
              onBlur={() => setSearchOpen(false)}
            />
          </div>
        </div>

        <div className="flex items-center gap-2">
          <button
            type="button"
            className="h-9 w-9 grid place-items-center text-muted-foreground hover:text-primary transition-colors relative"
          >
            <Bell className="h-4 w-4" />
            <span className="absolute top-2 right-2 h-1.5 w-1.5 bg-primary rounded-full animate-pulse-signal" />
          </button>
          <ActiveLink
            href="/settings"
            className="h-9 w-9 grid place-items-center text-muted-foreground hover:text-primary transition-colors"
          >
            <Settings className="h-4 w-4" />
          </ActiveLink>
          <ActiveLink
            href={currentUser?.username ? "/profile/" + currentUser.username : "/login?redirect=/profile"}
            className="h-9 w-9 grid place-items-center bg-secondary border border-border rounded-sm hover:border-primary transition-colors"
          >
            <User className="h-4 w-4 text-primary" />
          </ActiveLink>

          <button
            type="button"
            onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
            className="h-9 w-9 grid place-items-center text-muted-foreground hover:text-primary transition-colors"
            aria-label="Toggle theme"
          >
            {mounted ? (
              theme === "dark" ? (
                <Sun className="h-4 w-4" />
              ) : (
                <Moon className="h-4 w-4" />
              )
            ) : (
              <div className="h-4 w-4" />
            )}
          </button>
        </div>
      </div>
    </header>
  );
};
