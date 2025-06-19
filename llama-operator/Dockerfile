FROM golang:1.20-alpine AS build
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /manager main.go

FROM scratch
COPY --from=build /manager /manager
ENTRYPOINT ["/manager"]
