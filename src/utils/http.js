/**
 * Created by igor on 26.08.16.
 */

"use strict";
    
const helper = module.exports = {
    readRangeHeader: (range, totalLength) => {
        if (range == null || range.length == 0)
            return null;

        var array = range.split(/bytes=([0-9]*)-([0-9]*)/);
        var start = parseInt(array[1]);
        var end = parseInt(array[2]);
        var result = {
            Start: isNaN(start) ? 0 : start,
            End: isNaN(end) ? (totalLength - 1) : end
        };

        if (!isNaN(start) && isNaN(end)) {
            result.Start = start;
            result.End = totalLength - 1;
        }

        if (isNaN(start) && !isNaN(end)) {
            result.Start = totalLength - end;
            result.End = totalLength - 1;
        }

        return result;
    },
    
    streaming: (source, response, params) => {
        let responseHeaders = {};

        if (source && source.location && source.statusCode == 302) {
            responseHeaders['Location'] = source.location;
            responseHeaders['Content-Type'] = params.contentType;
            return helper.sendResponse(response, 302, responseHeaders);
        }

        if (params.dispositionName) {
            responseHeaders['Content-disposition'] = `attachment;  filename=${params.dispositionName}`;
        }

        if (!params.range) {
            responseHeaders['Content-Type'] = params.contentType;
            responseHeaders['Content-Length'] = params.totalLength || 0;
            responseHeaders['Accept-Ranges'] = 'bytes';
            return helper.sendResponse(response, 200, responseHeaders, source);
        }

        let start = params.range.Start,
            end = params.range.End
            ;

        if (start >= params.totalLength || end >= params.totalLength) {
            responseHeaders['Content-Range'] = 'bytes */' + params.totalLength; // File size.
            return helper.sendResponse(response, 416, responseHeaders, null);
        }

        responseHeaders['Content-Range'] = 'bytes ' + start + '-' + end + '/' + params.totalLength;
        responseHeaders['Content-Length'] = start == end ? 0 : (end - start + 1);
        responseHeaders['Content-Type'] = params.contentType;
        responseHeaders['Accept-Ranges'] = 'bytes';
        responseHeaders['Cache-Control'] = 'no-cache';
        return helper.sendResponse(response, 206, responseHeaders, source);
    },

    sendResponse: (response, responseStatus, responseHeaders, readable) => {
        response.writeHead(responseStatus, responseHeaders);
        // console.log(responseStatus, responseHeaders);

        if (readable == null) {
            response.end();
        }
        else
            readable.pipe(response);

        return null;
    }
};