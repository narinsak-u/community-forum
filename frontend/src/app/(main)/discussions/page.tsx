"use client";

import { useEffect, useRef, useState } from "react";
import { useRequireAuth } from "@/hooks/use-require-auth";
import { useMe } from "@/hooks/use-auth";
import { useChat, type ChatMessage, type OnlineUser } from "@/hooks/use-chat";
import { Skeleton } from "@/components/ui/skeleton";
import { Send } from "lucide-react";
import { timeAgo } from "@/lib/utils";

export default function DiscussionsPage() {
  const { requireAuth } = useRequireAuth();
  const { data: currentUser } = useMe();
  const {
    messages,
    onlineUsers,
    sendMessage,
    loadMore,
    isConnected,
    isLoading,
  } = useChat();

  useEffect(() => {
    requireAuth({ redirect: "/discussions" });
  }, [requireAuth]);

  const [input, setInput] = useState("");
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const messagesContainerRef = useRef<HTMLDivElement>(null);
  const [autoScroll, setAutoScroll] = useState(true);

  const handleSend = () => {
    const trimmed = input.trim();
    if (!trimmed) return;
    sendMessage(trimmed);
    setInput("");
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  useEffect(() => {
    if (autoScroll) {
      messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages, autoScroll]);

  const handleScroll = () => {
    const el = messagesContainerRef.current;
    if (!el) return;
    setAutoScroll(el.scrollTop + el.clientHeight >= el.scrollHeight - 100);
    if (el.scrollTop === 0) {
      loadMore();
    }
  };

  return (
    <div className="flex h-[calc(100vh-4rem)]">
      <div className="flex-1 flex flex-col min-w-0">
        <div className="border-b border-border/60 px-6 py-4">
          <h1 className="text-sm font-semibold text-foreground uppercase tracking-[0.12em]">
            # General Chat
          </h1>
        </div>

        <div
          ref={messagesContainerRef}
          onScroll={handleScroll}
          className="flex-1 overflow-y-auto px-6 py-4 space-y-4"
        >
          {isLoading ? (
            <div className="space-y-4 pt-4">
              {Array.from({ length: 5 }).map((_, i) => (
                <div key={i} className="flex gap-3">
                  <Skeleton className="h-8 w-8 rounded-full shrink-0" />
                  <div className="space-y-2 flex-1">
                    <Skeleton className="h-4 w-32" />
                    <Skeleton className="h-10 w-full max-w-md" />
                  </div>
                </div>
              ))}
            </div>
          ) : messages.length === 0 ? (
            <p className="text-sm text-muted-foreground text-center pt-12">
              No messages yet. Start the conversation!
            </p>
          ) : (
            messages.map((msg: ChatMessage) => {
              const isOwn = currentUser && msg.author?.id === currentUser.id;
              const initials = (msg.author?.username || "??")
                .replace("@", "")
                .slice(0, 2)
                .toUpperCase();
              return (
                <div
                  key={msg.id}
                  className={"flex gap-3 " + (isOwn ? "flex-row-reverse" : "")}
                >
                  <div className="h-8 w-8 rounded-full bg-gradient-signal grid place-items-center text-[10px] font-bold text-primary-foreground shrink-0">
                    {initials}
                  </div>
                  <div
                    className={
                      "space-y-1 min-w-0 " + (isOwn ? "text-right" : "")
                    }
                  >
                    <div
                      className={
                        "flex items-center gap-2 " +
                        (isOwn ? "flex-row-reverse" : "")
                      }
                    >
                      <span className="text-xs font-semibold text-foreground">
                        {msg.author?.username || "unknown"}
                      </span>
                      <span className="text-[10px] text-muted-foreground">
                        {msg.created_at ? timeAgo(msg.created_at) : ""}
                      </span>
                    </div>
                    <p
                      className={
                        "text-sm leading-relaxed " +
                        (isOwn ? "text-foreground/90" : "text-foreground/90")
                      }
                    >
                      {msg.content}
                    </p>
                  </div>
                </div>
              );
            })
          )}
          <div ref={messagesEndRef} />
        </div>

        <div className="border-t border-border/60 px-6 py-4">
          <div className="flex gap-3">
            <input
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder={isConnected ? "Type a message..." : "Connecting..."}
              disabled={!isConnected}
              className="flex-1 bg-secondary/40 border border-border rounded-sm px-4 py-2.5 text-sm text-foreground placeholder:text-muted-foreground/60 focus:outline-none focus:border-primary/40 disabled:opacity-50"
            />
            <button
              onClick={handleSend}
              disabled={!isConnected || !input.trim()}
              className="h-10 w-10 grid place-items-center bg-gradient-signal hover:opacity-90 text-primary-foreground rounded-sm disabled:opacity-40 disabled:cursor-not-allowed transition-all"
            >
              <Send className="h-4 w-4" />
            </button>
          </div>
        </div>
      </div>

      <div className="w-56 border-l border-border/60 md:flex-col shrink-0 hidden md:flex">
        <div className="border-b border-border/60 px-5 py-4 min-h-[52.8px]">
          <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
            Online{" "}
            <span className="text-foreground">— {onlineUsers.length}</span>
          </div>
        </div>
        <div className="flex-1 overflow-y-auto py-2">
          {onlineUsers.length === 0 ? (
            <p className="text-xs text-muted-foreground text-center pt-8">
              No users online
            </p>
          ) : (
            onlineUsers.map((u: OnlineUser) => {
              const initials = (u.username || "??")
                .replace("@", "")
                .slice(0, 2)
                .toUpperCase();
              return (
                <div
                  key={u.id}
                  className="flex items-center gap-3 px-5 py-2.5 hover:bg-secondary/30 transition-colors"
                >
                  <span className="h-2 w-2 rounded-full bg-success shrink-0" />
                  <div className="h-7 w-7 rounded-full bg-secondary grid place-items-center text-[9px] font-bold text-foreground shrink-0">
                    {initials}
                  </div>
                  <span className="text-xs text-foreground truncate">
                    {u.username || "user_" + u.id}
                  </span>
                </div>
              );
            })
          )}
        </div>
      </div>
    </div>
  );
}
