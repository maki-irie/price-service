FROM grafana/k6

COPY eei-test.js /

ENTRYPOINT ["k6", "run", "/eei-test.js"]