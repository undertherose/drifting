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




module.exports = router;