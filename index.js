import express from 'express';
import http from 'http';
import { Server } from "socket.io";
import { fileURLToPath } from 'url';

import Client from './mongo/mongo.js';

const __dirname = fileURLToPath(import.meta.url).replace('index.js', '');

const app = express();
const server = http.createServer(app);
const io = new Server(server, {
    cors: {
        origin: "*",
    }
});

var client = new Client();

await client.Connect();

app.get('/test', (req, res) => {
    res.sendFile(__dirname + '/static/index.html');
});

io.on('connection', (socket) => {
    // join a room, and create a new sketch
    // and send the sketch id to the client
    socket.on('join', (sketch_id) => {
        (async () => {
            let sketch = await client.GetSketch(sketch_id);

            if (sketch == null) {
                console.log("Creating new sketch");
                sketch_id = await client.NewSketch(sketch_id);
                console.log("New sketch id: " + sketch_id)
            }

            socket.join(sketch_id);

            console.log("Joined sketch: " + sketch_id);

            // emit the sketch to only the client that just joined
            socket.emit('sketch', sketch);
        })();
    });

    socket.on('line', (sketch_id, line) => {
        console.log("Line: " + line
            + " from sketch: " + sketch_id)
        // add the line to the database
        client.NewLine(sketch_id, line);
        socket.broadcast.to(sketch_id).emit('sync', line);
    });

    socket.on('clear', (sketch_id) => {
        client.ClearLines(sketch_id);
        io.to(sketch_id).emit('clear');
    });
});

server.listen(3000, () => {
    console.log('listening on *:3000');
});