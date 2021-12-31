FROM golang:1.16 AS build

RUN echo "Acquire::Check-Valid-Until \"false\";\nAcquire::Check-Date \"false\";" | cat > /etc/apt/apt.conf.d/10no--check-valid-until

RUN apt-get update && \
apt-get install -y build-essential && \
apt-get install -y software-properties-common && \
apt-get install -y curl git vim wget pkg-config ca-certificates && \
update-ca-certificates

ENV USER=appuser
ENV UID=10001

RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR $GOPATH/src/c3alert/
#COPY . .
COPY go ./go
COPY main.go .
COPY go.mod .
COPY go.sum .
RUN go mod download
#RUN go mod verify

ENV CGO_ENABLED=0 
ENV GO111MODULE=on

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags='-w -extldflags "-static"' -a -installsuffix cgo -o c3alert .
ENV CGO_ENABLED=0

FROM alpine:3.13

COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

# Copy our static executable
COPY --from=build /go/src/c3alert ./

# Use an unprivileged user.
USER appuser:appuser

EXPOSE 8080/tcp

# Run the c3alert server binary.
ENTRYPOINT ["./c3alert"]
