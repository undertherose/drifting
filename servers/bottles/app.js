const express = require('express');
const mongoose = require('mongoose');
const amqp = require('amqplib/callback_api')

const app = express();

const PORT = process.env.PORT;
const instanceName = process.env.NAME;
const dbURL = 'mongodb://mongo:27017/mydb';

const connectWithRetry = () => {
    console.log('MongoDB connection with retry');
    mongoose.connect(dbURL).then(() => {
      console.log('MongoDB is connected');
    }).catch(err => {
      console.log('MongoDB connection unsuccessful, retry after 2 seconds.')
      setTimeout(connectWithRetry, 2000);
    });
}

connectWithRetry();

//rabbitmq connection

amqp.connect('amqp://rabbitmq', (err, conn) => {
    if (err) {
        console.log('Failed to connect to rabbit ' + err);
        process.exit(1);
    }

    conn.createChannel((err, ch) => {
        if(err) {
            console.log('Failed to create a channel');
            process.exit(1);
        }

        console.log('channel created');
        ch.assertQueue('rabbitmq', {durable: true});
        app.set('ch', ch);

        // Channel.create({name: 'general'}).then((channel) => {
        //     let qPayLoad = {};
        //     qPayLoad.type = 'channel-new';
        //     qPayLoad.channel = channel;
        //     qPayLoad.userIDs = [];
        //     if(channel.private) {
        //         channel.members.map(mem => {
        //             qPayLoad.userIDs.push(mem.id);
        //         });
        //     }
        //     ch.sendToQueue(
        //         "MsgQueue",
        //         new Buffer(JSON.stringify(qPayLoad)),
        //         {persistent: true}
        //     );
        //     console.log('general created');
        // }).catch(err => {
        //     console.log(err);
        // });
    });
});

mongoose.Promise = global.Promise;

app.use((req, res, next) => {
    if(!req.get('X-User')) {
        res.status(401).send({message: 'No user header found'});
    } else {
        next();
    }
});

//Middleware to parse the request body as JSON
app.use(express.json());
app.use('/v1', require('./routes.js'));

app.listen(PORT, instanceName,() => {
    console.log('listening on ' + PORT);
    console.log('host: ' + instanceName);
});