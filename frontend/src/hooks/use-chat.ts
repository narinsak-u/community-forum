"use client";

import { useEffect, useRef, useState, useCallback } from "react";
import { toast } from "sonner";

export interface ChatMessage {
  id: number;
  content: string;
  author: { id: number; username: string; avatar: string };
  created_at: string;
}

export interface OnlineUser {
  id: number;
  username: string;
  avatar: string;
}

const WS_URL =
  (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080").replace(/^http/, "ws") +
  "/ws/chat";

interface UseChatReturn {
  messages: ChatMessage[];
  onlineUsers: OnlineUser[];
  sendMessage: (content: string) => void;
  loadMore: () => void;
  isConnected: boolean;
  isLoading: boolean;
  error: string | null;
}

export function useChat(): UseChatReturn {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [onlineUsers, setOnlineUsers] = useState<OnlineUser[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const retriesRef = useRef(0);
  const maxRetries = 3;
  const oldestIdRef = useRef<number | null>(null);

  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.close();
    }

    const ws = new WebSocket(WS_URL);
    const isCurrent = () => wsRef.current === ws;
    wsRef.current = ws;

    ws.onopen = () => {
      if (!isCurrent()) return;
      setIsConnected(true);
      retriesRef.current = 0;
    };

    ws.onmessage = (event) => {
      if (!isCurrent()) return;
      try {
        const data = JSON.parse(event.data);
        switch (data.type) {
          case "init":
            setMessages(data.messages || []);
            setOnlineUsers(data.users || []);
            setIsLoading(false);
            if (data.messages?.length) {
              oldestIdRef.current = data.messages[0].id;
            }
            break;
          case "message":
            setMessages((prev) => [...prev, data.message]);
            break;
          case "user_joined":
            setOnlineUsers((prev) => {
              if (prev.find((u) => u.id === data.user.id)) return prev;
              return [...prev, data.user];
            });
            break;
          case "user_left":
            setOnlineUsers((prev) =>
              prev.filter((u) => u.id !== data.user.id),
            );
            break;
          case "load_more":
            setMessages((prev) => [...data.messages, ...prev]);
            if (data.messages?.length) {
              oldestIdRef.current = data.messages[0].id;
            }
            break;
        }
      } catch {
        // ignore malformed messages
      }
    };

    ws.onclose = () => {
      if (!isCurrent()) return;
      setIsConnected(false);
      if (retriesRef.current < maxRetries) {
        retriesRef.current++;
        const delay = Math.pow(2, retriesRef.current - 1) * 1000;
        setTimeout(connect, delay);
      } else {
        setError("Could not connect to chat server");
        setIsLoading(false);
        toast.error("Chat connection failed. Please reload and try again.");
      }
    };
  }, []);

  useEffect(() => {
    connect();
    return () => {
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        wsRef.current.close();
      }
    };
  }, [connect]);

  const sendMessage = useCallback((content: string) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ type: "message", content }));
    }
  }, []);

  const loadMore = useCallback(() => {
    if (oldestIdRef.current && wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(
        JSON.stringify({ type: "load_more", before: oldestIdRef.current }),
      );
    }
  }, []);

  return { messages, onlineUsers, sendMessage, loadMore, isConnected, isLoading, error };
}
