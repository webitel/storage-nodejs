/**
 * Created by igor on 09.09.16.
 */

"use strict";

const google = require('googleapis'),
    fs = require('fs'),
    helper = require('./helper'),
    log = require(__appRoot + '/lib/log')(module),
    TYPE_ID = 3;

const SCOPES = [
    'https://www.googleapis.com/auth/plus.me',
    'https://www.googleapis.com/auth/drive'
];

module.exports = class GoogleStorage {
    constructor (config = {}, mask) {
        this.folderId = null;
        this.name = "gDrive";

        this._client_email = config.client_email;
        this._private_key = config.private_key;
        this._folderName = config.folderName;

        this._auth =  new google.auth.JWT(this._client_email, null, this._private_key, SCOPES, null);

        this._auth.authorize((err, tokens) => {
            if (err) {
                return log.error(err);
            }

            this.drive = google.drive({ version: 'v3', auth: this._auth});
            this.drive.files.list({
                q: `name='${this._folderName}' AND mimeType='application/vnd.google-apps.folder' AND trashed=false`,
                pageSize: 1,
                fields: 'files(id)',
                spaces: 'drive'
            }, (err, res) => {
                if (err)
                    return log.error(err);

                if (res.files.length == 0)
                    return this._createFolder(this._folderName);

                this.folderId = res.files[0].id;
                log.debug(`Set folder id: ${this.folderId} from name ${this._folderName}`);
            })
        });

    }

    _createFolder (name) {
        this.drive.files.create({
            resource: {
                name: name,
                mimeType: 'application/vnd.google-apps.folder'
            }
        }, (err, res) => {
            if (err)
                return log.error(err);
            log.debug(`Create folder name: ${name} (${res.id})`);
            this.folderId = res.id;
        })
    }

    get (fileDb, options, cb) {
        if (!this.folderId)
            return cb(new Error(`No initialize folder id.`));

        let stream = this.drive.files.get({
            fileId: fileDb.storageFileId,
            alt: "media"
        }).on('error', (err) => log.error(err));

        cb(null, stream);

    }

    getFilePath(domain, fileName) {

    }

    save (fileConf, option, cb) {
        if (!this.folderId)
            return cb(new Error(`No initialize folder id.`));

        this.drive.files.create({
            fields: 'id',
            resource: {
                name: fileConf.name,
                parents: [this.folderId]
            },
            media: {
                body: fs.createReadStream(fileConf.path),
                mimeType: fileConf.contentType
            }
        }, (err, res) => {
            if (err)
                return cb(err);

            log.trace(`Save file ${fileConf.name} to ${res.id}`);
            return cb(null, {
                path: fileConf.name,
                type: TYPE_ID,
                bucketName: this.folderId,
                storageFileId: res.id
            })
        })
    }

    del (fileDb, cb) {
        if (!this.folderId)
            return cb(new Error(`No initialize folder id.`));

        this.drive.files.delete({
            fileId: fileDb.storageFileId
        }, (err, res) => {
           if (err) {
                log.error(err);
                return cb(err);
            }
            log.trace(`Delete file ${fileDb.storageFileId}`);
            return cb(null, res)
        });
    }

    existsFile (fileDb, cb) {
        if (!this.folderId)
            return cb(new Error(`No initialize folder id.`));

        this.drive.files.get({
            fileId: fileDb.storageFileId
        }, (err, res) => {
            if (err && err.code === 404) {
                return cb(null, false)
            } else if (err) {
                log.error(err);
                return cb(err);
            }

            return cb(null, !!res)
        });
    }

    checkConfig (config = {}, mask) {
        return this._client_email == config.client_email && this._private_key == config.private_key
            && this._folderName == config.folderName;
    }
};