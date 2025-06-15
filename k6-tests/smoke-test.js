import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  scenarios: {
    smoke: {
      executor: 'constant-vus',
      vus: 1,
      duration: '30s',
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.01'], // less than 1% failed requests
    http_req_duration: ['p(95)<500'], // 95% of requests must complete below 500ms
  },
};

const BASE_URL = 'http://localhost:8080/api/v1';

export default function () {
  const res = http.get(`${BASE_URL}/tanks`);

  check(res, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(1);
}
