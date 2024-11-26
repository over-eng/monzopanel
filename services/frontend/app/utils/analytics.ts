import { MonzopanelSDK } from "@/lib/monzopanelsdk";
import { NEXT_PUBLIC_ANALYTICS_KEY, NEXT_PUBLIC_MONZOPANEL_API_HOST } from "../config";

export function createAnalyticsSingleton(distinctId: string) {
  return new MonzopanelSDK({
    writeKey: NEXT_PUBLIC_ANALYTICS_KEY,
    host: NEXT_PUBLIC_MONZOPANEL_API_HOST || "https://api.over-engineering.co.uk",
    distinctId: distinctId,
  });
}

let client: MonzopanelSDK | null = null;

export function getAnalyticsClient(distinctId: string) {
  if (!client) {
    client = createAnalyticsSingleton(distinctId);
  }
  return client;
}
