/**
 * Created by igor on 25.08.16.
 */

"use strict";

const log = require(__appRoot + '/lib/log')(module),
    recordingsService = require(__appRoot + '/services/recordings'),
    emailService = require(__appRoot + '/services/email'),
    fileService = require(__appRoot + '/services/file'),
    async = require('async')
    ;

module.exports = {
    addRoutes: addRoutes
};

/**
 * Adds routes to the api.
 */
function addRoutes(api) {
    api.put('/sys/formLoadFile?:id', saveFile); // deprecated
    api.put('/sys/records', saveFile);
}

function saveFile(req, res, next) {
    let uuid = req.query.id,
        name = req.query.name || 'tmp',
        type = req.query.type,
        domainName = req.query.domain,
        email = req.query.email,
        sendMail = email && email !== 'none',
        subject = req.query.email_sbj,
        text = req.query.email_msg
        ;

    fileService.requestStreamToCache( `${uuid}_${name}.${type || 'mp3'}`, req, (err, file) => {
        if (err)
            return next(err);

        file.domain = domainName;
        file.uuid = uuid;
        file.queryName = name;

        async.parallel(
            [
                (cb) => {
                    if (sendMail) {
                        let attachments = [{
                                path: file.path,
                                filename: file.name
                            }]
                            ;

                        if (!subject || subject === "none") {
                            if (type === 'pdf') {
                                subject = '[webitel] You have received a new fax';
                            } else {
                                subject = '[webitel] You have received a new call record file';
                            }
                        }

                        if (!text || text === "none") {
                            if (type === 'pdf') {
                                text = 'You have received a new fax from Webitel Fax Server\n\n--\nWebitel Cloud Platform';
                            } else {
                                text = 'You have received a new call record file from Webitel\n\n--\nWebitel Cloud Platform';
                            }
                        }

                        return emailService.send(
                            domainName,
                            {
                                to: email,
                                subject: subject,
                                text: text,
                                attachments: attachments
                            },
                            (err) => {
                                if (err)
                                    log.error(err);

                                cb(null)
                            }
                        )
                    } else {
                        cb();
                    }
                },

                (cb) => {
                    recordingsService.saveFile(file, req.query, cb);
                }
            ],
            (err) => {

                fileService.deleteFile(file.path, (err) => {
                    if (err)
                        log.error(err);
                });

                if (err)
                    return next(err);

                res.status(204).end();

            }
        );
        
    });
}