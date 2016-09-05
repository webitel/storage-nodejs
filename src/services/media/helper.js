/**
 * Created by igor on 31.08.16.
 */

"use strict";

const config = require(__appRoot + '/config'),
    providers = config.get('mediaFile:providers'),
    mediaFileConfig = config.get('mediaFile'),
    defProvider = config.get('mediaFile:defaultProvider'),
    mask = config.get('mediaFile:maskPath'),
    maxSize = config.get('mediaFile:maxFieldsSizeMB')
    ;

module.exports = {
    mask: mask,
    maxFileSize: (maxSize || 5) * 1024 * 1024,
    DEFAULT_PROVIDER_NAME: defProvider,
    DEFAULT_PROVIDERS_CONF: mediaFileConfig,
    getProviderConfigByName: (name) => {
        if (providers.hasOwnProperty(name))
            return providers[name];
    }
};