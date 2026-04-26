import { NavLink } from "@/components/NavLink";
import { Bell, Settings, Search, User } from "lucide-react";
import { Input } from "@/components/ui/input";

export const TopNav = () => {
  return (
    <header className="sticky top-0 z-40 border-b border-border/60 bg-background/85 backdrop-blur-xl">
      <div className="flex h-16 items-center gap-6 px-6">
        <NavLink to="/" className="flex items-center gap-2 group">
          <div className="h-8 w-8 grid place-items-center bg-gradient-signal rounded-sm font-display font-bold text-primary-foreground">
            M
          </div>
          <span className="font-display font-bold text-lg tracking-wide text-foreground group-hover:text-primary transition-colors">
            MIDNIGHT<span className="text-primary">FORGE</span>
          </span>
        </NavLink>

        <nav className="hidden md:flex items-center gap-1 ml-4">
          {[
            { to: "/", label: "Nexus" },
            { to: "/thread/architectural-shift", label: "Threads" },
            { to: "/profile", label: "Network" },
            { to: "/settings", label: "Terminal" },
          ].map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.to === "/"}
              className="px-3 py-1.5 text-xs uppercase tracking-[0.18em] text-muted-foreground hover:text-foreground transition-colors rounded-sm"
              activeClassName="!text-primary border-b-2 border-primary"
            >
              {item.label}
            </NavLink>
          ))}
        </nav>

        <div className="flex-1 max-w-md ml-auto relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-muted-foreground" />
          <Input
            placeholder="TERMINAL_SEARCH..."
            className="h-9 pl-9 bg-secondary/60 border-border/80 text-xs uppercase tracking-wider placeholder:text-muted-foreground/60 font-mono focus-visible:ring-primary/40"
          />
        </div>

        <div className="flex items-center gap-2">
          <button className="h-9 w-9 grid place-items-center text-muted-foreground hover:text-primary transition-colors relative">
            <Bell className="h-4 w-4" />
            <span className="absolute top-2 right-2 h-1.5 w-1.5 bg-primary rounded-full animate-pulse-signal" />
          </button>
          <NavLink to="/settings" className="h-9 w-9 grid place-items-center text-muted-foreground hover:text-primary transition-colors">
            <Settings className="h-4 w-4" />
          </NavLink>
          <NavLink to="/profile" className="h-9 w-9 grid place-items-center bg-secondary border border-border rounded-sm hover:border-primary transition-colors">
            <User className="h-4 w-4 text-primary" />
          </NavLink>
        </div>
      </div>
    </header>
  );
};
