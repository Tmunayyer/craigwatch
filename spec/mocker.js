export default {
    // this mocks the api interface used by vue on the $http property.
    // ive wrapped the function in a courier function to allow easy configuration
    // within test files
    api: function (config) {
        return function (url, options) {
            return new Promise((resolve, reject) => {
                if (config.shouldFail) {
                    reject(config.data);
                }

                resolve(config.data);
            });
        };
    }
};
