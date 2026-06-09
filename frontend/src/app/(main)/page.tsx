import HomeClient from "./HomeClient";
import { fetchJSON, apiUrl } from "@/lib/server-fetch";

export default async function HomePage() {
  const [featured, trending, threads] = await Promise.all([
    fetchJSON<any>(apiUrl("/threads/featured")),
    fetchJSON<any>(apiUrl("/threads/trending")),
    fetchJSON<any>(apiUrl("/threads?page=1&pageSize=5&sort=latest")),
  ]);

  return (
    <HomeClient
      initialFeatured={featured}
      initialTrending={trending}
      initialThreads={threads}
    />
  );
}
