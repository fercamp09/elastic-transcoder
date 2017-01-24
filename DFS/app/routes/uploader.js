var express = require('express');
var formidable = require('formidable');
var router = express.Router();
var mongoose = require('mongoose');

module.exports = router;

router.post('/files', function(req, res){
  var form = new formidable.IncomingForm();
  form.parse(req);
  //console.log(./);
  form.on('fileBegin', function (name, file){
      file.path = __dirname + '/../public/uploads/' + Date.now() + '_' + file.name;
      console.log(file.path);
  });
  form.on('file', function (name, file){
      console.log('Uploaded ' + file.name);
  });
  res.json({message: "se subio el archivo con exito"});

});

router.get('/files/:id', function(req, res){
  console.log("aqui");
  res.redirect('/uploads/' + req.params.id);

});
