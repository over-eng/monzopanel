import { EventData, MonzopanelSDK } from "@/lib/monzopanelsdk";
import { useCallback } from "react";

export function createAnalyticsSingleton() {
  return new MonzopanelSDK({
    writeKey: process.env.MONZOPANEL_PUBLIC_WRITE_KEY || "",
    host: process.env.MONZOPANEL_API_HOST || "https://api.over-engineering.co.uk",

    distinctId: typeof window !== undefined
      ? localStorage.getItem("user-distinct-id") || undefined 
      : undefined
  });
}

let client: MonzopanelSDK | null = null;

export function getAnalyticsClient() {
  if (!client) {
    client = createAnalyticsSingleton();
  }
  return client;
}

export function useAnalytics() {
    const track = useCallback((event: EventData) => {
        const client = getAnalyticsClient();
        client.track(event);
    }, []);
    return { track };
}
