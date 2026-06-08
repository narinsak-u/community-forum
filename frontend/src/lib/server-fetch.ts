const API_BASE =
  (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080") + "/api/v1";

export async function fetchJSON<T>(url: string): Promise<T | null> {
  try {
    const res = await fetch(url, { cache: "no-store" });
    if (!res.ok) return null;
    return res.json();
  } catch {
    return null;
  }
}

export function apiUrl(path: string): string {
  return `${API_BASE}${path}`;
}
