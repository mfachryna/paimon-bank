FROM golang:1.21

WORKDIR /go/src/app

# Copy the backend source code
COPY . .

# Build the backend
RUN go build -o backend ./cmd/main.go

# Set the environment variables
ENV DB_NAME=${DB_NAME}
ENV DB_PORT=${DB_PORT}
ENV DB_HOST=${DB_HOST}
ENV DB_USERNAME=${DB_USERNAME}
ENV DB_PASSWORD=${DB_PASSWORD}
ENV DB_PARAMS=${DB_PARAMS}
ENV JWT_SECRET=${JWT_SECRET}
ENV BCRYPT_SALT=${BCRYPT_SALT}
ENV S3_ID=${S3_ID}
ENV S3_SECRET_KEY=${S3_SECRET_KEY}
ENV S3_BUCKET_NAME=${S3_BUCKET_NAME}
ENV S3_REGION=${S3_REGION}
ENV ENV=${ENV}

# Set the command to run the backend
CMD ["./backend"]