const express = require('express');
const Courses = require('./models/course');

const router = express.Router();
const fetch = require("node-fetch");

////////////////////////// /courses ENDPOINT ROUTERS //////////////////////////

//create a course
router.post('/courses/create', (req, res) => {

    let user = JSON.parse(req.get('X-User'));
    let ch = req.app.get('ch');

    Courses.create({
        sln: req.body.sln,
        name: req.body.courseName,
    }).then(course => {
        let ta = {
            id: user.id,
            username: user.userName
        }
        course.ta.push(ta);
        course.save().then(() => {

            let qPayLoad = {};
            qPayLoad.type = 'course-new';
            qPayLoad.message = course;

            ch.sendToQueue(
                "rabbitmq",
                new Buffer.from(JSON.stringify(qPayLoad)),
                { persistent: true }
            );

            fetch("https://api.iqueue.zubinchopra.me/final/faq", {
                method: 'POST',
                headers: {
                    "Content-Type": "application/json",
                    'Authorization' : req.headers.authorization  
                },
                body: JSON.stringify({sln: req.body.sln})
            }).then(res => {
                //parses json data and returns a promise
                return res.json();

            }).catch(err => {
                console.log(err);
            });

            res.setHeader('Content-Type', 'application/json');
            res.status(201).send(course);
        }).catch(err => {
            console.log(err);
        })
    }).catch(err => {
        res.status(400).send({error: "couldn't create a course" + err});
    });

});

//gets courses for a specific user
router.get('/courses/me', (req, res) => {

    let user = JSON.parse(req.get('X-User'));

    Courses.find({ $or: [{ 'ta.id': user.id }, { 'students.id': user.id }] }).exec().then(courses => {
        res.setHeader('Content-Type', 'application/json');
        res.status(200).send(courses);
    }).catch(err => {
        res.send(500).send({error: "couldn't get courses for current user" + err });
    });

});

//get all courses
router.get('/courses/all', (req, res) => {

    Courses.find({}).exec().then(courses => {
        res.setHeader('Content-Type', 'application/json');
        res.status(200).send(courses);
    }).catch(err => {
        res.send(500).send({error: "couldn't get all existing courses" + err});
    });

});

//get course struct
router.get('/courses/:sln', (req, res) => {

    let user = JSON.parse(req.get('X-User'));

    //finding the course we're at from the db given course sln
    //returns course struct
    Courses.find({ "sln": req.params.sln }).exec().then(course => {

        if (!course) {
            return res.status(404).send({ error: "Course with given sln" + req.params.sln + "does not exist"});
        }

        res.setHeader('Content-Type', 'application/json');
        res.status(200).send(course);
    }).catch(err => {
        res.status(400).send({error: "couldn't get course for given sln" + err});
    });

});

//post (name, question) message to queue 
router.post('/courses/:sln', (req, res) => {

    let user = JSON.parse(req.get('X-User'));
    let ch = req.app.get('ch');

    if (!req.body.body || req.body.body.length === 0) {
        res.status(403).send({ error: "Cannot send an empty question" });
    }

    Courses.findOne({ "sln": req.params.sln }).exec().then(course => {

        if (!course) {
            return res.status(404).send({ error: "Course with given sln" + req.params.sln + "does not exist"});
        }

        if (course.isOpen === false) {
            res.status(400).send({ error: "office hours are not in session right now" });
        } else {

            let studentEnrolled = course.students.filter(student => student.id == user.id);

            if (studentEnrolled.length != 0) {
                let check = course.queuePosts.filter(post => post.creator.id === user.id);

                if (check.length === 0) {
                    let post = {
                        creator: user,
                        body: req.body.body,
                        createdAt: Date.now()
                    };

                    course.queuePosts.push(post);

                    course.save().then(() => {

                        let qPayLoad = {};
                        qPayLoad.type = 'post-new';
                        qPayLoad.message = post;

                        ch.sendToQueue(
                            "rabbitmq",
                            new Buffer.from(JSON.stringify(qPayLoad)),
                            { persistent: true }
                        );

                        res.setHeader('Content-Type', 'application/json');
                        res.status(200).send(post);
                    });
                } else {
                    res.status(400).send({ error: "you are already in the queue" });
                }
            } else {
                res.status(403).send({ error: "you are not enrolled in this class" });
            }
        }

    }).catch(err => {
        res.status(400).send({error: "couldn't add question to the course queuePosts " + err});
    });

});

