"use client";

import { Suspense, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import {
  ArrowRight,
  KeyRound,
  User,
  Mail,
  ShieldCheck,
  Fingerprint,
} from "lucide-react";
import { toast } from "sonner";
import { useSignin, useSignup } from "@/hooks/use-auth";

type Mode = "signin" | "signup";

function LoginForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [mode, setMode] = useState<Mode>(
    searchParams.get("mode") === "signup" ? "signup" : "signin",
  );
  const [persist, setPersist] = useState(false);
  const [agree, setAgree] = useState(false);
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const signin = useSignin();
  const signup = useSignup();

  const switchMode = (next: Mode) => {
    setMode(next);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (mode === "signup") {
      if (!agree) {
        toast.error("PROTOCOL_ACK required", {
          description: "Acknowledge the Forge Protocol to continue.",
        });
        return;
      }
      if (password !== confirmPassword) {
        toast.error("KEY_MISMATCH", {
          description: "Access codes do not match.",
        });
        return;
      }

      signup.mutate(
        { username, email, password },
        {
          onSuccess: () => {
            toast.success("IDENTITY_FORGED", {
              description: "Identity registered. You may now sign in.",
            });
            switchMode("signin");
          },
          onError: (err) => {
            toast.error("FORGE_FAILED", { description: err.message });
          },
        },
      );
    } else {
      signin.mutate(
        { login: email, password },
        {
          onSuccess: () => {
            toast.success("SESSION_INITIATED", {
              description: "Handshake complete. Routing to nexus...",
            });
            setTimeout(() => router.push("/"), 600);
          },
          onError: (err) => {
            toast.error("AUTH_FAILED", { description: err.message });
          },
        },
      );
    }
  };

  return (
    <div className="min-h-screen bg-background relative overflow-hidden">
      {/* Background decoration */}
      <div className="absolute inset-0 pointer-events-none">
        <div className="absolute -bottom-32 -left-32 h-96 w-96 border border-primary/10 rotate-45" />
        <div className="absolute top-20 right-10 h-2 w-2 bg-primary rounded-full animate-pulse-signal" />
        <div className="absolute inset-0 opacity-[0.04] [background-image:linear-gradient(hsl(var(--primary))_1px,transparent_1px),linear-gradient(90deg,hsl(var(--primary))_1px,transparent_1px)] [background-size:48px_48px]" />
      </div>

      {/* Header */}
      <header className="relative z-10 flex items-center justify-between px-8 py-6">
        <button
          onClick={() => router.push("/")}
          className="flex items-center gap-3 group"
        >
          <div className="h-10 w-10 grid place-items-center bg-gradient-signal rounded-sm font-display font-bold text-xl text-primary-foreground transition-transform group-hover:-rotate-6">
            M
          </div>
          <div className="font-display font-bold text-xl tracking-wide">
            MIDNIGHT <span className="text-primary">FORGE</span>
          </div>
        </button>
        <div className="hidden md:flex items-center gap-8 text-[10px] uppercase tracking-[0.2em] text-muted-foreground font-mono">
          <span className="flex items-center gap-2">
            <span className="h-1.5 w-1.5 rounded-full bg-success animate-pulse" />
            SYSTEM: <span className="text-foreground">ONLINE</span>
          </span>
          <span>
            AUTH_PROTOCOL: <span className="text-primary">V4.0.2</span>
          </span>
        </div>
      </header>

      {/* Side meta */}
      <div className="hidden lg:block absolute top-32 right-10 text-right text-[10px] font-mono text-muted-foreground/70 leading-relaxed tracking-wider">
        <div>
          // KERNEL_TYPE: <span className="text-foreground/80">MONOLITHIC</span>
        </div>
        <div>
          // ENCRYPTION: <span className="text-foreground/80">AES-256-GCM</span>
        </div>
        <div>
          // PROTOCOL:{" "}
          <span className="text-foreground/80">FORGE_STANDARD_V4</span>
        </div>
      </div>

      <main className="relative z-10 grid place-items-center px-4 py-10 min-h-[calc(100vh-180px)]">
        <div className="panel scanline relative w-full max-w-md p-8 md:p-10 animate-fade-up">
          {/* Corner brackets */}
          <span className="absolute top-2 right-2 h-4 w-4 border-t border-r border-primary/40" />
          <span className="absolute bottom-2 left-2 h-4 w-4 border-b border-l border-primary/40" />
          <span className="absolute top-2 left-2 h-4 w-4 border-t border-l border-primary/40" />
          <span className="absolute bottom-2 right-2 h-4 w-4 border-b border-r border-primary/40" />

          {/* Mode switcher */}
          <div className="grid grid-cols-2 mb-8 border border-border/80 rounded-sm overflow-hidden font-mono">
            <button
              type="button"
              onClick={() => switchMode("signin")}
              className={`relative h-10 text-[10px] uppercase tracking-[0.22em] transition-colors ${
                mode === "signin"
                  ? "bg-gradient-signal text-primary-foreground font-bold"
                  : "text-muted-foreground hover:text-foreground"
              }`}
            >
              <Fingerprint className="inline h-3 w-3 mr-1.5 -mt-0.5" />
              SIGN_IN
            </button>
            <button
              type="button"
              onClick={() => switchMode("signup")}
              className={`relative h-10 text-[10px] uppercase tracking-[0.22em] transition-colors border-l border-border/80 ${
                mode === "signup"
                  ? "bg-gradient-signal text-primary-foreground font-bold"
                  : "text-muted-foreground hover:text-foreground"
              }`}
            >
              <ShieldCheck className="inline h-3 w-3 mr-1.5 -mt-0.5" />
              SIGN_UP
            </button>
          </div>

          <div className="space-y-2 mb-7">
            <h1 className="heading-display text-3xl md:text-4xl text-foreground">
              {mode === "signin" ? "TERMINAL ACCESS" : "FORGE IDENTITY"}
            </h1>
            <p className="text-[11px] uppercase tracking-[0.18em] text-muted-foreground">
              {mode === "signin" ? (
                <>
                  SECURE_LINK_ESTABLISHED{" "}
                  <span className="text-primary">//</span> ENTER_CREDENTIALS
                </>
              ) : (
                <>
                  NEW_ARCHITECT_REGISTRATION{" "}
                  <span className="text-primary">//</span> CONFIGURE_KEYS
                </>
              )}
            </p>
          </div>

          <form className="space-y-5" onSubmit={handleSubmit}>
            <div className="space-y-2">
              {mode === "signup" && (
                <div className="animate-fade-up">
                  <label className="terminal-label block">USER_IDENTITY</label>
                  <div className="relative">
                    <User className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-primary/70" />
                    <Input
                      required
                      value={username}
                      onChange={(e) => setUsername(e.target.value)}
                      placeholder="Architect_Name"
                      className="pl-10 h-11 bg-transparent border-0 border-b border-border/80 rounded-none focus-visible:ring-0 focus-visible:border-primary text-foreground"
                    />
                  </div>
                </div>
              )}
            </div>

            <div className="space-y-2 animate-fade-up">
              <label className="terminal-label block">RELAY_ADDRESS</label>
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-primary/70" />
                <Input
                  required
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="architect@forge.net"
                  className="pl-10 h-11 bg-transparent border-0 border-b border-border/80 rounded-none focus-visible:ring-0 focus-visible:border-primary text-foreground"
                />
              </div>
            </div>

            <div className="space-y-2 animate-fade-up">
              <label className="terminal-label block">ACCESS_CODE</label>
              <div className="relative">
                <KeyRound className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-primary/70" />
                <Input
                  required
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="••••••••••••"
                  className="pl-10 h-11 bg-transparent border-0 border-b border-border/80 rounded-none focus-visible:ring-0 focus-visible:border-primary text-foreground"
                />
              </div>
              {mode === "signup" && (
                <div className="flex gap-1 pt-1">
                  <span className="h-0.5 flex-1 bg-primary/80" />
                  <span className="h-0.5 flex-1 bg-primary/50" />
                  <span className="h-0.5 flex-1 bg-border" />
                  <span className="h-0.5 flex-1 bg-border" />
                </div>
              )}
            </div>

            {mode === "signup" && (
              <div className="space-y-2 animate-fade-up">
                <label className="terminal-label block">CONFIRM_CODE</label>
                <div className="relative">
                  <KeyRound className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-primary/70" />
                  <Input
                    required
                    type="password"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    placeholder="••••••••••••"
                    className="pl-10 h-11 bg-transparent border-0 border-b border-border/80 rounded-none focus-visible:ring-0 focus-visible:border-primary text-foreground"
                  />
                </div>
              </div>
            )}

            {mode === "signin" ? (
              <div className="flex items-center justify-between">
                <label className="flex items-center gap-2 cursor-pointer">
                  <Checkbox
                    checked={persist}
                    onCheckedChange={(v) => setPersist(!!v)}
                    className="rounded-none border-border data-[state=checked]:bg-primary data-[state=checked]:border-primary"
                  />
                  <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                    RECOGNITION_PERSISTENCE
                  </span>
                </label>
                <button
                  type="button"
                  className="text-[10px] uppercase tracking-[0.18em] text-primary hover:text-primary-glow"
                >
                  LOST_KEY?
                </button>
              </div>
            ) : (
              <label className="flex items-start gap-2 cursor-pointer">
                <Checkbox
                  checked={agree}
                  onCheckedChange={(v) => setAgree(!!v)}
                  className="mt-0.5 rounded-none border-border data-[state=checked]:bg-primary data-[state=checked]:border-primary"
                />
                <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground leading-relaxed">
                  PROTOCOL_ACK <span className="text-primary">//</span> I accept
                  the Forge Manifesto and Encryption Terms.
                </span>
              </label>
            )}

            <Button
              type="submit"
              disabled={signin.isPending || signup.isPending}
              className="w-full h-12 bg-gradient-signal hover:opacity-90 text-primary-foreground font-bold uppercase tracking-[0.2em] text-xs rounded-sm group"
            >
              {signin.isPending || signup.isPending
                ? "PROCESSING..."
                : mode === "signin"
                  ? "INITIATE_SESSION"
                  : "FORGE_IDENTITY"}
              <ArrowRight className="h-4 w-4 ml-2 group-hover:translate-x-1 transition-transform" />
            </Button>

            {/* Divider */}
            <div className="relative py-1">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-border/60" />
              </div>
              <div className="relative flex justify-center">
                <span className="bg-card px-3 text-[10px] uppercase tracking-[0.22em] font-mono text-muted-foreground">
                  ALT_CHANNEL
                </span>
              </div>
            </div>

            <button
              type="button"
              onClick={() => switchMode("signup")}
              className="w-full h-11 border border-border/80 hover:border-primary/60 text-foreground font-mono text-[11px] uppercase tracking-[0.18em] rounded-sm transition-colors flex items-center justify-center gap-3"
            >
              <span className="h-2 w-2 rounded-full bg-success animate-pulse" />
              FORGE_NEW_IDENTITY
            </button>

            <div className="pt-3 border-t border-border/60">
              <p className="text-[10px] font-mono">
                <span className="text-primary font-bold">LOG:</span>{" "}
                <span className="text-muted-foreground">
                  {mode === "signin"
                    ? "LISTENING FOR HANDSHAKE ON PORT 8080"
                    : "AWAITING IDENTITY PAYLOAD // CHECKSUM PENDING"}
                </span>
                <span className="inline-block w-2 h-3 bg-primary ml-1 animate-blink align-middle" />
              </p>
            </div>

            <p className="text-center text-[10px] font-mono uppercase tracking-[0.18em] text-muted-foreground pt-1">
              {mode === "signin" ? (
                <>
                  NO_IDENTITY?{" "}
                  <button
                    type="button"
                    onClick={() => switchMode("signup")}
                    className="text-primary hover:text-primary-glow"
                  >
                    FORGE_ONE →
                  </button>
                </>
              ) : (
                <>
                  IDENTITY_EXISTS?{" "}
                  <button
                    type="button"
                    onClick={() => switchMode("signin")}
                    className="text-primary hover:text-primary-glow"
                  >
                    ← RETURN_TO_TERMINAL
                  </button>
                </>
              )}
            </p>
          </form>
        </div>

        <div className="mt-8 max-w-md w-full grid grid-cols-2 gap-4 text-[10px] uppercase tracking-[0.15em] text-muted-foreground/70 font-mono">
          <div>
            AUTHORIZED ACCESS ONLY. UNAUTHORIZED ENTRIES ARE LOGGED AND REPORTED
            TO NODE ADMINISTRATORS.
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
}

export default function LoginPage() {
  return (
    <Suspense fallback={<div className="min-h-screen bg-background" />}>
      <LoginForm />
    </Suspense>
  );
}
