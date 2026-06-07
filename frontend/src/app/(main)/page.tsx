import HomeClient from "./HomeClient";

const API_BASE =
  (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080") + "/api/v1";

async function fetchJSON<T>(url: string): Promise<T | null> {
  try {
    const res = await fetch(url, { next: { revalidate: 60 } });
    if (!res.ok) return null;
    return res.json();
  } catch {
    return null;
  }
}

export default async function HomePage() {
  const [featured, trending, threads] = await Promise.all([
    fetchJSON<any>(`${API_BASE}/threads/featured`),
    fetchJSON<any>(`${API_BASE}/threads/trending`),
    fetchJSON<any>(`${API_BASE}/threads?page=1&pageSize=5&sort=latest`),
  ]);

  return (
    <HomeClient
      initialFeatured={featured}
      initialTrending={trending}
      initialThreads={threads}
    />
  );
}
