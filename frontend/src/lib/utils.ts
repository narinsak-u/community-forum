import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function timeAgo(dateStr: string): string {
  const date = new Date(dateStr);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const seconds = Math.floor(diff / 1000);
  if (seconds < 60) return "JUST NOW";
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return minutes + " MIN AGO";
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return hours + "H AGO";
  const days = Math.floor(hours / 24);
  return days + "D AGO";
}

export function getInitials(username?: string): string {
  return (username || "??").replace("@", "").slice(0, 2).toUpperCase();
}

export function formatCount(n: number): string {
  if (n >= 1000) return (n / 1000).toFixed(1) + "K";
  return String(n);
}
