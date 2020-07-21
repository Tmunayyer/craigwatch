import API from "../src/api.js";

class FetchSpy {
    constructor() {
        this.callCount = 0;
    }

    fetch = (url, config) => {
        this.callCount++;
        const response = {
            json: function () {
                return new Promise(accept => accept(config.data));
            }
        };

        return new Promise((accept, reject) => {
            if (config.shouldFail) {
                return reject(config.data);
            }

            return accept(response);
        });
    };
}

describe("fetch", () => {
    const api = new API();
    it("should return data", async () => {
        const fetchSpy = new FetchSpy();
        global.fetch = fetchSpy.fetch;
        const config = {
            shouldFail: false,
            data: { testing: 123 }
        };
        const data = await api.fetch("testurl.something", config);

        expect(data.testing).toBe(123);
    });
});

describe("fetch_retry", () => {
    const api = new API({
        defaultTimeout: 100,
        tries: 2
    });
    it("should retry multiple times", async () => {
        const fetchSpy = new FetchSpy();
        global.fetch = fetchSpy.fetch;
        const config = {
            shouldFail: false,
            data: {
                retryme: true
            }
        };

        const shouldRetry = data => data.retryme;
        const data = await api.fetch_retry("testurl.something", config, shouldRetry);
        expect(fetchSpy.callCount).toBe(3);
    });

    it("should stop if shouldRetry func returns false", async () => {
        const fetchSpy = new FetchSpy();
        global.fetch = fetchSpy.fetch;
        const config = {
            shouldFail: false,
            data: {
                retryme: true
            }
        };

        const shouldRetry = () => false;
        const data = await api.fetch_retry("testurl.something", config, shouldRetry);
        expect(fetchSpy.callCount).toBe(1);
    });
});
