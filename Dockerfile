FROM golang:1.20 AS builder

# Create appuser.
ENV USER=report_handler
ENV UID=10001
ENV CGO_ENABLED=0
# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
  --disabled-password \    
  --gecos "" \    
  --home "/nonexistent" \    
  --shell "/sbin/nologin" \    
  --no-create-home \    
  --uid "${UID}" \    
  "${USER}"

COPY ./src /code
WORKDIR /code
RUN go build -ldflags "-s -w" -o csp

FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /code/csp /go/bin/csp
USER report_handler:report_handler
