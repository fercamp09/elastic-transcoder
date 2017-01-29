var express = require('express');
var formidable = require('formidable');
var router = express.Router();
var mongoose = require('mongoose');
var File = require('../models/File.js');
var fs = require("fs");

module.exports = router;

//Metodo que agrega uno o varios archivos a la base de datos
//Archivos van en form-data tipo file en el body
router.post('/files', function(req, res){
  var files = [];
  var form = new formidable.IncomingForm();
  var files_metadata = [];
  form.parse(req);

  form.on('fileBegin', function (name, file){
    var new_name = Date.now() + '_' + file.name;
    file.path = __dirname + '/../public/uploads/' + new_name;

    var file = new File({
      original_name: file.name,
      new_name: new_name
    });

    file.save(function(err, doc){
    		if (err) {
    			return next(err);
    		} else {
          files_metadata.push(doc);
          if (files.length == files_metadata.length) {
            res.json(files_metadata);
          }
    		}
    });

  });

  form.on('file', function(field, file) {
      files.push({"field": field, "file": file});
  })
  form.on('end', function() {
    console.log(files.length + ' archivos agregados');
  });

});

//Obtienes un archivo
router.get('/files/:id', function(req, res){
  File.findOne({_id:req.params.id}, function(err, doc){
		if (err) {
			return next(err);
		} else if (doc) { //validacion de la existencia metadata
      var path = __dirname + '/../public/uploads/' + doc.new_name;
      fs.access(path, function(err){ //validacion de la existencia de la data
        if (err) {
          res.json({message: 'this file does not exist'});
        } else {
          var readStream = fs.createReadStream(path);
    		  readStream.pipe(res);
        }
      });
		} else {
		    res.json({message: 'this file does not exist'});
		}
	});
});

//Obtienes una lista de todos los archivos del sistema
router.get('/files', function(req, res){
  File.find({}, function(err, docs){
		if (err) {
			return next(err);
		}
		res.json(docs);
	});
});


//Deputa la metadata del sistema, elimina los registros que no tienen un archivo asociado.
router.post('/files/depure', function(req, res){
  File.find({}, function(err, docs){
		if (err) {
			return next(err);
		} else {
      docs.forEach(function(doc){
        fs.access(__dirname + '/../public/uploads/' + doc.new_name, function(err){
            if (err) {
              File.remove({new_name: doc.new_name}, function(err){
                  if (err) {
                    console.log(err);
                  } else {
                    console.log(doc.new_name + ' deleted from metadata');
                  }
              });
            }
        });
      });
      res.json({message: "DFS has been depured"});
		}
	});
});

//Elimina un archivo
router.delete('/files/:id', function(req, res){
  File.findOne({_id:req.params.id}, function(err, doc){
		if(err){
			res.send(err);
		} else if (doc) {
      fs.unlinkSync(__dirname + '/../public/uploads/'+ doc.new_name);
      doc.remove();
      console.log('\n' + doc.new_name + ' has been deleted');
      res.json({message: 'this file has been deleted', file: doc});
    } else {
		    res.json({message: 'this file does not exist'});
		}
  });
});

//actualiza el archivo, se tiene que enviar solamente 1 archivo.
router.put('/files/:id', function(req, res){
  var form = new formidable.IncomingForm();
  form.parse(req);

  File.findOne({_id:req.params.id}, function(err, doc){
		if(err){
			res.send(err);
		} else if (doc) {
      console.log("aqui");
      form.on('fileBegin', function (name, file){
        console.log("aqui");
        var new_name = Date.now() + '_' + file.name; //Nombre del archivo actualizado
        var file_to_delete = doc.new_name; //Nombre del archivo antiguo
        file.path = __dirname + '/../public/uploads/' + new_name;
        doc.new_name = new_name;

        doc.save(function(err){
        		if (err) {
        			return next(err);
        		} else {
              fs.unlinkSync(__dirname + '/../public/uploads/'+ file_to_delete);
              console.log(file_to_delete + ' changed to ' + new_name);
              res.json({message: "You should send only one file", result: "File updated correctly"});
        		}
        });
      });
    } else {
		    res.json({message: 'this file does not exist'});
		}
  });
});
