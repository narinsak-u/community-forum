import Image from "next/image";
import { useMemo } from "react";
import { SectionLabel } from "@/components/forge/SectionLabel";
import { getInitials } from "@/lib/utils";

interface ParticipantsPanelProps {
  threadAuthor?: { id: number; username: string; avatar: string };
  comments?: { author?: { id: number; username: string; avatar: string }; replies?: any[] }[];
  repliesCount: number;
}

export function ParticipantsPanel({ threadAuthor, comments, repliesCount }: ParticipantsPanelProps) {
  const participants = useMemo(() => {
    const authorMap = new Map<number, { id: number; username: string; avatar: string }>();
    if (threadAuthor) {
      authorMap.set(threadAuthor.id, threadAuthor);
    }
    const collectAuthors = (cs: any[]) => {
      for (const c of cs || []) {
        if (c.author && !authorMap.has(c.author.id)) {
          authorMap.set(c.author.id, c.author);
        }
        collectAuthors(c.replies);
      }
    };
    collectAuthors(comments || []);
    return Array.from(authorMap.values()).slice(0, 5);
  }, [comments, threadAuthor]);

  return (
    <div className="panel p-5 space-y-3">
      <SectionLabel>Participants</SectionLabel>
      <div className="flex flex-wrap gap-1">
        {participants.map((p) => (
          <div
            key={p.id}
            className="h-8 w-8 rounded-full bg-secondary border border-border overflow-hidden"
            title={p.username}
          >
            {p.avatar ? (
              <Image
                src={p.avatar}
                alt={p.username}
                width={32}
                height={32}
                className="object-cover w-full h-full"
              />
            ) : (
              <div className="w-full h-full grid place-items-center text-[10px] font-bold text-primary">
                {getInitials(p.username)}
              </div>
            )}
          </div>
        ))}
        {repliesCount > 0 && (
          <div className="h-8 w-8 rounded-full bg-secondary border border-border grid place-items-center text-[10px] font-mono text-primary">
            +{repliesCount}
          </div>
        )}
      </div>
    </div>
  );
}