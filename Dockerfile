# Use a minimal Ubuntu image for the second stage
FROM ubuntu:22.04

# Set the Current Working Directory inside the container
WORKDIR /backend-engineer-assignment

# Copy the Pre-built binary file from the builder stage
COPY backend-engineer-assignment .

# Command to run the executable
CMD ["./backend-engineer-assignment"]
