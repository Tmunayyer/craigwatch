import { spinnerState } from "./components/Spinner.vue";

class API {
    constructor(timeout, tries) {
        this.defaultTimeout = timeout || 3000;
        this.tries = tries || 2;
    }
}

/**
 * fetch performs a single request and returns the body.
 * 
 * @param {string} url 
 * @param {object} options 
 */
API.fetch = async function (url, options) {
    spinnerState.setLoading(true);
    const response = await fetch(url, options);
    const body = await response.json();

    spinnerState.setLoading(false);
    return body;
};

/**
 * fetch_retry will use the shouldRetry function to inspect the body returned from
 *  the initial request. If true, it will begin retrying the
 *  request at increasing intervals. If it returns false, it will return the body.
 * 
 * @param {string} url 
 * @param {object} options 
 * @param {function} cb 
 */
API.fetch_retry = fetch_retry = async function (url, options, shouldRetry) {
    spinnerState.setLoading(true);

    let response;
    let body;

    for (let i = 0; i <= this.tries; i++) {
        await this.sleep(this.defaultTimeout * i);

        response = await fetch(url, options);
        body = await response.json();

        if (shouldRetry(body)) {
            continue;
        } else {
            break;
        }
    }

    spinnerState.setLoading(false);
    return body;
};

API.sleep = function (time) {
    return new Promise(accept => setTimeout(accept, time));
};

export default new API();

