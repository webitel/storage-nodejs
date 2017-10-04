/**
 * Created by I. Navrotskyj on 04.10.17.
 */

"use strict";

const tcpService = require(__appRoot + '/services/tcpDump');
const streaming = require(`${__appRoot}/utils/http`).streaming;

module.exports = {
    addRoutes: addRoutes
};

function addRoutes(api) {
    api.get('/api/v2/tcp_dump/:id/data', getFile);
}

function getFile(req, res, next) {
    const options = {
        id: req.params.id
    };

    tcpService.getFile(req.webitelUser, options, (err, response) => {
        if (err)
            return next(err);

        if (!response || !response.source)
            return next(`No source stream.`);

        return streaming(response.source, res, {
            dispositionName: req.params.id,
            totalLength: response.totalLength,
            contentType: response.contentType || "application/octet-stream"
        });
    })
}