"use client";

import { useState, useEffect } from "react";
import { ActiveLink } from "@/components/ActiveLink";
import {
  LayoutGrid,
  MessagesSquare,
  Shapes,
  BookOpen,
  Plus,
  ChevronLeft,
  ChevronRight,
  LogIn,
  LogOut,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import { useMe, useSignout } from "@/hooks/use-auth";

const STORAGE_KEY = "midnight-forge-sidebar";

const items = [
  { to: "/threads", label: "Dashboard", icon: LayoutGrid },
  { to: "/discussions", label: "Discussions", icon: MessagesSquare },
  { to: "/profile", label: "Categories", icon: Shapes },
  { to: "/settings", label: "Documentation", icon: BookOpen },
];

export const Sidebar = ({ showNewEntry = false }: { showNewEntry?: boolean }) => {
  const [collapsed, setCollapsed] = useState(false);

  useEffect(() => {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored === "true") setCollapsed(true);
  }, []);

  const toggle = () => {
    setCollapsed((prev) => {
      const next = !prev;
      localStorage.setItem(STORAGE_KEY, String(next));
      return next;
    });
  };

  const { data: user } = useMe();
  const isAuthenticated = !!user;
  const signout = useSignout();

  const handleSignout = () => {
    signout.mutate(undefined, {
      onSuccess: () => window.location.href = "/threads",
    });
  };

  return (
    <aside
      className={cn(
        "shrink-0 border-r border-border/60 bg-sidebar/40 min-h-[calc(100vh-4rem)] flex flex-col transition-all duration-300",
        collapsed ? "w-16" : "w-64",
      )}
    >
      <div className="flex items-center h-14 px-3 border-b border-border/60">
        {/*{!collapsed && (
          <span className="text-xs uppercase tracking-[0.2em] text-primary font-bold">
            NAV
          </span>
        )}*/}
        <button
          type="button"
          onClick={toggle}
          className={cn(
            "h-7 w-7 grid place-items-center text-muted-foreground hover:text-primary rounded-sm hover:bg-sidebar-accent transition-colors",
            collapsed ? "mx-auto" : "ml-auto",
          )}
          aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
        >
          {collapsed ? (
            <ChevronRight className="h-4 w-4" />
          ) : (
            <ChevronLeft className="h-4 w-4" />
          )}
        </button>
      </div>

      <nav className="p-3 space-y-1 flex-1">
        {items.map((item) => {
          const link = (
            <ActiveLink
              key={item.label}
              href={item.to}
              end={item.to === "/"}
              className={cn(
                "flex items-center gap-3 px-3 py-2.5 text-sm text-sidebar-foreground rounded-sm hover:bg-sidebar-accent transition-colors group",
                collapsed && "justify-center px-0",
              )}
              activeClassName="!bg-secondary !text-primary border-l-2 border-primary"
            >
              <item.icon className="h-4 w-4 shrink-0" />
              {!collapsed && <span>{item.label}</span>}
            </ActiveLink>
          );

          if (collapsed) {
            return (
              <Tooltip key={item.label} delayDuration={300}>
                <TooltipTrigger asChild>{link}</TooltipTrigger>
                <TooltipContent side="right" className="text-xs">
                  {item.label}
                </TooltipContent>
              </Tooltip>
            );
          }
          return link;
        })}

        {showNewEntry && (
          <div className={cn("pt-4", collapsed && "pt-3")}>
            <ActiveLink href="/create">
              <Button
                className={cn(
                  "w-full bg-gradient-signal hover:opacity-90 text-primary-foreground font-bold uppercase tracking-[0.15em] text-xs h-10 rounded-sm",
                  collapsed && "h-10 w-10 p-0 mx-auto",
                )}
              >
                <Plus className={cn("h-4 w-4", !collapsed && "mr-1")} />
                {!collapsed && "New_Entry"}
              </Button>
            </ActiveLink>
          </div>
        )}
      </nav>

      <div className="p-3 border-t border-border/60 space-y-1">
        {isAuthenticated && user ? (
          <>
            <ActiveLink
              href={"/profile/" + user.username}
              className={cn(
                "flex items-center gap-3 px-3 py-2 text-sm text-sidebar-foreground hover:text-foreground rounded-sm transition-colors",
                collapsed && "justify-center px-0",
              )}
            >
              <div className="h-6 w-6 rounded-sm bg-secondary border border-border grid place-items-center text-[10px] font-bold text-primary shrink-0">
                {user.username.replace("@", "").slice(0, 2).toUpperCase()}
              </div>
              {!collapsed && (
                <span className="truncate">{user.username}</span>
              )}
            </ActiveLink>
            {collapsed ? (
              <Tooltip delayDuration={300}>
                <TooltipTrigger asChild>
                  <button
                    type="button"
                    onClick={handleSignout}
                    className="flex items-center justify-center px-3 py-2 w-full text-sm text-muted-foreground hover:text-destructive rounded-sm transition-colors"
                  >
                    <LogOut className="h-4 w-4 shrink-0" />
                  </button>
                </TooltipTrigger>
                <TooltipContent side="right" className="text-xs">Log Out</TooltipContent>
              </Tooltip>
            ) : (
              <button
                type="button"
                onClick={handleSignout}
                className="flex items-center gap-3 px-3 py-2 w-full text-sm text-muted-foreground hover:text-destructive rounded-sm transition-colors"
              >
                <LogOut className="h-4 w-4 shrink-0" />
                <span>Log Out</span>
              </button>
            )}
          </>
        ) : (
          <ActiveLink
            href="/login?redirect=/threads"
            className={cn(
              "flex items-center gap-3 px-3 py-2 text-sm text-muted-foreground hover:text-foreground rounded-sm transition-colors",
              collapsed && "justify-center px-0",
            )}
          >
            <LogIn className="h-4 w-4 shrink-0" />
            {!collapsed && <span>Log In</span>}
          </ActiveLink>
        )}
      </div>
    </aside>
  );
};
