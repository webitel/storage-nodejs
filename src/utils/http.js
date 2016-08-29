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
    }
};