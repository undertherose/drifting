const mongoose = require('mongoose');
const Schema = mongoose.Schema;

const CourseSchema = new Schema({
    sln: {
        type: Number,
        required: true,
        unique: true
    },
    name: {
        type: String,
        required: true,
        unique: true
    },
    isOpen: {
        type: Boolean,
        default: false
    },
    ta: [{
        id: Number,
        username: String
    }],
    students: [{
        id: Number,
        username: String
    }],
    queuePosts: [{
        creator: {
            type: {
                id: {
                    type: Number
                },
                username: String,
                firstname: String,
                lastname: String,
                photourl: String
            }
        },
        body: {
            type: String,
            required: true
        },
        createdAt: {
            type: Date
        }
    }]

});

const Course = mongoose.model('course', CourseSchema);

module.exports = Course;