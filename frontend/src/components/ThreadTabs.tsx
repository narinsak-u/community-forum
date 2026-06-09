"use client";

import { Filter, LayoutGrid } from "lucide-react";

export interface Tab {
  label: string;
  value: string;
  sort: string;
}

interface ThreadTabsProps {
  tabs: Tab[];
  activeTab: string;
  onTabChange: (tab: Tab) => void;
}

export function ThreadTabs({ tabs, activeTab, onTabChange }: ThreadTabsProps) {
  return (
    <div className="flex items-center justify-between border-b border-border/60">
      <div className="flex gap-6">
        {tabs.map((tab) => (
          <button
            key={tab.value}
            onClick={() => onTabChange(tab)}
            className={`pb-3 text-xs uppercase tracking-[0.18em] transition-colors border-b-2 ${
              tab.value === activeTab
                ? "text-primary border-primary"
                : "text-muted-foreground border-transparent hover:text-foreground"
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>
      <div className="flex gap-2 pb-2">
        <button className="h-8 w-8 grid place-items-center bg-secondary/60 border border-border rounded-sm hover:border-primary/40 transition-colors">
          <Filter className="h-3.5 w-3.5 text-muted-foreground" />
        </button>
        <button className="h-8 w-8 grid place-items-center bg-secondary/60 border border-border rounded-sm hover:border-primary/40 transition-colors">
          <LayoutGrid className="h-3.5 w-3.5 text-muted-foreground" />
        </button>
      </div>
    </div>
  );
}
