var mongoose = require('mongoose');

var FileSchema = new mongoose.Schema({
	original_name: String,
  new_name: String
});

module.exports = mongoose.model('File', FileSchema);
