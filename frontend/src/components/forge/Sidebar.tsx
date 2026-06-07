"use client";

import { ActiveLink } from "@/components/ActiveLink";
import { LayoutGrid, MessagesSquare, Shapes, BookOpen, HelpCircle, Archive, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";

const items = [
  { to: "/", label: "Dashboard", icon: LayoutGrid },
  { to: "/thread/architectural-shift", label: "Discussions", icon: MessagesSquare },
  { to: "/profile", label: "Categories", icon: Shapes },
  { to: "/settings", label: "Documentation", icon: BookOpen },
];

const footerItems = [
  { to: "/settings", label: "Support", icon: HelpCircle },
  { to: "/settings", label: "Archive", icon: Archive },
];

export const Sidebar = ({ showNewEntry = false }: { showNewEntry?: boolean }) => {
  return (
    <aside className="w-64 shrink-0 border-r border-border/60 bg-sidebar/40 min-h-[calc(100vh-4rem)] flex flex-col">
      <div className="p-4 border-b border-border/60">
        <div className="flex items-center gap-3 p-2 rounded-sm bg-secondary/40">
          <div className="h-10 w-10 bg-gradient-signal grid place-items-center rounded-sm">
            <div className="h-4 w-4 bg-background/30 rotate-45" />
          </div>
          <div>
            <div className="text-sm font-bold text-foreground">Midnight Forge</div>
            <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">Technical Forum</div>
          </div>
        </div>
      </div>

      <nav className="p-3 space-y-1 flex-1">
        {items.map((item) => (
          <ActiveLink
            key={item.label}
            href={item.to}
            end={item.to === "/"}
            className="flex items-center gap-3 px-3 py-2.5 text-sm text-sidebar-foreground rounded-sm hover:bg-sidebar-accent transition-colors group"
            activeClassName="!bg-secondary !text-primary border-l-2 border-primary"
          >
            <item.icon className="h-4 w-4" />
            <span>{item.label}</span>
          </ActiveLink>
        ))}

        {showNewEntry && (
          <div className="pt-4">
            <ActiveLink href="/create">
              <Button className="w-full bg-gradient-signal hover:opacity-90 text-primary-foreground font-bold uppercase tracking-[0.15em] text-xs h-10 rounded-sm">
                <Plus className="h-4 w-4 mr-1" /> New_Entry
              </Button>
            </ActiveLink>
          </div>
        )}
      </nav>

      <div className="p-3 border-t border-border/60 space-y-1">
        {footerItems.map((item) => (
          <ActiveLink
            key={item.label}
            href={item.to}
            className="flex items-center gap-3 px-3 py-2 text-sm text-muted-foreground hover:text-foreground rounded-sm transition-colors"
          >
            <item.icon className="h-4 w-4" />
            <span>{item.label}</span>
          </ActiveLink>
        ))}
      </div>
    </aside>
  );
};