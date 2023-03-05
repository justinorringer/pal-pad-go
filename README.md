# PalPad - Server
PalPad is a collaborative drawing app ran with a Svelte frontend and this Node.js server.

The server works as a hub between clients for websocket events, piping sketch info to other clients and storing that info in its Mongo database. Thanks to socket.io, our app supports multiple rooms of users, with clients being able to access a room through a sharable link.

[Draw here](https://pal-pad-app.codymitchell.dev/)!

[PalPad UI](https://github.com/CodyWMitchell/pal-pad-ui)

For future work, here's our [Trello board](https://trello.com/b/SIDZ1Y4g/pal-pad).

## Websocket Events
### Joining a room
Whenever a client connects to the websocket, the frontend sends a "join" event. If the room ID sent is not already in the database, a new room is created. Either way, a `sketch` event is sent out to the clients with all the data for a sketch, so that a user gets all historic data.

### Drawing a line
At the end of a client drawing a line, the UI sends the line data to the server. The line data is formatted as:

```
{
  color: { r, g, b, a},
  points: [{
    x1, y1, x2, y2
  }, ...]
}
```

The array of points allows the UI to interpolate spaces between the points; this makes lines more full.

The line is saved to the `lines` collection and assigned an `_id`. That id is appended the `sketches` collection entry for the current sketch (the `sketches` collection is an index table of lines).

Finally, the line is synced to the other clients in that room by emitting that line to the `sync` event.

### Clearing a drawing
We are still working on how the user interacts with erasing their drawings. So far, the user can delete all lines in a `clear` event. This deteles all the lines from the `lines` collection related to the sketch from the `sketches` collection. Next, the server emits `clear` to all clients in that room.

The UI redraws the canvas.

## Development
### Local
Run `npm install` to install dependencies, then run `USERNAME=<> PASSWORD=<> npm start` to start the server. The username and password parameters refer to the how we login to our MongoDB.

View the app at `http://localhost:3000/test`. Enter a desired room ID and "draw" lines. This test mimicks how the Svelte app interacts with our server.

### Connecting your own MongoDB
If you are wanting to run the app with your own mongoDB, run an instance with podman. Here's a quick tutorial on [setting mongoDB up on Fedora](https://tecadmin.net/install-mongodb-on-fedora/).

Once running, access the database with the mongo cli. Run `mongo` in the command line, and you can run queries on collections and test databases.

Change the URL [here](https://github.com/justinorringer/pal-pad-server/blob/main/mongo/mongo.js#L5) to where your mongo DB is running; the default is `localhost:27017`. If the app can reach your database, it'll run successfully! Otherwise, it will panic and close.

## Developers
Made by [Justin](https://github.com/justinorringer) and [Cody](https://github.com/CodyWMitchell) during Red Hat's 2023 Q1 Hackathon.
