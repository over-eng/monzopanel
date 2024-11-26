import { useCallback, useEffect, useState } from "react";
import { v4 as uuidv4 } from 'uuid';
import { getAnalyticsClient } from "../utils/analytics";
import { EventData } from "@/lib/monzopanelsdk";

const STORAGE_KEY = "analytics_distinct_id";

function getStoredOrNewDistinctId(): string {
    // Check if running in a browser environment
    if (typeof window === 'undefined') {
        return uuidv4();
    }

    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored) {
        return stored;
    }

    const newId = uuidv4();
    localStorage.setItem(STORAGE_KEY, newId);
    return newId;
}

export function useAnalytics() {
    const [distinctId, setDistinctId] = useState<string>(() => getStoredOrNewDistinctId());

    useEffect(() => {
        // Re-run on mount to ensure consistency with SSR
        const id = getStoredOrNewDistinctId();
        if (id !== distinctId) {
            setDistinctId(id);
        }
    }, []);

    const track = useCallback((event: EventData) => {
        const client = getAnalyticsClient(distinctId);
        client.track(event);
    }, [distinctId]);

    return { track, distinctId };
}
