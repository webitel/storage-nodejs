/**
 * Created by igor on 25.08.16.
 */

"use strict";

const config = require(__appRoot + '/config'),
    providers = config.get('recordFile:providers'),
    recordFileConfig = config.get('recordFile'),
    defProvider = config.get('recordFile:defaultProvider'),
    mask = config.get('recordFile:maskPath'),
    cronJobDeleteOldFile = `${config.get('recordFile:cronJobDeleteOldFile')}`
    ;
    
module.exports = {
    mask: mask,
    cronJobDeleteOldFile: cronJobDeleteOldFile,
    DEFAULT_PROVIDER_NAME: defProvider,
    DEFAULT_PROVIDERS_CONF: recordFileConfig,
    getProviderConfigByName: (name) => {
        if (providers.hasOwnProperty(name))
            return providers[name];
    }
};