/**
 * Created by igor on 25.08.16.
 */

"use strict";

const fileService = require(__appRoot + '/services/recordings'),
    log = require(__appRoot + '/lib/log')(module),
    utilsHttp = require(__appRoot + '/utils/http')
    ;

module.exports = {
    addRoutes: addRoutes
};

/**
 * Adds routes to the api.
 */
function addRoutes(api) {
    api.get('/api/v2/files/:id', getFile);
}

const getFile = (req, res, next) => {
    let uuid = req.params.id,
        contentType = req.query.type || 'audio/mpeg',
        pathName = req.query.name,
        dispositionName = req.query['file_name']
        ;

    let params = {
        contentType: contentType,
        pathName: pathName,
        range: req.headers['range']
    };

    fileService.getFileFromUUID(
        uuid,
        params,
        (err, response) => {
            if (err)
                return next(err);

            if (!response || !response.source)
                return next(`No source stream.`);

            let responseHeaders = {},
                source = response.source;
            if (dispositionName) {
                responseHeaders['Content-disposition'] = `attachment;  filename=${dispositionName}`;
            }

            if (!params.range) {
                responseHeaders['Content-Type'] = contentType;
                responseHeaders['Content-Length'] = response.totalLength || 0;
                responseHeaders['Accept-Ranges'] = 'bytes';
                return sendResponse(res, 200, responseHeaders, source);
            }

            let start = params.range.Start,
                end = params.range.End
                ;

            if (start >= response.totalLength || end >= response.totalLength) {
                responseHeaders['Content-Range'] = 'bytes */' + response.totalLength; // File size.
                return sendResponse(res, 416, responseHeaders, null);
            }

            responseHeaders['Content-Range'] = 'bytes ' + start + '-' + end + '/' + response.totalLength;
            responseHeaders['Content-Length'] = start == end ? 0 : (end - start + 1);
            responseHeaders['Content-Type'] = contentType;
            responseHeaders['Accept-Ranges'] = 'bytes';
            responseHeaders['Cache-Control'] = 'no-cache';
            return sendResponse(res, 206, responseHeaders, source);
        }
    );
};

function sendResponse(response, responseStatus, responseHeaders, readable) {
    //console.dir(responseStatus);
    response.writeHead(responseStatus, responseHeaders);

    if (readable == null)
        response.end();
    else
        readable.pipe(response);

    return null;
}