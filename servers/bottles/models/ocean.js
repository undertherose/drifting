const mongoose = require('mongoose');
const Schema = mongoose.Schema;

const OceanSchema = new Schema({
    name: String, 
    bottles: [{
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
        },
        tags : [{
            tag: String
        }]
    }],
    tags: [{
        tagName: {
            type: String,
            required: true,
            unique: true
        }, 
        tagCount: {
            type: Number,
            required: true,
        },
        tagLastUpdate: {
            type: Date
        }
    }]
});

const Ocean = mongoose.model('ocean', OceanSchema);

module.exports = Ocean;