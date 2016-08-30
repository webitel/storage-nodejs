/**
 * Created by igor on 30.08.16.
 */

"use strict";

module.exports = function checkPermission (acl, resource, action, ignoreAllPerm) {
    if (!acl || !resource || !action || !(acl[resource] instanceof Array))
        return false;

    if (ignoreAllPerm)
        return !!~acl[resource].indexOf(action) && !~acl[resource].indexOf('*');

    return !!~acl[resource].indexOf(action) || !!~acl[resource].indexOf('*');
};