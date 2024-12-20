"use client";

import { useEffect } from "react";
import { useAnalytics } from "./hooks/useAnalytics";
import { useEventsForDistinctId } from "./hooks/useEventsForDistinctId";
import styles from "./page.module.css";
import Button from "./ui/button/button";
import EventTable from "./ui/eventtable/EventTable";
import { useEventsOvertime } from "./hooks/useEventsOvertime";
import EventsOvertimeChart from "./ui/EventsOvertimeChart/EventsOvertimeChart";

export default function Home() {
    const { track, distinctId } = useAnalytics();

    const handleTrackClick = () => {
        track({
            event: 'button_click',
            properties: {
                testing: "this is a test",
            }
        });
    };

    const { data: events, refetch } = useEventsForDistinctId(distinctId);
    // Set up automatic refresh: events for distinctId
    useEffect(() => {
        const intervalId = setInterval(() => {
            refetch();
        }, 1000); // 1 second
        
        return () => clearInterval(intervalId);
    }, [refetch]);
    
    const { data: overtimeData, refetch: refetchOvertime } = useEventsOvertime();
    useEffect(() => {
        const intervalId = setInterval(() => {
            refetchOvertime()
        }, 5000) // 5 seconds

        return () => clearInterval(intervalId);
    }, [refetchOvertime])

    return (
        <>
        <header className={styles.header_nav}>
            <div className={styles.header_container}>
                <h1 className={styles.logo}>
                    monzopanel
                </h1>
                <nav>
                    <Button>
                        <a href="https://www.linkedin.com/in/john-newman-336a82140/">Hire me</a>
                    </Button>
                </nav>
            </div>
        </header>
        
        <main>
            <section className={styles.blog_post_hero}>
                <div className={styles.blog_post_hero_metadata}>
                    <div>
                        <time className={styles.blog_post_hero_time} dateTime="2024-10-20">20 NOVEMBER 2024</time>
                    </div>
                    <Button>
                        <a href="https://github.com/over-eng/monzopanel">Github</a>
                    </Button>
                </div>
                <div className={styles.blog_post_hero_title_container}>
                    <h1 className={styles.post_title}>
                        An analytics pipeline using Monzo&apos;s core technologies
                    </h1>
                </div>
            </section>

            <section className={styles.blog_section}>
                <article className={styles.article}>
                    <p className={styles.article_p}>
                        This article details creating a working prototype to collect analytics data, similar to Mixpanel (hence Monzopanel).
                    </p>
                </article>
            </section>

            <section className={styles.analytics_interface}>
                <EventTable data={events || []}/>
                <div className={styles.events_overtime_container}>
                    <EventsOvertimeChart data={overtimeData || []}/>
                </div>
                <Button onClick={handleTrackClick}>
                    Track Me
                </Button>
            </section>
        </main>
        </>
    );
}
