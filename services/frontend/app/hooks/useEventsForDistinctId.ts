import { useQuery } from "@tanstack/react-query";
import { EventData } from "../ui/eventtable/EventTable";
import { NEXT_PUBLIC_MONZOPANEL_API_HOST, NEXT_PUBLIC_ANALYTICS_KEY } from "../config";

interface RawEventResponse {
    events: Array<{
        id: string;
        event: string;
        team_id: string;
        distinct_id: string;
        properties: {
            browser?: string;
        };
        client_timestamp: {
            seconds: number;
        };
        created_at: {
            seconds: number;
            nanos: number;
        };
        loaded_at: {
            seconds: number;
            nanos: number;
        };
    }>;
}

const transformEvents = (rawData: RawEventResponse): EventData[] => {
    return rawData.events.map(event => {
        const timestamp = new Date(event.client_timestamp.seconds * 1000).toISOString();
        
        const latency = (event.loaded_at.seconds * 1000 + event.loaded_at.nanos / 1000000) - 
                       (event.created_at.seconds * 1000 + event.created_at.nanos / 1000000);

        return {
            event: event.event,
            timestamp: timestamp,
            latency: Math.round(latency),
            browser: event.properties.browser ?? ""
        };
    });
};

export function useEventsForDistinctId(distinctId: string) {
    return useQuery({
        queryKey: ['events', distinctId],
        queryFn: () => fetchEvents(distinctId),
    });
}

async function fetchEvents(distinctId: string): Promise<EventData[]> {
    console.log("fetching some data")

    const params = new URLSearchParams({
        page_size: "10",
    });

    const response = await fetch(`${NEXT_PUBLIC_MONZOPANEL_API_HOST}/analytics/${distinctId}/events?${params}`, {
        headers: {
            Authorization: `Bearer ${NEXT_PUBLIC_ANALYTICS_KEY}`,
            'Content-Type': 'application/json'
        },
    });

    if (!response.ok) {
        throw new Error("error fetching analytics events");
    }

    const rawData: RawEventResponse = await response.json();
    return transformEvents(rawData);
}
