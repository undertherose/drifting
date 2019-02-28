const mongoose = require('mongoose');
const Schema = mongoose.Schema;

const OceanSchema = new Schema({
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
        },
        isPublic: {
            type: Boolean,
            required: true
        }
    }]
});

const Bottle = mongoose.model('bottle', BottleSchema);

module.exports = Bottle;