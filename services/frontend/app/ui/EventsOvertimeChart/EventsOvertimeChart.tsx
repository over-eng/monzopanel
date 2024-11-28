"use client";

import { EventCount } from "@/app/hooks/useEventsOvertime";
import GradientLineChart, { DataPoint } from "../OvertimeLineChart/OvertimeLineChart";

const EventsOvertimeChart = ({ data }: { data: EventCount[] }) => {
    const seriesData: DataPoint[] = data.map(series => ({
        name: series.timestamp.toString(),
        value: series.count,
    }));
    return <GradientLineChart data={seriesData} />;
};

export default EventsOvertimeChart;
