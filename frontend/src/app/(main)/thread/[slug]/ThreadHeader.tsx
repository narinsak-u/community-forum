import Link from "next/link";
import { Pencil } from "lucide-react";
import { Badge } from "@/components/ui/badge";

interface ThreadHeaderProps {
  title?: string;
  tagName?: string;
  authorUsername?: string;
  authorInitials?: string;
  createdAt?: string;
  viewCount: number;
  repliesCount: number;
  upvotes: number;
  downvotes: number;
  slug?: string;
  authorId?: number;
  currentUserId?: number;
}

export function ThreadHeader({
  title,
  tagName,
  authorUsername,
  authorInitials,
  createdAt,
  viewCount,
  repliesCount,
  upvotes,
  downvotes,
  slug,
  authorId,
  currentUserId,
}: ThreadHeaderProps) {
  const trustIndex =
    upvotes + downvotes > 0
      ? Math.round((upvotes / (upvotes + downvotes)) * 100) + "%"
      : "98%";

  const formatCount = (n: number) => {
    if (n >= 1000) return (n / 1000).toFixed(1) + "K";
    return String(n);
  };

  return (
    <>
      {/* breadcrumbs */}
      <nav className="flex items-center gap-2 text-[10px] uppercase tracking-[0.2em] text-muted-foreground">
        <span>Discussion</span>
        <span>/</span>
        <span>{tagName || "Technical"}</span>
        <span>/</span>
        <span className="text-primary">
          {title?.slice(0, 24) || "Thread"}...
        </span>
        {slug && authorId && currentUserId === authorId && (
          <Link
            href={`/thread/${slug}/edit`}
            className="ml-auto flex items-center gap-1 text-primary hover:text-primary-glow transition-colors"
          >
            <Pencil className="h-3 w-3" />
            EDIT
          </Link>
        )}
      </nav>

      <header className="space-y-4">
        {tagName && (
          <Badge className="bg-primary/15 text-primary border border-primary/40 rounded-sm uppercase text-[10px] tracking-[0.2em]">
            ● {tagName}
          </Badge>
        )}
        <h1 className="heading-display text-3xl md:text-5xl leading-tight text-foreground">
          {title || "Thread"}
        </h1>
        <div className="flex flex-wrap items-center justify-between gap-4 pt-2">
          <div className="flex items-center gap-3">
            <div className="h-10 w-10 rounded-full bg-gradient-signal grid place-items-center text-xs font-bold text-primary-foreground">
              {authorInitials}
            </div>
            <div>
              <div className="text-sm font-semibold text-foreground">
                {authorUsername || "@unknown"}
              </div>
              <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                {createdAt ? new Date(createdAt).toLocaleDateString() : ""}
              </div>
            </div>
          </div>
          <div className="flex gap-8 text-right">
            <div>
              <div className="text-lg font-display font-bold text-foreground">
                {formatCount(viewCount)}
              </div>
              <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                VIEWS
              </div>
            </div>
            <div>
              <div className="text-lg font-display font-bold text-foreground">
                {repliesCount}
              </div>
              <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                REPLIES
              </div>
            </div>
            <div>
              <div className="text-lg font-display font-bold text-primary">
                {trustIndex}
              </div>
              <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                TRUST_INDEX
              </div>
            </div>
          </div>
        </div>
      </header>
    </>
  );
}
