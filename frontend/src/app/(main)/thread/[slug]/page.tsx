import ThreadDetailClient from "./ThreadDetailClient";

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

export default async function ThreadPage({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;
  const thread = await fetchJSON<any>(`${API_BASE}/threads/${slug}`);
  return <ThreadDetailClient slug={slug} initialThread={thread} />;
}
