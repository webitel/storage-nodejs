/**
 * Created by igor on 25.08.16.
 */

"use strict";

const config = require(__appRoot + '/config'),
    providers = config.get('recordFile:providers'),
    recordFileConfig = config.get('recordFile'),
    defProvider = config.get('recordFile:defaultProvider'),
    mask = config.get('recordFile:maskPath')
    ;
    
module.exports = {
    mask: mask,
    DEFAULT_PROVIDER_NAME: defProvider,
    DEFAULT_PROVIDERS_CONF: recordFileConfig,
    getProviderConfigByName: (name) => {
        for (let conf of providers)
            if (conf.type === name)
                return conf
    },
    getPath: (domain, fileName) => {
        let date = new Date();
        let path = mask;
        path = domain + '/' + path.replace(/\$Y/g, date.getFullYear()).replace(/\$M/g, (date.getMonth() + 1)).
            replace(/\$D/g, date.getDate()).replace(/\$H/g, date.getHours()).
            replace(/\$m/g, date.getMinutes());

        if (fileName)
            return `${path}/${fileName}`;
        else return path;
    }
};