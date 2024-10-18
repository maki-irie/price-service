import http from 'k6/http';

export const options = {
  discardResponseBodies: true,

  scenarios: {
    contacts: {
      executor: 'ramping-arrival-rate',

      // Start iterations per `timeUnit`
      startRate: 50,

      // Start `startRate` iterations per minute
      timeUnit: '1s',

      // Pre-allocate necessary VUs.
      preAllocatedVUs: 50,

      stages: [
        // Start 50 TPS for 5min, warm-up.
        { target: 50, duration: '3m' },

        // Linearly ramp-up to 400 iterations per `timeUnit` over the following two minutes.
        { target: __ENV.TARGET_RATE || 400, duration: '2m' },

        // Keep 400 tps for 5 min
        { target: __ENV.TARGET_RATE || 400, duration: '5m' },

        // Linearly ramp-down to starting 50 iterations per `timeUnit` over the last two minutes.
        { target: 50, duration: '2m' },
      ],
    },
  },
};

const jwts=[
  "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpdGVtIjoiQU1EIFJ5emVuIDEgMjM1MFgiLCJ2YXRfaW5jbCI6ZmFsc2UsInF1YW50aXR5IjoxMn0.3Wf6vEJeJIxx8OG9GXO906cYD96ND1J_875aKltp_SldzD5CJ0uoyiH2AurhP8NFmP9hIrF4oewi2ZRI-xVriQ",
  "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpdGVtIjoiQU1EIFJ5emVuIDcgNzU2MVgiLCJ2YXRfaW5jbCI6dHJ1ZSwicXVhbnRpdHkiOjEyNH0.LEDoAmeVv4qLnmvM_bO4D3AfUhf65F5EC3AartsW96C3gkdgqQi6iLU8NOR7jrByoKbbJnQapNsmAKa5rN0sIw",
  "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpdGVtIjoiQU1EIFJ5emVuIDQgMTM0NlgiLCJ2YXRfaW5jbCI6dHJ1ZSwicXVhbnRpdHkiOjY0fQ.8x2pr969bcS_NwTprDSvaGU2VoP4amDfxwWKiQYgun1uSVG_7-zOdINyQxRZl_jA5cvszslWuRS0ib84wmQyYA",
  "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpdGVtIjoiQU1EIFJ5emVuIDcgNjQzOVgiLCJ2YXRfaW5jbCI6ZmFsc2UsInF1YW50aXR5Ijo1fQ.0umodDqKAx839AyrSVxcEh9G5xfF5S96QEjWtFlH6YUt2Dh6lMqmpGsqR-nOdX0SSzemxxmyHQQR1Frbx890bQ",
];

let target_ip = __ENV.TARGET_IP || '127.0.0.1:8080';

export default function () {
  http.get('http://'+target_ip+'/api/price?jwt='+jwts[Math.floor(Math.random() * 3)]);
}
