/**
 * Created by I. Navrotskyj on 27.02.18.
 */

"use strict";

const CodeError = require(__appRoot + '/lib/error');
const log = require(`${__appRoot}/lib/log`)(module);

function add(pool) {
    return {
        getLegByUuid: (uuid, domain, legName, cb) => {
            pool.query(
                legName === "b"
                    ? `SELECT event FROM cdr_b WHERE uuid = $1`
                    : `SELECT event FROM cdr_a WHERE uuid = $1`,
                [uuid],
                (err, res) => {
                    if (err)
                        return cb(err);

                    if (res && res.rowCount) {
                        let data = {};
                        try {
                            data = res.rows[0].event.toString('utf8')
                        } catch (e) {
                            log.error(e)
                        }
                        return cb(null, JSON.parse(data))
                    } else {
                        return cb(new CodeError(404, `Not found ${uuid}`));
                    }
                }
            )
        },

        removeLegA: (uuid, cb) => {
            pool.query(
                `
                WITH lega as (
                  DELETE 
                  FROM cdr_a 
                  WHERE uuid = $1
                  RETURNING uuid
                ), legb as (
                  DELETE 
                  FROM cdr_b 
                  WHERE parent_uuid = $1
                  RETURNING uuid
                ) 
                SELECT (select count(*) FROM lega) + (select count(*) FROM legb) as removed
                `,
                [uuid],
                (err, res) => {

                    if (err)
                        return cb(err);

                    if (res && res.rowCount) {
                        return cb(null, +res.rows[0].removed)
                    } else {
                        return cb(new CodeError(404, `Not found ${uuid}`));
                    }
                }
            )
        }
    }
}


module.exports = add;
