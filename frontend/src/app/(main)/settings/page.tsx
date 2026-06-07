"use client";

import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { Slider } from "@/components/ui/slider";
import {
  Shield,
  Bell,
  Eye,
  Terminal as TerminalIcon,
  RotateCcw,
} from "lucide-react";
import { toast } from "sonner";
import { useAuthStore } from "@/stores/auth-store";
import { useUpdateProfile } from "@/hooks/use-user";

const Settings = () => {
  const [twofa, setTwofa] = useState(true);
  const [alerts, setAlerts] = useState(true);
  const [direct, setDirect] = useState(true);
  const [digest, setDigest] = useState(false);
  const [accent, setAccent] = useState(0);

  const user = useAuthStore((s) => s.user);
  const updateProfile = useUpdateProfile(user?.username || "");

  const handleCommit = () => {
    updateProfile.mutate(
      {},
      {
        onSuccess: () => toast.success("CONFIG_COMMITTED"),
        onError: (err) =>
          toast.error("COMMIT_FAILED", { description: err.message }),
      },
    );
  };

  return (
    <div className="p-8 lg:p-10 space-y-8 animate-fade-up max-w-[1200px] ">
      <header className="space-y-2">
        <h1 className="heading-display text-5xl italic text-foreground tracking-tight">
          CONTROL PANEL
        </h1>
        <div className="text-[11px] uppercase tracking-[0.2em] text-muted-foreground">
          SYSTEM CONFIGURATION <span className="text-primary">/</span>{" "}
          {user?.username?.replace("@", "").toUpperCase() || "NODE"}
        </div>
      </header>

      {/* Top row: account + privacy */}
      <div className="grid lg:grid-cols-3 gap-6">
        <section className="lg:col-span-2 panel p-6 border-l-2 border-l-primary space-y-5">
          <div className="flex items-center gap-2">
            <Shield className="h-5 w-5 text-primary" />
            <h2 className="text-lg font-bold text-foreground">
              Account Security
            </h2>
          </div>
          <div className="space-y-2">
            <label className="terminal-label">EMAIL ADDRESS</label>
            <Input
              defaultValue={user?.email || "architect@midnight-forge.io"}
              className="bg-terminal border-border font-mono rounded-sm"
              disabled
            />
          </div>
          <div className="space-y-2">
            <label className="terminal-label">PASSWORD</label>
            <div className="grid sm:grid-cols-2 gap-3">
              <Button
                variant="outline"
                className="border-border bg-secondary/40 justify-between text-xs uppercase tracking-[0.18em] text-muted-foreground rounded-sm h-11"
              >
                CHANGE ACCESS KEY <span>→</span>
              </Button>
              <div className="flex items-center justify-between bg-secondary/40 border border-border rounded-sm px-4 h-11">
                <span className="text-xs uppercase tracking-[0.18em] text-foreground">
                  ENABLE 2FA AUTH
                </span>
                <Switch
                  checked={twofa}
                  onCheckedChange={setTwofa}
                  className="data-[state=checked]:bg-primary"
                />
              </div>
            </div>
          </div>
        </section>

        <section className="panel p-6 space-y-5">
          <div className="flex items-center gap-2">
            <Eye className="h-5 w-5 text-primary" />
            <h2 className="text-lg font-bold text-foreground">Privacy</h2>
          </div>
          <p className="text-sm text-foreground/80 leading-relaxed">
            Your node is currently{" "}
            <span className="text-primary font-mono">PUBLIC</span>. Metadata is
            indexed by the Forge.
          </p>
          <div className="space-y-2">
            <Button
              variant="outline"
              className="w-full border-border text-xs uppercase tracking-[0.18em] rounded-sm h-10"
            >
              EXPORT DATA.JSON
            </Button>
            <Button className="w-full bg-gradient-signal text-primary-foreground font-bold uppercase tracking-[0.18em] text-xs h-10 rounded-sm">
              SWITCH TO STEALTH
            </Button>
          </div>
        </section>
      </div>

      {/* Signals */}
      <section className="panel p-6 space-y-5">
        <div className="flex items-center gap-2">
          <Bell className="h-5 w-5 text-primary" />
          <h2 className="text-lg font-bold text-foreground">
            Signal Preferences
          </h2>
        </div>
        <div className="grid md:grid-cols-3 gap-4">
          {[
            {
              label: "System Alerts",
              desc: "Critical forge updates and node maintenance pings.",
              val: alerts,
              set: setAlerts,
            },
            {
              label: "Direct Links",
              desc: "Real-time alerts for incoming private comms.",
              val: direct,
              set: setDirect,
            },
            {
              label: "Daily Digest",
              desc: "Summary of activity across followed forums.",
              val: digest,
              set: setDigest,
            },
          ].map((s) => (
            <div
              key={s.label}
              className="bg-secondary/30 border border-border rounded-sm p-4 space-y-2"
            >
              <div className="flex items-center justify-between">
                <span className="text-sm font-semibold text-foreground">
                  {s.label}
                </span>
                <Switch
                  checked={s.val}
                  onCheckedChange={s.set}
                  className="data-[state=checked]:bg-primary"
                />
              </div>
              <p className="text-xs text-muted-foreground leading-relaxed">
                {s.desc}
              </p>
            </div>
          ))}
        </div>
      </section>

      {/* Appearance + preview */}
      <section className="panel p-6 grid lg:grid-cols-2 gap-8">
        <div className="space-y-6">
          <div className="flex items-center gap-2">
            <TerminalIcon className="h-5 w-5 text-primary" />
            <h2 className="text-lg font-bold text-foreground">
              Forge Appearance
            </h2>
          </div>

          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <span className="terminal-label">TERMINAL FONT SIZE</span>
              <span className="text-sm text-primary font-mono">14PX</span>
            </div>
            <Slider
              defaultValue={[40]}
              max={100}
              step={5}
              className="[&_[role=slider]]:bg-primary [&_[role=slider]]:border-primary"
            />
          </div>

          <div className="space-y-3">
            <span className="terminal-label">UI DENSITY</span>
            <div className="grid grid-cols-3 gap-2">
              {["COMPACT", "DEFAULT", "EDITORIAL"].map((d, i) => (
                <button
                  key={d}
                  className={`py-2.5 text-[10px] uppercase tracking-[0.18em] border rounded-sm ${
                    i === 0
                      ? "border-primary text-primary"
                      : "border-border text-muted-foreground hover:border-primary/40"
                  }`}
                >
                  {d}
                </button>
              ))}
            </div>
          </div>

          <div className="space-y-3">
            <span className="terminal-label">ACCENT INTENSITY</span>
            <div className="flex gap-2">
              {[
                "bg-primary-glow",
                "bg-primary",
                "bg-primary-deep",
                "bg-[hsl(14_50%_30%)]",
              ].map((c, i) => (
                <button
                  key={i}
                  onClick={() => setAccent(i)}
                  className={`h-10 w-12 rounded-sm transition-all ${c} ${accent === i ? "ring-2 ring-primary ring-offset-2 ring-offset-card scale-105" : "opacity-80 hover:opacity-100"}`}
                />
              ))}
            </div>
          </div>
        </div>

        {/* Terminal preview */}
        <div className="bg-terminal border border-border rounded-sm overflow-hidden flex flex-col">
          <div className="flex items-center justify-between px-3 py-2 border-b border-border/60">
            <div className="flex gap-1.5">
              <span className="h-2.5 w-2.5 rounded-full bg-destructive" />
              <span className="h-2.5 w-2.5 rounded-full bg-yellow-500" />
              <span className="h-2.5 w-2.5 rounded-full bg-success" />
            </div>
          </div>
          <div className="p-4 font-mono text-xs space-y-1.5 flex-1 text-foreground/90">
            <div className="text-muted-foreground">
              // PREVIEW_MODE: <span className="text-success">Active</span>
            </div>
            <div>
              <span className="text-primary">root@midnight-forge:</span>
              <span className="text-success">~$</span> update --visuals
            </div>
            <div className="text-muted-foreground">
              Applying neon-noir schemas...
            </div>
            <div>
              Density:{" "}
              <span className="text-primary">[################----] 80%</span>
            </div>
            <div>
              Status: <span className="text-success">SYNCED</span>
            </div>
            <div className="pt-2 flex items-center">
              <span className="inline-block w-2 h-3 bg-primary animate-blink" />
            </div>
          </div>
        </div>
      </section>

      {/* Footer actions */}
      <div className="flex items-center justify-between pt-4 border-t border-border/60">
        <span className="flex items-center gap-2 text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
          <RotateCcw className="h-3 w-3" /> LAST SYNC: 2026.04.26 14:32:01
        </span>
        <div className="flex gap-3">
          <Button
            variant="ghost"
            className="text-muted-foreground hover:text-foreground uppercase text-xs tracking-[0.18em]"
          >
            RESET DEFAULTS
          </Button>
          <Button
            onClick={handleCommit}
            disabled={updateProfile.isPending}
            className="bg-gradient-signal hover:opacity-90 text-primary-foreground font-bold uppercase tracking-[0.18em] text-xs rounded-sm"
          >
            {updateProfile.isPending ? "COMMITTING..." : "COMMIT CHANGES"}
          </Button>
        </div>
      </div>
    </div>
  );
};

export default Settings;
