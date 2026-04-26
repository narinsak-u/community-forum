import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import { ArrowRight, KeyRound, User } from "lucide-react";

const Login = () => {
  const navigate = useNavigate();
  const [persist, setPersist] = useState(false);

  return (
    <div className="min-h-screen bg-background relative overflow-hidden">
      {/* Background grid + scanline */}
      <div className="absolute inset-0 pointer-events-none">
        <div className="absolute -bottom-32 -left-32 h-96 w-96 border border-primary/10 rotate-45" />
        <div className="absolute top-20 right-10 h-2 w-2 bg-primary rounded-full animate-pulse-signal" />
      </div>

      {/* Header */}
      <header className="relative z-10 flex items-center justify-between px-8 py-6">
        <div className="flex items-center gap-3">
          <div className="h-10 w-10 grid place-items-center bg-gradient-signal rounded-sm font-display font-bold text-xl text-primary-foreground">M</div>
          <div className="font-display font-bold text-xl tracking-wide">
            MIDNIGHT <span className="text-primary">FORGE</span>
          </div>
        </div>
        <div className="hidden md:flex items-center gap-8 text-[10px] uppercase tracking-[0.2em] text-muted-foreground font-mono">
          <span className="flex items-center gap-2">
            <span className="h-1.5 w-1.5 rounded-full bg-success animate-pulse" />
            SYSTEM: <span className="text-foreground">ONLINE</span>
          </span>
          <span>AUTH_PROTOCOL: <span className="text-primary">V4.0.2</span></span>
        </div>
      </header>

      {/* Side meta */}
      <div className="hidden lg:block absolute top-32 right-10 text-right text-[10px] font-mono text-muted-foreground/70 leading-relaxed tracking-wider">
        <div>// KERNEL_TYPE: <span className="text-foreground/80">MONOLITHIC</span></div>
        <div>// ENCRYPTION: <span className="text-foreground/80">AES-256-GCM</span></div>
        <div>// PROTOCOL: <span className="text-foreground/80">FORGE_STANDARD_V4</span></div>
      </div>

      {/* Form */}
      <main className="relative z-10 grid place-items-center px-4 py-16 min-h-[calc(100vh-200px)]">
        <div className="panel scanline relative w-full max-w-md p-8 md:p-10 animate-fade-up">
          {/* Corner brackets */}
          <span className="absolute top-2 right-2 h-4 w-4 border-t border-r border-primary/40" />
          <span className="absolute bottom-2 left-2 h-4 w-4 border-b border-l border-primary/40" />

          <div className="space-y-2 mb-8">
            <h1 className="heading-display text-3xl md:text-4xl text-foreground">TERMINAL ACCESS</h1>
            <p className="text-[11px] uppercase tracking-[0.18em] text-muted-foreground">
              SECURE_LINK_ESTABLISHED <span className="text-primary">//</span> ENTER_CREDENTIALS
            </p>
          </div>

          <form
            className="space-y-6"
            onSubmit={(e) => { e.preventDefault(); navigate("/"); }}
          >
            <div className="space-y-2">
              <label className="terminal-label block">USER_IDENTITY</label>
              <div className="relative">
                <User className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-primary/70" />
                <Input
                  placeholder="Architect_Name"
                  className="pl-10 h-11 bg-transparent border-0 border-b border-border/80 rounded-none focus-visible:ring-0 focus-visible:border-primary text-foreground"
                />
              </div>
            </div>

            <div className="space-y-2">
              <label className="terminal-label block">ACCESS_CODE</label>
              <div className="relative">
                <KeyRound className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-primary/70" />
                <Input
                  type="password"
                  placeholder="••••••••••••"
                  className="pl-10 h-11 bg-transparent border-0 border-b border-border/80 rounded-none focus-visible:ring-0 focus-visible:border-primary text-foreground"
                />
              </div>
            </div>

            <div className="flex items-center justify-between">
              <label className="flex items-center gap-2 cursor-pointer">
                <Checkbox checked={persist} onCheckedChange={(v) => setPersist(!!v)} className="rounded-none border-border data-[state=checked]:bg-primary data-[state=checked]:border-primary" />
                <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">RECOGNITION_PERSISTENCE</span>
              </label>
              <button type="button" className="text-[10px] uppercase tracking-[0.18em] text-primary hover:text-primary-glow">
                LOST_KEY?
              </button>
            </div>

            <Button
              type="submit"
              className="w-full h-12 bg-gradient-signal hover:opacity-90 text-primary-foreground font-bold uppercase tracking-[0.2em] text-xs rounded-sm group"
            >
              INITIATE_SESSION <ArrowRight className="h-4 w-4 ml-2 group-hover:translate-x-1 transition-transform" />
            </Button>

            <div className="pt-4 border-t border-border/60">
              <p className="text-[10px] font-mono">
                <span className="text-primary font-bold">LOG:</span>{" "}
                <span className="text-muted-foreground">LISTENING FOR HANDSHAKE ON PORT 8080</span>
                <span className="inline-block w-2 h-3 bg-primary ml-1 animate-blink align-middle" />
              </p>
            </div>
          </form>
        </div>

        <div className="mt-10 max-w-md w-full grid grid-cols-2 gap-4 text-[10px] uppercase tracking-[0.15em] text-muted-foreground/70 font-mono">
          <div>
            AUTHORIZED ACCESS ONLY. UNAUTHORIZED ENTRIES ARE LOGGED AND REPORTED TO NODE ADMINISTRATORS.
          </div>
          <div className="text-right">
            FORGE_REF: <span className="text-foreground/80">MNF-0912</span>
          </div>
        </div>
      </main>

      <footer className="relative z-10 text-center pb-8 space-y-3">
        <div className="flex justify-center gap-8 text-[11px] uppercase tracking-[0.18em] text-muted-foreground">
          <span>Protocol</span>
          <span>Manifesto</span>
          <span>Support</span>
        </div>
        <div className="text-[10px] uppercase tracking-[0.2em] text-muted-foreground/60">
          © 2026 Midnight Forge // Encrypted Session
        </div>
      </footer>
    </div>
  );
};

export default Login;
