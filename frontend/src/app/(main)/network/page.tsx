"use client";

import Image from "next/image";
import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { Skeleton } from "@/components/ui/skeleton";
import type { UserProfile } from "@/lib/mock-data";

interface UsersResponse {
  users: UserProfile[];
}

export default function NetworkPage() {
  const { data, isLoading } = useQuery<UsersResponse>({
    queryKey: ["users"],
    queryFn: () => api.get<UsersResponse>("/users"),
  });

  const users = data?.users ?? [];

  return (
    <div className="px-8 py-10 max-w-[1400px] mx-auto space-y-10 animate-fade-up">
      <section className="space-y-4">
        <div className="flex items-end gap-4">
          <h1 className="heading-display text-5xl md:text-6xl text-foreground">
            Network <span className="text-primary">/</span>
          </h1>
        </div>
        <p className="text-sm text-muted-foreground tracking-wide font-mono">
          All nodes connected to The Lands Between
        </p>
      </section>

      <section>
        {isLoading ? (
          <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {Array.from({ length: 6 }).map((_, i) => (
              <Skeleton key={i} className="h-32 rounded-sm" />
            ))}
          </div>
        ) : users.length === 0 ? (
          <p className="text-sm text-muted-foreground text-center py-12">
            No nodes found
          </p>
        ) : (
          <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {users.map((user) => {
              const initials = user.username
                .replace("@", "")
                .slice(0, 2)
                .toUpperCase();
              return (
                <Link
                  key={user.id}
                  href={"/profile/" + user.username}
                  className="panel p-5 flex items-center gap-4 hover:border-primary/40 transition-colors group"
                >
                  <div className="h-12 w-12 rounded-full bg-gradient-signal grid place-items-center text-xs font-bold text-primary-foreground shrink-0 overflow-hidden">
                    {user.avatar ? (
                      <Image
                        src={user.avatar}
                        alt={user.username}
                        width={48}
                        height={48}
                        className="object-cover w-full h-full"
                      />
                    ) : (
                      initials
                    )}
                  </div>
                  <div className="min-w-0 space-y-1">
                    <div className="text-sm font-semibold text-foreground group-hover:text-primary transition-colors truncate">
                      {user.username}
                    </div>
                    <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground truncate">
                      {user.role}
                    </div>
                    {user.bio && (
                      <p className="text-xs text-muted-foreground/70 line-clamp-1">
                        {user.bio}
                      </p>
                    )}
                  </div>
                </Link>
              );
            })}
          </div>
        )}
      </section>
    </div>
  );
}