//edit course name (only TA privilage)
router.patch('/courses/:sln', (req, res) => {

    let user = JSON.parse(req.get('X-User'));
    let ch = req.app.get('ch');

    //course name to be updated must not be empty
    if (req.body.name.length === 0) {
        res.status(403).send({ error: "New course name cannot be empty" });
    }

    Courses.findOne({ "sln": req.params.sln }).exec().then(course => {

        if (!course) {
            return res.status(404).send({ error: "Course with given sln" + req.params.sln + "does not exist"});
        }

        //check that the user who is trying to edit 
        let check = course.ta.filter(ta => ta.id === user.id);

        if (check.length !== 0) {
            course.name = req.body.name;

            course.save().then(() => {

                let qPayLoad = {};
                qPayLoad.type = 'course-update';
                qPayLoad.course = course;

                ch.sendToQueue(
                    "rabbitmq",
                    new Buffer.from(JSON.stringify(qPayLoad)),
                    { persistent: true }
                );
                res.setHeader('Content-Type', 'application/json');
                res.status(200).send(course);

            });
        } else {
            res.status(403).send({ error: 'you are not authorized to edit a channel' });
        }

    }).catch(err => {
        res.status(404).send({error: "couldn't update class name " + err});
    });

});

//delete course
router.delete('/courses/:sln', (req, res) => {
    let user = JSON.parse(req.get('X-User'));
    let ch = req.app.get('ch');


    Courses.findOne({ "sln": req.params.sln }).then(course => {

        if (!course) {
            return res.status(404).send({ error: "Course with given sln" + req.params.sln + "does not exist"});
        }

        let checkCurrentTa = course.ta.filter(ta => ta.username === user.userName);

        if (checkCurrentTa.length !== 0) {

            Courses.findOneAndDelete({ "sln": req.params.sln, }, (err, response) => {

                let qPayLoad = {};
                qPayLoad.type = 'course-delete';
                qPayLoad.channel = response;
                ch.sendToQueue(
                    "rabbitmq",
                    new Buffer.from(JSON.stringify(qPayLoad)),
                    { persistent: true }
                );

                fetch("https://api.iqueue.zubinchopra.me/final/faq/"+req.params.sln, {
                    method: 'DELETE',
                    headers: {
                        "Content-Type": "application/json",
                        'Authorization' : req.headers.authorization  
                    }
                }).then(res => {
                    //parses json data and returns a promise
                    return res.json();
    
                }).catch(err => {
                    console.log(err);
                });

                res.status(200).send({ message: 'delete was successful' });

            });
        } else {
            res.status(400).send({ error: 'you are not a ta' });
        }

    }).catch(err => {
        res.status(400).send({error: "couldn't delete course" + err});
    });
});

////////////////////////// /courses/{sln}/tas ENDPOINT ROUTERS //////////////////////////

//add tas
router.post('/courses/:sln/tas', (req, res) => {

    let user = JSON.parse(req.get('X-User'));
    let ch = req.app.get('ch');

    //check that both of them get passed in
    if (!req.body.username) {
        res.status(400).send({ error: "username cannot be empty" });
        return;
    }

    //check both of them are not empty
    if (req.body.username.length === 0) {
        res.status(400).send({ error: "username cannot be empty" });
        return;
    }

    Courses.findOne({ "sln": req.params.sln }).then(course => {

        if (!course) {
            return res.status(404).send({ error: "Course with given sln" + req.params.sln + "does not exist"});
        }

        let checkCurrentTa = course.ta.filter(ta => ta.username === user.userName);

        if (checkCurrentTa.length !== 0) {

            //make sure that the requested user from reqbody is an existing, authorized person
            fetch("https://api.iqueue.zubinchopra.me/final/allusers").then(res => {
                //parses json data and returns a promise
                return res.json();
            }).then(data => {

                let userCheck = data.filter(person => person.userName === req.body.username);

                if (userCheck.length !== 0) {

                    let check = course.ta.filter(ta => ta.username === req.body.username);
                    if (check.length !== 0) {
                        res.status(400).send({ error: 'TA already has access to this course' });
                        return;
                    } else {
                        course.ta.push({ id: userCheck[0].id, username: req.body.username });
                        course.save().then(() => {
                            res.setHeader('Content-Type', 'application/json');
                            res.status(201).send(course);
                            

                            fetch("https://api.iqueue.zubinchopra.me/final/faq/"+req.params.sln+'/tas', {
                                method: 'POST',
                                headers: {
                                    "Content-Type": "application/json",
                                    'Authorization' : req.headers.authorization  
                                },
                                body: JSON.stringify({username: req.body.username})
                            }).then(res => {
                                //parses json data and returns a promise
                                return res.json();
                
                            }).catch(err => {
                                console.log(err);
                            });


                            let qPayLoad = {};
                            qPayLoad.type = 'ta-new';
                            qPayLoad.message = course;
                
                            ch.sendToQueue(
                                "rabbitmq",
                                new Buffer.from(JSON.stringify(qPayLoad)),
                                { persistent: true }
                            );
                        });

                    }

                } else {
                    res.status(400).send({ error: 'TA entered does not exist' });
                    return;
                }


            }).catch(err => {
                console.log(err);
            });

        } else {
            res.status(403).send({ error: 'you are not a TA' });
            return;
        }

    }).catch(err => {
        res.status(403).send({error: "couldn't add TA into the course" + err});
    });
});

