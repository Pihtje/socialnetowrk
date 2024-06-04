FROM golang:alpine as build
ENV CGO_ENABLED=1
# ENV GOPATH=/go
# ENV GOCdocker psACHE=/go-build

# Install tools required for project
RUN apk update && apk add --no-cache git g++ sqlite

# List project dependencies with go.mod and go.sum
# These layers are only re-built when Gopkg files are updated
WORKDIR /go/src/backend/
COPY go.mod go.sum ./
# Install library dependencies
RUN go mod download

# Copy the entire project and build it
# This layer is rebuilt when a file changes in the project directory
COPY . .
# RUN go build -o /bin/backend main.go
EXPOSE 8000
CMD [ "go","run","main.go" ]

# This results in a single layer image
# FROM alpine
# RUN apk update && apk add --no-cache sqlite
# COPY --from=build /bin/backend /bin/backend
# COPY database.db /bin/
# COPY db/migrations/sqlite db/migrations/sqlite
# EXPOSE 8000
# ENTRYPOINT ["/bin/backend"]
#CMD ["--help"]
