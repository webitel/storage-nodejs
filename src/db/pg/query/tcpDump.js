/**
 * Created by I. Navrotskyj on 04.10.17.
 */

"use strict";

const log = require(__appRoot + '/lib/log')(module),
    CodeError = require(__appRoot + '/lib/error');


const setFileConf = `
    INSERT INTO tcp_dump (id, meta_file, description)
    VALUES ($1, $2, $3)
    ON CONFLICT (id) DO UPDATE
        SET meta_file = $4;
`;

const getMetadata = `
    SELECT meta_file FROM tcp_dump WHERE id = $1
`;

function add(pool) {
    return {
        setFile: (id, metaFile, cb) => {
            pool.query(
                setFileConf,
                [+id, metaFile, `Auto create from old id ${id}`, metaFile],
                (err) => {
                    if (err)
                        return cb(err);

                    return cb()
                }
            )
        },

        getMetadata: (id, cb) => {
            pool.query(
                getMetadata,
                [+id],
                (err, res) => {
                    if (err)
                        return cb(err);

                    if (res && res.rowCount) {
                        return cb(null, res.rows[0].meta_file)
                    } else {
                        return cb(new CodeError(404, `Not found ${id}`));
                    }
                }
            )
        }
    }
}


module.exports = add;