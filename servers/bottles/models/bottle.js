const mongoose = require('mongoose');
const Schema = mongoose.Schema;

const BottleSchema = new Schema({
    name: String, 
    bottlePosts: [{
        creator: {
            type: {
                id: {
                    type: Number
                },
                username: String,
            }
        },
        body: {
            type: String,
            required: true
        },
        createdAt: {
            type: Date
        },
        editedAt: {
            type: Date
        }
    }]
});

const Bottle = mongoose.model('bottle', BottleSchema);

module.exports = Bottle;