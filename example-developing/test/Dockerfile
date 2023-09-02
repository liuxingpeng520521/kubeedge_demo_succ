# Based on ubuntu
FROM ubuntu:18.04
LABEL maintainers="KubeEdge Community Developer"
LABEL description="KubeEdge Raspi Counter App"

# Copy from build directory
COPY main /pi-counter-app

# Define default command
ENTRYPOINT ["/pi-counter-app"]

# Run the executable
CMD ["pi-counter-app"]
