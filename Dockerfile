    FROM golang:latest
    LABEL maintainer="klara.chess.school@gmail.com"
    RUN mkdir forum
    WORKDIR forum
    COPY . .
    RUN go build
    CMD ["./forum"]