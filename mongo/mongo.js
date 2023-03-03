import { MongoClient } from 'mongodb';

const username = process.env.USERNAME;
const password = process.env.PASSWORD;
const url = `mongodb+srv://${username}:${password}@palpad.dfyrb8d.mongodb.net/?retryWrites=true&w=majority`;
const db = "pal-pad";
const sketches_collection = "sketches";
const lines_collection = "lines";

console.log("MongoDB hosted at: " + url);

class Client {
    constructor() {
        this.client = new MongoClient(url);
    }

    async Connect() {
        try {
            await this.client.connect();
            console.log("Connected successfully to server");

            this.db = this.client.db(db);
        } catch (err) {
            console.log(err.stack);

            // if we can't connect to the database, we should exit
            process.exit(1);
        }
    }

    async NewSketch(id) {
        try {
            const collection = this.db.collection(sketches_collection);

            const result = await collection.insertOne({ _id: id, lines: [] });
            console.log(result);

            return await this.GetSketch(id);
        } catch (err) {
            console.log(err.stack);
        }
    }

    async GetSketch(id) {
        try {
            const s_collection = this.db.collection(sketches_collection);

            const sketch_table = await s_collection.findOne({ _id: id })
            
            if (sketch_table == null) {
                return null;
            }

            // now that we have the sketch_table,
            // we need to get the lines from the lines table
            // and return everything as a single object
            // sketch: { id, lines: [line1, line2, line3] }

            const l_collection = this.db.collection(lines_collection);

            const lines_array = await l_collection.find({ _id: { $in: sketch_table.lines } }).map((element, _) => {
                return element.line;
            }).toArray();

            let result = {
                id: sketch_table._id,
                lines: lines_array
            }
            return result;
        } catch (err) {
            console.log(err.stack);
        }
    }

    async indexLine(id, line_id) {
        try {
            const sketches = this.db.collection(sketches_collection);

            const result = await sketches.updateOne(
                { _id: id }, 
                { $push: { lines: line_id } });

            console.log(result);
        } catch (err) {
            console.log(err.stack);
        }
    }

    async NewLine(sketch_id, line) {
        try {
            const lines = this.db.collection(lines_collection);

            const result = await lines.insertOne({ line: line });

            this.indexLine(sketch_id, result.insertedId);
        } catch (err) {
            console.log(err.stack);
        }
    }

    async ClearLines(sketch_id) {
        try {
            const sketches = this.db.collection(sketches_collection);

            const result = await sketches.updateOne(
                { _id: sketch_id },
                { $set: { lines: [] } }
            );

            console.log(result);
        } catch (err) {
            console.log(err.stack);
        }
    }
}

export default Client;