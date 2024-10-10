import http from 'k6/http';

export const options = {
  vus: 200,
  duration: '5m',
};

const jwts=[
  "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpdGVtIjoiQU1EIFJ5emVuIDEgMjM1MFgiLCJ2YXRfaW5jbCI6ZmFsc2UsInF1YW50aXR5IjoxMn0.3Wf6vEJeJIxx8OG9GXO906cYD96ND1J_875aKltp_SldzD5CJ0uoyiH2AurhP8NFmP9hIrF4oewi2ZRI-xVriQ",
  "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpdGVtIjoiQU1EIFJ5emVuIDcgNzU2MVgiLCJ2YXRfaW5jbCI6dHJ1ZSwicXVhbnRpdHkiOjEyNH0.LEDoAmeVv4qLnmvM_bO4D3AfUhf65F5EC3AartsW96C3gkdgqQi6iLU8NOR7jrByoKbbJnQapNsmAKa5rN0sIw",
  "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpdGVtIjoiQU1EIFJ5emVuIDQgMTM0NlgiLCJ2YXRfaW5jbCI6dHJ1ZSwicXVhbnRpdHkiOjY0fQ.8x2pr969bcS_NwTprDSvaGU2VoP4amDfxwWKiQYgun1uSVG_7-zOdINyQxRZl_jA5cvszslWuRS0ib84wmQyYA",
  "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpdGVtIjoiQU1EIFJ5emVuIDcgNjQzOVgiLCJ2YXRfaW5jbCI6ZmFsc2UsInF1YW50aXR5Ijo1fQ.0umodDqKAx839AyrSVxcEh9G5xfF5S96QEjWtFlH6YUt2Dh6lMqmpGsqR-nOdX0SSzemxxmyHQQR1Frbx890bQ",
];

let target_ip='127.0.0.1:8080'

export default function () {
  http.get('http://'+target_ip+'/api/price?jwt='+jwts[Math.floor(Math.random() * 3)]);
}
