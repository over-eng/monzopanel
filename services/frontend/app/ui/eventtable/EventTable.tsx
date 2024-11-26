import React, { useState } from 'react';
import { ChevronUp, ChevronDown } from 'lucide-react';
import styles from './EventTable.module.css';

export interface EventData {
    event: string;
    timestamp: string;
    latency: number;
    browser: string;
}

interface EventTableProps {
    data: EventData[];
    initialSort?: {
        field: keyof EventData;
        direction: 'asc' | 'desc';
    };
    onSort?: (field: keyof EventData, direction: 'asc' | 'desc') => void;
}

const EventTable: React.FC<EventTableProps> = ({
        data,
        initialSort = { field: 'timestamp', direction: 'desc' },
        onSort,
    }) => {
    const [sortField, setSortField] = useState<keyof EventData>(initialSort.field);
    const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>(initialSort.direction);

    const sortData = (data: EventData[]) => {
        return [...data].sort((a, b) => {
        if (a[sortField] < b[sortField]) return sortDirection === 'asc' ? -1 : 1;
        if (a[sortField] > b[sortField]) return sortDirection === 'asc' ? 1 : -1;
        return 0;
        });
    };

    const handleSort = (field: keyof EventData) => {
        const newDirection = field === sortField && sortDirection === 'asc' ? 'desc' : 'asc';
        setSortField(field);
        setSortDirection(newDirection);
        onSort?.(field, newDirection);
    };

    const formatLatency = (latency: number) => {
        return `${latency}ms`;
    };

    const sortedData = sortData(data);

    return (
        <div className={styles.tableContainer}>
        <div className={styles.overflowWrapper}>
            <table className={styles.table}>
            <thead>
                <tr className={styles.headerRow}>
                {(['event', 'timestamp', 'latency', 'browser'] as const).map((field) => (
                    <th
                    key={field}
                    onClick={() => handleSort(field)}
                    className={styles.headerCell}
                    >
                    <div className={styles.headerContent}>
                        {field.charAt(0).toUpperCase() + field.slice(1)}
                        <span className={styles.sortIcon}>
                        {sortField === field && (
                            sortDirection === 'asc' ? (
                            <ChevronUp />
                            ) : (
                            <ChevronDown />
                            )
                        )}
                        </span>
                    </div>
                    </th>
                ))}
                </tr>
            </thead>
            <tbody>
                {sortedData.length === 0 ? (
                <tr>
                    <td colSpan={4} className={styles.emptyState}>
                    No data available
                    </td>
                </tr>
                ) : (
                sortedData.map((row, index) => (
                    <tr
                    key={index}
                    className={styles.tableRow}
                    >
                    <td className={styles.tableCell}>{row.event}</td>
                    <td className={styles.tableCell}>{row.timestamp}</td>
                    <td className={styles.tableCell}>{formatLatency(row.latency)}</td>
                    <td className={styles.tableCell}>{row.browser}</td>
                    </tr>
                ))
                )}
            </tbody>
            </table>
        </div>
        </div>
    );
};

export default EventTable;