//remove tas
router.delete('/courses/:sln/tas', (req, res) => {

    let user = JSON.parse(req.get('X-User'));
    let ch = req.app.get('ch');


    if (!req.body.username) {
        res.status(400).send({ error: "username cannot be empty" });
        return;
    }

    if (req.body.username.length === 0) {
        res.status(400).send({ error: "username cannot be empty" });
        return;
    }

    Courses.findOne({ "sln": req.params.sln }).then(course => {

        if (!course) {
            return res.status(404).send({ error: "Course with given sln" + req.params.sln + "does not exist"});
        }

        let checkCurrentTa = course.ta.filter(ta => ta.username == user.userName);

        if (checkCurrentTa.length !== 0) {

            let originalLength = course.ta.length;

            let ta = course.ta.filter(ta => ta.username != req.body.username);


            //if ta exists, remove it from the list
            if (ta.length != 0) {
                course.ta = ta;

                if (course.ta.length != originalLength) {

                    course.save().then(() => {

                        fetch("https://api.iqueue.zubinchopra.me/final/faq/"+req.params.sln+'/tas', {
                            method: 'DELETE',
                            headers: {
                                "Content-Type": "application/json",
                                'Authorization' : req.headers.authorization  
                            },
                            body: JSON.stringify({username: req.body.username})
                        }).then(res => {
                            //parses json data and returns a promise
                            return res.json();
            
                        }).catch(err => {
                            console.log(err);
                        });


                        let qPayLoad = {};
                        qPayLoad.type = 'ta-delete';
                        qPayLoad.message = course;
            
                        ch.sendToQueue(
                            "rabbitmq",
                            new Buffer.from(JSON.stringify(qPayLoad)),
                            { persistent: true }
                        );

                        res.status(200).send(course);
                    });
                } else {
                    res.status(400).send({ error: "TA does not exist in this course" });
                }
            } else {
                res.status(400).send({ error: "Could not remove the TA" });
            }
        } else {
            res.status(403).send({ error: 'you are not a TA' });
            return;
        }
    }).catch(err => {
        res.status(404).send({error: "couldn't remove TA from the course" + err});
    });
});


////////////////////////// /courses/{sln}/enrollment ENDPOINT ROUTERS //////////////////////////
router.post('/courses/:sln/enrollment', (req, res) => {

    let user = JSON.parse(req.get('X-User'));
    let ch = req.app.get('ch');

    Courses.findOne({ "sln": req.params.sln }).then(course => {

        if (!course) {
            return res.status(404).send({ error: "Course with given sln" + req.params.sln + "does not exist"});
        }

        let check = course.students.filter(student => student.username === user.userName);

        if (check.length !== 0) {
            res.status(400).send({ error: ' you are already enrolled in this course' });
            return;
        } else {
            course.students.push({ id: user.id, username: user.userName });
            course.save().then(() => {
                let qPayLoad = {};
                qPayLoad.type = 'student-enrolled';
                qPayLoad.message = course;

                ch.sendToQueue(
                    "rabbitmq",
                    new Buffer.from(JSON.stringify(qPayLoad)),
                    { persistent: true }
                );
                res.setHeader('Content-Type', 'application/json');
                res.status(200).send(course);
            });
        }

    }).catch(err => {
        res.status(403).send({error: "couldn't add student into the course" + err});
    });

});



router.delete('/courses/:sln/enrollment', (req, res) => {

    let user = JSON.parse(req.get('X-User'));
    let ch = req.app.get('ch');

    Courses.findOne({ "sln": req.params.sln }).then(course => {

        if (!course) {
            return res.status(404).send({ error: "Course with given sln" + req.params.sln + "does not exist"});
        }

        let student = course.students.filter(student => student.username !== user.userName);

        //if student exists, remove it from the list
        // compare to course.students.length since it's possible to have course with noone enrolled
        if (student.length !== course.students.length) {
            course.students = student;

            course.save().then(() => {
                let qPayLoad = {};
                qPayLoad.type = 'student-disenrolled';
                qPayLoad.message = Courses;

                ch.sendToQueue(
                    "rabbitmq",
                    new Buffer.from(JSON.stringify(qPayLoad)),
                    { persistent: true }
                );
                res.setHeader('Content-Type', 'application/json');
                res.status(200).send(Courses);
            });

        } else {
            res.status(400).send({ error: user.userName + ' is not enrolled in course so cant be removed' });
            return;
        }

    }).catch(err => {
        res.status(403).send({error: "couldn't remove student from the course" + err});
    });

});



