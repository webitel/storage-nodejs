const elasticService = require(__appRoot + '/services/elastic');

module.exports = {
    getElasticData: getElasticData
};

function getElasticData(webitelUser, options, res, next) {

    return elasticService.search(webitelUser, options, (err, result) => {
        if (err)
            return next(err);

        res.json(result);
    });
}