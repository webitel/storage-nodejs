/**
 * Created by igor on 29.08.16.
 */

"use strict";
    
module.exports = {
    LocalStorage: require('./local'),
    B2Storage: require('./b2'),
    S3Storage: require('./s3'),
    GDriveStorage: require('./gdrive'),
    helper: require('./helper')
};