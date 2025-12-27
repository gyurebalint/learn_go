import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    // The "Ramp Up" profile
    stages: [
        { duration: '5s', target: 10 },  // Wake up: 10 users
        { duration: '10s', target: 50 }, // Stress: 50 concurrent users
        { duration: '5s', target: 0 },   // Cooldown
    ],
};

export default function () {
    const res = http.get('http://localhost:3000/price?symbol=BTC');
    check(res, { 'status was 200': (r) => r.status == 200 });
    sleep(1);
}


// C:\Development\go\learn_go\crypto-aggregator>k6 run stress_test.js
//
//          /\      Grafana   /‾‾/
//     /\  /  \     |\  __   /  /
//    /  \/    \    | |/ /  /   ‾‾\
//   /          \   |   (  |  (‾)  |
//  / __________ \  |_|\_\  \_____/
//
// execution: local
// script: stress_test.js
// output: -
//
//     scenarios: (100.00%) 1 scenario, 50 max VUs, 50s max duration (incl. graceful stop):
// * default: Up to 50 looping VUs for 20s over 3 stages (gracefulRampDown: 30s, gracefulStop: 30s)
//
//
//
// █ TOTAL RESULTS
//
// checks_total.......: 367     17.535944/s
// checks_succeeded...: 100.00% 367 out of 367
// checks_failed......: 0.00%   0 out of 367
//
//     ✓ status was 200
//
// HTTP
// http_req_duration..............: avg=309.04ms min=249.29ms med=275.17ms max=747.06ms p(90)=305.19ms p(95)=710.13ms
// { expected_response:true }...: avg=309.04ms min=249.29ms med=275.17ms max=747.06ms p(90)=305.19ms p(95)=710.13ms
// http_req_failed................: 0.00% 0 out of 367
// http_reqs......................: 367   17.535944/s
//
// EXECUTION
// iteration_duration.............: avg=1.3s     min=1.24s    med=1.27s    max=1.74s    p(90)=1.3s     p(95)=1.71s
// iterations.....................: 367   17.535944/s
// vus............................: 7     min=2        max=49
// vus_max........................: 50    min=50       max=50
//
// NETWORK
// data_received..................: 52 kB 2.5 kB/s
// data_sent......................: 32 kB 1.5 kB/s
//
//
//
//
// running (20.9s), 00/50 VUs, 367 complete and 0 interrupted iterations
// default ✓ [======================================] 00/50 VUs  20s