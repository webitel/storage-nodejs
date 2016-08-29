/**
 * Created by igor on 29.08.16.
 */

"use strict";
    
module.exports = {
    getPath: (mask, domain, fileName) => {
        let date = new Date();
        let path = mask || '';
        path = path
            .replace(/\$DOMAIN/g, domain || "_unknown_")
            .replace(/\$Y/g, date.getFullYear())
            .replace(/\$M/g, (date.getMonth() + 1))
            .replace(/\$D/g, date.getDate())
            .replace(/\$H/g, date.getHours())
            .replace(/\$m/g, date.getMinutes())
        ;

        if (fileName)
            return `${path}/${fileName}`;
        else return path;
    }
};