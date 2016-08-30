/**
 * Created by igor on 29.08.16.
 */

"use strict";

const log = require(__appRoot + '/lib/log')(module),
    nodeMailer = require('nodemailer'),
    SMTPPool = require('nodemailer-smtp-pool'),
    CodeError = require(__appRoot + '/lib/error');

const PROVIDER = {
    smtp: SMTPPool
};
    
let Service = module.exports = {
    send (domain, mailOption, cb) {
        if (!domain)
            return cb(new CodeError(400, 'Domain is required.'));

        if (!mailOption)
            return cb(new CodeError(400, 'Option is required.'));

        application.DB._query.email.getByDomain(domain, (err, res) => {
            if (err)
                return cb(err);

            if (!res || !res.options || !res.provider)
                return cb(new CodeError(404, `Not found domain ${domain} configure email.`));

            if (!(PROVIDER[res.provider] instanceof Function))
                return cb(new CodeError(400, `Bad provider ${res.provider}`));

            mailOption['from'] = mailOption['from'] || res['from'];
            try {
                log.debug(`Send mail ${mailOption.to}`);
                var transporter = nodeMailer.createTransport(PROVIDER[res.provider](res.options));
                transporter.sendMail(
                    mailOption,
                    cb
                );
            } catch (e) {
                return cb(e);
            }
        })
    }
};