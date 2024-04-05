# Stage 1: Build stage
FROM golang:1.21.9 AS builder

WORKDIR /app

# Copy necessary files
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the backend
RUN CGO_ENABLED=0 GOOS=linux go build -o backend ./cmd/main.go

# Stage 2: Final stage
FROM alpine:latest

WORKDIR /app

ARG DB_NAME
ARG DB_PORT
ARG DB_HOST
ARG DB_USERNAME
ARG DB_PASSWORD
ARG DB_PARAMS
ARG BCRYPT_SALT
ARG JWT_SECRET
ARG S3_ID
ARG S3_SECRET_KEY
ARG S3_BUCKET_NAME
ARG S3_REGION
ARG PEM_FILE

ENV DB_NAME=$DB_NAME
ENV DB_PORT=$DB_PORT
ENV DB_HOST=$DB_HOST
ENV DB_USERNAME=$DB_USERNAME
ENV DB_PASSWORD=$DB_PASSWORD
ENV DB_PARAMS=$DB_PARAMS
ENV JWT_SECRET=$JWT_SECRET
ENV BCRYPT_SALT=$BCRYPT_SALT
ENV S3_ID=$S3_ID
ENV S3_SECRET_KEY=$S3_SECRET_KEY
ENV S3_BUCKET_NAME=$S3_BUCKET_NAME
ENV S3_REGION=$S3_REGION


COPY $PEM_FILE ap-southeast-1-bundle.pem
RUN chmod 600 ap-southeast-1-bundle.pem
COPY --from=builder /app/backend ./backend

# Set the command to run the backend
CMD ["./backend"]