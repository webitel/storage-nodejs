/**
 * Created by igor on 21.10.16.
 */

"use strict";

module.exports = api => {
    require('./acl').addRoutes(api);
    require('./cdr').addRoutes(api);
    require('./recordFile').addRoutes(api);
};