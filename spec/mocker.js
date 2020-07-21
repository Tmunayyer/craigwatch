class MockAPI {
    constructor(options) {
        this.options = options;
    }

    fetch() {
        return new Promise((resolve, reject) => {
            if (this.options.shouldFail) {
                reject(this.options.data);
            }

            resolve(this.options.data);
        });
    }
}

export default {
    // this mocks the api interface used by vue on the $http property.
    // ive wrapped the function in a courier function to allow easy configuration
    // within test files
    api: MockAPI
};
