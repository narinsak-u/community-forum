import ThreadsClient from "./ThreadsClient";
import { fetchJSON, apiUrl } from "@/lib/server-fetch";

export default async function ThreadsPage() {
  const [featured, trending, threads] = await Promise.all([
    fetchJSON<any>(apiUrl("/threads/featured")),
    fetchJSON<any>(apiUrl("/threads/trending")),
    fetchJSON<any>(apiUrl("/threads?page=1&pageSize=5&sort=latest")),
  ]);

  return (
    <ThreadsClient
      initialFeatured={featured}
      initialTrending={trending}
      initialThreads={threads}
    />
  );
}
