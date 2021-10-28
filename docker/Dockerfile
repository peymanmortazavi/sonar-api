FROM golang:alpine3.14 AS builder

WORKDIR /project
COPY . .

RUN go mod download
RUN go build -o sonar-api -ldflags='-s -w' cmd/sonar-api/sonar-api.go

FROM builder AS sonar-api

RUN addgroup sonar-api -S && \
	adduser sonar-api -G sonar-api -S -H -D

USER sonar-api

COPY --from=builder --chown=sonar-api /project/sonar-api /opt/sonar/sonar-api

CMD /opt/sonar/sonar-api
