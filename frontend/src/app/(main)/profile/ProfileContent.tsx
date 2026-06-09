"use client";

import Image from "next/image";
import { SectionLabel } from "@/components/forge/SectionLabel";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { FileText as FileIcon, Globe, Lock, Plus } from "lucide-react";
import { useUserProfile, useUserThreads } from "@/hooks/use-user";
import Link from "next/link";
import { timeAgo } from "@/lib/utils";

const vaultItems = [
  {
    icon: FileIcon,
    tag: "HARDWARE_LOGIC",
    title: "Neural Interface Schema v2",
    desc: "Complete architectural breakdown of the quantum-resistant neural gateway used...",
  },
  {
    icon: Globe,
    tag: "NETWORK_MAP",
    title: "Deep-Web Relay Topography",
    desc: "Visualizing the hidden relays within the forge-network to ensure 99.9% protocol...",
  },
  {
    icon: Lock,
    tag: "SECURITY_AUDIT",
    title: "Vault Entropy Benchmarks",
    desc: "Comparative analysis of random seed generators across the four vault sectors.",
  },
];

interface ProfileContentProps {
  username?: string;
  initialProfile?: any;
}

export function ProfileContent({
  username,
  initialProfile,
}: ProfileContentProps) {
  const { data: profile, isLoading } = useUserProfile(username || "");
  const { data: threadsData, isLoading: threadsLoading } = useUserThreads(
    username || "",
  );

  const currentProfile = profile || initialProfile;

  if (isLoading && !initialProfile) {
    return (
      <div className="p-8 lg:p-10 space-y-10">
        <Skeleton className="h-28 w-full" />
        <Skeleton className="h-48 w-full" />
      </div>
    );
  }

  const displayName = currentProfile?.username || "NODE_8829";
  const role = currentProfile?.role || "SENIOR ARCHITECT";
  const bio = currentProfile?.bio || "No bio available.";
  const stacks = currentProfile?.stacks || ["Rust", "Solidity", "Go"];

  return (
    <div className="p-8 lg:p-10 space-y-10 animate-fade-up">
      {/* Header */}
      <section className="grid md:grid-cols-[120px,1fr,auto] gap-6 items-start">
        <div className="h-28 w-28 panel overflow-hidden relative">
          {currentProfile?.avatar ? (
            <Image
              src={currentProfile.avatar}
              alt="Avatar"
              fill
              className="object-cover"
              sizes="112px"
            />
          ) : (
            <Image
              src="/images/forge-avatar.jpg"
              alt="Node avatar"
              fill
              className="object-cover"
              sizes="112px"
            />
          )}
        </div>
        <div className="space-y-3">
          <h1 className="heading-display text-5xl text-foreground uppercase">
            {displayName}
          </h1>
          <div className="text-xs uppercase tracking-[0.2em] text-primary font-mono">
            {role} // SYSTEMS INTEGRATOR
          </div>
          <div className="flex flex-wrap gap-2 pt-1">
            <Badge className="bg-primary/10 text-primary border border-primary/30 rounded-sm uppercase text-[10px] tracking-[0.18em]">
              ● CORE_CONTRIBUTOR
            </Badge>
            <Badge className="bg-primary/10 text-primary border border-primary/30 rounded-sm uppercase text-[10px] tracking-[0.18em]">
              ● SECURITY_AUDITOR
            </Badge>
          </div>
        </div>
        <div className="panel p-5 grid grid-cols-3 gap-4 min-w-[300px]">
          {[
            { v: "1.2k", l: "CONTRIBUTIONS" },
            { v: "98%", l: "TRUST INDEX" },
            { v: "L_04", l: "ACCESS LEVEL" },
          ].map((s) => (
            <div key={s.l}>
              <div className="text-2xl font-display font-bold text-foreground">
                {s.v}
              </div>
              <div className="text-[9px] uppercase tracking-[0.18em] text-muted-foreground">
                {s.l}
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Two-column body */}
      <section className="grid lg:grid-cols-[1fr,1.4fr] gap-10">
        {/* Identity column */}
        <div className="space-y-8">
          <div className="space-y-4">
            <SectionLabel>SYSTEM IDENTITY</SectionLabel>
            <p className="text-sm text-foreground/85 leading-relaxed">{bio}</p>
          </div>

          <div className="panel p-5 space-y-3">
            <div className="text-[10px] uppercase tracking-[0.18em] text-primary">
              TECHNICAL_STACK
            </div>
            <div className="flex flex-wrap gap-2">
              {stacks.map((stack) => (
                <span
                  key={stack}
                  className="px-2.5 py-1 bg-secondary border border-border text-xs text-foreground rounded-sm"
                >
                  {stack}
                </span>
              ))}
            </div>
          </div>

          <div className="space-y-4">
            <SectionLabel>CONTRIBUTION_STREAM</SectionLabel>
            {threadsLoading ? (
              <div className="space-y-4">
                <Skeleton className="h-12 w-full" />
                <Skeleton className="h-12 w-full" />
              </div>
            ) : threadsData?.threads?.length ? (
              threadsData.threads.map((c) => (
                <div
                  key={c.id}
                  className="flex justify-between items-start gap-4 border-b border-border/60 pb-3"
                >
                  <div className="space-y-1.5">
                    <Link
                      href={`/thread/${c.slug}`}
                      className="text-sm font-semibold text-foreground hover:text-primary transition-colors"
                    >
                      {c.title}
                    </Link>
                    <div className="flex gap-1.5">
                      {c.tags?.map((tg) => (
                        <span
                          key={tg.name}
                          className="px-1.5 py-0.5 bg-secondary border border-border text-[9px] uppercase tracking-[0.18em] text-muted-foreground rounded-sm"
                        >
                          {tg.name}
                        </span>
                      ))}
                    </div>
                  </div>
                  <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground shrink-0">
                    {c.created_at ? timeAgo(c.created_at) : ""}
                  </span>
                </div>
              ))
            ) : (
              <p className="text-sm text-muted-foreground italic">
                No recent contributions found.
              </p>
            )}
          </div>
        </div>

        {/* Vault grid */}
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <SectionLabel>DIGITAL VAULT</SectionLabel>
            <button className="text-[10px] uppercase tracking-[0.18em] text-primary hover:text-primary-glow">
              ACCESS_ALL_ARCHIVES
            </button>
          </div>
          <div className="grid sm:grid-cols-2 gap-4">
            {vaultItems.map((v) => (
              <div
                key={v.title}
                className="panel p-5 space-y-3 hover:border-primary/40 transition-colors group cursor-pointer"
              >
                <div className="aspect-[3/2] bg-terminal/80 border border-border/60 rounded-sm relative overflow-hidden">
                  <div className="absolute inset-0 grid place-items-center text-primary/30 group-hover:text-primary/60 transition-colors">
                    <v.icon className="h-12 w-12" />
                  </div>
                  <div
                    className="absolute inset-0 opacity-30"
                    style={{
                      backgroundImage:
                        "radial-gradient(circle at 30% 30%, hsl(var(--primary) / 0.4), transparent 60%)",
                    }}
                  />
                </div>
                <div className="text-[10px] uppercase tracking-[0.2em] text-primary">
                  {v.tag}
                </div>
                <div className="text-base font-bold text-foreground">
                  {v.title}
                </div>
                <div className="text-xs text-muted-foreground line-clamp-2">
                  {v.desc}
                </div>
              </div>
            ))}
            <button className="panel p-5 grid place-items-center min-h-[280px] border-dashed hover:border-primary/40 transition-colors group">
              <div className="text-center space-y-2">
                <div className="h-12 w-12 mx-auto bg-secondary border border-border rounded-sm grid place-items-center text-muted-foreground group-hover:text-primary group-hover:border-primary/40 transition-colors">
                  <Plus className="h-5 w-5" />
                </div>
                <div className="text-[10px] uppercase tracking-[0.2em] text-muted-foreground">
                  ADD_TO_VAULT
                </div>
              </div>
            </button>
          </div>
        </div>
      </section>
    </div>
  );
}
