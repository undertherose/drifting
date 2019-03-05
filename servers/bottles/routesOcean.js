const express = require('express');
const Oceans = require('./models/ocean');

const router = express.Router();
const fetch = require("node-fetch");

////////////////////////// /ocean ENDPOINT ROUTERS //////////////////////////
// Creating an Ocean
// Deleting an Ocean
// Posting to an Ocean
// Fetch Request for creating tags

//create an ocean
router.post("/ocean", (req, res) => {
    //let user = JSON.parse(req.get("X-User"));
    // check if user is an admin

    //fetch request for get user by id to double check it's it's type

    Oceans.create({
        name: req.body.name
    }).then(ocean => {
        //insert rabbitMQ stuff
        ocean.save().then(() => {
            res.setHeader("Content-Type", "application/json");
            res.status(201).send(ocean);
        }).catch(err => {
            console.log(err);
        });

    }).catch(err => {
        console.log(err);
    });
}).catch(err => {
    res.status(400).send({error : "couldn't create a ocean " + err});
});


//delete an ocean
router.delete("/ocean/:name", (req, res) => {
    //let user = JSON.parse(req.get("X-User"));
    //if (user.type == "admin") { //only be able to delete if if person is the admin
        Oceans.findOneAndDelete({"name" : req.params.name, }, (err, response) => {

        }).then(res => {
            return res.json();
        }).catch(err => {
            console.log(err);
        });
        res.status(200).send({message: "ocean " + req.params.name + " was sucessfully deleted"});
    //} else {
    //  res.status(400).send({error : "You are not authorized to delete"});    
    //}

}).catch(err => {
    res.status(400).send({error: "couldn't delete ocean named " + req.params.name })
});


module.exports = router;