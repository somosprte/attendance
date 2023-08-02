# Build the React application
FROM node:alpine AS node_builder
ADD ./frontend /frontend
WORKDIR /frontend/
RUN npm install
RUN npm run build

# Build the Go Application
FROM golang:latest AS builder
ENV GOPATH ""
RUN go env -w GOPROXY=direct
ADD . .
RUN go mod tidy
COPY --from=node_builder /frontend/build ./frontend/build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /main .

# Run
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /main ./
RUN chmod +x ./main
EXPOSE 8080
CMD ./main