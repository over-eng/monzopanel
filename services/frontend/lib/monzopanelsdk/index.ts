import { BatchQueue } from "./batchQueue";
import { v4 as uuidv4 } from "uuid";


export interface EventData {
    event: string;
    properties: Record<string, string>;
}

interface FullEventData extends EventData {
    distinct_id: string;
    client_timestamp: string;
}

export interface InitParams {
    writeKey: string;
    host: string;
    distinctId?: string;
    batchSize?: number;
    retries?: number;
}

export class MonzopanelSDK {
    private batchQueue: BatchQueue<FullEventData, number>;
    private host: string;
    private writeKey: string;
    private distinctId: string;

    constructor({
        writeKey,
        host,
        distinctId,
        retries = 5,
        batchSize = 10,
    }: InitParams) {
        this.distinctId = distinctId ?? uuidv4();
        this.host = host;
        this.writeKey = writeKey;
        this.batchQueue = new BatchQueue<FullEventData, number>({
            timeout: 500,
            batchSize: batchSize,
            setTimeout: (callback, timeout) => window.setTimeout(callback, timeout),
            clearTimeout: (timeoutHandle) => window.clearTimeout(timeoutHandle),
            baseDelay: 500,
            retries: retries,
            executeBatch: this.executeBatch.bind(this),
        });
    }

    private async executeBatch(tasks: FullEventData[]) {
        const url = `${this.host}/analytics/batch`        
        const headers = {
            authorization: `Bearer: ${this.writeKey}`,
            "Content-Type": "application/json",
        };
        const response = await fetch(url, {
            method: "POST",
            headers: headers,
            body: JSON.stringify(tasks),
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
    }

    public track(event: EventData) {
        const data: FullEventData = {
            distinct_id: this.distinctId,
            client_timestamp: new Date().toISOString(),
            ...event
        }
        this.batchQueue.submit(data);
    }

    public async flush() {
        await this.batchQueue.flush();
    }
}