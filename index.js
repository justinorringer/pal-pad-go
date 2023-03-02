const express = require('express');
const app = express();
const http = require('http');
const server = http.createServer(app);
const { Server } = require("socket.io");

const io = new Server(server);

// Using in memory storage for now
let lines = [];

app.get('/', (req, res) => {
    res.sendFile(__dirname + '/static/index.html');
});

io.on('connection', (socket) => {
    io.emit('sync', lines)

    socket.on('line', (line) => {
        lines.push(line);
        io.emit('sync', lines);
    });

    socket.on('clear', () => {
        lines = [];
        io.emit('sync', lines)
    });
});

server.listen(3000, () => {
    console.log('listening on *:3000');
});