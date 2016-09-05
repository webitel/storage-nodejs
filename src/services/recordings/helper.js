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
        if (providers.hasOwnProperty(name))
            return providers[name];
    }
};