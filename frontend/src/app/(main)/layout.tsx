import { TopNav } from "@/components/forge/TopNav";
import { Sidebar } from "@/components/forge/Sidebar";

export default function MainLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-background flex flex-col">
      <TopNav />
      <div className="flex flex-1 w-full">
        <Sidebar showNewEntry />
        <main className="flex-1 min-w-0">{children}</main>
      </div>
      <footer className="border-t border-border/60 mt-auto">
        <div className="px-8 py-6 flex flex-wrap items-center justify-between gap-4 text-[11px] uppercase tracking-[0.18em] text-muted-foreground">
          <div className="flex gap-6">
            <span>Manifesto</span>
            <span>Privacy</span>
            <span>Security</span>
          </div>
          <div>© 2026 Midnight Forge // Encrypted Session</div>
        </div>
      </footer>
    </div>
  );
}