const express = require('express');
const Oceans = require('./models/ocean');

const router = express.Router();
const fetch = require("node-fetch");

////////////////////////// /bottles ENDPOINT ROUTERS //////////////////////////

// create a bottle in the ocean
router.post("/ocean/:name", (req, res) => {
    if (req.body.body || req.body.body.length === 0) {
        res.status(403).send({ error: "Cannot posts an empty bottle" });
    }

    // searching mongo to find ocean with given name
    Oceans.findOne({ "name": req.params.name }).exec().then(ocean => {
        if (!ocean) {
            return res.status(404).send({ error: "Ocean named " + req.params.name + " was not found" });
        }

        //go through each tag in the req body
        // make post request to make all the tags
        // if the post is public, increase the tag count
        
        
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
    res.status(400).send({ error: "bottle couldn't be posted: " + err });
});



// update the bottle contents
router.patch("/ocean/:name/bottles/:id", (req, res) => {
    //get the xuser stuff

    // cannot have both the body and the tags in update be empty
    /*     if ((!req.body.body || req.body.body.length === 0) && (!req.body.tags || req.body.tags.length === 0)) {
            return res.status(403).send({error : "Cannot have no update body and no tags"});
        } */


    Oceans.findOne({ "name": req.params.name }).then(ocean => {
        if (!ocean) {
            //return not found mesage
        }

        let bottle = ocean.bottles.filter(bottle => bottle.__id == req.params.id);

        // check if you wrote the bottle

        if (!bottle) {
            //error message about non-existing bottle
        }


        // cannot have empty update
        if (req.body.body && req.body.body.length > 0) {
            bottle[0].body = req.body.body;
        }
    
        // cannot have no tags 
        if (req.body.tags && req.body.tags.length > 0) {


            if (bottle.isPublic) {
                
                // old tags == existing tags
                // remove all the old tags
                oldTags = bottle[0].tags;
                oldTags.forEach(function (t) {
                    // go through each of the tag in the original
                    // decrease the count for all the tags
                    currTag = bottle[0].tags.filter(tag => tag.tagName == t)
                    if (currTag.tagCount != 0) {
                        currTag.tagCount -= 1
                    }
                });
    
                newTags = req.body.tags;
                newTags.forEach(function(t) {
                    //go through all the new tags
                    // update the counts of all the tags
                    // udate the last update time


                    // if new tag: non-existent tag:
                    // > create a new tag and update the count
                    currTag = bottle[0].tags.filter(tag => tag.tagName == t)
                    if (currTag.tagCount != 0) {
                        currTag.tagCount += 1
                        currTag.tagLastUpdate = Date.now()
                    }
                });

            }




            


            // updating the tags
            bottle[0].tags = req.body.tags;

            ocean.save().then(() => {

            });

        }

        //if bottle.creator.id == currUser

    });

});

// deleting a bottle
router.delete("ocean/:name/bottles/:id", (req, res) => {
    Oceans.findOne({ "name": req.params.name }).exec().then(ocean => {

        if (!ocean) { // did not find a ocean with given name
            return res.status(404).send({ error: "ocean with given name was not found" })
        }

        // check if there is a bottle with that id
        let bottle = ocean.bottles.filter(b => b.__id == req.params.id);
        if (bottle.length == 0) {
            return res.status(400).send({ error: "No bottle found with that id" });
        }
        // check if you are a moderator account or that user


        let originalLenth = ocean.bottles.length;
        ocean.bottles = ocean.bottles.filter(bottle => bottle.__id != req.params.id); //filter out bottles that match the id

        if (originalLength != ocean.bottles.length) {

            ocean.save().then(() => {
                return res.status(200).send({ message: "bottle was sucessfully deleted" });
            });
        }

    }).catch(err => {
        return res.status(400).send({ error: "Unable to delete bottle: " + err });
    });
});





module.exports = router;