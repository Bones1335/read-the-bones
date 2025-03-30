FROM debian:stable-slim

COPY read-the-bones read-the-bones
COPY content/ content/ 
COPY static/ static/ 
COPY templates/ templates/ 

CMD ["./read-the-bones"]