////////////////////////// /queueposts/{queuepostID} ENDPOINT ROUTERS //////////////////////////

// update question from queue
router.patch('/courses/:sln/questions/:id', (req, res) => {

    let user = JSON.parse(req.get('X-User'));
    let ch = req.app.get('ch');

    if (!req.body.body) {
        res.status(400).send({ error: "Cannot send an empty update to question" });
    }

    if (req.body.body.length === 0) {
        res.status(400).send({ error: "Cannot send an empty update to question" });
    }

    Courses.findOne({ "sln": req.params.sln }).then(course => {

        if (!course) {
            return res.status(404).send({ error: "Course with given sln" + req.params.sln + "does not exist"});
        }

        let studentEnrolled = course.students.filter(student => student.id == user.id);

        if (studentEnrolled.length != 0) {
            let q = course.queuePosts.filter(queuePost => queuePost.creator.id == req.params.id);

            if (user.id == q[0].creator.id) {

                q[0].body = req.body.body;

                course.save().then(() => {
                    let qPayLoad = {};
                    qPayLoad.type = 'queuePost-update';
                    qPayLoad.message = course;
    
                    ch.sendToQueue(
                        "rabbitmq",
                        new Buffer.from(JSON.stringify(qPayLoad)),
                        { persistent: true }
                    );

                    res.status(200).send(q[0]);
                });

            } else {
                res.status(403).send({ error: "you are not authorized to update this post" });
            }
        } else {
            res.status(403).send({ error: "you are not enrolled in the class" });
        }

    }).catch(err => {
        res.status(400).send({error: "unable to update question in queue " + err});
    });

});


// delete question from queue
router.delete('/courses/:sln/questions/:id', (req, res) => {

    let user = JSON.parse(req.get('X-User'));
    let ch = req.app.get('ch');

    // find a course given the sln
    Courses.findOne({ "sln": req.params.sln }).exec().then(course => {

        if (!course) {
            return res.status(404).send({ error: "Course with given sln" + req.params.sln + "does not exist"});
        }

        // check if logged in user is a TA
        let check = course.ta.filter(ta => ta.id === user.id);

        let studentEnrolled = course.students.filter(student => student.id == user.id);

        // check if logged in user is question creator 
        // OR the TA, delete question 

        let queuePost = course.queuePosts.filter(queuePost => queuePost.creator.id == req.params.id);

        if (queuePost.length == 0) {
            res.status(400).send({ error: "no post to be deleted" });
            return;
        }

        if ((user.id == queuePost[0].creator.id && studentEnrolled.length != 0) || check.length !== 0) {

            let originalLength = course.queuePosts.length;

            course.queuePosts = course.queuePosts.filter(queuePost => queuePost.creator.id != req.params.id);

            if (originalLength != course.queuePosts.length) {
                let qPayLoad = {};
                qPayLoad.type = 'queuePost-delete';

                course.save().then(() => {
                    qPayLoad.message = course;

                    ch.sendToQueue(
                        "rabbitmq",
                        new Buffer.from(JSON.stringify(qPayLoad)),
                        { persistent: true }
                    );

                    res.status(200).send({ message: 'queue entry was deleted' });
                    return;
                });


            } else {
                res.status(400).send({ error: "post wasn't deleted" });
                return;
            }

        } else {
            res.status(403).send({ error: "you are not authorized to delete this post" });
            return;
        }
    }).catch(err => {
        res.status(400).send({error: "unable to delete question from queue " + err});
    });

});


////////////////////////////////handle open/close office hours///////////////////////////////

router.patch('/courses/:sln/status', (req, res) => {

    let user = JSON.parse(req.get('X-User'));
    let ch = req.app.get('ch');

    Courses.findOne({ "sln": req.params.sln }).exec().then(course => {

        if (!course) {
            return res.status(404).send({ error: "Course with given sln" + req.params.sln + "does not exist"});
        }

        //check that the user who is a TA that can open/close sessions
        let check = course.ta.filter(ta => ta.id === user.id);

        if (check.length !== 0) {

            if (course.isOpen === true) {
                course.isOpen = false;
                course.queuePosts = [];
            } else {
                course.isOpen = true;
            }

            course.save().then(() => {

                let qPayLoad = {};
                qPayLoad.type = 'course-status-update';
                qPayLoad.course = course;

                ch.sendToQueue(
                    "rabbitmq",
                    new Buffer.from(JSON.stringify(qPayLoad)),
                    { persistent: true }
                );
                res.setHeader('Content-Type', 'application/json');
                res.status(200).send(course);

            });

        } else {
            res.status(403).send({ error: 'you are not authorized to start/end office hour session' });
        }

    }).catch(err => {
        res.status(404).send({error: "couldn't open/close office hours" + err});
    });

});


module.exports = router;