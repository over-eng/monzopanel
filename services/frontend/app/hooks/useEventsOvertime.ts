import { useQuery } from "@tanstack/react-query";
import { NEXT_PUBLIC_ANALYTICS_KEY, NEXT_PUBLIC_MONZOPANEL_API_HOST } from "../config";

interface RawEventsOvertimeResponse {
    buckets: Array<{
        timestamp: string;
        count: number;
    }>;
}

export interface EventCount {
    timestamp: Date;
    count: number;
}

export function useEventsOvertime() {
    return useQuery({
        queryKey: ['overtime-events'],
        queryFn: () => fetchEvents(),
    });
}

function transformEvents(rawData: RawEventsOvertimeResponse): EventCount[] {
    return rawData.buckets.map(bucket => {
        return {
            timestamp: new Date(bucket.timestamp),
            count: bucket.count,
        }
    })
}

async function fetchEvents(): Promise<EventCount[]> {

    const response = await fetch(`${NEXT_PUBLIC_MONZOPANEL_API_HOST}/analytics/stats/events_overtime`, {
        headers: {
            Authorization: `Bearer ${NEXT_PUBLIC_ANALYTICS_KEY}`,
            'Content-Type': 'application/json'
        },
    });

    if (!response.ok) {
        throw new Error("error fetching analytics events");
    }

    const rawData: RawEventsOvertimeResponse = await response.json();
    return transformEvents(rawData);
}
