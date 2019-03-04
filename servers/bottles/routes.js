const express = require('express');
const Oceans = require('./models/ocean');

const router = express.Router();
const fetch = require("node-fetch");

////////////////////////// /bottles ENDPOINT ROUTERS //////////////////////////

//create an ocean
router.post("/ocean", (req, res) => {
    //let user = JSON.parse(req.get("X-User"));

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

// create a bottle in the ocean
router.post("/ocean/:name", (req, res) => {
    if (req.body.body || req.body.body.length === 0) {
        res.status(403).send({error : "Cannot posts an empty bottle"});
    }

    // searching mongo to find ocean with given name
    Oceans.findOne({ "name" : req.params.name}).exec().then(ocean => {
        if (!ocean) {
            return res.status(404).send({error : "Ocean named " + req.params.name + " was not found"});
        }

        // create a new bottle
        // make it so people can only edit on their personal page ????
        let bottle = {
            //creator: user,
            body: req.body.body,
            tags: req.body.tags,
            createdAt: Date.now(),
            isPublic: req.body.isPublic
        };

        // save the bottle in the specific ocean
        ocean.bottles.push(bottle);
        ocean.save().then(() => {
            res.setHeader("Content-Type", "application/json");
            res.status(200).send(bottle);
        });

    });
}).catch(err => {
    res.status(400).send({error: "bottle couldn't be posted: " + err});
});

// update the bottle contents
router.patch("/ocean/:name/bottles/:id", (req, res) => {
    if (req.body.body || req.body.body.length === 0) {
        res.status(403).send({error : "Cannot posts an empty bottle"});
    }

    Oceans.findOne({"name" : req.params.name}).then(ocean => {
        let b = ocean.bottles.filter(bottle =>  bottle.__id == req.params.id);
    });

    // add seperate API endpoint for getting user type
});

router.delete("ocean/:name/bottles/:id", (req, res) => {
    
})





module.exports = router;