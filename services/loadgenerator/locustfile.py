#!/usr/bin/python
import logging
import os

from datetime import datetime
from uuid import uuid4

import gevent

from locust import FastHttpUser, task

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('locust_requests.log'),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger(__name__)


AUTH_KEY = os.environ["AUTH_KEY"]

class APIUser(FastHttpUser):

    def on_start(self):
        self.distinct_id = uuid4()

    @task(3)
    def post_analytics_events(self):
        events = [
            {
                "distinct_id": str(self.distinct_id),
                "event": "locust-event",
                "client_timestamp": f"{datetime.now().isoformat()}Z",
                "properties": {
                    "browser": "locust"
                }
            }
            for _ in range(10)
        ]
        
        response = self.client.post(
            "/analytics/batch",
            json=events,
            headers={"Authorization": f"Bearer {AUTH_KEY}"},
        )

        logging.info(f"/analytics/batch: {response.status_code} {response.content}")

    @task
    def load_analytics_charts(self):
        def concurrent_request(url):
            response = self.client.get(
                url,
                headers={"Authorization": f"Bearer {AUTH_KEY}"},
            )
            logging.info(f"{url}: {response.status_code}")

        pool = gevent.pool.Pool()
        urls = [
            f"/analytics/distinct_id/{self.distinct_id}/events?page_size=10",
            "/analytics/stats/events_overtime",
        ]
        for url in urls:
            pool.spawn(concurrent_request, url)
        pool.join()
